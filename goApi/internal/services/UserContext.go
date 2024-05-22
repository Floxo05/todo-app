package services

import (
	"errors"
	"github.com/floxo05/todoapi/internal/types"
	"github.com/gin-gonic/gin"
)

type UserContext struct {
	UserRepository types.UserRepository
}

func NewUserContext(userRepository types.UserRepository) *UserContext {
	return &UserContext{UserRepository: userRepository}
}

func (u *UserContext) GetUserFromContext(c *gin.Context) (*types.User, error) {
	username, exists := c.Get("username")
	if !exists {
		return nil, errors.New("could not retrieve user")
	}

	return u.UserRepository.GetUserByUsername(username.(string))
}
