package dto

type ApplicationDocumentRes struct {
	ID            uint   `json:"id"`
	ApplicationID uint   `json:"application_id"`
	DocumentType  string `json:"document_type"`
	FilePath      string `json:"file_path"` // object_key
	PublicURL     string `json:"public_url"`
	CreatedAt     string `json:"created_at"`
}
