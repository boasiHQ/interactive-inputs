package config

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/boasihq/interactive-inputs/internal/errors"
    "github.com/boasihq/interactive-inputs/internal/fields"
    githubactions "github.com/sethvargo/go-githubactions"
)

type Config struct {

	// Title is the header that will be displayed at the top of the generated form
	Title string

	// Fields is the slice of fields that will be displayed in the generated form
	Fields *fields.Fields

	// Timeout is the timeout that will be used to manage how long the portal
	// will be available for users to use before it is automatically deactivated
	Timeout int

	// PortalHostMode determines how the portal is exposed (ngrok tunnel vs self-hosted)
	PortalHostMode string

	// SelfHostedListenAddress is the address the HTTP server should listen on when running
	// in self-hosted mode (for example, :8080 or 0.0.0.0:8080). This can then be placed behind
	// any load balancer such as an AWS ALB.
	SelfHostedListenAddress string

	// SelfHostedPublicURL is the URL that users (and notifications) should use to reach the portal
	// when operating in self-hosted mode.
	SelfHostedPublicURL string

	// NotifierSlackEnabled will be used to determine whether the Slack notifier
	// is enabled or not
	NotifierSlackEnabled bool

	// NotifierSlackToken is the token that will be used to make the Slack request
	// to send the message(s)
	NotifierSlackToken string

	// NotifierSlackThreadTs is the timestamp of the message to reply to in the thread
	NotifierSlackThreadTs string

	// NotifierSlackChannel is the channel that message(s) will be sent to
	NotifierSlackChannel string

	// NotifierSlackBotName is the name of the Slack bot that will we used
	// when sending notifications
	NotifierSlackBotName string

	// NotifierDiscordEnabled will be used to determine whether the Slack notifier
	// is enabled or not
	NotifierDiscordEnabled bool

	// NotifierDiscordThreadId is the ID of the Discord thread the message should be sent to
	// (as a threaded message)
	NotifierDiscordThreadId string

	// NotifierDiscordWebhook is the webhook that will be used to make the Discord request
	// to send the message(s)
	NotifierDiscordWebhook string

	// NotifierDiscordUsernameOverride is the username that will be used when sending
	//  the message(s)
	NotifierDiscordUsernameOverride string

	// GithubToken is the token that will be used to allow action to leverage the GitHub API
	GithubToken string

    // RunnerEndpointKey is a unique, URL-safe key used to namespace all
    // endpoints for a specific runner invocation. This ensures requests
    // (submit/cancel/uploads) are routed to the correct runner instance.
    RunnerEndpointKey string

    Action *githubactions.Action
}

const (
	// DefaultTimeout is the default timeout that will be used to manage how long the portal
	// will be available for users to use before it is automatically deactivated
	//
	// Defaults to 300 seconds (5 minutes)
	DefaultTimeout int = 300

	// DefaultSelfHostedListenAddress is the default bind address when the action runs in self-hosted mode
	DefaultSelfHostedListenAddress string = ":8080"
)

const (
    // PortalHostModeSelfHosted exposes the portal via a self-hosted HTTP listener
    PortalHostModeSelfHosted string = "self-hosted"
)

