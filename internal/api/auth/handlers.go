package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Mryashbhardwaj/marketAnalysis/internal/api/utils"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name   string `json:"Name"`
	Email  string `json:"Email"`
	Broker string `json:"Broker"`
	Photo  string `json:"Photo"`
}

func SignupHandler(db *sql.DB, photoUploadPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		broker := r.FormValue("broker")

		if name == "" || email == "" || password == "" || broker == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "All fields are required")
			return
		}

		var photoPath string
		file, handler, err := r.FormFile("photo")
		if err == nil {
			defer file.Close()
			photoPath, err = savePhoto(file, handler.Filename, photoUploadPath)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error saving photo")
				return
			}
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
			return
		}

		_, err = db.Exec(
			`INSERT INTO users (name, email, password, broker, photo_path) VALUES ($1, $2, $3, $4, $5)`,
			name, email, string(hashedPassword), broker, photoPath,
		)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Email may already be in use")
			return
		}

		signed, err := issueToken(db, email)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Could not issue token: "+err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    signed,
			Path:     "/",
			HttpOnly: true,
		})

		newUser := User{Name: name, Email: email, Broker: broker, Photo: photoPath}
		utils.RespondWithJSON(w, http.StatusCreated, newUser)
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var hashedPassword string
		var user User
		err := db.QueryRow("SELECT password, name, email, broker, photo_path FROM users WHERE email=$1", email).
			Scan(&hashedPassword, &user.Name, &user.Email, &user.Broker, &user.Photo)

		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		} else if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Server error")
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		signed, err := issueToken(db, email)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Could not issue token: "+err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    signed,
			Path:     "/",
			HttpOnly: true,
		})

		utils.RespondWithJSON(w, http.StatusOK, user)
	}
}

func LogoutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, tokenID, ok := Authorize(db, w, r)
		if !ok {
			return
		}
		_ = revokeToken(db, tokenID)

		cookie := &http.Cookie{Name: "auth_token", Value: "", Path: "/", MaxAge: -1, HttpOnly: true}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LogoutAllHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, _, ok := Authorize(db, w, r)
		if !ok {
			return
		}
		if err := revokeAllTokensForUser(db, email); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to logout from all devices: "+err.Error())
			return
		}

		cookie := &http.Cookie{Name: "auth_token", Value: "", Path: "/", MaxAge: -1, HttpOnly: true}
		http.SetCookie(w, cookie)

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Logged out from all devices"})
	}
}

func ProfileHandler(db *sql.DB, photoUploadPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, _, ok := Authorize(db, w, r)
		if !ok {
			return
		}

		if r.Method == http.MethodPost {
			if newPassword := r.FormValue("password"); newPassword != "" {
				hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
				if err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
					return
				}
				if _, err := db.Exec("UPDATE users SET password=$1 WHERE email=$2", string(hashed), email); err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update password: "+err.Error())
					return
				}
			}
			file, handler, err := r.FormFile("photo")
			if err == nil {
				defer file.Close()
				photoPath, err := savePhoto(file, handler.Filename, photoUploadPath)
				if err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save photo: "+err.Error())
					return
				}
				if _, err := db.Exec("UPDATE users SET photo_path=$1 WHERE email=$2", photoPath, email); err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update photo in DB: "+err.Error())
					return
				}
			}
			utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
			return
		}

		var user User
		err := db.QueryRow("SELECT email, name, broker, photo_path FROM users WHERE email=$1", email).
			Scan(&user.Email, &user.Name, &user.Broker, &user.Photo)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "User not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, user)
	}
}

func UploadHandler(db *sql.DB, tradebookUploadPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, _, ok := Authorize(db, w, r)
		if !ok {
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
			utils.RespondWithError(w, http.StatusBadRequest, "Error parsing form: "+err.Error())
			return
		}

		files := r.MultipartForm.File["csvfiles"]
		if len(files) == 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "No files uploaded")
			return
		}

		for _, header := range files {
			file, err := header.Open()
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Error opening file: "+err.Error())
				return
			}
			defer file.Close()

			relativePath, err := saveFile(file, header.Filename, tradebookUploadPath)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error saving file: "+err.Error())
				return
			}

			_, err = db.Exec(
				`INSERT INTO uploads (email, filename, file_path) VALUES ($1, $2, $3)`,
				email, header.Filename, relativePath,
			)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "DB error: "+err.Error())
				return
			}
		}

		// Instead of redirecting, send a success message
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Files uploaded successfully"})
	}
}

func ListFilesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, _, ok := Authorize(db, w, r)
		if !ok {
			return
		}

		rows, err := db.Query(
			`SELECT id, filename, file_path, uploaded_at 
			 FROM uploads 
			 WHERE email=$1 
			 ORDER BY uploaded_at DESC`, email)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "DB error: "+err.Error())
			return
		}
		defer rows.Close()

		type FileInfo struct {
			ID         int       `json:"id"`
			Filename   string    `json:"filename"`
			FilePath   string    `json:"file_path"`
			UploadedAt time.Time `json:"uploaded_at"`
		}

		var files []FileInfo
		for rows.Next() {
			var f FileInfo
			if err := rows.Scan(&f.ID, &f.Filename, &f.FilePath, &f.UploadedAt); err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Scan error: "+err.Error())
				return
			}
			files = append(files, f)
		}
		if err := rows.Err(); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Row iteration error: "+err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, files)
	}
}
