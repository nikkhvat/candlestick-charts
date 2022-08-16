package models

type Data struct {
	Symbol string  `json:"symbol"` // * The name of the currency pair, for example EURUSD or LTCUSD
	Ts     string  `json:"ts"`     // * TimeStamp
	Bid    float64 `json:"bid"`    // * Lowest prace
	Ask    float64 `json:"ask"`    // * Highest prace
	Mid    float64 `json:"mid"`    // * Arithmetic between bid and ask
}

type Candle struct {
	Open      float64 `json:"open"`      // * The value during the opening of the candle (the close value of the previous candle)
	Close     float64 `json:"close"`     // * The value during the closing of the candle (the next candle starts with this value)
	High      float64 `json:"high"`      // * The value where the candle was the maximum is the highest
	Low       float64 `json:"low"`       // * The value where the candle was the maximum is the lowest
	Timestamp int64   `json:"timestamp"` // * TimeStamp
}
