
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/runatlantis/atlantis/server/events/yaml/valid"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/runatlantis/atlantis/server/controllers"
	events_controllers "github.com/runatlantis/atlantis/server/controllers/events"
	"github.com/runatlantis/atlantis/server/controllers/templates"
	"github.com/runatlantis/atlantis/server/core/locking"
	"github.com/runatlantis/atlantis/server/events/yaml"
	"github.com/runatlantis/atlantis/server/static"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/logging"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

const (
	// LockViewRouteName is the named route in mux.Router for the lock view.
	// The route can be retrieved by this name, ex:
	//   mux.Router.Get(LockViewRouteName)
	LockViewRouteName = "lock-detail"
	// LockViewRouteIDQueryParam is the query parameter needed to construct the lock view
	// route. ex:
	//   mux.Router.Get(LockViewRouteName).URL(LockViewRouteIDQueryParam, "my id")
	LockViewRouteIDQueryParam = "id"

	// binDirName is the name of the directory inside our data dir where
	// we download binaries.
	BinDirName = "bin"

	// terraformPluginCacheDir is the name of the dir inside our data dir
	// where we tell terraform to cache plugins and modules.
	TerraformPluginCacheDirName = "plugin-cache"
)

// Server runs the Atlantis web server.
type Server struct {
	AtlantisVersion               string
	AtlantisURL                   *url.URL
	Router                        *mux.Router
	Port                          int
	//PreWorkflowHooksCommandRunner *events.DefaultPreWorkflowHooksCommandRunner
	//CommandRunner                 *events.DefaultCommandRunner
	Logger                        logging.SimpleLogging
	Locker                        locking.Locker
	ApplyLocker                   locking.ApplyLocker
	VCSEventsController           *events_controllers.VCSEventsController
	GithubAppController           *controllers.GithubAppController
	LocksController               *controllers.LocksController
	StatusController              *controllers.StatusController
	IndexTemplate                 templates.TemplateWriter
	LockDetailTemplate            templates.TemplateWriter
	SSLCertFile                   string
	SSLKeyFile                    string
	Drainer                       *events.Drainer
	WebAuthentication             bool
	WebUsername                   string
	WebPassword                   string
}

