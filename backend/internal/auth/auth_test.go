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
	database.ConnectTestDB()
	code := m.Run()
	os.Exit(code)
}

func TestRegister(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/register", strings.NewReader(`{"email":"test@test.com", "password":"password"}`))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Register)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the user was created in the database
	var user database.User
	database.DB.Where("email = ?", "test@test.com").First(&user)
	if user.Email != "test@test.com" {
		t.Errorf("user not created in database")
	}

	// Check the password was hashed
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password")); err != nil {
		t.Errorf("password not hashed correctly")
	}
}

func TestLogin(t *testing.T) {
	// Create a user to login with
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 8)
	user := database.User{Email: "test2@test.com", Password: string(hashedPassword)}
	database.DB.Create(&user)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"email":"test2@test.com", "password":"password"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Login)

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
	middlewareHandler := Middleware(protectedHandler)

	// Serve the request
	middlewareHandler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}