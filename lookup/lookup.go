package lookup

// Country returns the detected country based on dial code prefixes.
func Country(msisdn string) string {
	normalized := normalize(msisdn)
	if normalized == "" {
		return "Unknown"
	}
	if rule, _ := findCountryRule(normalized); rule != nil {
		return rule.Name
	}
	return "Unknown"
}
