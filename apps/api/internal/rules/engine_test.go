package rules

import (
	"testing"

	"github.com/pierre/event-driven-automation-platform/apps/api/internal/models"
)

func TestMatches(t *testing.T) {
	payload := map[string]any{
		"lead": map[string]any{
			"score":  82,
			"region": "US",
		},
	}

	conditions := []models.Condition{
		{Field: "lead.score", Operator: "gt", Value: 70},
		{Field: "lead.region", Operator: "eq", Value: "us"},
	}

	if !Matches(conditions, payload) {
		t.Fatal("expected conditions to match payload")
	}

	conditions[0].Value = 95
	if Matches(conditions, payload) {
		t.Fatal("expected conditions to fail")
	}
}
