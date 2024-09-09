package fields

import (
	"strings"

	"github.com/boasihq/interactive-inputs/internal/errors"
	"github.com/boasihq/interactive-inputs/internal/toolbox"
	"github.com/sethvargo/go-githubactions"
	"gopkg.in/yaml.v2"
)

var (

	// ValidFieldTypes  is a list of valid field types supported by the action.
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

// Fields is a struct that contains a list of Field structs, which represent the fields in a form to display to users.
// The Fields struct is typically used to define the structure and properties of the fields that will be displayed to users.
// Each Field in the Fields slice has a Label and a list of FieldProperties that define the display, type, and other characteristics of the field.
type Fields struct {
	Fields []Field `yaml:"fields"`
}

// Field represents a field in the Fields struct. It contains a label and a list of field properties.
type Field struct {
	Label      string          `yaml:"label"`
	Properties FieldProperties `yaml:"properties"`
}

// FieldProperties represents the properties of a field in the Fields struct.
// Display is the label to show the user for the field.
// Type is the type of the field, such as "text" or "options".
// Description is a description of the field to show the user.
// Choices is a list of options to display for the field if the Type is "options".
// Required indicates whether the field must be filled out.
// MaxLength is the maximum length of the field's value.
// DisableAutoCopySelection is whether the field should stop automatically coping the selected option to the clipboard (valid fields: select, multiselect).
type FieldProperties struct {
	Display                  string   `yaml:"display"`
	Type                     string   `yaml:"type"`
	Description              string   `yaml:"description"`
	Choices                  []string `yaml:"choices"`
	Required                 bool     `yaml:"required"`
	MaxLength                int      `yaml:"maxLength"`
	Placeholder              string   `yaml:"placeholder"`
	NumberMin                int      `yaml:"minNumber"`
	NumberMax                int      `yaml:"maxNumber"`
	DefaultValue             string   `yaml:"defaultValue"`
	ReadOnly                 bool     `yaml:"readOnly"`
	DisableAutoCopySelection bool     `yaml:"disableAutoCopySelection"`
	AcceptedFileTypes        []string `yaml:"acceptedFileTypes"`
}

// MarshalStringIntoValidFieldsStruct takes a YAML-formatted string representation of a Fields
// struct and unmarshals it into a valid Fields struct. If the unmarshaling is successful and
// the Fields struct contains at least one Field, the function returns a pointer to the Fields
// struct. If the unmarshaling fails or the Fields struct contains no Fields, the function
// returns an error.
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
