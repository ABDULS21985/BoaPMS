package helpers

import "encoding/base64"

// DefaultCompanyImageFallback is the fallback image path when no company image is available.
const DefaultCompanyImageFallback = "/uiassets/assets/office-building 1.svg"

// ToBase64DataURI converts a byte slice to a base64-encoded data URI string.
// Returns empty string if data is nil or contentType is empty.
func ToBase64DataURI(data []byte, contentType string) string {
	if len(data) == 0 || contentType == "" {
		return ""
	}
	return "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(data)
}

// ToCompanyImage converts a byte slice to a base64 data URI, or returns
// the fallback path if data is empty.
func ToCompanyImage(data []byte, contentType string, fallback string) string {
	if len(data) == 0 {
		return fallback
	}
	return ToBase64DataURI(data, contentType)
}
