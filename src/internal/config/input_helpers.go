package config

import (
	"strings"

	githubactions "github.com/sethvargo/go-githubactions"
)

// getInput fetches an input value from the GitHub Actions runtime.
// The upstream go-githubactions helper only replaces spaces with underscores
// when constructing the INPUT_* variable name. GitHub also emits variables
// where hyphens are converted to underscores, so we fall back to that form
// if the direct lookup returns empty.
func getInput(action *githubactions.Action, name string) string {
	val := action.GetInput(name)
	if val != "" {
		return val
	}

	if !strings.Contains(name, "-") {
		return val
	}

	alternate := strings.ReplaceAll(name, "-", "_")
	return action.GetInput(alternate)
}