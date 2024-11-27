package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-semaphore/internal/database"
	"strings"
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
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// connect to database
	_, err := database.InitializeMySqlConnection()
	logrus.Fatal(err)
}
