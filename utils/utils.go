package utils

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadDir(dir string) ([]string, error) {
	tradeFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var tradeFilesStrings []string
	for _, tf := range tradeFiles {
		tradeFilesStrings = append(tradeFilesStrings, dir+tf.Name())
	}
	return tradeFilesStrings, nil
}

func ReadCSV(tradeFiles []string) ([][]string, error) {
	var tradeFilesCombined [][]string
	for _, tf := range tradeFiles {
		file, err := os.ReadFile(tf)
		if err != nil {
			return nil, err
		}
		r := csv.NewReader(strings.NewReader(string(file)))
		k, err := r.ReadAll()
		if err != nil {
			return nil, err
		}
		tradeFilesCombined = append(tradeFilesCombined, k...)
	}
	return tradeFilesCombined, nil
}

func GetTimeRange(r *http.Request) (time.Time, time.Time, error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" {
		fromStr = "490147200000"
	}
	if toStr == "" {
		toStr = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}

	// Parse the "from" timestamp
	fromMilli, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid 'from' timestamp format. Expected epoch milliseconds")
	}
	from := time.Unix(fromMilli/1000, 0)

	// Parse the "to" timestamp
	toMilli, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid 'to' timestamp format. Expected epoch milliseconds")
	}
	to := time.Unix(toMilli/1000, 0)

	return from, to, nil
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
