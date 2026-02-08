package notifications

import "strings"

// Email notification templates with placeholder tokens.
// Use FormatTemplate to substitute: %NAME%, %REQUEST_NAME%, %ASSIGNED_DATE%, %SLA_HOURS%, %TREATED_DATE%.
const (
	SlaGenericNewRequest = `<p>Dear %NAME%</p>` +
		`<p>A request on %REQUEST_NAME% has been assigned to you on %ASSIGNED_DATE%.</p>` +
		`<p>Kindly logon to your account on Performance Management System to treat this request ` +
		`as soon as possible to avoid breaching SLA of %SLA_HOURS% hours after initiation.</p>` +
		`<p>Thank you, </br>CBN PMS></p>`

	GenericNewRequest = `<p>Dear %NAME%</p>` +
		`<p>A request on %REQUEST_NAME% has been assigned to you on %ASSIGNED_DATE%</p>` +
		`<p>Kindly logon to your account on Performance Management System to treat this request.</p>` +
		`<p>Thank you, </br>CBN PMS></p>`

	AssignerGenericNewRequest = `<p>Dear %NAME%</p>` +
		`<p>Your request on %REQUEST_NAME% has been initated on %ASSIGNED_DATE%</p>` +
		`<p>Thank you, </br>CBN PMS></p>`

	UpdateRequest = `<p>Dear %NAME%</p>` +
		`<p>You have treated the  request on %REQUEST_NAME% on %TREATED_DATE%</p>` +
		`<p>Thank you, </br>CBN PMS></p>`

	AssignerUpdateRequest = `<p>Dear %NAME%</p>` +
		`<p>Your request on %REQUEST_NAME% has been treated on %TREATED_DATE%</p>` +
		`<p>Thank you, </br>CBN PMS></p>`
)

// FormatTemplate replaces placeholder tokens in a template string.
// The replacements map keys should include the percent delimiters,
// e.g. {"%NAME%": "John Doe", "%REQUEST_NAME%": "Annual Review"}.
func FormatTemplate(template string, replacements map[string]string) string {
	if len(replacements) == 0 {
		return template
	}

	// Build old/new pairs for strings.NewReplacer.
	pairs := make([]string, 0, len(replacements)*2)
	for placeholder, value := range replacements {
		pairs = append(pairs, placeholder, value)
	}

	return strings.NewReplacer(pairs...).Replace(template)
}
