package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"lookup/lookup"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// 	html := `
	// <!DOCTYPE html>
	// <html>
	// <head>
	//     <meta charset="utf-8">
	//     <title>MSISDN Lookup</title>
	//     <script src="https://unpkg.com/htmx.org@2.0.0"></script>
	// </head>
	// <body>
	//     <h1>MSISDN Lookup</h1>
	//     <form hx-get="/lookup" hx-target="#result" hx-trigger="submit">
	//         <input type="text" name="msisdn" placeholder="+393383260866">
	//         <button type="submit">Lookup</button>
	//     </form>

	//     <pre id="result"></pre>
	// </body>
	// </html>`
	// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// 	fmt.Fprint(w, html)
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>MSISDN Lookup</title>
    <script src="https://unpkg.com/htmx.org@2.0.0"></script>
    <style>
        body {
            font-family: sans-serif;
            max-width: 600px;
            margin: 40px auto;
            padding: 0 16px;
        }
        h1 {
            margin-bottom: 0.5rem;
        }
        .card {
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 16px;
            margin-top: 16px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
        }
        input[type="text"] {
            width: 100%;
            padding: 8px;
            border-radius: 4px;
            border: 1px solid #ccc;
            margin-bottom: 8px;
        }
        button {
            padding: 8px 16px;
            border-radius: 4px;
            border: none;
            cursor: pointer;
            background: #007bff;
            color: white;
            font-weight: 600;
        }
        button:hover {
            background: #0056b3;
        }
        .muted {
            color: #777;
            font-size: 0.9rem;
        }
        .badge {
            display: inline-block;
            border-radius: 999px;
            padding: 2px 8px;
            font-size: 0.8rem;
            background: #eee;
            margin-left: 4px;
        }
        .valid {
            color: #0a7f2e;
        }
        .invalid {
            color: #b30000;
        }
        table.result-grid {
            width: 100%;
            border-collapse: collapse;
            margin-top: 16px;
        }
        table.result-grid th,
        table.result-grid td {
            border: 1px solid #e2e2e2;
            padding: 8px;
            text-align: left;
            font-size: 0.95rem;
        }
        table.result-grid th {
            background: #f8f8f8;
        }
    </style>
</head>
<body>
    <h1>MSISDN Lookup</h1>
    <p class="muted">Enter a number in the E.164 format (e.g. +393383260866) and click "Lookup".</p>

    <form hx-get="/lookup-view" hx-target="#result" hx-trigger="submit">
        <label for="msisdn">MSISDN</label>
        <input type="text" id="msisdn" name="msisdn" placeholder="+393383260866">
        <button type="submit">Lookup</button>
    </form>

    <div id="result"></div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// HTML snippet za rezultat (koristi HTMX)
func LookupViewHandler(w http.ResponseWriter, r *http.Request) {
	msisdn := r.URL.Query().Get("msisdn")
	if msisdn == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<div class="card"><strong>Error:</strong> missing <code>msisdn</code> parameter.</div>`)
		return
	}

	resp := lookup.LookupResponse{
		MSISDN:      msisdn,
		Country:     lookup.Country(msisdn),
		NumberType:  lookup.NumberType(msisdn),
		ValidLength: lookup.IsValidLength(msisdn),
		Operator:    lookup.Operator(msisdn),
	}

	validText := "No"
	validClass := "invalid"
	if resp.ValidLength {
		validText = "Yes"
		validClass = "valid"
	}

	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		jsonBytes = []byte(`{"error":"unable to format JSON"}`)
	}
	rawJSON := template.HTMLEscapeString(string(jsonBytes))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Ovo je HTML koji se ubacuje u <div id="result">
	fmt.Fprintf(w, `
<div class="card">
    <h2>Lookup result for %s</h2>
    <p class="muted">HLR-lite style analysis based purely on prefixes and number length.</p>
    <ul>
        <li><strong>MSISDN:</strong> %s</li>
        <li><strong>Country:</strong> %s</li>
        <li><strong>Number type:</strong> %s</li>
        <li><strong>Valid length:</strong> <span class="%s">%s</span></li>
        <li><strong>Operator:</strong> %s</li>
    </ul>
    <p class="muted">Raw JSON response returned by the API:</p>
    <pre>%s</pre>
    <table class="result-grid">
        <thead>
            <tr>
                <th>MSISDN</th>
                <th>Country</th>
                <th>Number type</th>
                <th>Valid length</th>
                <th>Operator</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
            </tr>
        </tbody>
    </table>
</div>
`, resp.MSISDN, resp.MSISDN, resp.Country, resp.NumberType, validClass, validText, resp.Operator,
		rawJSON, resp.MSISDN, resp.Country, resp.NumberType, validText, resp.Operator)
}
