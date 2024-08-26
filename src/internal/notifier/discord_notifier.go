package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boasihq/interactive-inputs/internal/errors"
	"github.com/sethvargo/go-githubactions"
)

// NewDiscordNotifierRequest is the request object for creating a new
// instance of a Discord Notifier.
type NewDiscordNotifierRequest struct {
	// Enabled whether the notifier is enabled or not
	Enabled bool

	// WebhookUrl is the URL of the Discord webhook
	WebhookUrl string

	// UsernameOverride is the username to override the username of the Discord webhook
	UsernameOverride string

	// ActionPkg represents the githubactions package
	ActionPkg *githubactions.Action

	// VerificationEndpoint is the endpoint override the endpoint
	// used to verify the token
	VerificationEndpoint string
}

// NewDiscordNotifier returns a new instance of a discord Notifier
func NewDiscordNotifier(r *NewDiscordNotifierRequest) Notifier {

	var username string = "Interactive Inputs"
	// by default, the verification endpoint is the same as the webhook url
	var verificationEndpoint string = r.WebhookUrl

	if r.VerificationEndpoint != "" {
		verificationEndpoint = r.VerificationEndpoint
	}

	if r.UsernameOverride != "" {
		username = r.UsernameOverride
	}

	return &DiscordNotifier{
		enabled:              r.Enabled,
		webhookUrl:           r.WebhookUrl,
		usernameOverride:     username,
		action:               r.ActionPkg,
		verificationEndpoint: verificationEndpoint,
	}
}

// DiscordNotifier is a struct that implements the Notifier interface
type DiscordNotifier struct {

	// enabled whether the notifier is enabled or not
	enabled bool

	// webhookUrl is the URL of the Discord webhook
	webhookUrl string

	// usernameOverride is the username to override the username of the Discord webhook
	usernameOverride string

	// action represents the githubactions package
	action *githubactions.Action

	// verificationEndpoint is the endpoint to verify the auth
	verificationEndpoint string
}

func (n *DiscordNotifier) Notify(title, message string) (string, error) {

	// Shape the message to be sent
	renderedMessage, err := n.renderStandardDiscordNofityMessage(title, message)
	if err != nil {
		return "", err
	}

	notificationMessage := DiscordPostMessageRequest{
		Username:  n.usernameOverride,
		AvatarUrl: "https://interactiveinputs.com/static/img/interactive-inputs-bg-black-text-white.png",
		Content:   renderedMessage,
	}

	notificationMessageBytes, err := json.Marshal(notificationMessage)
	if err != nil {
		n.action.Errorf("An error occured while shaping notification message. Message: %s", message)
		return "", err
	}
	requestBody := bytes.NewBuffer(notificationMessageBytes)

	// handle request to endpoint
	resp, err := http.NewRequest(http.MethodPost, n.webhookUrl, requestBody)
	if err != nil {
		n.action.Errorf("Error on response.\nError: %v", err)

		return "", err
	}

	resp.Header.Add("Content-type", "application/json")
	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		n.action.Errorf("An error occured while making call to verification endpoint. Error: %v", err)
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		n.action.Errorf("Unable to send the message to the provided Discord webhook. Status Code: %v", response.StatusCode)
		return "", errors.ErrFailedToSendMessageWithNotifier
	}

	n.action.Debugf("Successfully sent message to the provided Discord webhook.")

	return "", nil
}

// Verify checks the validity of the Discord webhook URL by making a GET request to the verification endpoint (the webhook).
// If the response status code is 200 OK or 401 Unauthorized, and the response message is empty, the webhook is considered valid.
// Otherwise, an error is returned indicating the reason for the verification failure.
func (n *DiscordNotifier) Verify() error {
	var verificationResponse DiscordVerifyResponse

	n.action.Debugf("Initiating the verification of the Discord webhook provided.")

	resp, err := http.NewRequest(http.MethodGet, n.verificationEndpoint, nil)
	if err != nil {
		n.action.Errorf("Error on response.\nError: %v", err)

		return err
	}

	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		n.action.Errorf("An error occured while making call to verification endpoint. Error: %v", err)

		return err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&verificationResponse)
	if err != nil {
		n.action.Errorf("Unexpected error while decoding response. Error: %v", err)
		return err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusUnauthorized {
		n.action.Errorf("Unexpected response from Discord endpoint. Status Code: %d", response.StatusCode)
		return errors.ErrUnexpectedDiscordVerificationStatusCode
	}

	if verificationResponse.Message != "" {
		n.action.Errorf("Unable to verify the the Discord webhook provided. Error: %v", verificationResponse.Message)
		return errors.ErrInvalidDiscordWebhookProvided
	}

	n.action.Debugf("Successfully verified the Discord webhook provided.")

	return nil
}

// Enabled returns whether the notifier is enabled or not
func (n *DiscordNotifier) Enabled() bool {
	return n.enabled
}

// renderStandardDiscordNofityMessage renders the standard Discord notification message.
func (n *DiscordNotifier) renderStandardDiscordNofityMessage(title, message string) (string, error) {

	// get action context
	actionCtx, err := n.action.Context()
	if err != nil {
		n.action.Errorf("Failed to get action context. Error: %v", err)
		return "", err
	}

	var optionalSentence string = "."

	if title != "" {
		optionalSentence = fmt.Sprintf(" - *`\"%s\"`*.", title)
	}

	defaultNotifyMessageFmt := "**`User Input Required`**" + `

Github user %s has kicked off a workflow that requires runtime input(s)%s *You can find out more by visiting the job at %s*.

%s


> *Powered by **[Interactive Inputs](https://interactiveinputs.com/)***`

	// build out url for action
	repoOwner, repoName := actionCtx.Repo()
	additionalContext := fmt.Sprintf(
		"%s/%s/actions/runs/%d",
		actionCtx.ServerURL,
		repoOwner+"/"+repoName,
		actionCtx.RunID,
	)

	return fmt.Sprintf(defaultNotifyMessageFmt, actionCtx.Actor, optionalSentence, additionalContext, message), nil
}
