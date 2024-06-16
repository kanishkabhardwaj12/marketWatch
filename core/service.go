package service

type broker interface {
	getCandleStickTrend(identifier, rangeKey string)               // from money control
	getSimpleTrend(identifier, rangeKey string)                    // from tickertape
	getSimpleTrendComparison(identifier []string, rangeKey string) // from tickertape
}
