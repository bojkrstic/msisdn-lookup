package lookup

import "testing"

func TestCountry(t *testing.T) {
	cases := []struct {
		msisdn string
		want   string
	}{
		{"+393383260866", "Italy"},
		{"+38164123456", "Serbia"},
		{"+306941234567", "Greece"},
		{"+41712345678", "Switzerland"},
	}

	for _, tc := range cases {
		if got := Country(tc.msisdn); got != tc.want {
			t.Fatalf("Country(%s) = %s, want %s", tc.msisdn, got, tc.want)
		}
	}
}

func TestNumberType(t *testing.T) {
	cases := []struct {
		msisdn string
		want   string
	}{
		{"+393383260866", "mobile"},
		{"+390636918899", "fixed"},
		{"+38164123456", "mobile"},
		{"+38111345678", "fixed"},
		{"+306941234567", "mobile"},
		{"+302112345678", "fixed"},
	}

	for _, tc := range cases {
		if got := NumberType(tc.msisdn); got != tc.want {
			t.Fatalf("NumberType(%s) = %s, want %s", tc.msisdn, got, tc.want)
		}
	}
}

func TestIsValidLength(t *testing.T) {
	cases := []struct {
		msisdn string
		want   bool
	}{
		{"+393383260866", true},
		{"+38164123456", true},
		{"+390636", false},
		{"+3816", false},
		{"+306941234567", true},
		{"+3069", false},
	}

	for _, tc := range cases {
		if got := IsValidLength(tc.msisdn); got != tc.want {
			t.Fatalf("IsValidLength(%s) = %v, want %v", tc.msisdn, got, tc.want)
		}
	}
}

func TestOperator(t *testing.T) {
	cases := []struct {
		msisdn string
		want   string
	}{
		{"+381601234567", "A1 Serbia (original range)"},
		{"+381621234567", "Yettel Serbia (original range)"},
		{"+381641234567", "Telekom Srbija (mts original range)"},
		{"+381671234567", "Globaltel Serbia (MVNO range)"},
		{"+393383260866", "TIM Italy (338 prefix)"},
		{"+393491234567", "Vodafone Italy (349 prefix)"},
		{"+390612345678", "Italy fixed (Rome 06)"},
		{"+41791234567", "Swisscom Mobile (079 prefix)"},
		{"+41761234567", "Sunrise UPC Switzerland (076 prefix)"},
		{"+41221234567", "Switzerland fixed (Geneva 22)"},
		{"+306971234567", "Cosmote Greece (697 prefix)"},
		{"+302310669985", "Greek fixed (OTE - Thessaloniki)"},
		{"+302109876543", "Greek fixed (OTE - Athens)"},
	}

	for _, tc := range cases {
		if got := Operator(tc.msisdn); got != tc.want {
			t.Fatalf("Operator(%s) = %s, want %s", tc.msisdn, got, tc.want)
		}
	}
}

func TestAnalyzeProvidesNormalizedView(t *testing.T) {
	resp := Analyze("+30 697 038 91 62")
	if resp.Normalized != "306970389162" {
		t.Fatalf("expected normalized digits to match, got %s", resp.Normalized)
	}
	if resp.E164 != "+306970389162" {
		t.Fatalf("expected E.164 format, got %s", resp.E164)
	}
	if !resp.Valid.DigitsOnly || !resp.Valid.KnownCountryCode || !resp.Valid.LengthOk {
		t.Fatalf("expected all validations to pass: %+v", resp.Valid)
	}
	if resp.CountryConfidence != "high" || resp.TypeConfidence != "medium" || resp.OperatorConfidence == "" {
		t.Fatalf("unexpected confidence payload: %+v", resp)
	}
}

func TestAnalyzeDetectsInvalidCharacters(t *testing.T) {
	resp := Analyze("+30/ 69A")
	if resp.Valid.DigitsOnly {
		t.Fatalf("expected digitsOnly=false for payload with invalid characters")
	}
	if resp.Country != "Unknown" {
		t.Fatalf("country should be unknown for incomplete prefix, got %s", resp.Country)
	}
}
