package rules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pierre/event-driven-automation-platform/apps/api/internal/models"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/utils"
)

func Matches(conditions []models.Condition, payload map[string]any) bool {
	for _, condition := range conditions {
		value, ok := utils.NestedValue(payload, condition.Field)
		if !ok {
			return false
		}
		if !compare(value, condition.Operator, condition.Value) {
			return false
		}
	}
	return true
}

func compare(left any, operator string, right any) bool {
	switch strings.ToLower(operator) {
	case "eq":
		return toString(left) == toString(right)
	case "neq":
		return toString(left) != toString(right)
	case "contains":
		return strings.Contains(strings.ToLower(toString(left)), strings.ToLower(toString(right)))
	case "gt":
		lf, lok := toFloat(left)
		rf, rok := toFloat(right)
		return lok && rok && lf > rf
	case "lt":
		lf, lok := toFloat(left)
		rf, rok := toFloat(right)
		return lok && rok && lf < rf
	default:
		return false
	}
}

func toString(v any) string {
	return strings.TrimSpace(strings.ToLower(fmt.Sprintf("%v", v)))
}

func toFloat(v any) (float64, bool) {
	s := fmt.Sprintf("%v", v)
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}
