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
			name: "successful - explicit ngrok mode",
			envMap: map[string]string{
				"INPUT_PORTAL-HOST-MODE": "ngrok",
				"INPUT_TITLE":            "Test Title",
				"INPUT_INTERACTIVE":      "fields:\n  - label: name\n    properties:\n      type: text\n",
				"INPUT_GITHUB-TOKEN":     "github-secret-token",
				"INPUT_NGROK-AUTHTOKEN":  "ngrok-secret-token",
			},
			expectedConfig: config.Config{
				Timeout:        300,
				Title:          "Test Title",
				PortalHostMode: config.PortalHostModeNgrok,
				NgrokAuthtoken: "ngrok-secret-token",
				GithubToken:    "github-secret-token",
				Fields: &fields.Fields{
					Fields: []fields.Field{{Label: "name", Properties: fields.FieldProperties{Type: "text"}}},
				},
				NotifierSlackToken:     "xoxb-secret-token",
				NotifierSlackChannel:   "#notifications",
				NotifierDiscordWebhook: "secret-webhook",
				RunnerEndpointKey:      "",
			},
			expectedOutput: "::debug::Ngrok mode active.\n::add-mask::ngrok-secret-token\n::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::debug::Title input provided: Test Title\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n",
			expectedError:  nil,
		},
		{
			name: "successful - explicit self-host alias",
			envMap: map[string]string{
				"INPUT_PORTAL-HOST-MODE":      "self-host",
				"INPUT_GITHUB-TOKEN":          "github-secret-token",
				"INPUT_SELFHOSTED-PUBLIC-URL": "https://portal.example.com",
				"INPUT_INTERACTIVE":           "fields:\n  - label: name\n    properties:\n      type: text\n",
			},
			expectedConfig: config.Config{
				Timeout:                 300,
				PortalHostMode:          config.PortalHostModeSelfHosted,
				SelfHostedPublicURL:     "https://portal.example.com",
				SelfHostedListenAddress: ":8080",
				GithubToken:             "github-secret-token",
				Fields: &fields.Fields{
					Fields: []fields.Field{{Label: "name", Properties: fields.FieldProperties{Type: "text"}}},
				},
				NotifierSlackToken:     "xoxb-secret-token",
				NotifierSlackChannel:   "#notifications",
				NotifierDiscordWebhook: "secret-webhook",
				RunnerEndpointKey:      "runner",
			},
			expectedOutput: "::debug::Self-hosted mode active.\n::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n",
			expectedError:  nil,
		},
		{
			name: "successful - explicit self-hosted mode",
			envMap: map[string]string{
				"INPUT_PORTAL-HOST-MODE":          "self-hosted",
				"INPUT_GITHUB-TOKEN":              "github-secret-token",
				"INPUT_SELFHOSTED-PUBLIC-URL":     "https://portal.example.com",
				"INPUT_SELFHOSTED-LISTEN-ADDRESS": ":9090",
				"INPUT_INTERACTIVE":               "fields:\n  - label: name\n    properties:\n      type: text\n",
			},
			expectedConfig: config.Config{
				Timeout:                 300,
				PortalHostMode:          config.PortalHostModeSelfHosted,
				SelfHostedPublicURL:     "https://portal.example.com",
				SelfHostedListenAddress: ":9090",
				GithubToken:             "github-secret-token",
				Fields: &fields.Fields{
					Fields: []fields.Field{{Label: "name", Properties: fields.FieldProperties{Type: "text"}}},
				},
				NotifierSlackToken:     "xoxb-secret-token",
				NotifierSlackChannel:   "#notifications",
				NotifierDiscordWebhook: "secret-webhook",
				RunnerEndpointKey:      "runner",
			},
			expectedOutput: "::debug::Self-hosted mode active.\n::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n",
			expectedError:  nil,
		},
		{
			name: "successful - self-hosted mode with underscore endpoint key",
			envMap: map[string]string{
				"INPUT_PORTAL-HOST-MODE":      "self-hosted",
				"INPUT_GITHUB-TOKEN":          "github-secret-token",
				"INPUT_SELFHOSTED-PUBLIC-URL": "https://portal.example.com",
				"INPUT_INTERACTIVE":           "fields:\n  - label: name\n    properties:\n      type: text\n",
				"INPUT_RUNNER_ENDPOINT_KEY":   "custom",
			},
			expectedConfig: config.Config{
				Timeout:                 300,
				PortalHostMode:          config.PortalHostModeSelfHosted,
				SelfHostedPublicURL:     "https://portal.example.com",
				SelfHostedListenAddress: ":8080",
				GithubToken:             "github-secret-token",
				RunnerEndpointKey:       "custom",
				Fields: &fields.Fields{
					Fields: []fields.Field{{Label: "name", Properties: fields.FieldProperties{Type: "text"}}},
				},
				NotifierSlackToken:     "xoxb-secret-token",
				NotifierSlackChannel:   "#notifications",
				NotifierDiscordWebhook: "secret-webhook",
			},
			expectedOutput: "::debug::Self-hosted mode active.\n::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::debug::runner-endpoint-key read from INPUT_RUNNER_ENDPOINT_KEY\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n",
			expectedError:  nil,
		},
		{
			name: "failed - portal host mode not provided",
			envMap: map[string]string{
				"INPUT_NGROK-AUTHTOKEN": "ngrok-secret-token",
			},
			expectedOutput: "::error::portal-host-mode is required and must be one of \"ngrok\", \"self-host\", \"self-hosted\".\n",
			expectedError:  errors.ErrNoHostingModeProvided,
		},
		{
			name: "failed - ngrok mode missing token",
			envMap: map[string]string{
				"INPUT_PORTAL-HOST-MODE": "ngrok",
				"INPUT_NGROK-AUTHTOKEN":  "",
			},
			expectedOutput: "::error::Ngrok authtoken must be provided when portal-host-mode is set to 'ngrok'.\n",
			expectedError:  errors.ErrNgrokAuthtokenNotProvided,
		},
		{
			name: "failed - self-hosted mode missing public url",
			envMap: map[string]string{
				"INPUT_PORTAL-HOST-MODE":      "self-hosted",
				"INPUT_SELFHOSTED-PUBLIC-URL": "",
			},
			expectedOutput: "::error::Self-hosted public URL must be provided when portal-host-mode is set to 'self-hosted'.\n",
			expectedError:  errors.ErrSelfHostedPublicURLMissing,
		},
		{
			name: "failed - invalid portal host mode",
			envMap: map[string]string{
				"INPUT_PORTAL-HOST-MODE": "invalid",
			},
			expectedOutput: "::error::Invalid portal-host-mode provided: invalid. Supported modes are \"ngrok\", \"self-host\", \"self-hosted\".\n",
			expectedError:  errors.ErrInvalidPortalHostModeProvided,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actionLog := bytes.NewBuffer(nil)

			getenv := func(key string) string {
				return test.envMap[key]
			}

			for k, v := range test.envMap {
				if k == "INPUT_RUNNER_ENDPOINT_KEY" {
					t.Setenv(k, v)
				}
			}

			action := githubactions.New(
				githubactions.WithWriter(actionLog),
				githubactions.WithGetenv(getenv),
			)

			cfg, inputsErr := config.NewFromInputs(action)

			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, inputsErr)
				assert.Equal(t, test.expectedOutput, actionLog.String())
			} else {
				assert.NoError(t, inputsErr)
				assert.Equal(t, test.expectedOutput, actionLog.String())
				cfg.Action = nil
				assert.EqualValues(t, test.expectedConfig, *cfg)
			}
		})
	}
}
