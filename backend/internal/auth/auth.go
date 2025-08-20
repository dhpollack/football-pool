package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/david/football-pool/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Debug("Login attempt for user:", "email", creds.Email)
	var user database.User
	if result := database.DB.Where("email = ?", creds.Email).First(&user); result.Error != nil {
		slog.Debug("User not found:", "email", creds.Email, "error", result.Error)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if creds.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Debug("User found. Stored password hash:", "hash", user.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		slog.Debug("Password comparison failed for user:", "email", creds.Email, "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	slog.Debug("Password comparison successful for user:", "email", creds.Email)

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: creds.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Read the request body into a byte slice for logging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Debug("Error reading request body:", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Restore the body for subsequent reads (json.NewDecoder)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	slog.Debug("Register request body:", "body", string(bodyBytes))

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		slog.Debug("Error decoding request body:", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var existingUser database.User
	if database.DB.Where("email = ?", creds.Email).First(&existingUser).Error == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Debug("Registering user:", "email", creds.Email)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		slog.Debug("Error hashing password:", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Debug("Hashed password:", "hash", string(hashedPassword))
	user := database.User{Name: creds.Name, Email: creds.Email, Password: string(hashedPassword), Role: "player"}
	if result := database.DB.Create(&user); result.Error != nil {
		slog.Debug("Error creating user:", "error", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Debug("User registered successfully:", "email", creds.Email)
	w.WriteHeader(http.StatusCreated)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
	})
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tknStr := ""

		// Check for token in Authorization header first
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tknStr = parts[1]
			}
		}

		// If not in Authorization header, check for token in cookie
		if tknStr == "" {
			c, err := r.Cookie("token")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: No token provided"})
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request: Invalid token cookie"})
				return
			}
			tknStr = c.Value
		}

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Invalid token signature"})
				return
			}
			if errors.Is(err, jwt.ErrTokenExpired) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Token expired"})
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request: Invalid token"})
			return
		}
		if !tkn.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Invalid token"})
			return
		}

		ctx := context.WithValue(r.Context(), EmailKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(EmailKey).(string)

		var user database.User
		if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Not Found: User not found"})
			return
		}

		if user.Role != "admin" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden: User is not an admin"})
			return
		}

		next.ServeHTTP(w, r)
	})
}