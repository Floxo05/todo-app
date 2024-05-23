package main

import (
	"github.com/floxo05/todoapi/internal/repository"
	"github.com/floxo05/todoapi/internal/routes"
	"github.com/floxo05/todoapi/internal/services"
	"github.com/floxo05/todoapi/internal/tools"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	// enable CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	r.Use(cors.New(config))

	db, err := tools.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	catRepo := repository.NewCategoryRepo(db)
	todoRepo := repository.NewTodoRepo(db, catRepo)
	userRepo := repository.NewUserRepo(db, todoRepo)
	userContextHelper := services.NewUserContext(userRepo)
	passwordHasher := services.NewPasswordHasher()

	todoRoute := routes.NewTodoRoute(todoRepo, userContextHelper)
	userRoute := routes.NewUserRoute(userRepo, passwordHasher, userContextHelper)
	tokenRoute := routes.NewTokenRoute()
	catRoute := routes.NewCategoryRoute(catRepo, userContextHelper)

	// register Routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.Use(routes.JWTAuthMiddleware())
		authRoutes.POST("/todo/create", todoRoute.CreateTodo)
		authRoutes.GET("/todos", todoRoute.GetTodos)
		authRoutes.PUT("/todo/:id", todoRoute.UpdateTodo)
		authRoutes.DELETE("/todo/:id", todoRoute.DeleteTodo)
		authRoutes.GET("/check-token", tokenRoute.CheckToken)
		authRoutes.POST("/share", userRoute.ShareToUser)
		authRoutes.POST("/category/create", catRoute.CreateCategory)
		authRoutes.GET("/categories", catRoute.GetCategories)
	}

	r.POST("/login", userRoute.Login)
	r.POST("/register", userRoute.Register)

	// Run the server
	r.Run(":8080")
}
