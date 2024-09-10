package main

import (
	auto_download "download/trump"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var port = "9999"

func main() {
	log.Info("Hello, moksha service is running!")

	cronJob := cron.New()
	cronJob.Start()

	log.Infof("init craw data at %s", "CRON_TZ=Asia/Ho_Chi_Minh 00 08 * * *")
	_, err := cronJob.AddFunc("CRON_TZ=Asia/Ho_Chi_Minh 00 08 * * *", func() {
		_, err := jobDownloadData()
		if err != nil {
			log.Error(err)
		}
	})
	if err != nil {
		log.Error("init job crawl data failed")
	}

	log.Fatal(Run())
}

func Run() error {
	r := mux.NewRouter()

	r.HandleFunc("/manual-download", downloadData).Methods("GET")

	http.Handle("/", r)
	log.Infof("moksha service is running at port %s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func jobDownloadData() ([]auto_download.DataForCSV, error) {
	dataFromMarker, err := crawlDataFromMarket()
	if err != nil {
		return nil, err
	}

	dataForCSV := convertDataForCSV(dataFromMarker)

	err = writeToFile(dataForCSV)
	if err != nil {
		return nil, err
	}
	return dataForCSV, nil
}

func downloadData(w http.ResponseWriter, r *http.Request) {
	rs, err := jobDownloadData()
	if err != nil {
		auto_download.ResponseError(w, err)
		return
	}

	auto_download.ResponseWithJSON(w, 200, rs)
}

func crawlDataFromMarket() (auto_download.MarketData, error) {
	now := time.Now()
	after30day := now.AddDate(0, -3, 0)

	url := fmt.Sprintf("%s%d", "https://clob.polymarket.com/prices-history?market=21742633143463906290569050155826241533067272736897614950488156847949938836455&fidelity=60&startTs=", after30day.Unix())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Errorf("err when post to market, resp: %+v", resp)
		return auto_download.MarketData{}, fmt.Errorf("error code: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	rs := auto_download.MarketData{}
	err := json.Unmarshal(body, &rs)
	if err != nil {
		return auto_download.MarketData{}, err
	}

	log.Infof("response: %+v", rs)

	return rs, nil
}

func convertDataForCSV(original auto_download.MarketData) []auto_download.DataForCSV {
	var result []auto_download.DataForCSV
	for _, d := range original.History {
		result = append(result, auto_download.DataForCSV{
			Date:      auto_download.ConvertSecondToDate(d.Time),
			StartTime: d.Time,
			Close:     d.Price,
		})
	}
	return result
}

func writeToFile(data []auto_download.DataForCSV) error {
	f, err := os.OpenFile("./data/trump_data.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Error("Failed to open log file")
		return err
	}
	defer f.Close()

	//delete old content
	err = f.Truncate(0)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(f)
	defer writer.Flush()

	header := []string{"date", "start_time", "close"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range data {
		var csvRow []string
		csvRow = append(csvRow, r.Date, fmt.Sprint(r.StartTime), fmt.Sprint(r.Close))
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
