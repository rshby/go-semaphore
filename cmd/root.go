package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-semaphore/config"
	"go-semaphore/internal/logger"
	"os"
)

var RootCMD = &cobra.Command{
	Use:   "go-semaphore",
	Short: "go-semaphore console",
	Long:  "this is go-semaphore console",
}

func init() {
	// load logrus at first
	logger.SetupLogger()

	// load config at first
	config.LoadConfig()
}

// Execute :nodoc:
func Execute() {
	if err := RootCMD.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
