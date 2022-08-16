package terminal

import (
	"forex/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func RouudTime(timestamp string, timeFrame int) int64 {

	inttime, err := strconv.Atoi(timestamp)

	if err != nil {
		return 0
	}

	tm := time.Unix(int64(inttime), 0)

	return tm.Truncate(time.Second * time.Duration(timeFrame)).Unix()
}

func DrowCandle(data []models.Data, timeFrame int) []models.Candle {
	var candles []models.Candle

	for index, dataImte := range data {
		if index == 0 {
			// * Creating the first candle
			candles = append(candles, models.Candle{
				Open:      dataImte.Mid,
				High:      dataImte.Mid,
				Low:       dataImte.Mid,
				Close:     dataImte.Mid,
				Timestamp: RouudTime(dataImte.Ts, timeFrame),
			})
		} else {
			last_index := len(candles) - 1

			inttime, err := strconv.Atoi(dataImte.Ts)

			if err != nil {
				return nil
			}

			condition := candles[last_index].Timestamp-int64(inttime) < (-1000 * int64(timeFrame))

			if !condition {
				candles[last_index].Close = dataImte.Mid

				if candles[last_index].High < dataImte.Mid {
					candles[last_index].High = dataImte.Mid
				}

				if candles[last_index].Low > dataImte.Mid {
					candles[last_index].Low = dataImte.Mid
				}
			} else {
				if candles[last_index].Close < candles[last_index].Open && candles[last_index].Low > candles[last_index].Close {
					candles[last_index].Low = candles[last_index].Close
				}

				if candles[last_index].Close > candles[last_index].Open && candles[last_index].High < candles[last_index].Close {
					candles[last_index].High = candles[last_index].Close
				}

				candles = append(candles, models.Candle{
					Open:      candles[last_index].Close,
					Close:     dataImte.Mid,
					High:      candles[last_index].Close,
					Low:       candles[last_index].Close,
					Timestamp: candles[last_index].Timestamp + (1000 * int64(timeFrame)),
				})
			}
		}
	}

	return candles
}

func Hostory(db *gorm.DB, symbol string, startdate, enddate int64, timeframe int) []models.Candle {

	var rawDataItems []models.Data

	db.Where("CAST(nullif(ts, '') as BIGINT) > ? and CAST(nullif(ts, '') as BIGINT) < ? and symbol = ?", startdate, enddate, symbol).Order("ts").Find(&rawDataItems)

	candles := DrowCandle(rawDataItems, timeframe)
	return candles
}

func ActualDate(db *gorm.DB, symbol string) models.Data {
	var rawDataItem models.Data

	db.Where("symbol=? ORDER BY ts DESC LIMIT 1;", symbol).Find(&rawDataItem)

	return rawDataItem
}
