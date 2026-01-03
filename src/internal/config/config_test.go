package config_test

import (
	"bytes"
	"testing"

	"github.com/boasihq/interactive-inputs/internal/config"
	"github.com/boasihq/interactive-inputs/internal/errors"
	"github.com/boasihq/interactive-inputs/internal/fields"
	githubactions "github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/assert"
)

func TestConfig_NewFromInputs(t *testing.T) {

	tests := []struct {
		name           string
		envMap         map[string]string
		expectedOutput string
		expectedConfig config.Config
		expectedError  error
	}{
		{
			name: "successful - created base config from input (default timeout)",
			envMap: map[string]string{
				"INPUT_TITLE":           "What name should be given to the barista?",
				"INPUT_INTERACTIVE":     "fields:\n  - label: name\n    properties:\n      display: name\n      type: text\n      description: Name of the user\n      maxLength: 20\n      required: false\n",
				"INPUT_GITHUB-TOKEN":    "github-secret-token",
				"INPUT_NGROK-AUTHTOKEN": "ngrok-secret-token",
			},
			expectedConfig: config.Config{
				Timeout: 300,
				Title:   "What name should be given to the barista?",
				Fields: &fields.Fields{
					Fields: []fields.Field{
						{
							Label: "name",
							Properties: fields.FieldProperties{
								Display:     "name",
								Type:        "text",
								Description: "Name of the user",
								MaxLength:   20,
								Required:    false,
							},
						},
					},
				},
				NotifierSlackEnabled:            false,
				NotifierSlackToken:              "xoxb-secret-token",
				NotifierSlackChannel:            "#notificatins",
				NotifierDiscordEnabled:          false,
				NotifierDiscordWebhook:          "secret-webhook",
				GithubToken:                     "github-secret-token",
				NgrokAuthtoken:                  "ngrok-secret-token",
				PortalHostMode:                  config.PortalHostModeNgrok,
				RunnerEndpointKey:               "runner",
			},
			expectedOutput: "::debug::Ngrok authtoken detected. Using Ngrok mode.\n::add-mask::ngrok-secret-token\n::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::debug::Title input provided: What name should be given to the barista?\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n::add-mask::ngrok-secret-token\n",
			expectedError:  nil,
		},
		{
			name: "successful - self-hosted mode (ngrok absent)",
			envMap: map[string]string{
				"INPUT_TITLE":                     "Self Hosted Test",
				"INPUT_INTERACTIVE":               "fields:\n  - label: name\n    properties:\n      type: text\n",
				"INPUT_GITHUB-TOKEN":              "github-secret-token",
				"INPUT_SELFHOSTED-PUBLIC-URL":     "https://portal.example.com",
				"INPUT_SELFHOSTED-LISTEN-ADDRESS": ":9090",
			},
			expectedConfig: config.Config{
				Timeout: 300,
				Title:   "Self Hosted Test",
				Fields: &fields.Fields{
					Fields: []fields.Field{
						{
							Label:      "name",
							Properties: fields.FieldProperties{Type: "text"},
						},
					},
				},
				PortalHostMode:          config.PortalHostModeSelfHosted,
				SelfHostedPublicURL:     "https://portal.example.com",
				SelfHostedListenAddress: ":9090",
				RunnerEndpointKey:       "runner",
				GithubToken:             "github-secret-token",
				NotifierSlackEnabled:    false,
				NotifierSlackToken:      "xoxb-secret-token",
				NotifierSlackChannel:    "#notificatins",
				NotifierDiscordEnabled:  false,
				NotifierDiscordWebhook:  "secret-webhook",
			},
			expectedOutput: "::debug::Self-hosted public URL detected. Using Self-Hosted mode.\n::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::debug::Title input provided: Self Hosted Test\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n::add-mask::\n",
			expectedError:  nil,
		},
		{
			name: "successful - ngrok priority (both provided)",
			envMap: map[string]string{
				"INPUT_NGROK-AUTHTOKEN":       "ngrok-secret",
				"INPUT_SELFHOSTED-PUBLIC-URL":  "https://ignored.com",
				"INPUT_GITHUB-TOKEN":          "github-secret-token",
				"INPUT_INTERACTIVE":           "fields:\n  - label: name\n    properties:\n      type: text\n",
			},
			expectedConfig: config.Config{
				Timeout: 300,
				Title:   "",
				Fields: &fields.Fields{
					Fields: []fields.Field{
						{
							Label:      "name",
							Properties: fields.FieldProperties{Type: "text"},
						},
					},
				},
				PortalHostMode:         config.PortalHostModeNgrok,
				NgrokAuthtoken:         "ngrok-secret",
				GithubToken:            "github-secret-token",
				SelfHostedPublicURL:    "https://ignored.com",
				RunnerEndpointKey:      "runner",
				NotifierSlackToken:     "xoxb-secret-token",
				NotifierSlackChannel:   "#notificatins",
				NotifierDiscordWebhook: "secret-webhook",
			},
			expectedOutput: "::debug::Ngrok authtoken detected. Using Ngrok mode.\n::add-mask::ngrok-secret\n::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n::add-mask::ngrok-secret\n",
			expectedError:  nil,
		},
		{
			name: "failed - missing both hosting modes",
			envMap: map[string]string{
				"INPUT_GITHUB-TOKEN": "github-secret-token",
				"INPUT_INTERACTIVE":  "fields:\n  - label: name\n    properties:\n      type: text\n",
			},
			expectedConfig: config.Config{},
			// Removed the timeout debug message because the function returns early
			expectedOutput: "::error::Configuration error: Either 'ngrok-authtoken' or 'selfhosted-public-url' must be provided.\n",
			expectedError:  errors.ErrNoHostingModeProvided,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actionLog := bytes.NewBuffer(nil)

			getenv := func(key string) string {
				return test.envMap[key]
			}

			action := githubactions.New(
				githubactions.WithWriter(actionLog),
				githubactions.WithGetenv(getenv),
			)

			cfg, inputsErr := config.NewFromInputs(action)

			if test.expectedError != nil {
				// For errors, check that the specific error is returned
				assert.Equal(t, test.expectedError, inputsErr)
				// For errors, check that the expected logs are present
				assert.Equal(t, test.expectedOutput, actionLog.String())
			} else {
				// For success, ensure no error and check full config
				assert.NoError(t, inputsErr)
				assert.Equal(t, test.expectedOutput, actionLog.String())
				cfg.Action = nil
				assert.EqualValues(t, test.expectedConfig, *cfg)
			}
		})
	}
}