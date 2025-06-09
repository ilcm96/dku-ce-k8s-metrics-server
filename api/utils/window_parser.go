package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type WindowSpec struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"` // "s", "m", "h"
}

// ParseWindow 는 윈도우 문자열을 파싱하여 WindowSpec을 반환합니다.
// 예: "30s", "5m", "2h"
func ParseWindow(window string) (*WindowSpec, error) {
	if window == "" {
		return nil, fmt.Errorf("window parameter is empty")
	}

	// 정규식: 숫자 + 단위
	re := regexp.MustCompile(`^(\d+)([smh])$`)
	matches := re.FindStringSubmatch(window)

	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid window format: %s (expected format: <number><unit>, e.g., 30s, 5m, 2h)", window)
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid window value: %s", matches[1])
	}

	if value <= 0 {
		return nil, fmt.Errorf("window value must be positive: %d", value)
	}

	unit := matches[2]

	// 단위별 최대값 검증
	if err := validateWindowLimits(value, unit); err != nil {
		return nil, err
	}

	return &WindowSpec{
		Value: value,
		Unit:  unit,
	}, nil
}

// validateWindowLimits 는 윈도우 크기 제한을 검증합니다.
func validateWindowLimits(value int, unit string) error {
	switch unit {
	case "s":
		if value > 3600 { // 최대 1시간
			return fmt.Errorf("maximum window for seconds is 3600 (1 hour), got: %d", value)
		}
	case "m":
		if value > 1440 { // 최대 24시간
			return fmt.Errorf("maximum window for minutes is 1440 (24 hours), got: %d", value)
		}
	case "h":
		if value > 168 { // 최대 7일
			return fmt.Errorf("maximum window for hours is 168 (7 days), got: %d", value)
		}
	default:
		return fmt.Errorf("unsupported window unit: %s", unit)
	}
	return nil
}

// ToDuration 는 WindowSpec을 time.Duration으로 변환합니다.
func (w *WindowSpec) ToDuration() time.Duration {
	switch w.Unit {
	case "s":
		return time.Duration(w.Value) * time.Second
	case "m":
		return time.Duration(w.Value) * time.Minute
	case "h":
		return time.Duration(w.Value) * time.Hour
	default:
		return 0
	}
}

// GetStartTime 는 종료 시간으로부터 시작 시간을 계산합니다.
func (w *WindowSpec) GetStartTime(endTime time.Time) time.Time {
	duration := w.ToDuration()
	return endTime.Add(-duration)
}

// String 는 WindowSpec의 문자열 표현을 반환합니다.
func (w *WindowSpec) String() string {
	return fmt.Sprintf("%d%s", w.Value, w.Unit)
}
