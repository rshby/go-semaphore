package cmd

import (
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-semaphore/config"
	"go-semaphore/internal/database"
	"gorm.io/gorm"
	"strconv"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: processMigration,
}

func init() {
	migrateCmd.PersistentFlags().Int("step", 0, "maximum migration steps")
	migrateCmd.PersistentFlags().String("direction", "up", "migration direction")
	RootCMD.AddCommand(migrateCmd)
}

// processMigration
func processMigration(cmd *cobra.Command, args []string) {
	logrus.Info("Process migration!")

	direction := cmd.Flag("direction").Value.String()
	stepStr := cmd.Flag("step").Value.String()
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		logrus.WithField("stepStr", stepStr).Fatal("Failed to parse step to int: ", err)
	}

	db, err := database.InitializeMySqlConnection()
	if err != nil {
		logrus.Fatalf("failed to connect database mysql : %s", err.Error())
	}

	migration(db, direction, step)
}

func migration(db *gorm.DB, direction string, step int) {
	var (
		n                  int
		migrationDirection = migrate.Up
	)

	migrations := &migrate.FileMigrationSource{
		Dir: "internal/database/migrations/",
	}

	migrate.SetTable("migrations")

	dbMysql, err := db.DB()
	defer dbMysql.Close()

	if err != nil {
		logrus.WithField("DatabaseDSN", config.MysqlDSN()).Fatal("Failed to connect database: ", err)
	}

	if direction == "down" {
		migrationDirection = migrate.Down
	}

	n, err = migrate.ExecMax(dbMysql, "mysql", migrations, migrationDirection, step)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"direction": direction,
		}).Fatal("Failed to migrate database: ", err)
	}

	logrus.Infof("Applied %d migrations!\n", n)
}
