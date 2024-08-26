package notifier

type DiscordVerifyResponse struct {

	// ApplicationId is the ID of the application that owns the webhook
	ApplicationId string `json:"application_id,omitempty"`

	// Avatar is the URL of the avatar of the webhook
	Avatar string `json:"avatar,omitempty"`

	// ChannelId is the ID of the channel that the webhook is for
	ChannelId string `json:"channel_id,omitempty"`

	// GuildId is the ID of the guild that the webhook is for
	GuildId string `json:"guild_id,omitempty"`

	// Id is the ID of the webhook
	Id string `json:"id,omitempty"`

	// Name is the name of the webhook
	Name string `json:"name,omitempty"`

	// Type is the type of the webhook
	Type int `json:"type,omitempty"`

	// Token is the token of the webhook
	Token string `json:"token,omitempty"`

	// Url is the URL of the webhook
	Url string `json:"url,omitempty"`

	// Failure

	// Message is the message of the error
	Message string `json:"message,omitempty"`

	// Code is the code of the error
	Code int `json:"code,omitempty"`
}
