package config

import (
	"strconv"

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

	// NotifierSlackEnabled will be used to determine whether the Slack notifier
	// is enabled or not
	NotifierSlackEnabled bool

	// NotifierSlackToken is the token that will be used to make the Slack request
	// to send the message(s)
	NotifierSlackToken string

	// NotifierSlackChannel is the channel that message(s) will be sent to
	NotifierSlackChannel string

	// NotifierSlackBotName is the name of the Slack bot that will we used
	// when sending notifications
	NotifierSlackBotName string

	// NotifierDiscordEnabled will be used to determine whether the Slack notifier
	// is enabled or not
	NotifierDiscordEnabled bool

	// NotifierDiscordWebhook is the webhook that will be used to make the Discord request
	// to send the message(s)
	NotifierDiscordWebhook string

	// NotifierDiscordUsernameOverride is the username that will be used when sending
	//  the message(s)
	NotifierDiscordUsernameOverride string

	// GithubToken is the token that will be used to allow action to leverage the GitHub API
	GithubToken string

	// NgrokAuthtoken is the authtoken that will be used to make Ngrok tunnels to host the
	// interactive inputs portals
	NgrokAuthtoken string

	Action *githubactions.Action
}

const (
	// DefaultTimeout is the default timeout that will be used to manage how long the portal
	// will be available for users to use before it is automatically deactivated
	//
	// Defaults to 300 seconds (5 minutes)
	DefaultTimeout int = 300
)

// NewFromInputs creates a new Config instance from the provided GitHub Actions inputs.
// It utilises the inputs from the GitHub Actions context, and returns a new Config
// instance with the parsed values.
// If the fields input is malformed and cannot be parsed into a valid Fields struct,
// it returns an ErrMalformedFieldsInputDataProvided error.
func NewFromInputs(action *githubactions.Action) (*Config, error) {

	var err error

	// handle input for fetching ngrok authtoken
	ngrokAuthtokenInput := action.GetInput("ngrok-authtoken")
	if ngrokAuthtokenInput == "" {
		action.Errorf("The ngrok-authtoken was not provided, this is needed before the action can be used")
		return nil, errors.ErrNgrokAuthtokenNotProvided
	}

	// handle input for fetching github token
	githubTokenInput := action.GetInput("github-token")
	if githubTokenInput == "" {
		action.Errorf("The github-token was not provided, this is needed before the action can be used")
		return nil, errors.ErrGithubTokenNotProvided
	}

	// handle input for fetching timeout
	var timeout int
	timeoutInput := action.GetInput("timeout")
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
	titleInput := action.GetInput("title")
	if titleInput != "" {
		action.Debugf("Title input provided: %s", titleInput)
	}

	// handle input for fetching interactive inputs portal fields if provided
	interactiveInput := action.GetInput("interactive")
	fields, err := fields.MarshalStringIntoValidFieldsStruct(interactiveInput, action)
	if err != nil {
		action.Errorf("Can't convert the 'fields' input to a valid fields config: %s", interactiveInput)
		return nil, errors.ErrMalformedFieldsInputDataProvided
	}

	// handle input for fetching slack notifier
	var notifierSlackToken string = "xoxb-secret-token"
	var notifierSlackChannel string = "#notificatins"
	var notifierSlackBotName string

	notifierSlackEnabledInput := action.GetInput("notifier-slack-enabled") == "true"
	if notifierSlackEnabledInput {

		notifierSlackTokenInput := action.GetInput("notifier-slack-token")
		if notifierSlackTokenInput == notifierSlackToken {
			action.Errorf("A valid Slack token was not provided, please provide a valid Slack token when enabling the Slack notifier")
			return nil, errors.ErrInvalidSlackTokenProvided
		}
		notifierSlackToken = notifierSlackTokenInput
		notifierSlackChannel = action.GetInput("notifier-slack-channel")
		notifierSlackBotName = action.GetInput("notifier-slack-bot")
	}

	// handle input for fetching discord notifier
	var notifierDiscordWebhook string = "secret-webhook"
	var notifierDiscordUsernameOverride string

	notifierDiscordEnabledInput := action.GetInput("notifier-discord-enabled") == "true"
	if notifierDiscordEnabledInput {

		notifierDiscordWebhookInput := action.GetInput("notifier-discord-webhook")
		if notifierDiscordWebhookInput == notifierDiscordWebhook {
			action.Errorf("A valid Discord webhook was not provided, please provide a valid Discord webhook when enabling the Discord notifier")
			return nil, errors.ErrInvalidDiscordWebhookProvided
		}

		notifierDiscordWebhook = notifierDiscordWebhookInput
		notifierDiscordUsernameOverride = action.GetInput("notifier-discord-username")
	}

	// handle masking of sensitive data
	action.AddMask(notifierSlackToken)
	action.AddMask(notifierDiscordWebhook)
	action.AddMask(githubTokenInput)
	action.AddMask(ngrokAuthtokenInput)

	c := Config{
		Title:   titleInput,
		Fields:  fields,
		Timeout: timeout,

		NgrokAuthtoken: ngrokAuthtokenInput,
		GithubToken:    githubTokenInput,

		NotifierSlackEnabled: notifierSlackEnabledInput,
		NotifierSlackToken:   notifierSlackToken,
		NotifierSlackChannel: notifierSlackChannel,
		NotifierSlackBotName: notifierSlackBotName,

		NotifierDiscordEnabled:          notifierDiscordEnabledInput,
		NotifierDiscordWebhook:          notifierDiscordWebhook,
		NotifierDiscordUsernameOverride: notifierDiscordUsernameOverride,

		Action: action,
	}
	return &c, nil
}
