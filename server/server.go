
package server

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/controller"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/logging"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Port                          int
	Logger                        logging.SimpleLogging
	App                           *iris.Application
	LocksController               *controllers.LocksController
	StatusController              *controllers.StatusController
	UsersController               *controllers.UsersController
	AdminController               *controllers.AdminController
	LoginController               *controllers.LoginController

	SSLCertFile                   string
	SSLKeyFile                    string
	SSLPort                       int
	Drainer                       *events.Drainer
	WebAuthentication             bool
	WebUsername                   string
	WebPassword                   string
}

type Config struct {
	AllowForkPRsFlag        string
	StarshipURLFlag         string
	StarshipVersion         string
	DefaultTFVersionFlag    string
	RepoConfigJSONFlag      string
	SilenceForkPRErrorsFlag string
}

func NewServer(userConfig UserConfig, config Config) (*Server, error) {
	logger, err := logging.NewStructuredLoggerFromLevel(userConfig.ToLogLevel())

	if err != nil {
		return nil, err
	}

	drainer := &events.Drainer{}
	db, err := db.NewDB(&db.DBConfig{
		MongoDBConnectionUri: userConfig.MongoDBConnectionUri,
		MongoDBName: userConfig.MongoDBName,
		MongoDBUserName: userConfig.MongoDBUserName,
		MongoDBPassword: userConfig.MongoDBPassword,
		MaxConnection: userConfig.MaxConnection,
		RootCmdLogPath: userConfig.RootCmdLogPath,
		RootSecret: userConfig.RootSecret,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize db.")
	}

	userController := &controllers.UsersController{
		Logger:            logger,
		Drainer:           drainer,
		DB:                db,
	}

	app := iris.New()

	if err != nil {
		return nil, err
	}

	return &Server{
		Port:                          userConfig.Port,
		Logger:                        logger,
		SSLKeyFile:                    userConfig.SSLKeyFile,
		SSLCertFile:                   userConfig.SSLCertFile,
		Drainer:                       drainer,
		UsersController:               userController,
		App:                           app,
	}, nil
}

func (s *Server) ControllersInitialize(){
	apiVer := "/api/v1"
	s.App.Get (apiVer + "/status", s.StatusController.Status)

	s.App.Get (apiVer + "/users/{userId:string}", s.UsersController.Get)
	s.App.Post(apiVer + "/users/create", s.UsersController.Create)
	s.App.Post(apiVer + "/users/create", s.UsersController.Delete)

	s.App.Get (apiVer + "/admin/users", s.AdminController.Users)
	s.App.Post(apiVer + "/login", s.LoginController.Login)
	s.App.Post(apiVer + "/logout", s.LoginController.Logout)
}

func (s *Server) Start() error {
	defer s.Logger.Flush()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	s.ControllersInitialize()

	go func() {
		s.Logger.Info("Starship-IaC started - listening on port %v", s.Port)

		var err error
		s.App.UseGlobal(checkToken)
		if s.SSLCertFile != "" && s.SSLKeyFile != "" {

			err = s.App.Run(iris.TLS(":"+string(s.SSLPort), s.SSLCertFile, s.SSLKeyFile))
		} else {
			err = s.App.Run(iris.Addr(":" + string(s.Port)))
		}

		if err != nil && err != http.ErrServerClosed {
			s.Logger.Err(err.Error())
		}
	}()
	<-stop

	s.Logger.Warn("Received interrupt. Waiting for in-progress operations to complete")
	s.waitForDrain()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint: vet

	if err := s.App.Shutdown(ctx); err != nil {
		return cli.NewExitError(fmt.Sprintf("while shutting down: %s", err), 1)
	}
	return nil
}

// waitForDrain blocks until draining is complete.
func (s *Server) waitForDrain() {
	drainComplete := make(chan bool, 1)
	go func() {
		s.Drainer.ShutdownBlocking()
		drainComplete <- true
	}()
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-drainComplete:
			s.Logger.Info("All in-progress operations complete, shutting down")
			return
		case <-ticker.C:
			s.Logger.Info("Waiting for in-progress operations to complete, current in-progress ops: %d", s.Drainer.GetStatus().InProgressOps)
		}
	}
}

