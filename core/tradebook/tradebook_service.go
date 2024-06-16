package tradebook_service

import (
	"fmt"

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

type MutualFundsTradebook struct {
	AllFunds             []FundName
	MutualFundsTradebook map[FundName][]MutualFundsTrade
}

var mutualFundsTradebook MutualFundsTradebook

func readTradeFiles(tradebookDir string) (map[FundName][]MutualFundsTrade, error) {
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
	return tradebook, nil
}

func BuildTradeBook(tradebookDir string) error {
	tradeMap, err := readTradeFiles(tradebookDir)
	if err != nil {
		return err
	}
	var trickers []FundName
	for fundName, _ := range tradeMap {
		trickers = append(trickers, fundName)
	}
	mutualFundsTradebook.MutualFundsTradebook = tradeMap
	mutualFundsTradebook.AllFunds = trickers
	fmt.Println("built cache ")
	fmt.Println(tradeMap)
	return nil
}

func GetMutualFundsList() []FundName {
	return mutualFundsTradebook.AllFunds
}
