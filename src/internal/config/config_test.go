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
		name   string
		preRun func()
		envMap map[string]string

		expectedOutput string
		expectedConfig config.Config
		expectedError  error
	}{
		{
			name: "successful - created base config from input (default timeout)",
			preRun: func() {
			},
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
				NotifierSlackBotName:            "",
				NotifierDiscordEnabled:          false,
				NotifierDiscordWebhook:          "secret-webhook",
				NotifierDiscordUsernameOverride: "",
				GithubToken:                     "github-secret-token",
				NgrokAuthtoken:                  "ngrok-secret-token",
			},
			expectedOutput: "::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::debug::Title input provided: What name should be given to the barista?\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n::add-mask::ngrok-secret-token\n",
			expectedError:  nil,
		},
		{
			name: "successful - created base config from input (specified timeout)",
			preRun: func() {
			},
			envMap: map[string]string{
				"INPUT_TITLE": "What name should be given to the barista?", "INPUT_TIMEOUT": "240",
				"INPUT_INTERACTIVE":     "fields:\n  - label: name\n    properties:\n      display: name\n      type: text\n      description: Name of the user\n      maxLength: 20\n      required: false\n",
				"INPUT_GITHUB-TOKEN":    "github-secret-token",
				"INPUT_NGROK-AUTHTOKEN": "ngrok-secret-token",
			},
			expectedConfig: config.Config{
				Timeout: 240,
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
				NotifierSlackBotName:            "",
				NotifierDiscordEnabled:          false,
				NotifierDiscordWebhook:          "secret-webhook",
				NotifierDiscordUsernameOverride: "",
				GithubToken:                     "github-secret-token",
				NgrokAuthtoken:                  "ngrok-secret-token",
			},
			expectedOutput: "::debug::Title input provided: What name should be given to the barista?\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n::add-mask::ngrok-secret-token\n",
			expectedError:  nil,
		},
		{
			name: "successful - created base config no debug messaging",
			preRun: func() {
			},
			envMap: map[string]string{
				"INPUT_TITLE":           "",
				"INPUT_INTERACTIVE":     "fields:\n  - label: name\n    properties:\n      display: name\n      type: text\n      description: Name of the user\n      maxLength: 20\n      required: false\n",
				"INPUT_GITHUB-TOKEN":    "github-secret-token",
				"INPUT_NGROK-AUTHTOKEN": "ngrok-secret-token",
			},
			expectedConfig: config.Config{
				Timeout: 300,
				Title:   "",
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
				NotifierSlackBotName:            "",
				NotifierDiscordEnabled:          false,
				NotifierDiscordWebhook:          "secret-webhook",
				NotifierDiscordUsernameOverride: "",
				GithubToken:                     "github-secret-token",
				NgrokAuthtoken:                  "ngrok-secret-token",
			},
			expectedOutput: "::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n::add-mask::ngrok-secret-token\n",
			expectedError:  nil,
		},
		{
			name: "successful - created base config from input square bracket array",
			preRun: func() {
			},
			envMap: map[string]string{
				"INPUT_TITLE":           "Where should application be deployed?",
				"INPUT_INTERACTIVE":     "fields:\n  - label: deployment-environment\n    properties:\n      display: Environment names\n      type: select\n      choices: ['option', 'option2', 'option3']\n",
				"INPUT_GITHUB-TOKEN":    "github-secret-token",
				"INPUT_NGROK-AUTHTOKEN": "ngrok-secret-token",
			},
			expectedConfig: config.Config{
				Timeout: 300,
				Title:   "Where should application be deployed?",
				Fields: &fields.Fields{
					Fields: []fields.Field{
						{
							Label: "deployment-environment",
							Properties: fields.FieldProperties{
								Display:  "Environment names",
								Type:     "select",
								Choices:  []string{"option", "option2", "option3"},
								Required: false,
							},
						},
					},
				},
				NotifierSlackEnabled:            false,
				NotifierSlackToken:              "xoxb-secret-token",
				NotifierSlackChannel:            "#notificatins",
				NotifierSlackBotName:            "",
				NotifierDiscordEnabled:          false,
				NotifierDiscordWebhook:          "secret-webhook",
				NotifierDiscordUsernameOverride: "",
				Action:                          nil,
				GithubToken:                     "github-secret-token",
				NgrokAuthtoken:                  "ngrok-secret-token",
			},
			expectedOutput: "::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::debug::Title input provided: Where should application be deployed?\n::add-mask::xoxb-secret-token\n::add-mask::secret-webhook\n::add-mask::github-secret-token\n::add-mask::ngrok-secret-token\n",
			expectedError:  nil,
		},
		{
			name: "failed - invalid inputs fields passed",
			preRun: func() {
			},
			envMap: map[string]string{
				"INPUT_TITLE":           "Where should application be deployed?",
				"INPUT_INTERACTIVE":     "cas:\n  - label: deployment-environment\n    properties:\n      display: Environment names\n      type: select\n      choices: ['option', 'option2', 'option3']\n",
				"INPUT_GITHUB-TOKEN":    "github-secret-token",
				"INPUT_NGROK-AUTHTOKEN": "ngrok-secret-token",
			},
			expectedConfig: config.Config{},
			expectedOutput: "::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::debug::Title input provided: Where should application be deployed?\n::error::No fields provided\n::error::Can't convert the 'fields' input to a valid fields config: cas:%0A  - label: deployment-environment%0A    properties:%0A      display: Environment names%0A      type: select%0A      choices: ['option', 'option2', 'option3']\n",
			expectedError:  errors.ErrMalformedFieldsInputDataProvided,
		},
		{
			name: "failed - unsupported type passed",
			preRun: func() {
			},
			envMap: map[string]string{
				"INPUT_TITLE": "Where should application be deployed?", "INPUT_INTERACTIVE": "fields:\n  - label: deployment-environment\n    properties:\n      display: Environment names\n      type: options\n      choices: ['option', 'option2', 'option3']\n",
				"INPUT_GITHUB-TOKEN":    "github-secret-token",
				"INPUT_NGROK-AUTHTOKEN": "ngrok-secret-token",
			},
			expectedConfig: config.Config{
				Timeout: 300,
				Title:   "Where should application be deployed?",
				Fields: &fields.Fields{
					Fields: []fields.Field{
						{
							Label: "deployment-environment",
							Properties: fields.FieldProperties{
								Display:  "Environment names",
								Type:     "options",
								Choices:  []string{"option", "option2", "option3"},
								Required: false,
							},
						},
					},
				},
				NotifierSlackEnabled:            false,
				NotifierSlackToken:              "xoxb-secret-token",
				NotifierSlackChannel:            "#notificatins",
				NotifierSlackBotName:            "",
				NotifierDiscordEnabled:          false,
				NotifierDiscordWebhook:          "secret-webhook",
				NotifierDiscordUsernameOverride: "",
				Action:                          nil,
				GithubToken:                     "github-secret-token",
				NgrokAuthtoken:                  "ngrok-secret-token",
			},
			expectedOutput: "::debug::The timeout was not provided, will use the default timeout of 300 seconds\n::debug::Title input provided: Where should application be deployed?\n::error::Invalid field type 'options' provided for field 'deployment-environment'. Valid field types are: text, textarea, number, boolean, select, multiselect\n::error::Can't convert the 'fields' input to a valid fields config: fields:%0A  - label: deployment-environment%0A    properties:%0A      display: Environment names%0A      type: options%0A      choices: ['option', 'option2', 'option3']\n",
			expectedError:  errors.ErrMalformedFieldsInputDataProvided,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			actionLog := bytes.NewBuffer(nil)

			test.preRun()

			getenv := func(key string) string {
				return test.envMap[key]
			}

			action := githubactions.New(
				githubactions.WithWriter(actionLog),
				githubactions.WithGetenv(getenv),
			)

			cfg, inputsErr := config.NewFromInputs(action)
			if inputsErr != nil {
				assert.Equal(t, test.expectedOutput, actionLog.String())
				assert.Equal(t, test.expectedError, inputsErr)
			}

			if inputsErr == nil {
				assert.NotNil(t, cfg.Action)
				assert.Equal(t, test.expectedOutput, actionLog.String())

				// Make config's action nil for comparison
				cfg.Action = nil

				assert.EqualValues(t, test.expectedConfig, *cfg)
			}

		})
	}
}
