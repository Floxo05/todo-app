package routes

import (
	"github.com/floxo05/todoapi/internal/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type TodoRoute struct {
	todoRepository    types.TodoRepository
	userContextHelper types.UserContextInterface
}

func NewTodoRoute(todoRepository types.TodoRepository, userContextHelper types.UserContextInterface) *TodoRoute {
	return &TodoRoute{todoRepository: todoRepository, userContextHelper: userContextHelper}
}

func (t *TodoRoute) GetTodos(c *gin.Context) {
	user, err := t.userContextHelper.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	todos, err := t.todoRepository.GetAllTodosByUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (t *TodoRoute) CreateTodo(c *gin.Context) {
	user, err := t.userContextHelper.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req types.CreateTodoRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	//validate the request
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'title' must not be empty"})
		return
	}

	todo := types.Todo{Title: req.Title, OwnerID: user.ID, Completed: false, CreatedAt: time.Now()}
	err = t.todoRepository.CreateTodo(&todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (t *TodoRoute) UpdateTodo(c *gin.Context) {
	user, err := t.userContextHelper.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req types.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Title == nil || req.Completed == nil || req.ID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'id', 'title' and 'completed' are required"})
		return
	}

	if *req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'title' must not be empty"})
		return
	}

	todo := types.Todo{ID: *req.ID, Title: *req.Title, Completed: *req.Completed, Category: *req.Category}
	todo.Category.CreatedUserId = user.ID

	err = t.todoRepository.UpdateTodoById(&todo, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (t *TodoRoute) DeleteTodo(c *gin.Context) {
	user, err := t.userContextHelper.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var todoId int
	todoId, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	todo := types.Todo{ID: todoId}
	err = t.todoRepository.DeleteTodoById(&todo, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
