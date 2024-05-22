package routes

import (
	"database/sql"
	"github.com/floxo05/todoapi/internal/models"
	"github.com/floxo05/todoapi/internal/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRoute struct {
	// depedencies
	db       *sql.DB
	userRepo *UserRepo
}

func NewUserRoute(db *sql.DB, userRepo *UserRepo) *UserRoute {
	return &UserRoute{db: db, userRepo: userRepo}
}

func (u *UserRoute) Login(c *gin.Context) {
	var req types.AuthRequest
	var user models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

type UserRepo interface {
	FindUserByUsername(username string) (models.User, error)
}
