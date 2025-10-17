package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var tokenInactivityWindow = 24 * time.Hour
var tokenAbsoluteExpiry = 7 * 24 * time.Hour

type Claims struct {
	Email   string `json:"email"`
	TokenID string `json:"tid"`
	jwt.RegisteredClaims
}

func issueToken(db *sql.DB, email string) (string, error) {
	tokenID := uuid.NewString()
	expiration := time.Now().Add(tokenAbsoluteExpiry)
	claims := &Claims{
		Email:   email,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	_, err = db.Exec(
		`INSERT INTO tokens (id, email, created_at, last_activity) VALUES ($1, $2, now(), now())`,
		tokenID, email,
	)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func validateAndRefreshToken(r *http.Request, db *sql.DB) (string, string, error) {
	tokenStr := extractTokenFromRequest(r)
	if tokenStr == "" {
		return "", "", errors.New("missing token")
	}

	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !parsed.Valid {
		return "", "", errors.New("invalid token")
	}

	if claims.TokenID == "" || claims.Email == "" {
		return "", "", errors.New("malformed token claims")
	}

	var lastActivity time.Time
	err = db.QueryRow(`SELECT last_activity FROM tokens WHERE id=$1 AND email=$2`,
		claims.TokenID, claims.Email).Scan(&lastActivity)
	if err == sql.ErrNoRows {
		return "", "", errors.New("token not found / revoked")
	} else if err != nil {
		return "", "", err
	}

	if lastActivity.Add(tokenInactivityWindow).Before(time.Now()) {
		_, _ = db.Exec(`DELETE FROM tokens WHERE id=$1`, claims.TokenID)
		return "", "", errors.New("token expired by inactivity")
	}

	_, err = db.Exec(`UPDATE tokens SET last_activity = now() WHERE id=$1`, claims.TokenID)
	if err != nil {
		return "", "", err
	}

	return claims.Email, claims.TokenID, nil
}

func extractTokenFromRequest(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth != "" {
		parts := strings.Fields(auth)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return parts[1]
		}
	}

	if cookie, err := r.Cookie("auth_token"); err == nil {
		return cookie.Value
	}
	return ""
}

func revokeToken(db *sql.DB, tokenID string) error {
	_, err := db.Exec(`DELETE FROM tokens WHERE id=$1`, tokenID)
	return err
}

func revokeAllTokensForUser(db *sql.DB, email string) error {
	_, err := db.Exec(`DELETE FROM tokens WHERE email=$1`, email)
	return err
}

func Authorize(db *sql.DB, w http.ResponseWriter, r *http.Request) (string, string, bool) {
	email, tokenID, err := validateAndRefreshToken(r, db)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return "", "", false
	}
	return email, tokenID, true
}
