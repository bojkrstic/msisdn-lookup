package lookup

// LookupResponse represents a full MSISDN analysis payload.
type LookupResponse struct {
	Input              string   `json:"input"`
	Normalized         string   `json:"normalized"`
	E164               string   `json:"e164"`
	Country            string   `json:"country"`
	NumberType         string   `json:"numberType"`
	Operator           string   `json:"operator"`
	Valid              Validity `json:"valid"`
	CountryConfidence  string   `json:"countryConfidence"`
	TypeConfidence     string   `json:"typeConfidence"`
	OperatorConfidence string   `json:"operatorConfidence"`
	Explain            Explain  `json:"explain"`
}

// Validity captures lightweight client-side style validations.
type Validity struct {
	DigitsOnly       bool `json:"digitsOnly"`
	KnownCountryCode bool `json:"knownCountryCode"`
	LengthOk         bool `json:"lengthOk"`
}

// Explain contains human readable rules that were applied.
type Explain struct {
	Country  string `json:"country"`
	Type     string `json:"type"`
	Operator string `json:"operator"`
}
