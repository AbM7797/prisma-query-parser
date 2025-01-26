package handler

import (
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/usecase"
	"github.com/AbM7797/prisma-query-parser/pkg/parser"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	useCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: useCase}
}

func (h *UserHandler) Find(c *gin.Context) {
	filters := parser.ProcessArgs(c.Request.URL.Query())
	users, err := h.useCase.Find(c.Request.Context(), filters)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, users)
}

func (h *UserHandler) Create(c *gin.Context) {
	var user domain.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.useCase.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, createdUser)
}
