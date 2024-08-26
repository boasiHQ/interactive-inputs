package notifier

// DiscordPostMessageRequest represents the request for sending a message using the
// Discord webhook
type DiscordPostMessageRequest struct {

	// AvatarUrl is the URL of the icon to use for the bot
	AvatarUrl string `json:"avatar_url,omitempty"`

	// Username is the username of the bot sending the message
	Username string `json:"username,omitempty"`

	// Content is the message to send
	Content string `json:"content,omitempty"`

	// Embeds is a list of embeds to send in the message
	Embeds []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {

	// Titles is the title of the embed
	Title string `json:"title,omitempty"`

	// Description is the description of the embed
	Description string `json:"description,omitempty"`

	// Color is the color of the embed
	// More info: https://birdie0.github.io/discord-webhooks-guide/structure/embed/color.html
	Color string `json:"color,omitempty"`

	// Url is the link for the embed's title
	// requires "title" to be set
	Url string `json:"url,omitempty"`

	// Timestamp is the user specified timestamp of the embed
	// expected format: "YYYY-MM-DDTHH:MM:SS.MSSZ"
	Timestamp string `json:"timestamp,omitempty"`

	// Thumbnail is the thumbnail of the embed
	Thumbnail DiscordEmbedThumbnail `json:"thumbnail,omitempty"`

	// Image is the image of the embed
	Image DiscordEmbedImage `json:"image,omitempty"`

	// Footer is the footer of the embed
	Footer DiscordEmbedFooter `json:"footer,omitempty"`

	// Author is the author of the embed
	Author DiscordEmbedAuthor `json:"author,omitempty"`

	// Fields is a list of fields to send in the embed
	Fields []DiscordEmbedField `json:"fields,omitempty"`
}

type DiscordEmbedField struct {
	// Name is the name of the field
	Name string `json:"name,omitempty"`

	// Value is the value of the field
	Value string `json:"value,omitempty"`

	// Inline is whether the field should be displayed inline
	Inline bool `json:"inline,omitempty"`
}

type DiscordEmbedAuthor struct {

	// Name is the name of the author
	Name string `json:"name,omitempty"`

	// Url is the url to use for the author's name
	Url string `json:"url,omitempty"`

	// IconUrl is the url to use for the author's avatar
	IconUrl string `json:"icon_url,omitempty"`
}

type DiscordEmbedFooter struct {

	// Text is the text of the footer
	Text string `json:"text,omitempty"`

	// IconUrl is the URL of the icon to use for the footer
	IconUrl string `json:"icon_url,omitempty"`
}

type DiscordEmbedImage struct {
	// Url is the URL of the image
	Url string `json:"url,omitempty"`
}

type DiscordEmbedThumbnail struct {
	// Url is the URL of the image to use for the thumbnail
	Url string `json:"url,omitempty"`
}
