package tradebook_service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Mryashbhardwaj/marketAnalysis/models"
	"github.com/Mryashbhardwaj/marketAnalysis/utils"
)

type EquityTrade struct {
	Isin               string
	Symbol             string
	TradeDate          string
	Exchange           string
	Segment            string
	Series             string
	TradeType          string
	Auction            string
	Quantity           string
	Price              string
	TradeId            string
	OrderId            string
	OrderExecutionTime string
}

type MutualFundsTrade struct {
	Isin               string
	TradeDate          string
	Exchange           string
	Segment            string
	Series             string
	TradeType          string
	Auction            string
	Quantity           string
	Price              string
	TradeID            string
	OrderID            string
	OrderExecutionTime string
}

type FundName string

func (f FundName) String() string {
	return string(f)
}

type ISIN string
type ScriptName string

type MutualFundsTradebook struct {
	AllFunds             map[FundName]ISIN
	ISINToFundName       map[ISIN]FundName
	MutualFundsTradebook map[FundName][]MutualFundsTrade
}

type EquityTradebook struct {
	AllScripts      []ScriptName
	EquityTradebook map[ScriptName][]EquityTrade
}

var mutualFundsTradebook MutualFundsTradebook

func (m MutualFundsTradebook) GetFundNameFromISIN(k ISIN) FundName {
	return m.ISINToFundName[k]
}

var equityTradebook EquityTradebook

var shareHistory = make(map[ScriptName][]models.EquityPriceData)
var mutualFundsHistory = make(map[ISIN][]models.MFPriceData)

func readMFTradeFiles(tradebookDir string) (map[FundName][]MutualFundsTrade, map[FundName]ISIN, error) {
	// to remove duplidate trade ids
	tradeSet := make(map[string]struct{})
	allFunds := make(map[FundName]ISIN)

	tradeFiles, err := utils.ReadDir(tradebookDir)
	if err != nil {
		return nil, nil, err
	}
	tradebookCSV, err := utils.ReadCSV(tradeFiles)
	if err != nil {
		return nil, nil, err
	}
	tradebook := make(map[FundName][]MutualFundsTrade)
	for _, record := range tradebookCSV {
		if record[0] == "symbol" {
			continue
		}
		if _, ok := tradeSet[record[10]]; ok {
			fmt.Println("found duplicate trade")
			continue
		}
		tradeSet[record[10]] = struct{}{}

		symbol := FundName(record[0])
		if _, ok := tradebook[symbol]; !ok {
			tradebook[symbol] = []MutualFundsTrade{}
		}
		allFunds[symbol] = ISIN(record[1])
		tradebook[symbol] = append(tradebook[symbol], MutualFundsTrade{
			Isin:               record[1],
			TradeDate:          record[2],
			Exchange:           record[3],
			Segment:            record[4],
			Series:             record[5],
			TradeType:          record[6],
			Auction:            record[7],
			Quantity:           record[8],
			Price:              record[9],
			TradeID:            record[10],
			OrderID:            record[11],
			OrderExecutionTime: record[12],
		})

	}
	return tradebook, allFunds, nil
}

func readEquityTradeFiles(tradebookDir string) (map[ScriptName][]EquityTrade, error) {
	// to remove duplidate trade ids
	tradeSet := make(map[string]struct{})

	tradeFiles, err := utils.ReadDir(tradebookDir)
	if err != nil {
		return nil, err
	}
	tradebookCSV, err := utils.ReadCSV(tradeFiles)
	if err != nil {
		return nil, err
	}
	tradebook := make(map[ScriptName][]EquityTrade)
	for _, record := range tradebookCSV {
		if record[1] == "symbol" {
			continue
		}
		if _, ok := tradeSet[record[10]]; ok {
			fmt.Println("found duplicate trade")
			continue
		}
		tradeSet[record[10]] = struct{}{}

		symbol := ScriptName(record[0])
		if _, ok := tradebook[symbol]; !ok {
			tradebook[symbol] = []EquityTrade{}
		}
		tradebook[symbol] = append(tradebook[symbol], EquityTrade{
			Symbol:             record[1],
			TradeDate:          record[2],
			Exchange:           record[3],
			Segment:            record[4],
			TradeType:          record[6],
			Quantity:           record[8],
			Price:              record[9],
			OrderExecutionTime: record[12],
		})
	}
	return tradebook, nil
}

func BuildMFTradeBook(tradebookDir string) error {
	tradeMap, allFunds, err := readMFTradeFiles(tradebookDir)
	if err != nil {
		return err
	}
	mutualFundsTradebook.MutualFundsTradebook = tradeMap
	mutualFundsTradebook.AllFunds = allFunds
	mutualFundsTradebook.ISINToFundName = make(map[ISIN]FundName)
	for k, v := range allFunds {
		mutualFundsTradebook.ISINToFundName[v] = k
	}
	return nil
}

