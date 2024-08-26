package fields_test

import (
	"bytes"
	"testing"

	"github.com/boasihq/interactive-inputs/internal/fields"
	"github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/assert"
)

func TestMarshalStringIntoValidFieldsStruct(t *testing.T) {
	tests := []struct {
		name           string
		fieldsString   string
		expectedError  bool
		expectedField  *fields.Fields
		expectedOutput string
	}{
		{
			name:          "success - strip  special chars in label string",
			fieldsString:  "fields:\n  - label: na\\me\n    properties:\n      display: name\n      type: text\n      description: Name of the user\n      maxLength: 20\n      required: false\n",
			expectedError: false,
			expectedField: &fields.Fields{
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
			expectedOutput: "",
		},
		{
			name:          "success - valid YAML string parsed",
			fieldsString:  "fields:\n  - label: name\n    properties:\n      display: name\n      type: text\n      description: Name of the user\n      maxLength: 20\n      required: false\n  - label: age\n    properties:\n      display: age\n      type: text\n      description: Age of the user\n      maxLength: 20\n      required: false\n  - label: city\n    properties:\n      display: city\n      type: text\n      description: City of the user\n      maxLength: 20\n      required: false",
			expectedError: false,
			expectedField: &fields.Fields{
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
					{
						Label: "age",
						Properties: fields.FieldProperties{

							Display:     "age",
							Type:        "text",
							Description: "Age of the user",
							MaxLength:   20,
							Required:    false,
						},
					},
					{
						Label: "city",
						Properties: fields.FieldProperties{

							Display:     "city",
							Type:        "text",
							Description: "City of the user",
							MaxLength:   20,
							Required:    false,
						},
					},
				},
			},
		},
		{
			name:          "Empty string",
			fieldsString:  "",
			expectedError: true,
			expectedField: &fields.Fields{
				Fields: []fields.Field{},
			},
			expectedOutput: "::error::No fields provided\n",
		},
		{
			name:           "Invalid YAML string",
			fieldsString:   "invalid: yaml: :",
			expectedError:  true,
			expectedField:  &fields.Fields{},
			expectedOutput: "::error::Unmarshalling field(s): yaml: mapping values are not allowed in this context\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actionLog := bytes.NewBuffer(nil)

			getenv := func(key string) string {

				mapper := map[string]string{
					"INPUT_EXAMPLE": "getenv",
				}
				return mapper[key]
			}

			action := githubactions.New(
				githubactions.WithWriter(actionLog),
				githubactions.WithGetenv(getenv),
			)

			result, err := fields.MarshalStringIntoValidFieldsStruct(tt.fieldsString, action)

			assert.Equal(t, tt.expectedOutput, actionLog.String())

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedField, result)
			}
		})
	}
}
