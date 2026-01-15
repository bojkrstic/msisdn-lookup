package lookup

import (
	"strings"
	"unicode"
)

type normalizedPayload struct {
	digits     string
	digitsOnly bool
}

func normalize(msisdn string) string {
	return normalizeDetailed(msisdn).digits
}

func normalizeDetailed(msisdn string) normalizedPayload {
	trimmed := strings.TrimSpace(msisdn)
	var builder strings.Builder
	builder.Grow(len(trimmed))
	digitsOnly := true

	for _, r := range trimmed {
		switch {
		case unicode.IsDigit(r):
			builder.WriteRune(r)
		case r == '+' && builder.Len() == 0:
			// skip leading plus
		case unicode.IsSpace(r):
			// ignore whitespace completely
		case r == '-' || r == '(' || r == ')' || r == '.':
			// treated as permissive formatting characters
		default:
			digitsOnly = false
		}
	}

	digits := builder.String()
	if strings.HasPrefix(digits, "00") {
		digits = digits[2:]
	}

	return normalizedPayload{digits: digits, digitsOnly: digitsOnly}
}
