package portal

// UploadToPortalResponse represents the response for uploading files to the portal
type UploadToPortalResponse struct {

	// Status represents the status of the upload
	Status string `json:"status"`

	// UploadedFiles represents the list of files uploaded successfully
	UploadedFiles []string `json:"uploaded_files,omitempty"`

	// FailedFiles represents the list of files that failed to upload
	FailedFiles []string `json:"failed_files,omitempty"`
}

// ResetUploadResponse represents the response for resetting the upload
type ResetUploadResponse struct {
	// Status represents the status of the reset
	Status string `json:"status"`

	// DeletedFiles represents the list of files deleted
	DeletedFiles []string `json:"deleted_files,omitempty"`

	// FailedFiles represents the list of files that failed to delete
	FailedFiles []string `json:"failed_files,omitempty"`

	// TotalFilesToDelete represents the total number of files that were to be deleted
	TotalFilesToDelete int `json:"total_files_to_delete"`

	// TotalFilesDeleted represents the total number of files that were deleted
	TotalFilesDeleted int `json:"total_files_deleted"`
}
