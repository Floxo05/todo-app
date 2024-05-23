package repository

import (
	"database/sql"
	"errors"
	"github.com/floxo05/todoapi/internal/types"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (c *CategoryRepo) UpsertCategory(category *types.Category) error {

	categoryFromDb, err := c.GetCategoryFromDB(category)
	if err != nil {
		return err
	}

	var res sql.Result
	if categoryFromDb == nil {
		// if the category does not exist, insert it

		res, err = c.db.Exec("INSERT INTO categories (title, created_user_id) VALUES (?, ?)", category.Title, category.CreatedUserId)
		if err != nil {
			return err
		}

		var categoryID int64
		categoryID, err = res.LastInsertId()
		if err != nil {
			return err
		}
		category.ID = int(categoryID)
	} else {
		// if the category exists, update it
		res, err = c.db.Exec("UPDATE categories SET title = ? WHERE id = ? AND created_user_id = ?", category.Title, category.ID, category.CreatedUserId)
		if err != nil {
			return err
		}

		// if affected rows is 0, it means the category does not exist
		var affectedRows int64
		affectedRows, err = res.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("category could not be updated")
		}
	}

	return nil
}

func (c *CategoryRepo) GetCategoryFromDB(category *types.Category) (*types.Category, error) {
	var res *sql.Rows
	var err error

	// if there is an id, get the category by id, otherwise get the category by title and created_user_id
	if category.ID != 0 {
		res, err = c.db.Query("SELECT id, title, created_user_id FROM categories WHERE id = ?", category.ID)
	} else {
		res, err = c.db.Query("SELECT id, title, created_user_id FROM categories WHERE title = ? AND created_user_id = ?", category.Title, category.CreatedUserId)
	}

	if err != nil {
		return nil, err
	}
	defer res.Close()

	var newCategory types.Category
	for res.Next() {
		err = res.Scan(&newCategory.ID, &newCategory.Title, &newCategory.CreatedUserId)
		if err != nil {
			return nil, err
		}
	}

	return &newCategory, nil
}

func (c *CategoryRepo) GetCategoryByID(id int) (*types.Category, error) {
	res, err := c.db.Query("SELECT id, title, created_user_id FROM categories WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var category types.Category
	for res.Next() {
		err = res.Scan(&category.ID, &category.Title, &category.CreatedUserId)
		if err != nil {
			return nil, err
		}
	}

	return &category, nil
}
