package dto

type Attachment struct {
	ID       uint   `json:"id"`
	FileUrl  string `json:"file_url"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
}
