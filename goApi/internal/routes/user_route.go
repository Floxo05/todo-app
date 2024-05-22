package routes

import (
	"github.com/floxo05/todoapi/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
	"unicode"
)

type UserRoute struct {
	userRepository    types.UserRepository
	passwordHasher    types.PasswordHasherInterface
	userContextHelper types.UserContextInterface
}

func NewUserRoute(
	userRepo types.UserRepository,
	passwordHasher types.PasswordHasherInterface,
	userContextHelper types.UserContextInterface) *UserRoute {
	return &UserRoute{
		userRepository:    userRepo,
		passwordHasher:    passwordHasher,
		userContextHelper: userContextHelper}
}

func (u *UserRoute) Login(c *gin.Context) {
	var req types.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	}

	user, err := u.userRepository.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = u.passwordHasher.ComparePasswords(user.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	err = sendJWTToken(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (u *UserRoute) Register(c *gin.Context) {
	var req types.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 3 characters long"})
		return
	}

	if !validatePassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character"})
		return
	}

	hashedPassword, err := u.passwordHasher.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := types.User{Username: req.Username, Password: hashedPassword}
	err = u.userRepository.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = sendJWTToken(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (u *UserRoute) ShareToUser(c *gin.Context) {
	user, err := u.userContextHelper.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req types.ShareToUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'username' must not be empty"})
		return
	}

	if req.TodoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'id' must not be empty"})
		return
	}

	var shareUser *types.User
	shareUser, err = u.userRepository.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve shareUser"})
		return
	}

	err = u.userRepository.ShareTodoWithUser(req.TodoID, user, shareUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo shared successfully"})
}

func sendJWTToken(c *gin.Context, req types.AuthRequest) error {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = req.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	signedString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{"token": signedString})
	return nil
}

func validatePassword(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
