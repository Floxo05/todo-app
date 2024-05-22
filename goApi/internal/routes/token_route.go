package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TokenRoute struct{}

func NewTokenRoute() *TokenRoute {
	return &TokenRoute{}
}

func (t *TokenRoute) CheckToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}
