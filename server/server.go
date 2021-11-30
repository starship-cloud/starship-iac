
package server

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/runatlantis/atlantis/server/controllers/templates"
	"github.com/starship-cloud/starship-iac/server/controller"
	//"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/core/locking"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/logging"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server runs the Atlantis web server.
type Server struct {
	Port                          int
	Logger                        logging.SimpleLogging
	//GithubAppController           *controllers.GithubAppController
	//LocksController               *controllers.LocksController

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

// Config holds config for server that isn't passed in by the user.
type Config struct {
	AllowForkPRsFlag        string
	StarshipURLFlag         string
	StarshipVersion         string
	DefaultTFVersionFlag    string
	RepoConfigJSONFlag      string
	SilenceForkPRErrorsFlag string
}

// WebhookConfig is nested within UserConfig. It's used to configure webhooks.
type WebhookConfig struct {
	// Event is the type of event we should send this webhook for, ex. apply.
	Event string `mapstructure:"event"`
	// WorkspaceRegex is a regex that is used to match against the workspace
	// that is being modified for this event. If the regex matches, we'll
	// send the webhook, ex. "production.*".
	WorkspaceRegex string `mapstructure:"workspace-regex"`
	// Kind is the type of webhook we should send, ex. slack.
	Kind string `mapstructure:"kind"`
	// Channel is the channel to send this webhook to. It only applies to
	// slack webhooks. Should be without '#'.
	Channel string `mapstructure:"channel"`
}

// NewServer returns a new server. If there are issues starting the server or
// its dependencies an error will be returned. This is like the main() function
// for the server CLI command because it injects all the dependencies.
func NewServer(userConfig UserConfig, config Config) (*Server, error) {
	logger, err := logging.NewStructuredLoggerFromLevel(userConfig.ToLogLevel())

	if err != nil {
		return nil, err
	}

	drainer := &events.Drainer{}
	db, err := db.NewDB(userConfig.DBConfig)

	if err != nil {
		return nil, err
	}
	var lockingClient locking.Locker
	var applyLockingClient locking.ApplyLocker
	if userConfig.DisableRepoLocking {
		lockingClient = locking.NewNoOpLocker()
	} else {
		//lockingClient = locking.NewClient(boltdb)
	}
	//applyLockingClient = locking.NewApplyClient(boltdb, userConfig.DisableApply)
	workingDirLocker := events.NewDefaultWorkingDirLocker()

	var workingDir events.WorkingDir = &events.FileWorkspace{
		DataDir:       userConfig.DataDir,
		CheckoutMerge: userConfig.CheckoutStrategy == "merge",
	}

	deleteLockCommand := &events.DefaultDeleteLockCommand{
		Locker:           lockingClient,
		Logger:           logger,
		WorkingDir:       workingDir,
		WorkingDirLocker: workingDirLocker,
		DB:               db,
	}

	locksController := &controllers.LocksController{
		Locker:             lockingClient,
		ApplyLocker:        applyLockingClient,
		Logger:             logger,
		LockDetailTemplate: templates.LockTemplate,
		WorkingDir:         workingDir,
		WorkingDirLocker:   workingDirLocker,
		DB:                 db,
		DeleteLockCommand:  deleteLockCommand,
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
		WebAuthentication:             userConfig.WebBasicAuth,
		WebUsername:                   userConfig.WebUsername,
		WebPassword:                   userConfig.WebPassword,
		LocksController:               locksController,
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

	// Ensure server gracefully drains connections when stopped.
	stop := make(chan os.Signal, 1)
	// Stop on SIGINTs and SIGTERMs.
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

