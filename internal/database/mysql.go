package database

import (
	"github.com/sirupsen/logrus"
	"go-semaphore/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitializeMySqlConnection is function to initialize connection with database mysql
func InitializeMySqlConnection() (*gorm.DB, error) {
	db, err := openMySqlConnection(config.MysqlDSN())
	if err != nil {
		logrus.Fatalf("failed connect to database mysql : %s", err.Error())
	}

	return db, nil
}

// openMySqlConnection is function to connect mysql with gorm
func openMySqlConnection(dsn string) (*gorm.DB, error) {
	logrus.Infof("database MYSQL DSN : [%s]", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	conn, err := db.DB()
	if err != nil {
		return nil, err
	}

	// set database pooling
	conn.SetMaxOpenConns(config.MysqlMaxOpenConns())
	conn.SetMaxIdleConns(config.MysqlMaxIdleConns())
	conn.SetConnMaxLifetime(config.MysqlConnMaxLifetime())

	logrus.Infof("success connected to database MYSQL [%s:%d]", config.MysqlHost(), config.MysqlPort())

	return db, nil
}
