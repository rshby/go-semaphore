package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-semaphore/config"
	"go-semaphore/internal/database"
	"go-semaphore/internal/middleware"
	"go-semaphore/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var runServer = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  `This subcommand start the server`,
	Run:   server,
}

func init() {
	RootCMD.AddCommand(runServer)
}

func server(cmd *cobra.Command, args []string) {
	var appPort = config.AppPort()
	if len(args) > 0 {
		port, err := strconv.Atoi(args[0])
		if err != nil {
			logrus.Fatalf("cant parse port to int : %v", err)
		}

		appPort = port
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// connect to database
	db, err := database.InitializeMySqlConnection()
	if err != nil {
		logrus.Fatal(err)
	}

	// init service
	customerService := InitCustomerService(db)

	// create grpc server
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.WithErrorInterceptor()))

	// register service
	pb.RegisterCustomerServiceServer(grpcServer, customerService)

	reflection.Register(grpcServer)

	// create channel
	sigChan := make(chan os.Signal, 1)
	chanErr := make(chan error, 1)
	quitChan := make(chan struct{}, 1)

	// notify signal intterupt
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-sigChan:
				logrus.Infof("receive interupt signal")
				gracefullShutdown(grpcServer)
				quitChan <- struct{}{}
				return
			case err := <-chanErr:
				logrus.Infof("receive error signal")
				logrus.Error(err)
				gracefullShutdown(grpcServer)
				quitChan <- struct{}{}
				return
			}
		}
	}()

	// spawn a goroutine to run the gRPC server, allowing it to be stopped gracefully
	go func() {
		var err error
		listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.AppHost(), appPort))
		if err != nil {
			chanErr <- err
			return
		}

		logrus.Infof("running grpc server on [%s:%d] ⏳", config.AppHost(), appPort)
		if err = grpcServer.Serve(listen); err != nil {
			chanErr <- err
			return
		}
	}()

	_ = <-quitChan
	close(quitChan)
	close(chanErr)
	close(sigChan)

	logrus.Info("grpc server exiting ⚠️")
}

// gracefullShutdown is function to stop grpc server with gracefull or force stop
func gracefullShutdown(grpcServer *grpc.Server) {
	if grpcServer != nil {
		var (
			wg   = &sync.WaitGroup{}
			done = make(chan struct{}, 1)
		)

		// spawn a goroutine : close server with gracefull
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			grpcServer.GracefulStop()
			done <- struct{}{}
			close(done)
		}(wg)

		// spawn a goroutine : force close if more than 5 seconds
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			for {
				select {
				case <-done:
					logrus.Infof("stop server with gracefull ❎")
					return
				case <-time.After(5 * time.Second):
					grpcServer.Stop()
					logrus.Infof("force stop server 🚨")
					return
				}
			}
		}(wg)

		wg.Wait()
	}
}
