package model

// Attachment is a file uploaded to a task.
type Attachment struct {
	Base
	TaskID      uint   `gorm:"index;not null" json:"task_id"`
	FileName    string `gorm:"size:255;not null" json:"file_name"`
	FilePath    string `gorm:"size:512;not null" json:"file_path"`
	Size        int64  `json:"size"`
	ContentType string `gorm:"size:128" json:"content_type"`
	UploadedBy  uint   `gorm:"not null" json:"uploaded_by"`
}
