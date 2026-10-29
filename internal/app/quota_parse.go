package app

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"
)

// QuotaWindow preserves only reported or safely derived values.
type QuotaWindow struct {
	Type            string     `json:"type"`
	Percentage      float64    `json:"percentage"`
	Used            *float64   `json:"used,omitempty"`
	Limit           *float64   `json:"limit,omitempty"`
	Remaining       *float64   `json:"remaining,omitempty"`
	DurationMinutes *int64     `json:"duration_minutes,omitempty"`
	ResetsAt        *time.Time `json:"resets_at,omitempty"`
}

func quotaNumber(row map[string]json.RawMessage, key string) (*float64, error) {
	raw, ok := row[key]
	if !ok {
		return nil, nil
	}
	if string(raw) == "null" {
		return nil, errors.New("null numeric field")
	}
	var n float64
	if err := json.Unmarshal(raw, &n); err != nil || math.IsNaN(n) || math.IsInf(n, 0) || n < 0 {
		return nil, errors.New("invalid numeric field")
	}
	return &n, nil
}

func parseQuotaWindows(body []byte) ([]QuotaWindow, error) {
	var envelope struct {
		Success *bool           `json:"success"`
		Code    json.RawMessage `json:"code"`
		Data    struct {
			Limits []json.RawMessage `json:"limits"`
		} `json:"data"`
	}
	if json.Unmarshal(body, &envelope) != nil || envelope.Success != nil && !*envelope.Success || len(envelope.Data.Limits) == 0 || len(envelope.Data.Limits) > 1024 {
		return nil, errors.New("invalid quota response")
	}
	if len(envelope.Code) > 0 {
		var code any
		if json.Unmarshal(envelope.Code, &code) != nil {
			return nil, errors.New("invalid quota code")
		}
		if code != float64(0) && code != float64(200) && code != "0" && code != "200" {
			return nil, errors.New("quota rejected")
		}
	}
	windows := make([]QuotaWindow, 0, len(envelope.Data.Limits))
	for _, raw := range envelope.Data.Limits {
		var row map[string]json.RawMessage
		if json.Unmarshal(raw, &row) != nil {
			continue
		}
		var kind string
		if json.Unmarshal(row["type"], &kind) != nil || kind != "CREDIT_LIMIT" && kind != "TOKENS_LIMIT" && kind != "TIME_LIMIT" {
			continue
		}
		window, err := parseQuotaWindow(row, kind)
		if err != nil {
			return nil, err
		}
		windows = append(windows, window)
	}
	if len(windows) == 0 {
		return nil, errors.New("no usable quota windows")
	}
	return windows, nil
}

func parseQuotaWindow(row map[string]json.RawMessage, kind string) (QuotaWindow, error) {
	used, e1 := quotaNumber(row, "currentValue")
	limit, e2 := quotaNumber(row, "usage")
	percentage, e3 := quotaNumber(row, "percentage")
	remaining, e4 := quotaNumber(row, "remaining")
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		return QuotaWindow{}, errors.New("invalid quota row")
	}
	if percentage == nil && used != nil && limit != nil && *limit > 0 {
		value := *used / *limit * 100
		percentage = &value
	}
	if percentage == nil || math.IsInf(*percentage, 0) {
		return QuotaWindow{}, errors.New("missing quota percentage")
	}
	if remaining == nil && used != nil && limit != nil {
		value := math.Max(0, *limit-*used)
		remaining = &value
	}
	window := QuotaWindow{Type: kind, Percentage: *percentage, Used: used, Limit: limit, Remaining: remaining}
	if err := applyQuotaTiming(row, &window); err != nil {
		return QuotaWindow{}, err
	}
	return window, nil
}

func applyQuotaTiming(row map[string]json.RawMessage, window *QuotaWindow) error {
	unit, err := quotaNumber(row, "unit")
	if err != nil {
		return err
	}
	count, err := quotaNumber(row, "number")
	if err != nil {
		return err
	}
	if unit != nil && count != nil && *count > 0 && math.Trunc(*count) == *count {
		minutes := int64(0)
		switch *unit {
		case 3:
			minutes = 60
		case 4:
			minutes = 1440
		case 6:
			minutes = 10080
		}
		if minutes > 0 && *count <= float64(int64(^uint64(0)>>1)/minutes) {
			value := minutes * int64(*count)
			window.DurationMinutes = &value
		}
	}
	reset, err := quotaNumber(row, "nextResetTime")
	if err != nil {
		return err
	}
	if reset != nil && *reset > 0 && *reset <= float64(int64(^uint64(0)>>1)) && math.Trunc(*reset) == *reset {
		ms := int64(*reset)
		date := time.Unix(ms/1000, (ms%1000)*int64(time.Millisecond)).UTC()
		if date.Year() >= 1 && date.Year() <= 9999 {
			window.ResetsAt = &date
		}
	}
	return nil
}

// quotaCode recognizes provider authentication codes without exposing response data.
func quotaCode(body []byte) string {
	var envelope struct {
		Code json.RawMessage `json:"code"`
	}
	if json.Unmarshal(body, &envelope) != nil {
		return ""
	}
	var code string
	if json.Unmarshal(envelope.Code, &code) == nil {
		return code
	}
	var number int
	if json.Unmarshal(envelope.Code, &number) == nil {
		return strconv.Itoa(number)
	}
	return ""
}
