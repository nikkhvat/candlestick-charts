# Forex Data 🚀

forex (golang) - Microservice that works with Streaming API Key
- ✅ Getting real-time data (websocket)
- ✅ Procy web socket server for data 
- ✅ Save data to db (PostgreSQL)
- ✅ History rest request (Generate Candlestick chart) (GET) /?startDate /?endDate /?symbol /?timeframe
- ✅ Get up-to-date data now /?symbol

Raw socket 
```json
{
  "symbol": "LTCUSD",
  "ts": "1660694052026",
  "bid": 61.348,
  "ask": 61.353,
  "mid": 61.350502
}
```

Candlestick chart Data
```json
[{
  "open": 61.5284602,
  "close": 61.5655,
  "high": 61.600502,
  "low": 61.5284602,
  "timestamp": 1660571731680
}]
```



## Requests

### Request example getting historical data

GET -> `https://domain/forex/hostory?start_date=1660567220494&end_date=1660567787724&symbol=ETHUSD&timeframe=60`
- end_date - end date (timestamp)
- start_date - start date (timestamp)
- symbol - currency pair
- timeframe - in seconds

### Example getting up-to-date data now by symbol

GET -> `https://domain/forex/actual?symbol=ETHUSD`
- symbol - currency pair

### Example ping

GET -> `https://domain/forex/ping`

### Example socket

WS -> `wss://domain:3495/forexws`

## Run

in `.env` file:

```.env
key=qwerft_wedrftg
nominals=EURCAD,EURUSD,GBPAUD,AUDCHF,CHFJPY,EURNZD
```

>key - Streaming API Key

>nominals - Сomma-separated list of currency pair

run on port: `3495`

Before launching, create a database in PostgreSql
```sql
CREATE DATABASE forex;
```