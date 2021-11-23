
package server

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/runatlantis/atlantis/server/controllers"
	"github.com/runatlantis/atlantis/server/controllers/templates"
	"github.com/runatlantis/atlantis/server/events/yaml"
	"github.com/runatlantis/atlantis/server/events/yaml/valid"
	"github.com/starship-cloud/starship-iac/api"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/logging"
	"github.com/urfave/cli"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Server runs the Atlantis web server.
type Server struct {
	StarshipVersion               string
	StarshipURL                   *url.URL
	Port                          int
	Logger                        logging.SimpleLogging
	GithubAppController           *controllers.GithubAppController
	LocksController               *controllers.LocksController
	StatusController              *controllers.StatusController
	IndexTemplate                 templates.TemplateWriter
	LockDetailTemplate            templates.TemplateWriter
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

	if userConfig.GithubUser != "" || userConfig.GithubAppID != 0 {
		if userConfig.GithubUser != "" {
			//githubCredentials = &vcs.GithubUserCredentials{
			//	User:  userConfig.GithubUser,
			//	Token: userConfig.GithubToken,
			//}
		} else if userConfig.GithubAppID != 0 && userConfig.GithubAppKeyFile != "" {
		} else if userConfig.GithubAppID != 0 && userConfig.GithubAppKey != "" {
		}

		//var err error
		//githubClient, err = vcs.NewGithubClient(userConfig.GithubHostname, githubCredentials, logger, userConfig.VCSStatusName)
		//if err != nil {
		//	return nil, err
		//}
	}

	if userConfig.WriteGitCreds {
		home, err := homedir.Dir()
		if err != nil {
			return nil, errors.Wrap(err, "getting home dir to write ~/.git-credentials file")
		}
		if userConfig.GithubUser != "" {
			//if err := events.WriteGitCreds(userConfig.GithubUser, userConfig.GithubToken, userConfig.GithubHostname, home, logger, false); err != nil {
			//	return nil, err
			//}
			fmt.Println(home)
		}
		if userConfig.GitlabUser != "" {
		}
	}

	parsedURL, err := ParseURL(userConfig.StarshipURL)
	if err != nil {
		return nil, errors.Wrapf(err,
			"parsing --%s flag %q", config.StarshipURLFlag, userConfig.StarshipURL)
	}
	validator := &yaml.ParserValidator{}

	globalCfg := valid.NewGlobalCfgFromArgs(
		valid.GlobalCfgArgs{
			AllowRepoCfg:       userConfig.AllowRepoConfig,
			MergeableReq:       userConfig.RequireMergeable,
			ApprovedReq:        userConfig.RequireApproval,
			UnDivergedReq:      userConfig.RequireUnDiverged,
			PolicyCheckEnabled: userConfig.EnablePolicyChecksFlag,
		})
	if userConfig.RepoConfig != "" {
		globalCfg, err = validator.ParseGlobalCfg(userConfig.RepoConfig, globalCfg)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing %s file", userConfig.RepoConfig)
		}
	} else if userConfig.RepoConfigJSON != "" {
		globalCfg, err = validator.ParseGlobalCfgJSON(userConfig.RepoConfigJSON, globalCfg)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing --%s", config.RepoConfigJSONFlag)
		}
	}

	drainer := &events.Drainer{}

	return &Server{
		StarshipVersion:               config.StarshipVersion,
		StarshipURL:                   parsedURL,
		Port:                          userConfig.Port,
		Logger:                        logger,
		IndexTemplate:                 templates.IndexTemplate,
		LockDetailTemplate:            templates.LockTemplate,
		SSLKeyFile:                    userConfig.SSLKeyFile,
		SSLCertFile:                   userConfig.SSLCertFile,
		Drainer:                       drainer,
		WebAuthentication:             userConfig.WebBasicAuth,
		WebUsername:                   userConfig.WebUsername,
		WebPassword:                   userConfig.WebPassword,
	}, nil
}

func (s *Server) Start() error {
	defer s.Logger.Flush()

	// Ensure server gracefully drains connections when stopped.
	stop := make(chan os.Signal, 1)
	// Stop on SIGINTs and SIGTERMs.
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	app := api.Init()

	go func() {
		s.Logger.Info("Starship-IaC started - listening on port %v", s.Port)

		var err error
		if s.SSLCertFile != "" && s.SSLKeyFile != "" {
			err = app.Run(iris.TLS(":" + string(s.SSLPort), s.SSLCertFile, s.SSLKeyFile))
		} else {
			err = app.Run(iris.Addr(":" + string(s.Port)))
		}

		if err != nil && err != http.ErrServerClosed {
			s.Logger.Err(err.Error())
		}
	}()
	<-stop

	s.Logger.Warn("Received interrupt. Waiting for in-progress operations to complete")
	s.waitForDrain()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint: vet

	if err := app.Shutdown(ctx); err != nil {
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

// Healthz returns the health check response. It always returns a 200 currently.
//func (s *Server) Healthz(w http.ResponseWriter, _ *http.Request) {
//	data, err := json.MarshalIndent(&struct {
//		Status string `json:"status"`
//	}{
//		Status: "ok",
//	}, "", "  ")
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprintf(w, "Error creating status json response: %s", err)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(data) // nolint: errcheck
//}

// ParseURL parses the user-passed atlantis URL to ensure it is valid
// and we can use it in our templates.
// It removes any trailing slashes from the path so we can concatenate it
// with other paths without checking.
func ParseURL(u string) (*url.URL, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	if !(parsed.Scheme == "http" || parsed.Scheme == "https") {
		return nil, errors.New("http or https must be specified")
	}
	// We want the path to end without a trailing slash so we know how to
	// use it in the rest of the program.
	parsed.Path = strings.TrimSuffix(parsed.Path, "/")
	return parsed, nil
}
