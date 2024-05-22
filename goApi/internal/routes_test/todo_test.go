package routes_test

import (
	"github.com/floxo05/todoapi/internal/models"
	"github.com/floxo05/todoapi/internal/routes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTodos(t *testing.T) {
	t.Run("Test Get Todos", func(t *testing.T) {

		var testcases = []struct {
			header string
		}{
			{"Bearer XXX"},
		}

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		router.GET("/todos", routes.GetTodos)

		// Mock the GetUserFromContext function to bypass the JWT token validation
		routes.GetUserFromContext = func(c *gin.Context) (models.User, error) {
			return models.User{Username: "testUser"}, nil
		}

		for _, tc := range testcases {
			t.Run("Test Get Todos", func(t *testing.T) {
				req, err := http.NewRequest(http.MethodGet, "/todos", nil)
				if err != nil {
					t.Fatal(err)
				}

				req.Header.Set("Authorization", tc.header)

				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				if status := rr.Code; status != http.StatusOK {
					t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
				}
			})
		}
	})
}
