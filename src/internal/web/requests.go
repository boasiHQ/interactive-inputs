package webui

import "github.com/boasihq/interactive-inputs/internal/fields"

// CreateInteractiveInputsPortalRequest is the request that will
// be used to create a new interactive inputs portal.
type CreateInteractiveInputsPortalRequest struct {
	// Title is the header that will be displayed at the top of the generated form
	Title string

	// Fields is the slice of fields that will be displayed in the generated form
	Fields *fields.Fields

	// RepoOwner is the name that will feature in the portal's title
	RepoOwner string

	// Timeout is how long the portal will be available for users to use before it is
	// automatically deactivated
	Timeout string
}
