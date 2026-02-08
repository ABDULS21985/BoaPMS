package utils

import "encoding/base64"

// ToBase64DataURL converts a byte slice to a base64-encoded data URL with the
// specified content type. The format is:
//
//	data:<contentType>;base64,<encoded-data>
//
// Returns an empty string if data is nil or empty.
func ToBase64DataURL(data []byte, contentType string) string {
	if len(data) == 0 {
		return ""
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return "data:" + contentType + ";base64," + encoded
}

// ToBase64ImageOrDefault converts a byte slice to a base64 image data URL.
// If data is nil or empty, it returns the provided fallbackURL instead.
// This is useful for displaying user avatars or profile images with a default
// placeholder when no image data is available.
func ToBase64ImageOrDefault(data []byte, contentType string, fallbackURL string) string {
	if len(data) == 0 {
		return fallbackURL
	}
	return ToBase64DataURL(data, contentType)
}
