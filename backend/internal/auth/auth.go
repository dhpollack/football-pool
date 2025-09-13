// Package auth provides authentication and authorization middleware for the football pool application.
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

	"github.com/david/football-pool/internal/api"
	"github.com/david/football-pool/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

var jwtKey = []byte("my_secret_key")

// Auth provides authentication and authorization services.
type Auth struct {
	db *database.Database
}

// NewAuth creates a new Auth instance with the provided database connection.
func NewAuth(db *database.Database) *Auth {
	return &Auth{db: db}
}

// Claims represents JWT token claims for user authentication.
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Login handles user authentication and issues JWT tokens.
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var creds api.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Debug("Login attempt for user:", "email", creds.Email)
	var user database.User
	if result := a.db.GetDB().Where("email = ?", creds.Email).First(&user); result.Error != nil {
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
	if err := json.NewEncoder(w).Encode(api.LoginResponse{
		Token: tokenString,
		User: api.UserResponse{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}); err != nil {
		slog.Debug("Error encoding token response:", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Register handles new user registration.
func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var creds api.RegisterRequest
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
	if a.db.GetDB().Where("email = ?", creds.Email).First(&existingUser).Error == nil {
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

	if creds.Role == nil || *creds.Role == "" {
		role := "player"
		creds.Role = &role
	}

	user := database.User{Name: creds.Name, Email: creds.Email, Password: string(hashedPassword), Role: *creds.Role}
	if err := validate.Struct(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: err.Error()})
		return
	}

	if result := a.db.GetDB().Create(&user); result.Error != nil {
		slog.Debug("Error creating user:", "error", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create associated Player record for the new user
	player := database.Player{
		UserID:  user.ID,
		Name:    creds.Name,
		Address: "", // Default empty address
	}
	if err := validate.Struct(player); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: err.Error()})
		return
	}
	if result := a.db.GetDB().Create(&player); result.Error != nil {
		slog.Debug("Error creating player:", "error", result.Error)
		// If player creation fails, clean up the user record
		a.db.GetDB().Delete(&user)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Debug("User and player registered successfully:", "email", creds.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(api.RegisterResponse{
		Message: "User registered successfully",
	}); err != nil {
		slog.Debug("Error encoding registration response:", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Logout handles user logout by clearing authentication cookies.
func (a *Auth) Logout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
	})
}

// Middleware provides authentication middleware for HTTP handlers.
func (a *Auth) Middleware(next http.Handler) http.Handler {
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
			tknStr = a.getTokenFromCookie(w, r)
			if tknStr == "" {
				return
			}
		}

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(_ *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			a.handleJWTError(w, err)
			return
		}
		if !tkn.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(api.ErrorResponse{
				Error: "Unauthorized: Invalid token",
			}); err != nil {
				slog.Debug("Error encoding error response:", "error", err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), EmailKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware provides authorization middleware for admin-only endpoints.
func (a *Auth) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(EmailKey).(string)

		var user database.User
		if result := a.db.GetDB().Where("email = ?", email).First(&user); result.Error != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(api.ErrorResponse{
				Error: "Not Found: User not found",
			}); err != nil {
				slog.Debug("Error encoding error response:", "error", err)
			}
			return
		}

		if user.Role != "admin" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(api.ErrorResponse{
				Error: "Forbidden: User is not an admin",
			}); err != nil {
				slog.Debug("Error encoding error response:", "error", err)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getTokenFromCookie extracts the token from the cookie and handles any errors.
func (a *Auth) getTokenFromCookie(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie("token")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(api.ErrorResponse{
				Error: "Unauthorized: No token provided",
			}); err != nil {
				slog.Debug("Error encoding error response:", "error", err)
			}
			return ""
		}
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Bad Request: Invalid token cookie",
		}); err != nil {
			slog.Debug("Error encoding error response:", "error", err)
		}
		return ""
	}
	return c.Value
}

// handleJWTError handles JWT validation errors with appropriate HTTP responses.
func (a *Auth) handleJWTError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	if errors.Is(err, jwt.ErrSignatureInvalid) {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Unauthorized: Invalid token signature",
		}); err != nil {
			slog.Debug("Error encoding error response:", "error", err)
		}
		return
	}
	if errors.Is(err, jwt.ErrTokenExpired) {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(api.ErrorResponse{
			Error: "Unauthorized: Token expired",
		}); err != nil {
			slog.Debug("Error encoding error response:", "error", err)
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(api.ErrorResponse{
		Error: "Bad Request: Invalid token",
	}); err != nil {
		slog.Debug("Error encoding error response:", "error", err)
	}
}
