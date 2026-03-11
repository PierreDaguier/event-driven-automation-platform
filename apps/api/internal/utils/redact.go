package utils

import "strings"

var sensitiveFields = []string{"password", "secret", "token", "authorization", "api_key", "apikey"}

func RedactPayload(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		lower := strings.ToLower(k)
		if isSensitive(lower) {
			out[k] = "***redacted***"
			continue
		}
		switch typed := v.(type) {
		case map[string]any:
			out[k] = RedactPayload(typed)
		case []any:
			out[k] = redactSlice(typed)
		default:
			out[k] = v
		}
	}
	return out
}

func redactSlice(items []any) []any {
	out := make([]any, 0, len(items))
	for _, v := range items {
		switch typed := v.(type) {
		case map[string]any:
			out = append(out, RedactPayload(typed))
		case []any:
			out = append(out, redactSlice(typed))
		default:
			out = append(out, v)
		}
	}
	return out
}

func isSensitive(key string) bool {
	for _, needle := range sensitiveFields {
		if strings.Contains(key, needle) {
			return true
		}
	}
	return false
}
