package runner

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/boasihq/interactive-inputs/internal/config"
	"github.com/boasihq/interactive-inputs/internal/errors"
	"github.com/boasihq/interactive-inputs/internal/notifier"
	"github.com/boasihq/interactive-inputs/internal/portal"
	webui "github.com/boasihq/interactive-inputs/internal/web"
	"github.com/gorilla/mux"
	"golang.ngrok.com/ngrok"
	nconfig "golang.ngrok.com/ngrok/config"
)

func InvokeAction(ctx context.Context, ctxCancel context.CancelFunc, cfg *config.Config, embeddedContent fs.FS, embeddedContentFilePathPrefix string) error {

	defer ctxCancel()

	var githubActionWorkingDir string = os.Getenv("GITHUB_WORKSPACE")
	var isRunningLocal bool = os.Getenv("IAIP_LOCAL_RUN") != ""

	if githubActionWorkingDir == "" {
		cfg.Action.Errorf("GITHUB_WORKSPACE not found")
		return errors.ErrGitHubWorkspaceEnvVarIsMissing
	}

	// TODO: Get the source job's url that's calling the
	// action so that we can link to it. and send users
	// back to it

	/// Notifiers
	slackNotifier := notifier.NewSlackNotifier(&notifier.NewSlackNotifierRequest{})
	discordNotifier := notifier.NewDiscordNotifier(&notifier.NewDiscordNotifierRequest{})

	if cfg.NotifierSlackEnabled {
		slackNotifier = notifier.NewSlackNotifier(&notifier.NewSlackNotifierRequest{
			Enabled:   cfg.NotifierSlackEnabled,
			Token:     cfg.NotifierSlackToken,
			Channel:   cfg.NotifierSlackChannel,
			BotName:   cfg.NotifierSlackBotName,
			ActionPkg: cfg.Action,
			ThreadTs:  cfg.NotifierSlackThreadTs,
		})

		verifiedSlackNotifierErr := slackNotifier.Verify()
		if verifiedSlackNotifierErr != nil {
			cfg.Action.Errorf("Slack Notifier Verification Failed")
			return verifiedSlackNotifierErr
		}

		cfg.Action.Debugf("Slack Notifier Verification Succeeded")
	}

	if cfg.NotifierDiscordEnabled {
		discordNotifier = notifier.NewDiscordNotifier(&notifier.NewDiscordNotifierRequest{
			Enabled:          cfg.NotifierDiscordEnabled,
			WebhookUrl:       cfg.NotifierDiscordWebhook,
			UsernameOverride: cfg.NotifierDiscordUsernameOverride,
			ActionPkg:        cfg.Action,
			ThreadId:         cfg.NotifierDiscordThreadId,
		})

		verifiedDiscordNotifierErr := discordNotifier.Verify()
		if verifiedDiscordNotifierErr != nil {
			cfg.Action.Errorf("Discord Notifier Verification Failed")
			return verifiedDiscordNotifierErr
		}
		cfg.Action.Debugf("Discord Notifier Verification Succeeded")
	}

	/// Handlers
	uiHandler := webui.NewWebAppHandler(&webui.NewWebAppHandlerRequest{
		EmbeddedContent:               embeddedContent,
		EmbeddedContentFilePathPrefix: embeddedContentFilePathPrefix,
		Config:                        cfg,
	})

	portalEventHandler := portal.NewHandler(cfg.Action, isRunningLocal, embeddedContent, embeddedContentFilePathPrefix, cfg.GithubToken)

	/// Routes
	r := mux.NewRouter()

	portal.AttachRoutes(&portal.AttachRoutesRequest{
		Router:                        r,
		PortalEventHandler:            portalEventHandler,
		UiHandler:                     uiHandler,
		EmbeddedContent:               embeddedContent,
		EmbeddedContentFilePathPrefix: embeddedContentFilePathPrefix,
		ActionPkg:                     cfg.Action,
	})

	/// Server
	serverDone := make(chan error, 1)
	serverInitMessageTmpl := "Your Interactive Inputs portal is reachable at: %s"
	notifierSlackEnterInputMessageTmpl := "<%s|*Enter required input*>"
	notifierDiscordEnterInputMessageTmpl := "[**Enter required input**](%s)"
	universalNotifierFailedToSelfHost := "A failure has occurred while starting/running your self-hosted portal: %v"

	// TODO: Add a flag to enable/disable the ngrok tunnel respsective
	// of whether the action is running locally or not
	if !isRunningLocal {
		ln, err := ngrok.Listen(ctx,
			nconfig.HTTPEndpoint(),
			ngrok.WithAuthtoken(cfg.NgrokAuthtoken),
		)
		if err != nil {
			return err
		}

		serverInitMessage := fmt.Sprintf(serverInitMessageTmpl, ln.URL())

		cfg.Action.Noticef(serverInitMessage)

		if slackNotifier.Enabled() {
			_, err := slackNotifier.Notify(cfg.Title, fmt.Sprintf(notifierSlackEnterInputMessageTmpl, ln.URL()))
			if err != nil {
				cfg.Action.Errorf("Slack Notifier Notification Failed: %v", err)
				return err
			}
		}

		if discordNotifier.Enabled() {
			_, err := discordNotifier.Notify(cfg.Title, fmt.Sprintf(notifierDiscordEnterInputMessageTmpl, ln.URL()))
			if err != nil {
				cfg.Action.Errorf("Discord Notifier Notification Failed: %v", err)
				return err
			}
		}

		go func() {
			// server logic
			if err := http.Serve(ln, r); err != nil {
				serverErrorMessage := fmt.Sprintf(universalNotifierFailedToSelfHost, err)

				cfg.Action.Errorf(serverErrorMessage)
				if slackNotifier.Enabled() {
					_, err := slackNotifier.Notify(cfg.Title, serverErrorMessage)
					if err != nil {
						cfg.Action.Errorf("Slack Notifier Notification Failed: %v", err)
					}
				}

				if discordNotifier.Enabled() {
					_, err := discordNotifier.Notify(cfg.Title, serverErrorMessage)
					if err != nil {
						cfg.Action.Errorf("Discord Notifier Notification Failed: %v", err)
					}
				}

				serverDone <- err
			}
			serverDone <- ln.CloseWithContext(ctx)
		}()

	} else {
		localPort := ":8080"
		server := &http.Server{Addr: localPort, Handler: r}
		completeLocalUrl := fmt.Sprintf("http://localhost%s", localPort)
		serverInitMessage := fmt.Sprintf(serverInitMessageTmpl, completeLocalUrl)

		cfg.Action.Noticef(serverInitMessage)
		if slackNotifier.Enabled() {
			_, err := slackNotifier.Notify(cfg.Title, fmt.Sprintf(notifierSlackEnterInputMessageTmpl, completeLocalUrl))
			if err != nil {
				cfg.Action.Errorf("Slack Notifier Notification Failed: %v", err)
				return err
			}
		}

		if discordNotifier.Enabled() {
			_, err := discordNotifier.Notify(cfg.Title, fmt.Sprintf(notifierDiscordEnterInputMessageTmpl, completeLocalUrl))
			if err != nil {
				cfg.Action.Errorf("Discord Notifier Notification Failed: %v", err)
				return err
			}
		}

		go func() {
			// server logic
			if err := server.ListenAndServe(); err != nil {
				serverErrorMessage := fmt.Sprintf(universalNotifierFailedToSelfHost, err)

				cfg.Action.Errorf(serverErrorMessage)
				if slackNotifier.Enabled() {
					_, err := slackNotifier.Notify(cfg.Title, serverErrorMessage)
					if err != nil {
						cfg.Action.Errorf("Slack Notifier Notification Failed: %v", err)
					}
				}

				if discordNotifier.Enabled() {
					_, err := discordNotifier.Notify(cfg.Title, serverErrorMessage)
					if err != nil {
						cfg.Action.Errorf("Discord Notifier Notification Failed: %v", err)
					}
				}

				serverDone <- err
			}
			serverDone <- server.Shutdown(ctx)
		}()
	}

	select {
	case err := <-serverDone:
		return handlePrettierTimeoutErrorMessage(err, cfg.Timeout)
	case <-ctx.Done():
		// Timeout occurred
		ctxCancel() // Ensure all resources are cleaned up

		return handlePrettierTimeoutErrorMessage(ctx.Err(), cfg.Timeout)
	}

}

// handlePrettierTimeoutErrorMessage is a helper function that prints a nicer error message
// when the context deadline is exceeded. Otherwise, it returns the original error.
func handlePrettierTimeoutErrorMessage(err error, timeout int) error {
	// Print nicer timeout message
	if err != nil && err.Error() == "context deadline exceeded" {
		//nolint:go-staticcheck
		return fmt.Errorf("Your session has expired (timed out) due to inactivity for %d seconds", timeout)
	}

	return err
}