func BuildEquityTradeBook(tradebookDir string) error {
	tradeMap, err := readEquityTradeFiles(tradebookDir)
	if err != nil {
		return err
	}
	var trickers []ScriptName
	for fundName, _ := range tradeMap {
		trickers = append(trickers, fundName)
	}
	equityTradebook.EquityTradebook = tradeMap
	equityTradebook.AllScripts = trickers
	return nil
}

func persistInFile(symbol string, trend interface{}) error {
	fileContent, err := json.Marshal(trend)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("./data/trends/equity/%s.json", symbol)
	return os.WriteFile(fileName, fileContent, os.ModePerm)
}

func persistMFInFile(symbol string, trend interface{}) error {
	fileContent, err := json.Marshal(trend)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("./data/trends/mutual_funds/%s.json", symbol)
	return os.WriteFile(fileName, fileContent, os.ModePerm)
}

func fetchTradeHistories(script ScriptName) ([]models.EquityPriceData, error) {
	k, err := getFromMoneyControll(script)
	if err != nil {
		return nil, err
	}

	candlePoints := make([]models.EquityPriceData, len(k.T))

	for i, timeStamp := range k.T {
		candlePoints[i] = models.EquityPriceData{
			Close:      k.C[i],
			High:       k.H[i],
			Volume:     k.V[i],
			Open:       k.O[i],
			Low:        k.L[i],
			Timestamps: time.Unix(timeStamp, 0),
		}
	}
	return candlePoints, nil
}

func BuildPriceHistoryCache() error {
	var errorList []string
	for _, symbol := range equityTradebook.AllScripts {
		history, err := fetchTradeHistories(symbol)
		if err != nil {
			fmt.Printf("error fetching history for %s, err:%s", symbol, err.Error())
			errorList = append(errorList, fmt.Sprintf("error fetching history for %s, err:%s", symbol, err.Error()))
			continue
		}
		shareHistory[symbol] = history
		err = persistInFile(string(symbol), history)
		if err != nil {
			fmt.Printf("error persisting history for %s, err:%s", symbol, err.Error())
			errorList = append(errorList, fmt.Sprintf("error persisting history for %s, err:%s", symbol, err.Error()))
			continue
		}
	}
	if len(errorList) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(errorList, "\n"))
}

func BuildMFPriceHistoryCache() error {
	var errorList []string
	for name, isin := range mutualFundsTradebook.AllFunds {
		history, err := GetMFHistoryFromMoneyControll(string(isin))
		if err != nil {
			fmt.Printf("error fetching history for MF %s, err:%s", name, err.Error())
			errorList = append(errorList, fmt.Sprintf("error fetching history for %s, err:%s", isin, err.Error()))
			continue
		}
		mutualFundsHistory[isin] = history
		err = persistMFInFile(string(isin), history)
		if err != nil {
			fmt.Printf("error persisting history for %s, err:%s", isin, err.Error())
			errorList = append(errorList, fmt.Sprintf("error persisting history for %s, err:%s", isin, err.Error()))
			continue
		}
	}
	if len(errorList) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(errorList, "\n"))
}

func buildEquityCacheFromFile(symbol ScriptName) ([]models.EquityPriceData, error) {
	fileName := fmt.Sprintf("./data/trends/equity/%s.json", symbol)
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	trend := []models.EquityPriceData{}
	err = json.Unmarshal(fileContent, &trend)
	return trend, err
}

func buildFundsCacheFromFile(symbol ISIN) ([]models.MFPriceData, error) {
	fileName := fmt.Sprintf("./data/trends/mutual_funds/%s.json", symbol)
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	trend := []models.MFPriceData{}
	err = json.Unmarshal(fileContent, &trend)
	return trend, err
}

func BuildPriceHistoryCacheFromFile() error {
	for _, symbol := range equityTradebook.AllScripts {
		history, err := buildEquityCacheFromFile(symbol)
		if err != nil {
			fmt.Printf("error fetching history from file for %s, err:%s\n", symbol, err.Error())
			continue
		}
		shareHistory[symbol] = history
	}

	for _, symbol := range mutualFundsTradebook.AllFunds {
		history, err := buildFundsCacheFromFile(symbol)
		if err != nil {
			fmt.Printf("error fetching history from file for %s, err:%s\n", symbol, err.Error())
			continue
		}
		mutualFundsHistory[symbol] = history
	}
	return nil
}
