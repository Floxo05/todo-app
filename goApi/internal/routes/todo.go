package routes

import (
	"errors"
	"github.com/floxo05/todoapi/internal/models"
	"github.com/floxo05/todoapi/internal/tools"
	"github.com/floxo05/todoapi/internal/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func CreateTodo(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		return
	}

	var req types.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	//validate the request
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'title' must not be empty"})
		return
	}

	// get database
	db, err := tools.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to the database"})
		return
	}

	// create todo
	todo := models.Todo{Title: req.Title}
	todo, err = upsertTodoById(todo, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo"})
		return
	}

	_, err = db.Exec("INSERT INTO user_todos (user_id, todo_id) VALUES (?, ?)", user.ID, todo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user-todo relationship"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo_id": todo.ID, "title": req.Title, "completed": false})
}

func GetTodos(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		return
	}

	todos := getAllTodosByUser(user)
	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func GetTodoById(c *gin.Context) {
	_, err := GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		return
	}

	var todoId int
	todoId, err = strconv.Atoi(c.Param("id"))

	todo, err := getTodoById(todoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func UpdateTodoById(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		return
	}

	var req types.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// validate the request
	if req.Title == nil || req.Completed == nil || req.ID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'id', 'title' and 'completed' are required"})
		return
	}

	if *req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'title' must not be empty"})
		return
	}

	var todo models.Todo
	todo.ID = int64(*req.ID)
	todo.Title = *req.Title
	todo.Completed = *req.Completed

	todo, err = upsertTodoById(todo, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func DeleteTodoById(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		return
	}

	var todoId int
	todoId, err = strconv.Atoi(c.Param("id"))

	if err := deleteTodoById(todoId, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

func deleteTodoById(todoId int, user models.User) error {
	db, err := tools.InitDB()
	if err != nil {
		return err
	}

	// Check if the user has access to the todo
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM user_todos ut JOIN todos t ON t.id = ut.todo_id WHERE ut.user_id = ? AND ut.todo_id = ? AND t.owner_id = ?", user.ID, todoId, user.ID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("user does not have access to delete the todo")
	}

	// delete todo assignments
	_, err = db.Exec("DELETE FROM user_todos WHERE todo_id = ?", todoId)
	if err != nil {
		return err
	}

	// delete todo
	_, err = db.Exec("DELETE FROM todos WHERE id = ?", todoId)
	return err
}

// ##################################################
// ##################################################
// ##################################################

func upsertTodoById(todo models.Todo, user models.User) (models.Todo, error) {
	db, err := tools.InitDB()
	if err != nil {
		return todo, err
	}

	if todo.ID == 0 {
		// insert
		res, err := db.Exec("INSERT INTO todos (title, completed, owner_id) VALUES (?, ?, ?)", todo.Title, false, user.ID)
		if err != nil {
			return todo, err
		}

		todo.ID, err = res.LastInsertId()
		if err != nil {
			return todo, err
		}

		return todo, nil
	} else {
		// check if user has access to the todo
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM user_todos WHERE user_id = ? AND todo_id = ?", user.ID, todo.ID).Scan(&count)
		if err != nil {
			return todo, err
		}

		if count == 0 {
			return todo, errors.New("user does not have access to update the todo")
		}

		// update
		_, err = db.Exec("UPDATE todos SET title = ?, completed = ? WHERE id = ?", todo.Title, todo.Completed, todo.ID)
		if err != nil {
			return todo, err
		}

		return todo, nil
	}
}

func getAllTodosByUser(user models.User) []models.Todo {
	db, err := tools.InitDB()
	if err != nil {
		return nil
	}

	rows, err := db.Query("SELECT t.id, t.title, t.completed, t.created_at FROM todos t JOIN user_todos ut ON t.id = ut.todo_id WHERE ut.user_id = ?", user.ID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		var createdAt string
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &createdAt)
		if err != nil {
			return nil
		}

		// Parse the string into a time.Time type
		todo.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil
		}

		todos = append(todos, todo)
	}

	return todos
}

func getTodoById(todoId int) (models.Todo, error) {
	db, err := tools.InitDB()
	if err != nil {
		return models.Todo{}, err
	}

	var todo models.Todo
	var createdAt string
	err = db.QueryRow("SELECT id, title, completed, created_at FROM todos WHERE id = ?", todoId).Scan(&todo.ID, &todo.Title, &todo.Completed, &createdAt)
	if err != nil {
		return models.Todo{}, err
	}

	// Parse the string into a time.Time type
	todo.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return models.Todo{}, err
	}

	return todo, nil
}
