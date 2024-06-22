package models

import "time"

type TradeType int

const (
	TradeTypeBuy = iota
	TradeTypeSell
)

type TradePoint struct {
	Price     float32
	Date      time.Time
	TradeType TradeType
}

type MoneyControlResponse struct {
	C []float32
	H []float32
	L []float32
	O []float32
	T []int64
	V []float32
}

type CandlePoint struct {
	Timestamps    time.Time
	Open          float32
	Close         float32
	High          float32
	Low           float32
	Volume        float32
	PercentChange float32
}

func NewTradePoint(price float32) *TradePoint {

	return &TradePoint{
		TradeType: TradeTypeBuy,
	}
}

type Ticker struct {
	Symbol     string
	Low52Week  float32
	High52Week float32
}
