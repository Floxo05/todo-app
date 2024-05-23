package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/floxo05/todoapi/internal/types"
	"reflect"
	"testing"
	"time"
)

func TestTodoRepo_NewTodoRepo(t *testing.T) {
	t.Run("should return a new TodoRepo", func(t *testing.T) {
		// Act
		repo := NewTodoRepo(nil, nil)

		// Assert
		expectedType := "*repository.TodoRepo"
		if reflect.TypeOf(repo).String() != expectedType {
			t.Errorf("Expected type %s, but got %s", expectedType, reflect.TypeOf(repo))
		}
	})
}

func TestTodoRepo_GetAllTodosByUser(t *testing.T) {
	t.Run("Test GetAllTodosByUser", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		// case 1
		t.Run("should return a list of todos", func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"id", "title", "completed", "created_at", "owner_id"}).
				AddRow(1, "Test Todo", false, "2022-01-01 00:00:00", 1).
				AddRow(2, "Test Todo 2", true, "2022-01-01 00:00:00", 2)

			mock.ExpectQuery("^SELECT (.+) FROM todos").WillReturnRows(rows)

			repo := NewTodoRepo(db, nil)

			// Act
			todos, err := repo.GetAllTodosByUser(&types.User{ID: 1})

			// Assert
			if err != nil {
				t.Errorf("Expected error to be nil, but got %s", err.Error())
			}

			if todos == nil {
				t.Errorf("Expected todos to be a list, but got nil")
			}

			if len(todos) != 2 {
				t.Errorf("Expected 2 todos, but got %d", len(todos))
			}
		})

		// case 2
		t.Run("should return an error", func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"id", "title", "completed", "created_at", "owner_id"}).
				AddRow(1, "Test Todo", false, "", 1)

			mock.ExpectQuery("^SELECT (.+) FROM todos").WillReturnRows(rows)

			repo := NewTodoRepo(db, nil)

			// Act
			_, err := repo.GetAllTodosByUser(&types.User{ID: 1})

			// Assert
			if err == nil {
				t.Errorf("Expected an error, but got nil")
			}
		})

		// case 3
		t.Run("should return nil if list is empty", func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"id", "title", "completed", "created_at", "owner_id"})

			mock.ExpectQuery("^SELECT (.+) FROM todos").WillReturnRows(rows)

			repo := NewTodoRepo(db, nil)

			// Act
			todos, err := repo.GetAllTodosByUser(&types.User{ID: 1})

			// Assert
			if err != nil {
				t.Errorf("Expected error to be nil, but got %s", err.Error())
			}

			if todos != nil {
				t.Errorf("Expected todos to be nil")
			}
		})
	})
}

func TestTodoRepo_CreateTodo(t *testing.T) {
	t.Run("Test CreateTodo", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		// case 1
		t.Run("should create a new todo", func(t *testing.T) {
			mock.ExpectExec("^INSERT INTO todos").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("^INSERT INTO user_todos").WillReturnResult(sqlmock.NewResult(1, 1))

			repo := NewTodoRepo(db, nil)

			// Act
			date, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 00:00:00")
			err := repo.CreateTodo(&types.Todo{
				Title:     "Test Todo",
				Completed: false,
				CreatedAt: date,
				OwnerID:   1,
			})

			// Assert
			if err != nil {
				t.Errorf("Expected error to be nil, but got %s", err.Error())
			}
		})

		// case 2
		t.Run("should return an error", func(t *testing.T) {
			mock.ExpectExec("^INSERT INTO todos").WillReturnError(err)

			repo := NewTodoRepo(db, nil)

			// Act
			date, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 00:00:00")
			err := repo.CreateTodo(&types.Todo{
				Title:     "Test Todo",
				Completed: false,
				CreatedAt: date,
				OwnerID:   1,
			})

			// Assert
			if err == nil {
				t.Errorf("Expected an error, but got nil")
			}
		})
	})
}

func TestTodoRepo_UpdateTodoById(t *testing.T) {
	t.Run("Test UpdateTodoById", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		// case 1
		t.Run("should update a todo", func(t *testing.T) {
			mock.ExpectExec("^UPDATE todos").WillReturnResult(sqlmock.NewResult(1, 1))

			repo := NewTodoRepo(db, nil)

			// Act
			date, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 00:00:00")
			err := repo.UpdateTodoById(&types.Todo{
				ID:        1,
				Title:     "Test Todo",
				Completed: false,
				CreatedAt: date,
				OwnerID:   1,
			}, &types.User{ID: 1})

			// Assert
			if err != nil {
				t.Errorf("Expected error to be nil, but got %s", err.Error())
			}
		})

		// case 2
		t.Run("should return an error", func(t *testing.T) {
			mock.ExpectExec("^UPDATE todos").WillReturnError(err)

			repo := NewTodoRepo(db, nil)

			// Act
			date, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 00:00:00")
			err := repo.UpdateTodoById(&types.Todo{
				ID:        1,
				Title:     "Test Todo",
				Completed: false,
				CreatedAt: date,
				OwnerID:   1,
			}, &types.User{ID: 1})

			// Assert
			if err == nil {
				t.Errorf("Expected an error, but got nil")
			}
		})
	})
}
