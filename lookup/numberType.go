package lookup

func NumberType(msisdn string) string {
	normalized := normalize(msisdn)
	if normalized == "" {
		return "unknown"
	}

	country, prefix := findCountryRule(normalized)
	if country == nil {
		return "unknown"
	}

	local := normalized[len(prefix):]
	numberType, _ := resolveType(local, country)
	return numberType
}
