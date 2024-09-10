package auto_download

type MarketData struct {
	History []History `json:"history"`
}

type History struct {
	Time  int64   `json:"t"`
	Price float64 `json:"p"`
}

type DataForCSV struct {
	Date      string  `json:"date"`
	StartTime int64   `json:"start_time"`
	Close     float64 `json:"close"`
}
