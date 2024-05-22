package routes

import (
	"bytes"
	"encoding/json"
	"github.com/floxo05/todoapi/internal/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var (
	router *gin.Engine
)

func setup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	router.POST("/login", Login)
	router.POST("/register", Register)
}

func TestLogin(t *testing.T) {
	setup(t)

	// Testfall 1: Leerer Anforderungskörper
	resp := performRequest(router, http.MethodPost, "/login", []byte{})
	if status := resp.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Testfall 2: Ungültiger Benutzername oder Passwort
	reqBody := types.AuthRequest{
		Username: "invalidUser",
		Password: "invalidPassword",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	resp = performRequest(router, http.MethodPost, "/login", reqBodyBytes)
	if status := resp.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Testfall 3: Erfolgreiche Anmeldung

	reqBody = types.AuthRequest{
		Username: "testUser",
		Password: "testPassword",
	}
	reqBodyBytes, _ = json.Marshal(reqBody)
	resp = performRequest(router, http.MethodPost, "/login", reqBodyBytes)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Überprüfen Sie, ob die erwartete Antwort zurückgegeben wurde
	expected := `{"token":".+"}`
	if match, _ := regexp.MatchString(expected, resp.Body.String()); !match {
		t.Errorf("Handler returned unexpected body: got %v want %v", resp.Body.String(), expected)
	}
}

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
