package utils

import (
	"encoding/csv"
	"os"
	"strings"
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
