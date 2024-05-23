package repository

import (
	"database/sql"
	"errors"
	"github.com/floxo05/todoapi/internal/types"
	"time"
)

type TodoRepo struct {
	db           *sql.DB
	categoryRepo types.CategoryRepository
}

func NewTodoRepo(db *sql.DB, repo types.CategoryRepository) *TodoRepo {
	return &TodoRepo{db: db, categoryRepo: repo}
}

func (t *TodoRepo) GetAllTodosByUser(user *types.User) ([]types.Todo, error) {
	rows, err := t.db.Query(`
		SELECT 
			t.id, t.title, t.completed, t.created_at, t.owner_id, t.category_id
		FROM todos t 
		    JOIN user_todos ut ON t.id = ut.todo_id 
		WHERE ut.user_id = ?`, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []types.Todo
	for rows.Next() {
		var todo types.Todo
		var category types.Category
		var createdAt string
		var categoryID sql.NullInt64
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &createdAt, &todo.OwnerID, &categoryID)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			var cat *types.Category
			cat, err = t.categoryRepo.GetCategoryByID(int(categoryID.Int64))
			if err != nil {
				return nil, err
			}
			category = *cat
		}

		// Parse the string into a time.Time type
		todo.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, err
		}

		todo.Category = category

		todos = append(todos, todo)
	}

	return todos, nil
}

func (t *TodoRepo) CreateTodo(todo *types.Todo) error {
	res, err := t.db.Exec("INSERT INTO todos (title, completed, created_at, owner_id) VALUES (?, ?, ?, ?)", todo.Title, todo.Completed, todo.CreatedAt, todo.OwnerID)
	if err != nil {
		return err
	}

	todoID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	todo.ID = int(todoID)

	_, err = t.db.Exec("INSERT INTO user_todos (user_id, todo_id) VALUES (?, ?)", todo.OwnerID, todo.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoRepo) UpdateTodoById(todo *types.Todo, user *types.User) error {
	var count int
	var err error

	// check if user has access to the todo
	err = t.db.QueryRow("SELECT COUNT(*) FROM user_todos WHERE user_id = ? AND todo_id = ?", user.ID, todo.ID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("user does not have access to update the todo")
	}

	err = t.categoryRepo.UpsertCategory(&todo.Category)
	if err != nil {
		return err
	}

	_, err = t.db.Exec("UPDATE todos SET title = ?, completed = ?, category_id = ? WHERE id = ?", todo.Title, todo.Completed, todo.Category.ID, todo.ID)

	// update
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoRepo) DeleteTodoById(todo *types.Todo, user *types.User) error {
	// check if user is the owner of the todo
	isOwner, err := t.IsOwner(todo, user)
	if err != nil {
		return err
	}

	if !isOwner {
		return errors.New("user does not have right to delete the todo")
	}

	// delete association
	_, err = t.db.Exec("DELETE FROM user_todos where todo_id = ?", todo.ID)
	if err != nil {
		return err
	}

	_, err = t.db.Exec("DELETE FROM todos WHERE id = ?", todo.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoRepo) IsOwner(todo *types.Todo, user *types.User) (bool, error) {
	var count int
	err := t.db.QueryRow("SELECT COUNT(*) FROM todos WHERE id = ? AND owner_id = ?", todo.ID, user.ID).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
