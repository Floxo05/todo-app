package repository

import (
	"database/sql"
	"errors"
	"github.com/floxo05/todoapi/internal/types"
)

type UserRepo struct {
	db       *sql.DB
	todoRepo types.TodoRepository
}

func NewUserRepo(db *sql.DB, repo types.TodoRepository) *UserRepo {
	return &UserRepo{db: db, todoRepo: repo}
}

func (u *UserRepo) GetUserByUsername(username string) (*types.User, error) {
	var user types.User
	err := u.db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return &types.User{}, err
	}

	return &user, nil
}

func (u *UserRepo) CreateUser(user *types.User) error {
	_, err := u.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) ShareTodoWithUser(todoID int, user *types.User, shareUser *types.User) error {
	todo := types.Todo{ID: todoID}
	isOwner, err := u.todoRepo.IsOwner(&todo, user)
	if err != nil {
		return err
	}

	if !isOwner {
		return errors.New("user does not have access to share the todo")
	}

	_, err = u.db.Exec("INSERT INTO user_todos (todo_id, user_id) VALUES (?, ?)", todoID, shareUser.ID)
	if err != nil {
		return err
	}

	return nil
}
