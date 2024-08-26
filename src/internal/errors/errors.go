package errors

import "errors"

var (
	// ErrMalformedFieldsInputDataProvided is returned when the input data provided is malformed
	// and cannot be parsed into a valid Fields struct
	ErrMalformedFieldsInputDataProvided = errors.New("MalformedFieldsInputDataProvided")

	// ErrNoFieldsProvided is returned when no fields are detected in the input data provided
	// after parsing the input data into a Fields struct
	ErrNoFieldsProvided = errors.New("NoFieldsProvided")

	// ErrInvalidFieldTypeProvided is returned when the type provided in the input data provided
	// is not one of the supported field types
	ErrInvalidFieldTypeProvided = errors.New("InvalidFieldTypeProvided")

	// ErrGitHubWorkspaceEnvVarIsMissing is returned when the GitHub Workspace
	// environment variable is not set.
	ErrGitHubWorkspaceEnvVarIsMissing = errors.New("GitHubWorkspaceEnvVarIsMissing")

	// ErrInvalidLabelProvided is returned when the label provided in the input data cannot
	// be converted to kebab case
	ErrInvalidLabelProvided = errors.New("InvalidLabelProvided")

	// ErrInvalidTimeoutValueProvided is returned when the timeout provided cannot be converted
	// to an integer
	ErrInvalidTimeoutValueProvided = errors.New("InvalidTimeoutValueProvided")

	// ErrInvalidSlackTokenProvided is returned when the Slack token provided is not valid
	ErrInvalidSlackTokenProvided = errors.New("InvalidSlackTokenProvided")

	// ErrUnexpectedSlackVerificationStatusCode is returned when the Slack verification endpoint
	// returns a status code other than 200
	ErrUnexpectedSlackVerificationStatusCode = errors.New("UnexpectedSlackVerificationStatusCode")

	// ErrInvalidDiscordWebhookProvided is returned when the Discord webhook provided is not valid
	ErrInvalidDiscordWebhookProvided = errors.New("InvalidDiscordWebhookProvided")

	// ErrUnexpectedDiscordVerificationStatusCode is returned when the Discord verification endpoint
	// returns a status code other than 200 or 401
	ErrUnexpectedDiscordVerificationStatusCode = errors.New("UnexpectedDiscordVerificationStatusCode")

	// ErrFailedToSendMessageWithNotifier is returned when the notifier is unable to send a message
	ErrFailedToSendMessageWithNotifier = errors.New("FailedToSendMessageWithNotifier")

	//ErrNgrokAuthtokenNotProvided is returned when the ngrok authtoken is not provided
	ErrNgrokAuthtokenNotProvided = errors.New("NgrokAuthtokenNotProvided")

	// ErrGithubTokenNotProvided is returned when the github token is not provided
	ErrGithubTokenNotProvided = errors.New("GithubTokenNotProvided")
)
