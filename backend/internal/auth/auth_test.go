package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/david/football-pool/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	// Database setup is now handled in individual tests
	code := m.Run()
	os.Exit(code)
}

func TestRegister(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	auth := NewAuth(db)

	// Create a request to pass to our handler.
	req, err := http.NewRequest("POST", "/register", strings.NewReader(`{"email":"test@test.com", "password":"password"}`))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.Register)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the user was created in the database
	var user database.User
	db.GetDB().Where("email = ?", "test@test.com").First(&user)
	if user.Email != "test@test.com" {
		t.Errorf("user not created in database")
	}

	// Check the password was hashed
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password")); err != nil {
		t.Errorf("password not hashed correctly")
	}

	// Check the JSON response content
	expectedResponse := `{"message":"User registered successfully"}`
	if strings.TrimSpace(rr.Body.String()) != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedResponse)
	}

	// Check the Content-Type header
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type header not set correctly: got %v want %v", ct, "application/json")
	}
}

func TestLogin(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	auth := NewAuth(db)

	// Create a user to login with
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 8)
	user := database.User{Email: "test2@test.com", Password: string(hashedPassword)}
	db.GetDB().Create(&user)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"email":"test2@test.com", "password":"password"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the cookie is set
	cookie := rr.Result().Cookies()[0]
	if cookie.Name != "token" {
		t.Errorf("cookie not set correctly")
	}

	// Check the token is valid
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !tkn.Valid {
		t.Errorf("token not valid")
	}
	if claims.Email != "test2@test.com" {
		t.Errorf("email not set correctly in token")
	}
}

func TestMiddleware(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	auth := NewAuth(db)

	// Create a handler to be protected by the middleware
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a request with a valid token
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: "test@test.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create the middleware handler
	middlewareHandler := auth.Middleware(protectedHandler)

	// Serve the request
	middlewareHandler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAdminMiddleware(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	auth := NewAuth(db)

	// Create a handler to be protected by the middleware
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a user with admin role
	user := database.User{Email: "admin@test.com", Password: "password", Role: "admin"}
	db.GetDB().Create(&user)

	// Create a request with a valid token for the admin user
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: "admin@test.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create the middleware chain
	finalHandler := auth.Middleware(auth.AdminMiddleware(protectedHandler))

	// Serve the request
	finalHandler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestLoginErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	auth := NewAuth(db)

	// Create a user to login with
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 8)
	user := database.User{Email: "test3@test.com", Password: string(hashedPassword)}
	db.GetDB().Create(&user)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "Invalid JSON",
			body:           `{"email":"test3@test.com", "password":"password"`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "User not found",
			body:           `{"email":"notfound@test.com", "password":"password"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing password",
			body:           `{"email":"test3@test.com", "password":""}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Incorrect password",
			body:           `{"email":"test3@test.com", "password":"wrongpassword"}`,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/login", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(auth.Login)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestRegisterErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	auth := NewAuth(db)

	// Create a user to conflict with
	user := database.User{Email: "test4@test.com", Password: "password"}
	db.GetDB().Create(&user)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "Invalid JSON",
			body:           `{"email":"test5@test.com", "password":"password"`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "User already exists",
			body:           `{"email":"test4@test.com", "password":"password"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/register", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(auth.Register)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	auth := NewAuth(db)

	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.Logout)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	cookie := rr.Result().Cookies()[0]
	if cookie.Name != "token" {
		t.Errorf("cookie not set correctly")
	}
	if cookie.Value != "" {
		t.Errorf("cookie value not cleared")
	}
}

func TestMiddlewareErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	auth := NewAuth(db)

	// Create a handler to be protected by the middleware
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		cookie         *http.Cookie
		expectedStatus int
	}{
		{
			name:           "No token cookie",
			cookie:         nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid token",
			cookie: &http.Cookie{
				Name:  "token",
				Value: "invalid-token",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Expired token",
			cookie: &http.Cookie{
				Name: "token",
				Value: func() string {
					expirationTime := time.Now().Add(-5 * time.Minute)
					claims := &Claims{
						Email: "test@test.com",
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(expirationTime),
						},
					}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					tokenString, _ := token.SignedString(jwtKey)
					return tokenString
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			rr := httptest.NewRecorder()
			middlewareHandler := auth.Middleware(protectedHandler)
			middlewareHandler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestAdminMiddlewareErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	auth := NewAuth(db)

	// Create a handler to be protected by the middleware
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a user with player role
	user := database.User{Email: "player@test.com", Password: "password", Role: "player"}
	db.GetDB().Create(&user)

	tests := []struct {
		name           string
		email          string
		expectedStatus int
	}{
		{
			name:           "User not found",
			email:          "notfound@test.com",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "User is not an admin",
			email:          "player@test.com",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with a valid token for the user
			expirationTime := time.Now().Add(5 * time.Minute)
			claims := &Claims{
				Email: tt.email,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(expirationTime),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, _ := token.SignedString(jwtKey)

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})

			// Create a ResponseRecorder
			rr := httptest.NewRecorder()

			// Create the middleware chain
			finalHandler := auth.Middleware(auth.AdminMiddleware(protectedHandler))

			// Serve the request
			finalHandler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
