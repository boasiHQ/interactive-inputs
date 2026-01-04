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
	
	// BasePath is the path prefix to reach this runner's portal, e.g. "/run-123"
    BasePath string

    // BalloonData holds per-field suggestion values for the scroll balloon UI
    // map key: field label; value: list of suggestions
    BalloonData map[string][]string

    // PreOutput holds a small read-only output (e.g. previous step result) to
    // display above a field. Keyed by field label.
    PreOutput map[string]struct{ Title, Value string }

}
