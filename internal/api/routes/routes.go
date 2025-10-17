package routes

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Mryashbhardwaj/marketAnalysis/internal/api/auth"
	"github.com/Mryashbhardwaj/marketAnalysis/internal/api/handlers"
	"github.com/gorilla/mux"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = filepath.Join(h.staticPath, path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// File does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func SetupRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	photoUploadPath := getEnv("PHOTO_UPLOAD_PATH", "./photo_uploads")
	tradebookUploadPath := getEnv("TRADEBOOK_UPLOAD_PATH", "./tradebook_uploads")

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/signup", auth.SignupHandler(db, photoUploadPath)).Methods("POST")
	api.HandleFunc("/login", auth.LoginHandler(db)).Methods("POST")
	api.HandleFunc("/logout", auth.LogoutHandler(db)).Methods("GET")
	api.HandleFunc("/logout_all", auth.LogoutAllHandler(db)).Methods("POST")
	api.HandleFunc("/profile", auth.ProfileHandler(db, photoUploadPath)).Methods("GET", "POST")
	api.HandleFunc("/upload", auth.UploadHandler(db, tradebookUploadPath)).Methods("POST")
	api.HandleFunc("/listFiles", auth.ListFilesHandler(db)).Methods("GET")

	router.HandleFunc("/equity/list", handlers.GetEquityList).Methods("GET")
	router.HandleFunc("/equity/trend", handlers.GetTrend).Methods("GET")
	router.HandleFunc("/equity/trend/compare", handlers.GetTrendComparison).Methods("GET")
	router.HandleFunc("/equity/history/refresh", handlers.RefreshPriceHistory).Methods("GET")
	router.HandleFunc("/equity/breakdown", handlers.GetEqBreakdown).Methods("GET")

	router.HandleFunc("/mutual_funds/list", handlers.GetMutualFundsList).Methods("GET")
	router.HandleFunc("/mutual_funds/positions", handlers.GetMFPositions).Methods("GET")
	router.HandleFunc("/mutual_funds/trend", handlers.GetMFTrend).Methods("GET")
	router.HandleFunc("/mutual_funds/summary", handlers.GetMFSummary).Methods("GET")
	router.HandleFunc("/mutual_funds/trend/compare", handlers.GetMFGrowthComparison).Methods("GET")
	router.HandleFunc("/mutual_funds/history/refresh", handlers.RefreshMFPriceHistory).Methods("GET")

	router.PathPrefix("/photo_uploads/").Handler(http.StripPrefix("/photo_uploads/", http.FileServer(http.Dir(photoUploadPath))))
	router.PathPrefix("/tradebook_uploads/").Handler(http.StripPrefix("/tradebook_uploads/", http.FileServer(http.Dir(tradebookUploadPath))))

	spa := spaHandler{staticPath: "dist", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	return router
}
