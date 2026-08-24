package domain

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type CheckOutcome string

const (
	OutcomePass    CheckOutcome = "pass"
	OutcomeFail    CheckOutcome = "fail"
	OutcomeUnknown CheckOutcome = "unknown"
)

func EvaluateItem(item InspectionItem, value string) (CheckOutcome, string) {
	value = strings.TrimSpace(value)
	if value == "" {
		if item.Required {
			return OutcomeFail, "required value is missing"
		}
		return OutcomeUnknown, "optional value is missing"
	}
	if item.Min == nil && item.Max == nil {
		return OutcomePass, "value recorded"
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(number) || math.IsInf(number, 0) {
		return OutcomeFail, fmt.Sprintf("value must be a finite numeric for %s", item.Name)
	}
	if item.Min != nil && number < *item.Min {
		return OutcomeFail, fmt.Sprintf("value %.2f is below minimum %.2f", number, *item.Min)
	}
	if item.Max != nil && number > *item.Max {
		return OutcomeFail, fmt.Sprintf("value %.2f is above maximum %.2f", number, *item.Max)
	}
	return OutcomePass, "value is within threshold"
}

func ValidateTaskResults(template InspectionTemplate, results []TaskResult) error {
	byID := make(map[string]TaskResult, len(results))
	for _, result := range results {
		if _, exists := byID[result.ItemID]; exists {
			return fmt.Errorf("%w: duplicate result %s", ErrInvalid, result.ItemID)
		}
		byID[result.ItemID] = result
	}
	for _, item := range template.Items {
		result, ok := byID[item.ID]
		if !ok {
			if item.Required {
				return fmt.Errorf("%w: missing required item %s", ErrInvalid, item.ID)
			}
			continue
		}
		outcome, _ := EvaluateItem(item, result.Value)
		if result.Passed != nil {
			if outcome == OutcomeFail && *result.Passed {
				return fmt.Errorf("%w: result marked passed despite threshold failure", ErrInvalid)
			}
			if outcome == OutcomePass && !*result.Passed {
				return fmt.Errorf("%w: result marked failed despite threshold pass", ErrInvalid)
			}
		}
		delete(byID, item.ID)
	}
	for itemID := range byID {
		return fmt.Errorf("%w: unknown result item %s", ErrInvalid, itemID)
	}
	return nil
}

func NextRun(now time.Time, trigger TriggerType, intervalHours int) time.Time {
	if trigger == TriggerEvent || intervalHours <= 0 {
		return now
	}
	return now.Add(time.Duration(intervalHours) * time.Hour)
}
