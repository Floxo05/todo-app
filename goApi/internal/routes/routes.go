package routes

import (
	"errors"
	"github.com/floxo05/todoapi/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserContextGetter interface {
	GetUserFromContext(c *gin.Context) (models.User, error)
}

func SendInternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func GetUserFromContext(c *gin.Context) (models.User, error) {
	username, exists := c.Get("username")
	if !exists {
		return models.User{}, errors.New("Could not retrieve user")
	}

	return GetUserByUsername(username.(string))
}
