package tradebook_service

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	MC "github.com/Mryashbhardwaj/marketAnalysis/clients/moneyControl"
	"github.com/Mryashbhardwaj/marketAnalysis/models"
	"github.com/Mryashbhardwaj/marketAnalysis/utils"
	"github.com/pkg/errors"
)

type MutualFundsTrade struct {
	Isin               string
	TradeDate          time.Time
	Exchange           string
	Segment            string
	Series             string
	TradeType          string
	Auction            string
	Quantity           float64
	Price              float64
	TradeID            string
	OrderID            string
	OrderExecutionTime string
}

func (m MutualFundsTrade) GetTime() time.Time {
	return m.TradeDate
}
func (m MutualFundsTrade) GetPrice() float64 {
	return m.Price
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

// explain the purpose of this function
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
	// avoid magic numbers, instead of 10 it should be index of trade_id
	for _, record := range tradebookCSV {
		if record[0] == "symbol" { //such handling shoul dnot be needed, instead handle using skip header
			continue
		}
		if _, ok := tradeSet[record[10]]; ok {
			continue
		}
		tradeSet[record[10]] = struct{}{}

		symbol := FundName(record[0])
		isin := ISIN(record[1])
		allFunds[symbol] = isin
		if _, ok := tradebook[isin]; !ok {
			tradebook[isin] = []MutualFundsTrade{}
		}

		tradeTime, err := time.Parse(time.DateOnly, record[2])
		if err != nil {
			fmt.Println(err.Error())
		}

		quantityString := record[8]
		quantity, err := strconv.ParseFloat(quantityString, 64)
		if err != nil {
			fmt.Println(err.Error())
		}

		priceString := record[9]
		price, err := strconv.ParseFloat(priceString, 64)
		if err != nil {
			fmt.Println(err.Error())
		}

		tradebook[isin] = append(tradebook[isin], MutualFundsTrade{
			Isin:               record[1],
			TradeDate:          tradeTime,
			Exchange:           record[3],
			Segment:            record[4],
			Series:             record[5],
			TradeType:          record[6],
			Auction:            record[7],
			Quantity:           quantity,
			Price:              price,
			TradeID:            record[10],
			OrderID:            record[11],
			OrderExecutionTime: record[12],
		})

	}

	for isin, _ := range tradebook {
		sort.Slice(tradebook[isin], func(i, j int) bool {
			return tradebook[isin][i].TradeDate.Before(tradebook[isin][j].TradeDate)
		})
	}

	return tradebook, allFunds, nil
}

func BuildMFTradeBook(tradebookDir string) error {
	tradeMap, allFunds, err := readMFTradeFiles(tradebookDir)
	if err != nil {
		return errors.Wrap(err, "unable to read MF trade file")
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
		return errors.Wrap(err, "unable to persist MF trade file")
	}
	fileName := fmt.Sprintf("./data/trends/MF/%s.json", symbol)
	return os.WriteFile(fileName, fileContent, os.ModePerm)
}

func BuildMFPriceHistoryCache() error {
<<<<<<< Updated upstream:core/tradebook/MF_cache.go
	var errorList []string
	for name, isin := range mutualFundsTradebook.AllFunds {
		history, err := MC.GetMFHistoryFromMoneyControll(string(isin))
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
=======
	var (
		mu        sync.Mutex
		errorList []string
		wg        sync.WaitGroup
		semaphore = make(chan struct{}, 5)
	)

	for name, isin := range MutualFundsTradebookCache.AllFunds {
		wg.Add(1)

		go func(name FundName, isin ISIN) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			fmt.Printf("Processing MF: %s\n", name)

			history, err := MC.GetMFHistoryFromMoneyControll(string(isin))
			if err != nil {
				mu.Lock()
				errorList = append(errorList, fmt.Sprintf("error fetching history for MF %s: %s", name, err.Error()))
				mu.Unlock()
				return
			}

			mutualFundsHistory[isin] = history

			if err := persistMFInFile(string(isin), history); err != nil {
				mu.Lock()
				errorList = append(errorList, fmt.Sprintf("error persisting history for MF %s: %s", name, err.Error()))
				mu.Unlock()
			}
		}(name, isin)
>>>>>>> Stashed changes:internal/domain/service/MF_cache.go
	}

	wg.Wait()

	if len(errorList) > 0 {
		return fmt.Errorf(strings.Join(errorList, "\n"))
	}
<<<<<<< Updated upstream:core/tradebook/MF_cache.go
	return fmt.Errorf(strings.Join(errorList, "\n"))
=======
	return nil
>>>>>>> Stashed changes:internal/domain/service/MF_cache.go
}

func buildFundsCacheFromFile(symbol ISIN) ([]models.MFPriceData, error) {
	fileName := fmt.Sprintf("./data/trends/MF/%s.json", symbol)
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	trend := []models.MFPriceData{}
	err = json.Unmarshal(fileContent, &trend)
	return trend, err
}

func BuildMFPriceHistoryCacheFromFile() error {
	for _, isin := range mutualFundsTradebook.AllFunds {
		history, err := buildFundsCacheFromFile(isin)
		if err != nil {
			fmt.Printf("error reading MF cache for %s: %s\n", isin, err)
			continue
		}
		mutualFundsHistory[isin] = history
	}
	return nil
}
