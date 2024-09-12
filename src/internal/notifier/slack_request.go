package notifier

// SlackChatPostMessageRequest represents the requesr for the Slack API chat.postMessage endpoint
type SlackChatPostMessageRequest struct {

	// Channel is the ID of the channel that the message should be sent to
	Channel string `json:"channel,omitempty"`

	// ThreadTs is the timestamp of the message to reply to
	ThreadTs string `json:"thread_ts,omitempty"`

	// IconUrl is the URL of the icon to use for the bot
	IconUrl string `json:"icon_url,omitempty"`

	// UnfurlMedia is whether to enable unfurling of media content.
	UnfurlMedia bool `json:"unfurl_media,omitempty"`

	// UnfurlLinks is whether to enable unfurling of primarily text-based content.
	UnfurlLinks bool `json:"unfurl_links,omitempty"`

	// Username is the username of the bot sending the message
	Username string `json:"username,omitempty"`

	// Blocks is a list of blocks to send in the message
	Blocks []SlackBlock `json:"blocks,omitempty"`
}

// SlackBlock represents a block in a Slack message
type SlackBlock struct {
	// Type is the type of the block
	Type string `json:"type,omitempty"`

	// BlockText represets the text of the block
	Text *BlockText `json:"text,omitempty"`
}

// BLockText represents the text of a block
type BlockText struct {
	// Type is the type of the text
	Type string `json:"type,omitempty"`

	// Text is the text of the block
	Text string `json:"text,omitempty"`
}
