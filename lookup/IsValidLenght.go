package lookup

func IsValidLength(msisdn string) bool {
	normalized := normalize(msisdn)
	if normalized == "" {
		return false
	}

	if country, _ := findCountryRule(normalized); country != nil {
		return withinLength(normalized, country)
	}

	return false
}
