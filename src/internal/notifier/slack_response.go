package notifier

// SlackVerifyResponse represents the response from the Slack API auth.test endpoint
type SlackVerifyResponse struct {

	// Ok is a boolean indicating whether the verification was successful or not
	Ok bool `json:"ok,omitempty"`

	// Error is a string representing the reason if the verification was not successful
	// i.e. not_authed
	Error string `json:"error,omitempty"`
}

// SlackChatPostMessageResponse represents the response from the Slack API chat.postMessage endpoint
type SlackChatPostMessageResponse struct {

	// Ok is a boolean indicating whether the message was sent successfully or not
	Ok bool `json:"ok,omitempty"`

	// Success

	// Channel is the ID of the channel that the message was sent to
	Channel string `json:"channel,omitempty"`

	// Ts is the timestamp of the message
	Ts string `json:"ts,omitempty"`

	// Message is the message that was sent
	Message *SlackMessage `json:"message,omitempty"`

	// Failure

	// Error is a string representing the reason if the message was not sent successfully
	// i.e. too_many_attachments
	Error string `json:"error,omitempty"`
}

// SlackMessage represents the message that was sent
type SlackMessage struct {

	// Text is the text of the message
	Text string `json:"text,omitempty"`

	// Username is the username of the bot sending the message
	Username string `json:"username,omitempty"`

	// Attachments is a list of attachments
	Attachments []SlackAttachment `json:"attachments,omitempty"`

	// Type is the type of the message
	Type string `json:"type,omitempty"`

	// Subtype is the subtype of the message
	Subtype string `json:"subtype,omitempty"`

	// Ts is the timestamp of the message
	Ts string `json:"ts,omitempty"`
}

// SlackAttachment represents the attachment of the message
type SlackAttachment struct {

	// Text is the text of the attachment
	Text string `json:"text,omitempty"`

	// Id is the ID of the attachment
	Id int `json:"id,omitempty"`

	// Fallback is the fallback text of the attachment
	Fallback string `json:"fallback,omitempty"`
}
