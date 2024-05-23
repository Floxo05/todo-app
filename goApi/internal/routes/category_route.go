package routes

import (
	"github.com/floxo05/todoapi/internal/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CategoryRoute struct {
	categoryRepository types.CategoryRepository
	userContextHelper  types.UserContextInterface
}

func NewCategoryRoute(categoryRepo types.CategoryRepository, userContextHelper types.UserContextInterface) *CategoryRoute {
	return &CategoryRoute{categoryRepository: categoryRepo, userContextHelper: userContextHelper}
}

func (cr *CategoryRoute) CreateCategory(c *gin.Context) {
	user, err := cr.userContextHelper.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req types.CreateCategoryRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	//validate the request
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'title' must not be empty"})
		return
	}

	category := types.Category{Title: req.Title, CreatedUserId: user.ID}
	err = cr.categoryRepository.UpsertCategory(&category)
}
