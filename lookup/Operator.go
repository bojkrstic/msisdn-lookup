package lookup

// Operator returns the operator guess based on prefix rules.
func Operator(msisdn string) string {
	normalized := normalize(msisdn)
	if normalized == "" {
		return "Unknown"
	}

	if op, _ := resolveOperator(normalized); op != "" {
		return op
	}

	return "Unknown"
}
