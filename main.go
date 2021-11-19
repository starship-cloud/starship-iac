package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/starship-cloud/starship-iac/cmd"
	"github.com/starship-cloud/starship-iac/server/logging"
)

const starshipVersion = "0.0.1"

func main() {
	v := viper.New()

	logger, err := logging.NewStructuredLogger()

	if err != nil {
		panic(fmt.Sprintf("unable to initialize logger. %s", err.Error()))
	}

	// We're creating commands manually here rather than using init() functions
	// (as recommended by cobra) because it makes testing easier.
	server := &cmd.ServerCmd{
		ServerCreator:   &cmd.DefaultServerCreator{},
		Viper:           v,
		StarshipVersion: starshipVersion,
		Logger:          logger,
	}
	version := &cmd.VersionCmd{StarshipVersion: starshipVersion}
	cmd.RootCmd.AddCommand(server.Init())
	cmd.RootCmd.AddCommand(version.Init())
	cmd.Execute()
}
