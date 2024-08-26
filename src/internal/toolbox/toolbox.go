package toolbox

import (
	"fmt"
	"regexp"
	"strings"
)

// StringConvertToKebabCase returns a string in kebab case format
func StringConvertToKebabCase(text string) (string, error) {

	// Trim string
	text = strings.TrimSpace(text)

	// Remove all the special characters
	reg, err := regexp.Compile(`[^a-zA-Z0-9\\s]+`)
	if err != nil {
		return "", err
	}
	cleanedText := reg.ReplaceAllString(text, " ")

	// Make sure it's lower case and remove double space
	cleanedText = strings.ToLower(
		StringRemoveMultiSpace(
			cleanedText,
		),
	)

	// Remove Spaces for hyphen
	cleanedText = strings.ReplaceAll(cleanedText, " ", "-")

	return cleanedText, nil
}

// StringRemoveSpecialCharactersWith removes all special characters from the given text
// string and replaces them with the provided replaceValue string.
func StringRemoveSpecialCharactersWith(text, replaceValue string) string {

	// Todo: figure out the best way to include ` in the special char group
	reg, err := regexp.Compile(`[!"'#%&,:;<>=@{}~\$\(\)\*\+\/\\\?\[\]\^\|]+`)
	if err != nil {
		return ""
	}
	cleanedText := reg.ReplaceAllString(text, replaceValue)
	return cleanedText
}

// StringConvertToSnakeCase subsitutes all instances of a space with an underscore
func StringConvertToSnakeCase(s string) string {

	s = StringRemoveMultiSpace(s)

	return strings.Replace(s, " ", "_", -1)
}

// StringStandardisedToUpper returns a string with no explicit spacing strategy
// that is all uppercase and standardised.
func StringStandardisedToUpper(s string) string {
	s = strings.ToUpper(s)

	return StringRemoveMultiSpace(strings.TrimSpace(strings.ReplaceAll(s, "’", "'")))
}

// StringStandardisedToLower returns a string with no explicit spacing strategy
// that is all lowercase and standardised.
func StringStandardisedToLower(s string) string {
	s = strings.ToLower(s)

	return StringRemoveMultiSpace(strings.TrimSpace(strings.ReplaceAll(s, "’", "'")))
}

// StringRemoveMultiSpace subsitutes all multispace with single space
func StringRemoveMultiSpace(s string) string {
	multipleSpaceRegex := regexp.MustCompile(`\s\s+`)

	return multipleSpaceRegex.ReplaceAllString(s, " ")
}

// StringInSlice checks to see if string is within slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// SecondsToMinutes converts the given number of seconds to a string representation
// in the format "minutes:seconds". For example, 125 seconds would be converted
// to the string "2:05".
func SecondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	str := fmt.Sprintf("%d:%d", minutes, seconds)
	return str
}
