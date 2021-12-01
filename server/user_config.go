package server

import (
	"github.com/starship-cloud/starship-iac/server/logging"
)

type UserConfig struct {
	LogLevel                   string `mapstructure:"log-level"`
	Port                       int    `mapstructure:"port"`
	SSLPort                    int    `mapstructure:"ssl-port"`
	SSLCertFile                string          `mapstructure:"ssl-cert-file"`
	SSLKeyFile                 string          `mapstructure:"ssl-key-file"`

	MongoDBConnectionUri string `mapstructure:"mongodburi"`
	MongoDBName          string `mapstructure:"mongodbname"`
	MongoDBUserName      string `mapstructure:"mongodbusername"`
	MongoDBPassword      string `mapstructure:"mongodbpassword"`
	MaxConnection        int    `mapstructure:"maxconnection"`
	RootCmdLogPath       string `mapstructure:"rootcmdlogpath"`
	RootSecret           string `mapstructure:"rootsecret"`
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
