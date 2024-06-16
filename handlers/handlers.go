package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	tradebook_service "github.com/Mryashbhardwaj/marketAnalysis/core/tradebook"
	"github.com/Mryashbhardwaj/marketAnalysis/models"
)

// var users = []models.User{
// 	{ID: 1, Name: "John Doe", Email: "john.doe@example.com"},
// 	{ID: 2, Name: "Jane Smith", Email: "jane.smith@example.com"},
// }

func GetTrend(w http.ResponseWriter, r *http.Request) {
	fileBytes, _ := os.ReadFile("./eichermotors_moneycontroll.json")
	k := models.MoneyControlRequest{}
	json.Unmarshal(fileBytes, &k)
	fmt.Println(len(k.T))

	candlePoints := make([]models.CandlePoint, len(k.T))

	for i, timeStamp := range k.T {
		candlePoints[i] = models.CandlePoint{
			Close:      k.C[i],
			High:       k.H[i],
			Volume:     k.V[i],
			Open:       k.O[i],
			Low:        k.L[i],
			Timestamps: time.Unix(timeStamp, 0),
		}
	}
	processedData, _ := json.Marshal(candlePoints)
	w.Write(processedData)
}

func GetMutualFundsList(w http.ResponseWriter, r *http.Request) {
	mfList := tradebook_service.GetMutualFundsList()
	processedData, _ := json.Marshal(mfList)
	w.Write(processedData)
}

// func GetUsers(w http.ResponseWriter, r *http.Request) {
// 	json.NewEncoder(w).Encode(users)
// }

// func GetUser(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	for _, item := range users {
// 		if string(item.ID) == params["id"] {
// 			json.NewEncoder(w).Encode(item)
// 			return
// 		}
// 	}
// 	http.NotFound(w, r)
// }

// func CreateUser(w http.ResponseWriter, r *http.Request) {
// 	var user models.User
// 	_ = json.NewDecoder(r.Body).Decode(&user)
// 	user.ID = len(users) + 1
// 	users = append(users, user)
// 	json.NewEncoder(w).Encode(user)
// }

// func UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	for index, item := range users {
// 		if string(item.ID) == params["id"] {
// 			var user models.User
// 			_ = json.NewDecoder(r.Body).Decode(&user)
// 			user.ID = item.ID
// 			users[index] = user
// 			json.NewEncoder(w).Encode(user)
// 			return
// 		}
// 	}
// 	http.NotFound(w, r)
// }
