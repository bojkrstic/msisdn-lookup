package lookup

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type ruleSet struct {
	Countries []CountryRule `json:"countries"`
}

type CountryRule struct {
	Name          string         `json:"name"`
	Codes         []string       `json:"codes"`
	MinLength     int            `json:"minLength"`
	MaxLength     int            `json:"maxLength"`
	TypeRules     []TypeRule     `json:"typeRules"`
	OperatorRules []OperatorRule `json:"operatorRules"`
}

type TypeRule struct {
	Prefix      string `json:"prefix"`
	Type        string `json:"type"`
	Explanation string `json:"explanation"`
}

type OperatorRule struct {
	Prefix      string `json:"prefix"`
	Operator    string `json:"operator"`
	Explanation string `json:"explanation"`
}

type operatorMetadata struct {
	Name        string
	Explanation string
}

var (
	loadOnce             sync.Once
	loadErr              error
	countryByPrefix      map[string]*CountryRule
	maxCountryPrefixLen  int
	operatorByPrefix     map[string]operatorMetadata
	maxOperatorPrefixLen int
)

func init() {
	loadOnce.Do(func() {
		loadErr = loadRuleData()
	})
	if loadErr != nil {
		panic(loadErr)
	}
}

func loadRuleData() error {
	data, err := os.ReadFile(resolveRulesPath())
	if err != nil {
		return fmt.Errorf("lookup: unable to load rules: %w", err)
	}

	var set ruleSet
	if err := json.Unmarshal(data, &set); err != nil {
		return fmt.Errorf("lookup: unable to parse rules: %w", err)
	}

	tmpCountryByPrefix := make(map[string]*CountryRule)
	tmpOperatorByPrefix := make(map[string]operatorMetadata)
	tmpMaxCountryPrefixLen := 0
	tmpMaxOperatorPrefixLen := 0

	for i := range set.Countries {
		country := &set.Countries[i]
		for _, code := range country.Codes {
			if code == "" {
				continue
			}
			tmpCountryByPrefix[code] = country
			if l := len(code); l > tmpMaxCountryPrefixLen {
				tmpMaxCountryPrefixLen = l
			}
		}
		for _, opRule := range country.OperatorRules {
			if opRule.Prefix == "" {
				continue
			}
			tmpOperatorByPrefix[opRule.Prefix] = operatorMetadata{
				Name:        opRule.Operator,
				Explanation: opRule.Explanation,
			}
			if l := len(opRule.Prefix); l > tmpMaxOperatorPrefixLen {
				tmpMaxOperatorPrefixLen = l
			}
		}
	}

	if len(tmpCountryByPrefix) == 0 {
		return errors.New("lookup: no country prefixes loaded")
	}

	countryByPrefix = tmpCountryByPrefix
	operatorByPrefix = tmpOperatorByPrefix
	maxCountryPrefixLen = tmpMaxCountryPrefixLen
	maxOperatorPrefixLen = tmpMaxOperatorPrefixLen
	return nil
}

func resolveRulesPath() string {
	if envPath := os.Getenv("LOOKUP_RULES_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
	}

	candidates := []string{
		"rules.json",
		filepath.Join("lookup", "rules.json"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// fallback to executable directory
	if exe, err := os.Executable(); err == nil {
		if dir := filepath.Dir(exe); dir != "" {
			candidate := filepath.Join(dir, "rules.json")
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}

	// default to repo-relative path to trigger useful error later
	return filepath.Join("lookup", "rules.json")
}
