package lookup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
)

type batchResponse struct {
	Results []LookupResponse `json:"results"`
	Table   string           `json:"table"`
}

// BatchHandler performs multi lookup on newline separated input.
func BatchHandler(w http.ResponseWriter, r *http.Request) {
	if !enforceAccess(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "batch endpoint expects POST", http.StatusMethodNotAllowed)
		return
	}

	msisdns, err := parseBatchBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results := make([]LookupResponse, 0, len(msisdns))
	for _, value := range msisdns {
		results = append(results, Analyze(value))
	}

	resp := batchResponse{
		Results: results,
		Table:   renderBatchTable(results),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func parseBatchBody(r *http.Request) ([]string, error) {
	reader := io.LimitReader(r.Body, 1<<20) // 1MB cap for safety
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	payload := strings.TrimSpace(string(raw))
	if payload == "" {
		return nil, errors.New("empty batch payload")
	}

	ct := r.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		var list []string
		if err := json.Unmarshal(raw, &list); err != nil {
			return nil, errors.New("invalid JSON array payload")
		}
		return normalizeBatchList(list), nil
	}

	lines := strings.Split(payload, "\n")
	return normalizeBatchList(lines), nil
}

func normalizeBatchList(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

func renderBatchTable(results []LookupResponse) string {
	if len(results) == 0 {
		return `<div class="muted">No inputs processed.</div>`
	}

	var b strings.Builder
	b.WriteString(`<table class="result-grid"><thead><tr>`)
	b.WriteString("<th>#</th><th>Input</th><th>E.164</th><th>Country</th><th>Type</th><th>Operator</th><th>Valid</th>")
	b.WriteString("</tr></thead><tbody>")

	for idx, res := range results {
		country := template.HTMLEscapeString(res.Country)
		numberType := template.HTMLEscapeString(res.NumberType)
		operator := template.HTMLEscapeString(res.Operator)
		input := template.HTMLEscapeString(res.Input)
		e164 := template.HTMLEscapeString(res.E164)
		valid := "No"
		if res.Valid.DigitsOnly && res.Valid.KnownCountryCode && res.Valid.LengthOk {
			valid = "Yes"
		}
		b.WriteString("<tr>")
		b.WriteString(fmt.Sprintf("<td>%d</td>", idx+1))
		b.WriteString("<td>" + input + "</td>")
		b.WriteString("<td>" + e164 + "</td>")
		b.WriteString("<td>" + country + "</td>")
		b.WriteString("<td>" + numberType + "</td>")
		b.WriteString("<td>" + operator + "</td>")
		b.WriteString("<td>" + valid + "</td>")
		b.WriteString("</tr>")
	}

	b.WriteString("</tbody></table>")
	return b.String()
}