// NewFromInputs creates a new Config instance from the provided GitHub Actions inputs.
// It utilises the inputs from the GitHub Actions context, and returns a new Config
// instance with the parsed values.
// If the fields input is malformed and cannot be parsed into a valid Fields struct,
// it returns an ErrMalformedFieldsInputDataProvided error.
func NewFromInputs(action *githubactions.Action) (*Config, error) {

	var err error

    // self-hosted is the only supported mode; read input for compatibility and warn if not supported
    portalHostModeRaw := strings.TrimSpace(strings.ToLower(getInput(action, "portal-host-mode")))
    if portalHostModeRaw != "" && portalHostModeRaw != PortalHostModeSelfHosted {
        action.Warningf("Ignoring unsupported portal-host-mode '%s'. Only '%s' is supported; using self-hosted.", portalHostModeRaw, PortalHostModeSelfHosted)
    }
    portalHostModeInput := PortalHostModeSelfHosted

	selfHostedListenAddress := getInput(action, "selfhosted-listen-address")
	if strings.TrimSpace(selfHostedListenAddress) == "" {
		selfHostedListenAddress = DefaultSelfHostedListenAddress
	}
	selfHostedPublicURL := strings.TrimSpace(getInput(action, "selfhosted-public-url"))

    if portalHostModeInput == PortalHostModeSelfHosted && selfHostedPublicURL == "" {
        action.Errorf("The selfhosted-public-url input must be provided when portal-host-mode is set to '%s'", PortalHostModeSelfHosted)
        return nil, errors.ErrSelfHostedPublicURLMissing
    }

    // handle input for fetching github token
    githubTokenInput := getInput(action, "github-token")
    if githubTokenInput == "" {
        action.Errorf("The github-token was not provided, this is needed before the action can be used")
        return nil, errors.ErrGithubTokenNotProvided
    }

	// handle input for fetching timeout
	var timeout int
	timeoutInput := getInput(action, "timeout")
	if timeoutInput == "" {
		timeout = DefaultTimeout
		action.Debugf("The timeout was not provided, will use the default timeout of %d seconds", DefaultTimeout)
	}
	if timeoutInput != "" {
		timeout, err = strconv.Atoi(timeoutInput)
		if err != nil {
			action.Fatalf("Cannot convert the 'timeout' input (%s) to an int!", timeoutInput)
			return nil, errors.ErrInvalidTimeoutValueProvided
		}
	}

	// handle input for fetching form title if provided
	titleInput := getInput(action, "title")
	if titleInput != "" {
		action.Debugf("Title input provided: %s", titleInput)
	}

	// handle input for fetching interactive inputs portal fields if provided
	interactiveInput := getInput(action, "interactive")
	fields, err := fields.MarshalStringIntoValidFieldsStruct(interactiveInput, action)
	if err != nil {
		action.Errorf("Can't convert the 'fields' input to a valid fields config: %s", interactiveInput)
		return nil, errors.ErrMalformedFieldsInputDataProvided
	}

    // handle input for fetching slack notifier
    var notifierSlackToken string = "xoxb-secret-token"
    var notifierSlackChannel string = "#notificatins"
    var notifierSlackBotName string
    var notifierSlackThreadTs string

	notifierSlackEnabledInput := getInput(action, "notifier-slack-enabled") == "true"
	if notifierSlackEnabledInput {

		notifierSlackTokenInput := getInput(action, "notifier-slack-token")
		if notifierSlackTokenInput == notifierSlackToken {
			action.Errorf("A valid Slack token was not provided, please provide a valid Slack token when enabling the Slack notifier")
			return nil, errors.ErrInvalidSlackTokenProvided
		}
		notifierSlackToken = notifierSlackTokenInput
		notifierSlackChannel = getInput(action, "notifier-slack-channel")
		notifierSlackBotName = getInput(action, "notifier-slack-bot")
		notifierSlackThreadTs = strings.TrimSpace(getInput(action, "notifier-slack-thread-ts"))
	}

    // handle input for fetching discord notifier
    var notifierDiscordWebhook string = "secret-webhook"
    var notifierDiscordUsernameOverride string
    var notifierDiscordThreadId string

	notifierDiscordEnabledInput := getInput(action, "notifier-discord-enabled") == "true"
	if notifierDiscordEnabledInput {

		notifierDiscordWebhookInput := getInput(action, "notifier-discord-webhook")
		if notifierDiscordWebhookInput == notifierDiscordWebhook {
			action.Errorf("A valid Discord webhook was not provided, please provide a valid Discord webhook when enabling the Discord notifier")
			return nil, errors.ErrInvalidDiscordWebhookProvided
		}

		notifierDiscordWebhook = notifierDiscordWebhookInput
		notifierDiscordUsernameOverride = getInput(action, "notifier-discord-username")
		notifierDiscordThreadId = strings.TrimSpace(getInput(action, "notifier-discord-thread-id"))
	}

    // handle masking of sensitive data
    action.AddMask(notifierSlackToken)
    action.AddMask(notifierDiscordWebhook)
    action.AddMask(githubTokenInput)

    // derive or accept a runner endpoint key for namespacing endpoints
    // this allows multiple runner instances to coexist under unique paths
    runnerEndpointKey := strings.Trim(getInput(action, "runner-endpoint-key"), "/ ")
    if runnerEndpointKey == "" {
        // Try to derive from the GitHub context if available
        if ctx, err := action.Context(); err == nil && ctx.RunID != 0 {
            sha := strings.TrimSpace(ctx.SHA)
            if sha != "" {
                // Use short SHA for readability
                short := sha
                if len(short) > 8 {
                    short = sha[:8]
                }
                runnerEndpointKey = fmt.Sprintf("input-%s-%d", short, ctx.RunID)
            } else {
                runnerEndpointKey = fmt.Sprintf("run-%d", ctx.RunID)
            }
        } else {
            // Safe, deterministic placeholder suitable for tests/local
            runnerEndpointKey = "runner"
        }
    }

    c := Config{
        Title:                   titleInput,
        Fields:                  fields,
        Timeout:                 timeout,
        PortalHostMode:          portalHostModeInput,
        SelfHostedListenAddress: selfHostedListenAddress,
        SelfHostedPublicURL:     selfHostedPublicURL,

        GithubToken: githubTokenInput,

        RunnerEndpointKey: runnerEndpointKey,

        NotifierSlackEnabled:  notifierSlackEnabledInput,
        NotifierSlackToken:    notifierSlackToken,
        NotifierSlackChannel:  notifierSlackChannel,
        NotifierSlackBotName:  notifierSlackBotName,
        NotifierSlackThreadTs: notifierSlackThreadTs,

		NotifierDiscordEnabled:          notifierDiscordEnabledInput,
		NotifierDiscordWebhook:          notifierDiscordWebhook,
		NotifierDiscordUsernameOverride: notifierDiscordUsernameOverride,
		NotifierDiscordThreadId:         notifierDiscordThreadId,

		Action: action,
    }
    return &c, nil
}