package routes

import (
	"errors"
	"github.com/floxo05/todoapi/internal/models"
	"github.com/floxo05/todoapi/internal/tools"
	"github.com/floxo05/todoapi/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"unicode"
)

func Login(c *gin.Context) {
	var req types.AuthRequest
	var user models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := tools.InitDB()
	if err != nil {
		SendInternalServerError(c, errors.New("could not connect to the database"))
		return
	}

	err = db.QueryRow("SELECT username, password FROM users WHERE username = ?", req.Username).Scan(&user.Username, &user.Password)
	if err != nil {
		SendInternalServerError(c, errors.New("username does not exist"))
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
		return
	}

	sendJWTToken(c, req, "User logged in successfully")
}

func Register(c *gin.Context) {
	var req types.AuthRequest
	var err error

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// validate the request
	if len(req.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 3 characters long"})
		return
	}

	if !validatePassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character"})
		return
	}

	var user models.User
	user.Username = req.Username
	user.Password, err = tools.HashPassword(req.Password)
	if err != nil {
		SendInternalServerError(c, err)
		return
	}

	// save the user to the database
	err = saveUser(user)
	if err != nil {
		SendInternalServerError(c, err)
		return
	}

	sendJWTToken(c, req, "User registered successfully")
}

func CheckToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}

func ShareToUser(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve shareUser"})
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

	var shareUser models.User
	shareUser, err = GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve shareUser"})
		return
	}

	err = shareTodoToUser(req.TodoID, user, shareUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo shared successfully"})
}

func shareTodoToUser(id int, user models.User, shareUser models.User) error {
	db, err := tools.InitDB()
	if err != nil {
		return err
	}

	// check if user is the owner of the todo
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM user_todos ut JOIN todos t ON t.id = ut.todo_id WHERE ut.user_id = ? AND ut.todo_id = ? AND t.owner_id = ?", user.ID, id, user.ID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("user does not have access to delete the todo")
	}

	_, err = db.Exec("INSERT INTO user_todos (todo_id, user_id) VALUES (?, ?)", id, shareUser.ID)
	if err != nil {
		return err
	}

	return nil
}

// #############################
// #############################
// #############################

func saveUser(user models.User) error {
	db, err := tools.InitDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		ok := errors.As(err, &mysqlErr)
		if ok && mysqlErr.Number == 1062 {
			return errors.New("username already exists")
		}
		return err
	}

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

func sendJWTToken(c *gin.Context, req types.AuthRequest, message string) {
	token, err := tools.GetJWTToken(req)
	if err != nil {
		SendInternalServerError(c, errors.New("could not generate JWT token"))
		return
	}

	c.JSON(http.StatusOK, types.AuthResponse{Message: message, Token: token})
}

func GetUserByUsername(username string) (models.User, error) {
	db, err := tools.InitDB()
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	err = db.QueryRow("SELECT id, username FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
