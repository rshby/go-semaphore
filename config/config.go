package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

// LoadConfig is method to load config file from yml
func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("failed to load config : %s", err.Error())
	}

	// iterate all config keys and set them as env variable
	for _, k := range viper.AllKeys() {
		key := strings.ToUpper(strings.Replace(k, ".", "_", -1))
		_ = os.Setenv(key, viper.GetString(k))
	}
}

func AppPort() int {
	return viper.GetInt("app.port")
}

func MysqlHost() string {
	return viper.GetString("mysql.host")
}

func MysqlPort() int {
	return viper.GetInt("mysql.port")
}

func MysqlUser() string {
	return viper.GetString("mysql.user")
}

func MysqlPassword() string {
	return viper.GetString("mysql.password")
}

func MysqlName() string {
	return viper.GetString("mysql.name")
}

func MysqlTimezone() string {
	return viper.GetString("mysql.timezone")
}

func MysqlMaxIdleConns() int {
	return viper.GetInt("mysql.max_idle_conns")
}

func MysqlMaxOpenConns() int {
	return viper.GetInt("mysql.max_open_conns")
}

func MysqlConnMaxLifetime() time.Duration {
	maxLifetime := viper.GetString("mysql.conn_max_lifetime")
	duration, err := time.ParseDuration(maxLifetime)
	if err != nil {
		logrus.Error(err)
		return time.Duration(1 * time.Hour)
	}

	return duration
}

func MysqlDSN() string {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s`,
		MysqlUser(),
		MysqlPassword(),
		MysqlHost(),
		MysqlPort(),
		MysqlName(),
		MysqlTimezone(),
	)
	return dsn
}
