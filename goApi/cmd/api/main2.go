package main

import (
	"github.com/floxo05/todoapi/internal/routes"
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
		log.Fatal("could not connect to the database")
	}

	// create routes
	userRoute := routes.NewUserRoute(db, UserRepo)

	// register Routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.Use(routes.JWTAuthMiddleware())
		authRoutes.POST("/todo/create", routes.CreateTodo)
		authRoutes.GET("/todos", routes.GetTodos)
		authRoutes.GET("/todo/:id", routes.GetTodoById)
		authRoutes.PUT("/todo/:id", routes.UpdateTodoById)
		authRoutes.DELETE("/todo/:id", routes.DeleteTodoById)
		authRoutes.GET("/check-token", routes.CheckToken)
		authRoutes.POST("/share", routes.ShareToUser)
	}

	r.POST("/login", routes.Login)
	r.POST("/register", routes.Register)

	// Run the server
	r.Run(":8080")
}