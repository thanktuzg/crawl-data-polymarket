package auto_download

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	PatternDate = "01/02/2006 15:04"
)

func ConvertSecondToDate(second int64) string {
	tm := time.Unix(second, 0).UTC()
	return tm.Format(PatternDate)
}

func ResponseWithJSON(w http.ResponseWriter, httpStatusCode int, data interface{}) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write(resp)
	return
}

func ResponseError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
