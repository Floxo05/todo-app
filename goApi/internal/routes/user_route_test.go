package routes

import (
	"bytes"
	"errors"
	"github.com/floxo05/todoapi/internal/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserRoute_Login(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	loginRoute := NewUserRoute(&mockUserRepository{}, &mockPasswordHasher{}, &mockUserContextHelper{})
	router.POST("/login", loginRoute.Login)

	testcases := []struct {
		Body             []byte
		expectedResponse int
	}{
		{Body: nil, expectedResponse: http.StatusBadRequest},
		{Body: []byte(`{"Username": "", "Password": ""}`), expectedResponse: http.StatusNotFound},
		{Body: []byte(`{"Username": "test", "Password": ""}`), expectedResponse: http.StatusUnauthorized},
		{Body: []byte(`{"Username": "test", "Password": "Test1234!"}`), expectedResponse: http.StatusOK},
	}

	for _, tc := range testcases {
		t.Run("Test Login", func(t *testing.T) {
			// Act
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(tc.Body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			if tc.expectedResponse != w.Code {
				t.Errorf("Expected status code 200, but got %d", w.Code)
			}
		})
	}
}

func TestUserRoute_Register(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	loginRoute := NewUserRoute(&mockUserRepository{}, &mockPasswordHasher{}, &mockUserContextHelper{})
	router.POST("/register", loginRoute.Register)

	testcases := []struct {
		Body             []byte
		expectedResponse int
	}{
		{Body: nil, expectedResponse: http.StatusBadRequest},
		{Body: []byte(`{"Username": "", "Password": ""}`), expectedResponse: http.StatusBadRequest},
		{Body: []byte(`{"Username": "test", "Password": ""}`), expectedResponse: http.StatusBadRequest},
		{Body: []byte(`{"Username": "test", "Password": "Test1234!"}`), expectedResponse: http.StatusOK},
	}

	for _, tc := range testcases {
		t.Run("Test Register", func(t *testing.T) {
			// Act
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(tc.Body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			if tc.expectedResponse != w.Code {
				t.Errorf("Expected status code 200, but got %d", w.Code)
			}
		})
	}
}

/////////////////////////////////////////////

type mockUserRepository struct{}

func (m *mockUserRepository) GetUserByUsername(username string) (*types.User, error) {
	if username == "test" {
		return &types.User{
			Username: "test",
			Password: "hashedPassword",
		}, nil
	}

	return nil, errors.New("error")
}

func (m *mockUserRepository) CreateUser(user *types.User) error {
	return nil
}

func (m *mockUserRepository) ShareTodoWithUser(todoID int, user *types.User, shareUser *types.User) error {
	return nil
}

type mockPasswordHasher struct{}

func (m *mockPasswordHasher) HashPassword(password string) (string, error) {
	if password == "Test1234!" {
		return "hashedPassword", nil
	}

	return "", errors.New("error")
}

func (m *mockPasswordHasher) ComparePasswords(hashedPassword string, password string) error {
	hPw, _ := m.HashPassword(password)
	if hashedPassword != hPw {
		return errors.New("error")
	}

	return nil
}

type mockUserContextHelper struct{}

func (m *mockUserContextHelper) GetUserFromContext(c *gin.Context) (*types.User, error) {
	return &types.User{}, nil
}
