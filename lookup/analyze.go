package lookup

import (
	"fmt"
	"strings"
)

const (
	confidenceHigh   = "high"
	confidenceMedium = "medium"
	confidenceLow    = "low"
)

// Analyze performs full lookup with metadata/explanations.
func Analyze(msisdn string) LookupResponse {
	norm := normalizeDetailed(msisdn)
	normalized := norm.digits
	e164 := ""
	if normalized != "" {
		e164 = "+" + normalized
	}

	resp := LookupResponse{
		Input:              msisdn,
		Normalized:         normalized,
		E164:               e164,
		Country:            "Unknown",
		NumberType:         "unknown",
		Operator:           "Unknown",
		Valid:              Validity{DigitsOnly: norm.digitsOnly},
		CountryConfidence:  confidenceHigh,
		TypeConfidence:     confidenceMedium,
		OperatorConfidence: confidenceLow,
	}

	if normalized == "" {
		resp.Explain.Country = "Country: missing digits after normalization"
		resp.Explain.Type = "Type: unable to evaluate without digits"
		resp.Explain.Operator = "Operator: unable to evaluate without digits"
		return resp
	}

	if country, prefix := findCountryRule(normalized); country != nil {
		resp.Country = country.Name
		resp.Valid.KnownCountryCode = true
		resp.Valid.LengthOk = withinLength(normalized, country)
		resp.Explain.Country = fmt.Sprintf("Country: +%s -> %s (country code %s)", prefix, country.Name, prefix)

		local := normalized[len(prefix):]
		resp.NumberType, resp.Explain.Type = resolveType(local, country)
	} else {
		resp.Explain.Country = "Country: prefix not in rules"
		resp.Explain.Type = "Type: country unknown so range can't be interpreted"
	}

	if op, explanation := resolveOperator(normalized); op != "" {
		resp.Operator = op
		resp.Explain.Operator = explanation
	} else {
		resp.Explain.Operator = "Operator guess: no matching prefix rule"
	}

	return resp
}

func findCountryRule(msisdn string) (*CountryRule, string) {
	if maxCountryPrefixLen == 0 {
		return nil, ""
	}
	for l := maxCountryPrefixLen; l >= 1; l-- {
		if len(msisdn) < l {
			continue
		}
		prefix := msisdn[:l]
		if rule, ok := countryByPrefix[prefix]; ok {
			return rule, prefix
		}
	}
	return nil, ""
}

func withinLength(msisdn string, country *CountryRule) bool {
	length := len(msisdn)
	if country.MinLength > 0 && length < country.MinLength {
		return false
	}
	if country.MaxLength > 0 && length > country.MaxLength {
		return false
	}
	return true
}

func resolveType(local string, country *CountryRule) (string, string) {
	for _, rule := range country.TypeRules {
		if rule.Prefix == "" {
			continue
		}
		if strings.HasPrefix(local, rule.Prefix) {
			return rule.Type, fmt.Sprintf("Type: %s -> %s", rule.Prefix, rule.Explanation)
		}
	}

	for _, rule := range country.TypeRules {
		if rule.Prefix == "" {
			if rule.Explanation != "" {
				return rule.Type, fmt.Sprintf("Type fallback: %s", rule.Explanation)
			}
			return rule.Type, "Type fallback rule applied"
		}
	}

	return "unknown", "Type: no matching rules"
}

func resolveOperator(msisdn string) (string, string) {
	if maxOperatorPrefixLen == 0 {
		return "", ""
	}
	for l := maxOperatorPrefixLen; l >= 1; l-- {
		if len(msisdn) < l {
			continue
		}
		prefix := msisdn[:l]
		if op, ok := operatorByPrefix[prefix]; ok {
			explanation := op.Explanation
			if explanation == "" {
				explanation = fmt.Sprintf("Prefix %s matches %s", prefix, op.Name)
			}
			return op.Name, fmt.Sprintf("Operator guess: %s", explanation)
		}
	}
	return "", ""
}
