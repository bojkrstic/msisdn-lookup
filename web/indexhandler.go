package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"lookup/lookup"
	"net/http"
)

type validationCheck struct {
	Label  string
	Passed bool
	Icon   string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>MSISDN Lookup</title>
    <script src="https://unpkg.com/htmx.org@2.0.0"></script>
    <style>
        :root {
            --accent: #0a84ff;
            --muted: #6b7280;
            --border: #e5e7eb;
            --card-bg: #fff;
            --error: #dc2626;
            --warning: #fbbf24;
            --success: #16a34a;
        }
        body {
            font-family: system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
            max-width: 960px;
            margin: 0 auto;
            padding: 40px 16px 80px;
            background: #f5f5f5;
            color: #111;
        }
        h1 {
            margin-bottom: 0.5rem;
        }
        .muted {
            color: var(--muted);
            font-size: 0.95rem;
        }
        pre {
            max-height: 250px;
            overflow: auto;
            background: #0f172a;
            color: #e2e8f0;
            padding: 12px;
            border-radius: 8px;
            font-size: 0.85rem;
        }
        .card {
            background: var(--card-bg);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 20px;
            margin-top: 20px;
            box-shadow: 0 8px 30px rgba(0,0,0,0.06);
        }
        .card h2 {
            margin-top: 0;
        }
        label {
            font-weight: 600;
            display: block;
            margin-bottom: 6px;
        }
        input[type="text"], textarea {
            width: 100%;
            padding: 10px 12px;
            border-radius: 8px;
            border: 1px solid var(--border);
            font-size: 1rem;
            box-sizing: border-box;
            margin-bottom: 12px;
        }
        textarea {
            min-height: 160px;
            resize: vertical;
            font-family: monospace;
        }
        button, .secondary-btn {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            padding: 10px 16px;
            background: var(--accent);
            color: #fff;
            border-radius: 8px;
            border: none;
            font-weight: 600;
            cursor: pointer;
        }
        button[disabled], .secondary-btn[disabled] {
            opacity: 0.5;
            cursor: not-allowed;
        }
        button:hover:not([disabled]), .secondary-btn:hover:not([disabled]) {
            background: #0561c9;
        }
        .secondary-btn {
            background: #e5e7eb;
            color: #111;
        }
        .layout {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 20px;
        }
        .badge {
            display: inline-flex;
            align-items: center;
            gap: 4px;
            padding: 2px 10px;
            border-radius: 999px;
            font-size: 0.85rem;
            font-weight: 600;
        }
        .badge.mobile { background: rgba(10,132,255,0.1); color: #0369a1; }
        .badge.fixed { background: rgba(22,163,74,0.1); color: #15803d; }
        .badge.invalid { background: rgba(220,38,38,0.1); color: #991b1b; }
        .confidence-pill {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            font-weight: 600;
            padding: 2px 10px;
            border-radius: 999px;
        }
        .confidence-pill.high { background: rgba(22,163,74,0.1); color: #15803d; }
        .confidence-pill.medium { background: rgba(251,191,36,0.15); color: #a16207; }
        .confidence-pill.low { background: rgba(248,113,113,0.15); color: #b91c1c; }
        .checks {
            list-style: none;
            padding-left: 0;
            margin: 0;
        }
        .checks li {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 6px 0;
            border-bottom: 1px solid var(--border);
            font-size: 0.95rem;
        }
        .checks li:last-child { border-bottom: none; }
        .checks .icon { font-size: 1.2rem; }
        .result-grid {
            width: 100%;
            border-collapse: collapse;
            margin-top: 12px;
            font-size: 0.93rem;
        }
        .result-grid th, .result-grid td {
            border: 1px solid var(--border);
            padding: 8px;
            text-align: left;
        }
        .result-grid th { background: #f8fafc; }
        .mcc-mnc {
            display: flex;
            gap: 16px;
            font-size: 0.85rem;
            color: var(--muted);
            font-weight: 600;
            margin-top: 4px;
            flex-wrap: wrap;
        }
        .copy-btn {
            background: transparent;
            color: var(--accent);
            border: 1px solid transparent;
            padding: 4px 8px;
            border-radius: 6px;
            font-size: 0.85rem;
        }
        .copy-btn:hover { border-color: var(--accent); }
        details {
            margin-top: 16px;
        }
        .recent-list {
            list-style: none;
            padding-left: 0;
            margin: 0;
        }
        .recent-item {
            width: 100%;
            justify-content: flex-start;
            margin-bottom: 8px;
            background: rgba(5,97,201,0.1);
            color: #0369a1;
            border: none;
        }
        .alert {
            padding: 10px 12px;
            border-radius: 8px;
            margin-top: 12px;
            font-weight: 600;
        }
        .alert.error {
            background: rgba(220,38,38,0.12);
            color: #991b1b;
        }
        .actions {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
        }
    </style>
</head>
<body>
    <h1>MSISDN Lookup</h1>
    <p class="muted">HLR-lite style enrichment for Serbia, Italy, Switzerland, Greece and friends. Prefix-driven rules, instant explanations, exportable results.</p>

    <div class="layout">
        <div>
            <div class="card">
                <h2>Single lookup</h2>
                <form id="single-form" hx-get="lookup-view" hx-target="#result" hx-trigger="submit">
                    <label for="msisdn">MSISDN</label>
                    <input type="text" id="msisdn" name="msisdn" placeholder="+30 697 038 91 62" autocomplete="off">
                    <button type="submit">Lookup</button>
                </form>
            </div>
            <div class="card">
                <h2>Recent lookups</h2>
                <p class="muted">Stored in your browser (latest 10). Click to reuse.</p>
                <ul id="recent-items" class="recent-list"></ul>
            </div>
        </div>
        <div>
            <div id="result"></div>
        </div>
    </div>

    <div class="card">
        <h2>Batch lookup</h2>
        <p class="muted">Paste newline-separated MSISDNs or upload from CSV. The backend responds with JSON + an HTML table; you can also export the JSON here.</p>
        <form id="batch-form">
            <label for="batch-input">Numbers</label>
            <textarea id="batch-input" placeholder="+41761234567\n+38163111222\n+393491234567"></textarea>
            <div class="actions">
                <button type="submit" id="run-batch">Run batch</button>
                <button type="button" id="export-json" class="secondary-btn" disabled>Copy JSON</button>
                <button type="button" id="export-csv" class="secondary-btn" disabled>Export CSV</button>
            </div>
        </form>
        <div id="batch-result" class="card" style="margin-top:16px; display:none;"></div>
    </div>

    <script>
        (function() {
            const storageKey = 'lookupRecent';
            const recentList = document.getElementById('recent-items');
            const msisdnInput = document.getElementById('msisdn');
            const batchForm = document.getElementById('batch-form');
            const batchResult = document.getElementById('batch-result');
            const exportJsonBtn = document.getElementById('export-json');
            const exportCsvBtn = document.getElementById('export-csv');
            let lastBatchResults = [];

            document.addEventListener('click', (evt) => {
                const copyBtn = evt.target.closest('[data-copy]');
                if (copyBtn) {
                    const value = copyBtn.getAttribute('data-copy') || '';
                    const defaultLabel = copyBtn.dataset.defaultLabel || 'Copy';
                    copyToClipboard(value).then(() => {
                        copyBtn.textContent = 'Copied!';
                        setTimeout(() => copyBtn.textContent = defaultLabel, 1200);
                    }).catch(() => {
                        alert('Clipboard blocked by browser. Please copy manually.');
                    });
                }
                const recentBtn = evt.target.closest('.recent-item');
                if (recentBtn) {
                    msisdnInput.value = recentBtn.getAttribute('data-value');
                    msisdnInput.focus();
                }
            });

            document.addEventListener('htmx:afterSwap', (evt) => {
                if (evt.target.id !== 'result') {
                    return;
                }
                const card = evt.target.querySelector('.result-card');
                if (!card) {
                    return;
                }
                const payload = card.dataset.json ? JSON.parse(card.dataset.json) : null;
                if (payload) {
                    pushRecent(payload);
                    renderRecent();
                }
            });

            batchForm.addEventListener('submit', async (evt) => {
                evt.preventDefault();
                const payload = document.getElementById('batch-input').value;
                if (!payload.trim()) {
                    alert('Please paste at least one MSISDN');
                    return;
                }
                const res = await fetch('batch', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'text/plain'
                    },
                    body: payload
                });
                if (!res.ok) {
                    const text = await res.text();
                    alert('Batch error: ' + text);
                    return;
                }
                const data = await res.json();
                lastBatchResults = data.results || [];
                exportJsonBtn.disabled = lastBatchResults.length === 0;
                exportCsvBtn.disabled = lastBatchResults.length === 0;
                batchResult.style.display = 'block';
                batchResult.innerHTML = data.table;
            });

            exportJsonBtn.addEventListener('click', () => {
                if (!lastBatchResults.length) return;
                copyToClipboard(JSON.stringify(lastBatchResults, null, 2)).then(() => {
                    exportJsonBtn.textContent = 'JSON copied';
                    setTimeout(() => exportJsonBtn.textContent = 'Copy JSON', 1500);
                }).catch(() => alert('Clipboard blocked by browser.'));
            });

            exportCsvBtn.addEventListener('click', () => {
                if (!lastBatchResults.length) return;
                const headers = ['input','e164','country','numberType','operator'];
                const rows = lastBatchResults.map((row) => headers.map((key) => '"' + (row[key] || '').replaceAll('"','""') + '"').join(','));
                const csv = [headers.join(','), ...rows].join('\n');
                copyToClipboard(csv).then(() => {
                    exportCsvBtn.textContent = 'CSV copied';
                    setTimeout(() => exportCsvBtn.textContent = 'Export CSV', 1500);
                }).catch(() => alert('Clipboard blocked by browser.'));
            });

            function copyToClipboard(value) {
                return new Promise((resolve, reject) => {
                    if (navigator.clipboard && window.isSecureContext) {
                        navigator.clipboard.writeText(value).then(resolve).catch(reject);
                        return;
                    }
                    const textarea = document.createElement('textarea');
                    textarea.value = value;
                    textarea.style.position = 'fixed';
                    textarea.style.opacity = '0';
                    document.body.appendChild(textarea);
                    textarea.focus();
                    textarea.select();
                    try {
                        const successful = document.execCommand('copy');
                        document.body.removeChild(textarea);
                        successful ? resolve() : reject();
                    } catch (err) {
                        document.body.removeChild(textarea);
                        reject(err);
                    }
                });
            }

            function pushRecent(payload) {
                const list = JSON.parse(window.localStorage.getItem(storageKey) || '[]');
                const filtered = list.filter(item => item.input !== payload.input);
                filtered.unshift({
                    input: payload.input,
                    country: payload.country,
                    operator: payload.operator
                });
                window.localStorage.setItem(storageKey, JSON.stringify(filtered.slice(0, 10)));
            }

            function renderRecent() {
                const list = JSON.parse(window.localStorage.getItem(storageKey) || '[]');
                recentList.innerHTML = '';
                if (!list.length) {
                    recentList.innerHTML = '<li class="muted">No history yet.</li>';
                    return;
                }
                list.forEach((item) => {
                    const li = document.createElement('li');
                    li.innerHTML = '<button type="button" class="recent-item" data-value="' + item.input + '">' +
                        item.input + ' · ' + (item.country || 'Unknown') + '</button>';
                    recentList.appendChild(li);
                });
            }

            renderRecent();
        })();
    </script>
</body>
</html>

`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// LookupViewHandler renders a richly formatted card for HTMX swaps.
func LookupViewHandler(w http.ResponseWriter, r *http.Request) {
	msisdn := r.URL.Query().Get("msisdn")
	if msisdn == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<div class="card alert error"><strong>Error:</strong> missing <code>msisdn</code> parameter.</div>`)
		return
	}

	resp := lookup.Analyze(msisdn)
	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<div class="card alert error">Unable to format JSON response.</div>`)
		return
	}

	compactJSON, err := json.Marshal(resp)
	if err != nil {
		compactJSON = jsonBytes
	}

	rawJSON := template.HTMLEscapeString(string(jsonBytes))
	dataJSON := template.HTMLEscapeString(string(compactJSON))

	checks := []validationCheck{
		{Label: "Digits only", Passed: resp.Valid.DigitsOnly, Icon: "🔢"},
		{Label: "Known country code", Passed: resp.Valid.KnownCountryCode, Icon: "🌍"},
		{Label: "Length OK", Passed: resp.Valid.LengthOk, Icon: "📏"},
	}

	confidenceBadges := map[string]string{
		"high":   "✅ High",
		"medium": "🟡 Medium",
		"low":    "🟠 Low",
	}

	badge := func(level string) string {
		if val, ok := confidenceBadges[level]; ok {
			return val
		}
		return level
	}

	validAlert := ""
	if !resp.Valid.KnownCountryCode {
		validAlert = `<div class="alert error">Unknown country code. We can't map this prefix.</div>`
	}

	normalized := resp.Normalized
	if normalized == "" {
		normalized = "—"
	}

	numberTypeClass := "badge invalid"
	switch resp.NumberType {
	case "mobile":
		numberTypeClass = "badge mobile"
	case "fixed":
		numberTypeClass = "badge fixed"
	}
	numberTypeBadge := template.HTML(fmt.Sprintf(`<span class="%s">%s</span>`, numberTypeClass, template.HTMLEscapeString(resp.NumberType)))

	mcc := resp.MCC
	if mcc == "" {
		mcc = "N/A"
	}
	mnc := resp.MNC
	if mnc == "" {
		mnc = "N/A"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
<div class="card result-card" data-json='%s'>
    <h2>Lookup result for %s</h2>
    <p class="muted">Normalized presentation + explanations for confidence and operator guess.</p>
    <ul>
        <li><strong>Input:</strong> %s</li>
        <li><strong>Normalized (digits only):</strong> %s</li>
        <li><strong>E.164 canonical:</strong> %s <button type="button" class="copy-btn" data-copy="%s" data-default-label="Copy">Copy</button></li>
        <li><strong>Country:</strong> %s</li>
        <li><strong>Number type:</strong> %s</li>
        <li><strong>Operator guess:</strong> %s<div class="mcc-mnc"><span>MCC: %s</span><span>MNC: %s</span></div></li>
    </ul>
    <ul class="checks">
        %s
    </ul>
    %s
    <div style="margin-top:16px;">
        <strong>Confidence</strong>
        <div>Country: <span class="confidence-pill high">%s</span></div>
        <div>Type: <span class="confidence-pill medium">%s</span></div>
        <div>Operator: <span class="confidence-pill low">%s</span></div>
    </div>
    <details>
        <summary>How we decided</summary>
        <ul>
            <li>%s</li>
            <li>%s</li>
            <li>%s</li>
        </ul>
    </details>
    <div class="json-block">
        <div class="json-header">
            <p class="muted" style="margin:0;">Raw JSON response</p>
            <button type="button" class="copy-btn" data-copy='%s' data-default-label="Copy JSON">Copy JSON</button>
        </div>
        <pre>%s</pre>
    </div>
</div>
`, dataJSON,
		template.HTMLEscapeString(resp.Input),
		template.HTMLEscapeString(resp.Input),
		template.HTMLEscapeString(normalized),
		template.HTMLEscapeString(resp.E164), template.HTMLEscapeString(resp.E164),
		template.HTMLEscapeString(resp.Country),
		numberTypeBadge,
		template.HTMLEscapeString(resp.Operator),
		template.HTMLEscapeString(mcc),
		template.HTMLEscapeString(mnc),
		renderChecks(checks),
		validAlert,
		template.HTMLEscapeString(badge(resp.CountryConfidence)),
		template.HTMLEscapeString(badge(resp.TypeConfidence)),
		template.HTMLEscapeString(badge(resp.OperatorConfidence)),
		template.HTMLEscapeString(resp.Explain.Country),
		template.HTMLEscapeString(resp.Explain.Type),
		template.HTMLEscapeString(resp.Explain.Operator),
		dataJSON,
		rawJSON)
}

func renderChecks(checks []validationCheck) string {
	var out string
	for _, c := range checks {
		class := "badge fixed"
		status := "OK"
		if !c.Passed {
			class = "badge invalid"
			status = "Needs attention"
		}
		out += fmt.Sprintf(`<li><span class="icon">%s</span><span>%s</span><span class="%s">%s</span></li>`, c.Icon, c.Label, class, status)
	}
	return out
}
