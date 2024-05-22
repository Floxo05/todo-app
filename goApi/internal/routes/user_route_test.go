package routes

import (
	"bytes"
	"encoding/json"
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
		Username   string
		Password   string
		shouldPass bool
	}{
		{Username: "", Password: "", shouldPass: false},
		{Username: "test", Password: "", shouldPass: false},
		{Username: "", Password: "test", shouldPass: false},
		{Username: "test", Password: "Test1234!", shouldPass: true},
	}

	for _, tc := range testcases {
		t.Run("Test Login", func(t *testing.T) {
			// Act
			reqBodyBytes, _ := json.Marshal(tc)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBodyBytes))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			if tc.shouldPass && w.Code != http.StatusOK {
				t.Errorf("Expected status code 200, but got %d", w.Code)
			}
			if !tc.shouldPass && w.Code == http.StatusOK {
				t.Errorf("Expected failure status code, but got %d", w.Code)
			}
		})
	}
}

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
	if password == "" {
		return "", errors.New("error")

	}
	return "hashedPassword", nil
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
