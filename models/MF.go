package models

import (
	"time"
)

type MFPriceData struct {
	Timestamps    time.Time
	Price         float32
	PercentChange float32
}

type MFHoldingsData struct {
	Timestamps     time.Time
	TotalValue     float64
	TotalUnitsHeld float64
	Transaction    float64
}

func (m MFPriceData) GetTime() time.Time {
	return m.Timestamps
}
