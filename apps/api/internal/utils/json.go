package utils

import "strings"

func NestedValue(payload map[string]any, dottedPath string) (any, bool) {
	parts := strings.Split(dottedPath, ".")
	var current any = payload
	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		v, exists := m[part]
		if !exists {
			return nil, false
		}
		current = v
	}
	return current, true
}
