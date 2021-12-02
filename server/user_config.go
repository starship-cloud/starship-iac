package server

import (
	"github.com/starship-cloud/starship-iac/server/logging"
)

type UserConfig struct {
	LogLevel      string `mapstructure:"log-level"`
	Port          int    `mapstructure:"port"`
	SSLPort       int    `mapstructure:"ssl-port"`
	SSLCertFile   string `mapstructure:"ssl-cert-file"`
	SSLKeyFile    string `mapstructure:"ssl-key-file"`
	SkipAuthToken bool   `mapstructure:"skip-auth-token"`

	MongoDBConnectionUri string `mapstructure:"mongodb-connection-uri"`
	MongoDBName          string `mapstructure:"mongodb-name"`
	MongoDBUserName      string `mapstructure:"mongodb-username"`
	MongoDBPassword      string `mapstructure:"mongodb-password"`
	MaxConnection        int    `mapstructure:"mongodb-max-connection"`
	RootCmdLogPath       string `mapstructure:"mongodb-root-cmd-logpath"`
	RootSecret           string `mapstructure:"mongodb-root-secret"`
}

func (u UserConfig) ToLogLevel() logging.LogLevel {
	switch u.LogLevel {
	case "debug":
		return logging.Debug
	case "info":
		return logging.Info
	case "warn":
		return logging.Warn
	case "error":
		return logging.Error
	}
	return logging.Info
}
