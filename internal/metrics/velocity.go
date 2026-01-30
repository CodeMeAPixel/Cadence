package metrics

import (
	"fmt"
	"time"
)

type VelocityMetrics struct {
	LOCPerMinute float64
}

func CalculateVelocity(loc int64, timeDelta time.Duration) (*VelocityMetrics, error) {
	velocity, err := CalculateVelocityPerMinute(loc, timeDelta)
	if err != nil {
		return nil, err
	}
	return &VelocityMetrics{
		LOCPerMinute: velocity,
	}, nil
}

func CalculateVelocityPerMinute(loc int64, timeDelta time.Duration) (float64, error) {
	if timeDelta <= 0 {
		return 0, fmt.Errorf("invalid time delta: %v (must be positive)", timeDelta)
	}
	return float64(loc) / timeDelta.Minutes(), nil
}
