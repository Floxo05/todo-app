package types

import (
	"github.com/gin-gonic/gin"
	"time"
)

type UserRepository interface {
	GetUserByUsername(username string) (*User, error)
	CreateUser(user *User) error
	ShareTodoWithUser(todoID int, user *User, shareUser *User) error
}

type TodoRepository interface {
	GetAllTodosByUser(user *User) ([]Todo, error)
	CreateTodo(todo *Todo) error
	UpdateTodoById(todo *Todo, user *User) error
	DeleteTodoById(todo *Todo, user *User) error
	IsOwner(todo *Todo, user *User) (bool, error)
}

type CategoryRepository interface {
	UpsertCategory(category *Category) error
	GetCategoryFromDB(category *Category) (*Category, error)
	GetCategoryByID(id int) (*Category, error)
}

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	OwnerID   int       `json:"owner_id"`
	Category  Category  `json:"category"`
}

type Category struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	CreatedUserId int    `json:"created_user_id"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserContextInterface interface {
	GetUserFromContext(c *gin.Context) (*User, error)
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type CreateCategoryRequest struct {
	Title string `json:"title"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateTodoRequest struct {
	ID        *int    `json:"id"`
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

type ShareToUserRequest struct {
	Username string `json:"username"`
	TodoID   int    `json:"id"`
}

type PasswordHasherInterface interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword, password string) error
}