// Config holds config for server that isn't passed in by the user.
type Config struct {
	AllowForkPRsFlag        string
	AtlantisURLFlag         string
	AtlantisVersion         string
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
			/*privateKey, err := os.ReadFile(userConfig.GithubAppKeyFile)
			if err != nil {
				return nil, err
			}
			githubCredentials = &vcs.GithubAppCredentials{
				AppID:    userConfig.GithubAppID,
				Key:      privateKey,
				Hostname: userConfig.GithubHostname,
				AppSlug:  userConfig.GithubAppSlug,
			}*/
			//githubAppEnabled = true
		} else if userConfig.GithubAppID != 0 && userConfig.GithubAppKey != "" {
			//githubCredentials = &vcs.GithubAppCredentials{
			//	AppID:    userConfig.GithubAppID,
			//	Key:      []byte(userConfig.GithubAppKey),
			//	Hostname: userConfig.GithubHostname,
			//	AppSlug:  userConfig.GithubAppSlug,
			//}
			//githubAppEnabled = true
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
			//if err := events.WriteGitCreds(userConfig.GitlabUser, userConfig.GitlabToken, userConfig.GitlabHostname, home, logger, false); err != nil {
			//	return nil, err
			//}
		}
	}

	parsedURL, err := ParseAtlantisURL(userConfig.AtlantisURL)
	if err != nil {
		return nil, errors.Wrapf(err,
			"parsing --%s flag %q", config.AtlantisURLFlag, userConfig.AtlantisURL)
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

	underlyingRouter := mux.NewRouter()

	drainer := &events.Drainer{}

	return &Server{
		AtlantisVersion:               config.AtlantisVersion,
		AtlantisURL:                   parsedURL,
		Router:                        underlyingRouter,
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

// Start creates the routes and starts serving traffic.
func (s *Server) Start() error {
	s.Router.HandleFunc("/", s.Index).Methods("GET").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.URL.Path == "/" || r.URL.Path == "/index.html"
	})
	s.Router.HandleFunc("/healthz", s.Healthz).Methods("GET")
	s.Router.HandleFunc("/status", s.StatusController.Get).Methods("GET")
	s.Router.PathPrefix("/static/").Handler(http.FileServer(&assetfs.AssetFS{Asset: static.Asset, AssetDir: static.AssetDir, AssetInfo: static.AssetInfo}))
	s.Router.HandleFunc("/events", s.VCSEventsController.Post).Methods("POST")
	s.Router.HandleFunc("/github-app/exchange-code", s.GithubAppController.ExchangeCode).Methods("GET")
	s.Router.HandleFunc("/github-app/setup", s.GithubAppController.New).Methods("GET")
	s.Router.HandleFunc("/apply/lock", s.LocksController.LockApply).Methods("POST").Queries()
	s.Router.HandleFunc("/apply/unlock", s.LocksController.UnlockApply).Methods("DELETE").Queries()
	s.Router.HandleFunc("/locks", s.LocksController.DeleteLock).Methods("DELETE").Queries("id", "{id:.*}")
	s.Router.HandleFunc("/lock", s.LocksController.GetLock).Methods("GET").
		Queries(LockViewRouteIDQueryParam, fmt.Sprintf("{%s}", LockViewRouteIDQueryParam)).Name(LockViewRouteName)
	n := negroni.New(&negroni.Recovery{
		Logger:     log.New(os.Stdout, "", log.LstdFlags),
		PrintStack: false,
		StackAll:   false,
		StackSize:  1024 * 8,
	}, NewRequestLogger(s))
	n.UseHandler(s.Router)

	defer s.Logger.Flush()

	// Ensure server gracefully drains connections when stopped.
	stop := make(chan os.Signal, 1)
	// Stop on SIGINTs and SIGTERMs.
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{Addr: fmt.Sprintf(":%d", s.Port), Handler: n}
	go func() {
		s.Logger.Info("Atlantis started - listening on port %v", s.Port)

		var err error
		if s.SSLCertFile != "" && s.SSLKeyFile != "" {
			err = server.ListenAndServeTLS(s.SSLCertFile, s.SSLKeyFile)
		} else {
			err = server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			s.Logger.Err(err.Error())
		}
	}()
	<-stop

	s.Logger.Warn("Received interrupt. Waiting for in-progress operations to complete")
	s.waitForDrain()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint: vet
	if err := server.Shutdown(ctx); err != nil {
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

// Index is the / route.
func (s *Server) Index(w http.ResponseWriter, _ *http.Request) {
	locks, err := s.Locker.List()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Could not retrieve locks: %s", err)
		return
	}

	var lockResults []templates.LockIndexData
	for id, v := range locks {
		lockURL, _ := s.Router.Get(LockViewRouteName).URL("id", url.QueryEscape(id))
		lockResults = append(lockResults, templates.LockIndexData{
			// NOTE: must use .String() instead of .Path because we need the
			// query params as part of the lock URL.
			LockPath:      lockURL.String(),
			RepoFullName:  v.Project.RepoFullName,
			PullNum:       v.Pull.Num,
			Path:          v.Project.Path,
			Workspace:     v.Workspace,
			Time:          v.Time,
			TimeFormatted: v.Time.Format("02-01-2006 15:04:05"),
		})
	}

	applyCmdLock, err := s.ApplyLocker.CheckApplyLock()
	s.Logger.Info("Apply Lock: %v", applyCmdLock)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Could not retrieve global apply lock: %s", err)
		return
	}

	applyLockData := templates.ApplyLockData{
		Time:          applyCmdLock.Time,
		Locked:        applyCmdLock.Locked,
		TimeFormatted: applyCmdLock.Time.Format("02-01-2006 15:04:05"),
	}
	//Sort by date - newest to oldest.
	sort.SliceStable(lockResults, func(i, j int) bool { return lockResults[i].Time.After(lockResults[j].Time) })

	err = s.IndexTemplate.Execute(w, templates.IndexData{
		Locks:           lockResults,
		ApplyLock:       applyLockData,
		AtlantisVersion: s.AtlantisVersion,
		CleanedBasePath: s.AtlantisURL.Path,
	})
	if err != nil {
		s.Logger.Err(err.Error())
	}
}

func mkSubDir(parentDir string, subDir string) (string, error) {
	fullDir := filepath.Join(parentDir, subDir)
	if err := os.MkdirAll(fullDir, 0700); err != nil {
		return "", errors.Wrapf(err, "unable to create dir %q", fullDir)
	}

	return fullDir, nil
}

// Healthz returns the health check response. It always returns a 200 currently.
func (s *Server) Healthz(w http.ResponseWriter, _ *http.Request) {
	data, err := json.MarshalIndent(&struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating status json response: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data) // nolint: errcheck
}

// ParseAtlantisURL parses the user-passed atlantis URL to ensure it is valid
// and we can use it in our templates.
// It removes any trailing slashes from the path so we can concatenate it
// with other paths without checking.
func ParseAtlantisURL(u string) (*url.URL, error) {
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
