package util

import "strings"

// Username allowed characters: 0-9 a-z(all lower case)
func AsUsername(input string) (output string) {
	// make the username all lowercase
	input = strings.ToLower(input)
	outputBuilder := strings.Builder{}

	for _, ch := range input {
		dec := int(ch)

		// 0-9
		if dec >= 48 && dec <= 57 {
			outputBuilder.WriteRune(ch)
		}

		// a-z
		if dec >= 97 && dec <= 122 {
			outputBuilder.WriteRune(ch)
		}

	}

	return outputBuilder.String()
}

// Only allow 0-9, a-z , A-Z and persian and arabic letters
// It also trim spaces from begging and end of the input
func AsClearText(input string) (output string) {
	input = strings.Trim(input, " ")

	outputBuilder := strings.Builder{}

	for _, ch := range input {
		dec := int(ch)

		// 0-9
		if dec >= 48 && dec <= 57 {
			outputBuilder.WriteRune(ch)
			continue
		}

		// space
		if dec == 32 /*space*/ || dec == 46 /*.*/ || dec == 45 /*-*/ || dec == 47 /*/*/ || dec == 58 /*:*/ {
			outputBuilder.WriteRune(ch)
			continue
		}

		// A-Z
		if dec >= 65 && dec <= 90 {
			outputBuilder.WriteRune(ch)
			continue
		}

		// a-z
		if dec >= 97 && dec <= 122 {
			outputBuilder.WriteRune(ch)
			continue
		}

		// Persian and Arabic letters
		if dec >= 1567 && dec <= 1919 {
			outputBuilder.WriteRune(ch)
			continue
		}

		outputBuilder.WriteRune(ch)

	}

	ArToFa(outputBuilder.String())

	return outputBuilder.String()
}

// AsCode sanitize input to be a acceptable value to set as Code properties
func AsCode(input string) (output string) {
	input = strings.Trim(input, " ")
	outputBuilder := strings.Builder{}

	for _, ch := range input {
		dec := int(ch)

		// 0-9
		if dec >= 48 && dec <= 57 {
			outputBuilder.WriteRune(ch)
			continue
		}

		// space
		if dec == 32 {
			outputBuilder.WriteRune(ch)
			continue
		}

		// A-Z
		if dec >= 65 && dec <= 90 {
			outputBuilder.WriteRune(ch)
			continue
		}

		// a-z
		if dec >= 97 && dec <= 122 {
			outputBuilder.WriteRune(ch)
			continue
		}

	}

	return outputBuilder.String()
}

func AsPermissionKey(input string) string {
	input = strings.Trim(input, " ")
	outputBuilder := strings.Builder{}

	for _, ch := range input {
		dec := int(ch)

		// 0-9
		if dec >= 48 && dec <= 57 {
			outputBuilder.WriteRune(ch)
			continue
		}

		// _
		if dec == 95 {
			outputBuilder.WriteRune(ch)
			continue
		}

		// a-z
		if dec >= 97 && dec <= 122 {
			outputBuilder.WriteRune(ch)
			continue
		}

	}

	return outputBuilder.String()
}

func AsNumeric(input string) string {
	input = strings.Trim(input, " ")
	outputBuilder := strings.Builder{}

	for _, ch := range input {
		dec := int(ch)
		// 0-9
		if dec >= 48 && dec <= 57 {
			outputBuilder.WriteRune(ch)
			continue
		}
	}

	return outputBuilder.String()
}

func ArToFa(ar string) string {
	s := strings.Trim(ar, " ")
	s = strings.ReplaceAll(s, "ي", "ی")
	s = strings.ReplaceAll(s, "ك", "ک")
	s = strings.ReplaceAll(s, "د", "د")

	s = TransformFa2En(s)

	return s
}
