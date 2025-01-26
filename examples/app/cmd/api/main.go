package main

import (
	"log"

	"github.com/AbM7797/prisma-query-parser/examples/app/internal/handler"
	prisma "github.com/AbM7797/prisma-query-parser/examples/app/internal/infrastructure"
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/repository"
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/service"
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/usecase"
	"github.com/AbM7797/prisma-query-parser/examples/app/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found. Using system environment variables instead.")
	}
	client := prisma.InitPrismaClient()
	defer client.Prisma.Disconnect()
	tm := types.NewTypeMapper()

	userRepo := repository.NewUserRepository(client, tm)
	userService := service.NewUserService(userRepo)
	userUseCase := usecase.NewUserUseCase(userService)
	userHandler := handler.NewUserHandler(userUseCase)
	postRepo := repository.NewPostRepository(client, tm)
	postService := service.NewPostService(postRepo)
	postUseCase := usecase.NewPostUseCase(postService)
	postHandler := handler.NewPostHandler(postUseCase)
	commentRepo := repository.NewCommentRepository(client, tm)
	commentService := service.NewCommentService(commentRepo)
	commentUseCase := usecase.NewCommentUseCase(commentService)
	commentHandler := handler.NewCommentHandler(commentUseCase)

	r := gin.Default()
	r.GET("/users", userHandler.Find)
	r.POST("/users", userHandler.Create)
	r.GET("/posts", postHandler.Find)
	r.POST("/posts", postHandler.Create)
	r.GET("/comments", commentHandler.Find)
	r.POST("/comments", commentHandler.Create)

	r.Run(":8888")
}
