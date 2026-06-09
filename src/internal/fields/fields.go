package fields

import (
	"strings"

	"github.com/boasihq/interactive-inputs/internal/errors"
	"github.com/boasihq/interactive-inputs/internal/toolbox"
	"github.com/sethvargo/go-githubactions"
	"gopkg.in/yaml.v2"
)

var (

	// ValidFieldTypes is the list of field types supported by the portal.
	ValidFieldTypes = []string{
		"text",
		"textarea",
		"number",
		"boolean",
		"select",
		"multiselect",
		"file",
		"multifile",
	}
)

// Fields is the top-level YAML shape used by the interactive input definition.
type Fields struct {
	// Fields is the list of input controls that will be displayed to users.
	Fields []Field `yaml:"fields"`
}

// Field represents one input control that will be rendered in the portal.
type Field struct {
	// Label is the input identifier used for form names, output keys, and upload cache mappings.
	Label string `yaml:"label"`

	// Properties define the display, validation, and behaviour settings for the field.
	Properties FieldProperties `yaml:"properties"`
}

// FieldProperties controls how a field is displayed, validated, and processed
// by the portal UI.
type FieldProperties struct {
	// Display is the label to show the user for the field.
	Display string `yaml:"display"`

	// Type is the type of the field, such as "text" or "select".
	Type string `yaml:"type"`

	// Description is the help text to show the user for the field.
	Description string `yaml:"description"`

	// Choices is the list of options to display when Type is "select" or "multiselect".
	Choices []string `yaml:"choices"`

	// Required indicates whether the field must be filled out.
	Required bool `yaml:"required"`

	// MaxLength is the maximum length of the field's value.
	MaxLength int `yaml:"maxLength"`

	// Placeholder is the hint text displayed before the user enters a value.
	Placeholder string `yaml:"placeholder"`

	// NumberMin is the minimum allowed value for number fields.
	NumberMin int `yaml:"minNumber"`

	// NumberMax is the maximum allowed value for number fields.
	NumberMax int `yaml:"maxNumber"`

	// DefaultValue is the initial value displayed in the field.
	DefaultValue string `yaml:"defaultValue"`

	// ReadOnly indicates whether the field should display its value without allowing edits.
	ReadOnly bool `yaml:"readOnly"`

	// DisableAutoCopySelection stops automatically copying the selected option to the clipboard.
	DisableAutoCopySelection bool `yaml:"disableAutoCopySelection"`

	// AcceptedFileTypes limits file and multifile uploads to the listed MIME types or file extensions.
	AcceptedFileTypes []string `yaml:"acceptedFileTypes"`
}

// MarshalStringIntoValidFieldsStruct parses the YAML `interactive` input into
// a normalised Fields value.
//
// It validates that at least one field exists, each field type is supported,
// labels can be converted to kebab case, and labels remain unique after
// normalisation.
func MarshalStringIntoValidFieldsStruct(fieldsString string, action *githubactions.Action) (*Fields, error) {
	var fields Fields
	var detectedFieldLabels []string = make([]string, 0)
	fields.Fields = make([]Field, 0)

	err := yaml.Unmarshal([]byte(fieldsString), &fields)
	if err != nil {

		action.Errorf("Unmarshalling field(s): %s", err)
		return &fields, err
	}

	if len(fields.Fields) == 0 {
		action.Errorf("No fields provided")
		return nil, errors.ErrNoFieldsProvided
	}

	for i, field := range fields.Fields {
		if !toolbox.StringInSlice(
			toolbox.StringStandardisedToLower(field.Properties.Type),
			ValidFieldTypes,
		) {
			action.Errorf(
				"Invalid field type '%s' provided for field '%s'. Valid field types are: %s",
				field.Properties.Type,
				field.Label,
				strings.Join(ValidFieldTypes, ", "),
			)

			return nil, errors.ErrInvalidFieldTypeProvided
		}

		// make sure label is camel case
		labelKebabCase, err := toolbox.StringConvertToKebabCase(
			toolbox.StringRemoveSpecialCharactersWith(field.Label, ""),
		)
		if err != nil || labelKebabCase == "" {
			action.Errorf("Invalid label provided - '%s' is not kebab case compatible", field.Label)
			return nil, errors.ErrInvalidLabelProvided
		}
		fields.Fields[i].Label = labelKebabCase

		// make sure the type is lower case
		fields.Fields[i].Properties.Type = toolbox.StringStandardisedToLower(field.Properties.Type)

		// check if the field label has already been detected
		if toolbox.StringInSlice(field.Label, detectedFieldLabels) {
			action.Errorf("Duplicate field label detected: '%s'", field.Label)
			return nil, errors.ErrDuplicateFieldLabelDetected
		}

		// add the field label to the detected field labels
		detectedFieldLabels = append(detectedFieldLabels, field.Label)
	}

	return &fields, nil
}
