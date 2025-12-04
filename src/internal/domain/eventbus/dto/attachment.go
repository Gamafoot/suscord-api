package dto

type Attachment struct {
	ID       uint   `json:"id"`
	FilePath string `json:"file_path"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
}
