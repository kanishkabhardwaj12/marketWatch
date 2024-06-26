package models

import (
	"math"
	"time"
)

type TimeGetter interface {
	GetTime() time.Time
}
type MFPriceData struct {
	Timestamps    time.Time
	Price         float32
	PercentChange float32
}

func (m MFPriceData) GetTime() time.Time {
	return m.Timestamps
}

type EquityPriceData struct {
	Timestamps    time.Time
	Open          float32
	Close         float32
	High          float32
	Low           float32
	Volume        float32
	PercentChange float32
}

func (e EquityPriceData) GetTime() time.Time {
	return e.Timestamps
}

type Ticker struct {
	Symbol     string
	Low52Week  float32
	High52Week float32
}

func MomentBinarySearch[V TimeGetter](timestamps []V, target time.Time) int {
	left, right := 0, len(timestamps)-1
	nearestIndex := -1
	minDiff := math.MaxInt64

	for left <= right {
		mid := left + (right-left)/2

		// Check if the target is present at mid
		if timestamps[mid].GetTime().Equal(target) {
			return mid
		}

		// Update the nearest index if the current difference is smaller
		diff := absDuration(timestamps[mid].GetTime().Sub(target))
		if diff < time.Duration(minDiff) {
			minDiff = int(diff)
			nearestIndex = mid
		}

		// If the target is greater, ignore the left half
		if timestamps[mid].GetTime().Before(target) {
			left = mid + 1
		} else {
			// If the target is smaller, ignore the right half
			right = mid - 1
		}
	}
	return nearestIndex
}

// absDuration is a helper function to calculate the absolute value of a time.Duration.
func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
