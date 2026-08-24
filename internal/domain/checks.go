package domain

import (
	"fmt"
	"sort"
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
	if err != nil {
		return OutcomeFail, fmt.Sprintf("value must be numeric for %s", item.Name)
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
	sort.Slice(template.Items, func(i, j int) bool {
		return template.Items[i].ID < template.Items[j].ID
	})
	sort.Slice(results, func(i, j int) bool {
		return results[i].ItemID < results[j].ItemID
	})
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
		if outcome == OutcomeFail && result.Passed != nil && *result.Passed {
			return fmt.Errorf("%w: result marked passed despite threshold failure", ErrInvalid)
		}
	}
	return nil
}

func NextRun(now time.Time, trigger TriggerType, intervalHours int) time.Time {
	if trigger == TriggerEvent || intervalHours <= 0 {
		return now
	}
	return now.Add(time.Duration(intervalHours) * time.Hour)
}
