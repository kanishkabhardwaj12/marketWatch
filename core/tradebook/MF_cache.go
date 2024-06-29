package tradebook_service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Mryashbhardwaj/marketAnalysis/core/client"
	"github.com/Mryashbhardwaj/marketAnalysis/models"
	"github.com/Mryashbhardwaj/marketAnalysis/utils"
)

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

type MutualFundsTradebook struct {
	AllFunds             map[FundName]ISIN
	ISINToFundName       map[ISIN]FundName
	MutualFundsTradebook map[ISIN][]MutualFundsTrade
}

var mutualFundsTradebook MutualFundsTradebook

func (m MutualFundsTradebook) GetFundNameFromISIN(k ISIN) FundName {
	return m.ISINToFundName[k]
}

var mutualFundsHistory = make(map[ISIN][]models.MFPriceData)

func readMFTradeFiles(tradebookDir string) (map[ISIN][]MutualFundsTrade, map[FundName]ISIN, error) {
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
	tradebook := make(map[ISIN][]MutualFundsTrade)
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
		isin := ISIN(record[1])
		allFunds[symbol] = isin
		if _, ok := tradebook[isin]; !ok {
			tradebook[isin] = []MutualFundsTrade{}
		}
		tradebook[isin] = append(tradebook[isin], MutualFundsTrade{
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

func persistMFInFile(symbol string, trend interface{}) error {
	fileContent, err := json.Marshal(trend)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("./data/trends/mutual_funds/%s.json", symbol)
	return os.WriteFile(fileName, fileContent, os.ModePerm)
}

func BuildMFPriceHistoryCache() error {
	var errorList []string
	for name, isin := range mutualFundsTradebook.AllFunds {
		history, err := client.GetMFHistoryFromMoneyControll(string(isin))
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
