package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boasihq/interactive-inputs/internal/errors"
	"github.com/sethvargo/go-githubactions"
)

// NewSlackNotifierRequest is the request object for creating a new
// instance of a Slack Notifier.
type NewSlackNotifierRequest struct {
	// Enabled whether the notifier is enabled or not
	Enabled bool

	// Token is the Slack API token
	Token string

	// Channel is the Slack channel to send notifications to
	Channel string

	// BotName is the name of the Slack bot that will send notifications
	BotName string

	// ActionPkg represents the githubactions package
	ActionPkg *githubactions.Action

	// VerifificationEndpoint is the endpoint to verify the token
	VerificationEndpoint string

	// ThreadTs is the timestamp of the message to reply to in the thread
	ThreadTs string
}

// NewSlackNotifier returns a new instance of a Slack Notifier
func NewSlackNotifier(r *NewSlackNotifierRequest) Notifier {

	var botName string = "Interactive Inputs"
	var verificationEndpoint string = "https://slack.com/api/auth.test"

	if r.VerificationEndpoint != "" {
		verificationEndpoint = r.VerificationEndpoint
	}

	if r.BotName != "" {
		botName = r.BotName
	}

	return &SlackNotifier{
		enabled:              r.Enabled,
		token:                r.Token,
		channel:              r.Channel,
		botName:              botName,
		action:               r.ActionPkg,
		verificationEndpoint: verificationEndpoint,
		threadTs:             r.ThreadTs,
	}
}

// SlackNotifier is a struct that implements the Notifier interface
type SlackNotifier struct {

	// enabled whether the notifier is enabled or not
	enabled bool

	// token is the Slack API token
	token string

	// channel is the Slack channel to send notifications to
	channel string

	// botName is the name of the Slack bot that will send notifications
	botName string

	// action represents the githubactions package
	action *githubactions.Action

	// verificationEndpoint is the endpoint to verify the auth
	verificationEndpoint string

	// threadTs is the timestamp of the message to reply to in the thread
	threadTs string
}

// Notify sends a notification to the Slack channel
func (n *SlackNotifier) Notify(title, message string) (string, error) {

	var notificationResponse SlackChatPostMessageResponse
	var slackPostChatMessageUrl string = "https://slack.com/api/chat.postMessage"

	// Shape the message to be sent
	renderedMessage, err := n.renderStandardSlackNofityMessage(title, message)
	if err != nil {
		return "", err
	}

	notificationMessage := SlackChatPostMessageRequest{
		Channel:  n.channel,
		Username: n.botName,
		IconUrl:  "https://interactiveinputs.com/static/img/interactive-inputs-no-bg-text-black.png",
		Blocks: []SlackBlock{
			{
				Type: "section",
				Text: &BlockText{
					Type: "mrkdwn",
					Text: renderedMessage,
				},
			},
		},
	}

	// Check if thread ts provided
	if n.threadTs != "" {
		notificationMessage.ThreadTs = n.threadTs
	}

	notificationMessageBytes, err := json.Marshal(notificationMessage)
	if err != nil {
		n.action.Errorf("An error occured while shaping notification message. Message: %s", message)
		return "", err
	}
	requestBody := bytes.NewBuffer(notificationMessageBytes)

	// handle request to endpoint
	resp, err := http.NewRequest(http.MethodPost, slackPostChatMessageUrl, requestBody)
	if err != nil {
		n.action.Errorf("Error on response.\nError: %v", err)

		return "", err
	}

	resp.Header.Add("Content-type", "application/json")
	resp.Header.Add("Authorization", "Bearer "+n.token)
	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		n.action.Errorf("An error occured while making call to verification endpoint. Error: %v", err)
		return "", err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&notificationResponse)
	if err != nil {
		n.action.Errorf("Unexpected error while decoding response. Error: %v", err)
		return "", err
	}

	if !notificationResponse.Ok {
		n.action.Errorf("Unable to send the message to the provided Slack channel. Error: %v", notificationResponse.Error)
		return "", errors.ErrFailedToSendMessageWithNotifier
	}

	n.action.Debugf("Successfully sent message to the provided Slack channel.")

	return notificationResponse.Ts, nil
}

// Verify checks if the Slack token provided is valid by making a call to the Slack API's auth.test endpoint.
// If the token is valid, it returns nil. If the token is invalid, it returns an error.
func (n *SlackNotifier) Verify() error {
	var verificationResponse SlackVerifyResponse

	n.action.Debugf("Initiating the verification of the Slack token provided.")

	resp, err := http.NewRequest(http.MethodGet, n.verificationEndpoint, nil)
	if err != nil {
		n.action.Errorf("Error on response.\nError: %v", err)

		return err
	}

	resp.Header.Add("Authorization", "Bearer "+n.token)
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

	if response.StatusCode != http.StatusOK {
		n.action.Errorf("Unexpected response from Slack endpoint. Status Code: %d", response.StatusCode)
		return errors.ErrUnexpectedSlackVerificationStatusCode
	}

	if !verificationResponse.Ok {
		n.action.Errorf("Unable to verify the the Slack token provided. Error: %v", verificationResponse.Error)
		return errors.ErrInvalidSlackTokenProvided
	}

	n.action.Debugf("Successfully verified the Slack token provided.")

	return nil
}

// Enabled returns whether the notifier is enabled or not
func (n *SlackNotifier) Enabled() bool {
	return n.enabled
}

// renderStandardSlackNofityMessage renders the standard Slack notification message.
func (n *SlackNotifier) renderStandardSlackNofityMessage(title, message string) (string, error) {

	// get action context
	actionCtx, err := n.action.Context()
	if err != nil {
		n.action.Errorf("Failed to get action context. Error: %v", err)
		return "", err
	}

	var optionalSentence string = ""

	if title != "" {
		optionalSentence = fmt.Sprintf("*Title:* _`\"%s\"`_ | ", title)
	}

	defaultNotifyMessageFmt := "*`User Input Required`*" + `

%s<%s|*Go to run*>
*Initiator:* %s


%s
`

	// build out url for action
	repoOwner, repoName := actionCtx.Repo()
	additionalContext := fmt.Sprintf(
		"%s/%s/actions/runs/%d",
		actionCtx.ServerURL,
		repoOwner+"/"+repoName,
		actionCtx.RunID,
	)

	return fmt.Sprintf(defaultNotifyMessageFmt, optionalSentence, additionalContext, actionCtx.Actor, message), nil
}
