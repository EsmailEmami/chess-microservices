package util

import "strings"

// search the input for persian and arabic numbers and transform them to English numbers
// Any other characters will be left unchanged
func TransformFa2En(input string) string {
	outputBuilder := strings.Builder{}

	for _, ch := range input {
		dec := int(ch)

		// Arabic 0-9
		if dec >= 1632 && dec <= 1641 {
			outputBuilder.WriteRune(rune(dec - 1584))
			continue
		}

		// Persian 0-9
		if dec >= 1776 && dec <= 1785 {
			outputBuilder.WriteRune(rune(dec - 1728))
			continue
		}

		// otherwise, do not change anything
		outputBuilder.WriteRune(ch)

	}

	return outputBuilder.String()
}
