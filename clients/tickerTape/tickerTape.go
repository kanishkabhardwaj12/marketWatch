package tickerTape

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	summaryUrl = "https://api.tickertape.in/mutualfunds/%s/summary"
)

var ttSymbolISINMap = make(map[TtSymbol]string)

type TtSymbol string

func (t TtSymbol) String() string {
	return string(t)
}

func (t TtSymbol) GetMFSummary() (*MFSummary, error) {
	url := fmt.Sprintf(summaryUrl, t)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	k := MFSummaryResp{}
	err = json.Unmarshal(body, &k)

	return &k.Data.Meta, err
}
