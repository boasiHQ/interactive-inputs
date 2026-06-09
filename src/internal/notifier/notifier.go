package notifier

// Notifier sends and verifies optional portal notifications.
type Notifier interface {

	// Notify sends a notification to respective integration
	// returning the id that can be used for future thread communication
	// if the integration supports it.
	Notify(title, message string) (string, error)

	// Verify checks the connection to the integration.
	Verify() error

	// Enabled returns whether the notifier is enabled or not
	Enabled() bool
}
