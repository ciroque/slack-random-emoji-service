package server

type SlackEmojiResponse struct {
	ResponseType string              `json:"response_type""`
	Text         string              `json:"text""`
	Attachments  []map[string]string `json:"attachments""`
}
