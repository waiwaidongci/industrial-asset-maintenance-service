package domain

import "testing"

func TestEvaluateItemRejectsNonFiniteNumbers(t *testing.T) {
	min := 0.0
	item := InspectionItem{ID: "pressure", Name: "pressure", Min: &min}
	for _, value := range []string{"NaN", "+Inf", "-Inf"} {
		outcome, _ := EvaluateItem(item, value)
		if outcome != OutcomeFail {
			t.Fatalf("EvaluateItem(%q) = %q, want fail", value, outcome)
		}
	}
}

func TestValidateTaskResultsRejectsUnknownItems(t *testing.T) {
	template := InspectionTemplate{Items: []InspectionItem{{ID: "known", Required: false}}}
	if err := ValidateTaskResults(template, []TaskResult{{ItemID: "ghost", Value: "ok"}}); err == nil {
		t.Fatal("unknown item was accepted")
	}
}

func TestValidateTaskResultsRejectsInconsistentPassFlag(t *testing.T) {
	template := InspectionTemplate{Items: []InspectionItem{{ID: "known", Required: false}}}
	passed := false
	if err := ValidateTaskResults(template, []TaskResult{{ItemID: "known", Value: "ok", Passed: &passed}}); err == nil {
		t.Fatal("inconsistent pass flag was accepted")
	}
}
