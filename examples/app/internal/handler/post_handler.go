package handler

import (
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/usecase"
	"github.com/AbM7797/prisma-query-parser/pkg/parser"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	useCase usecase.PostUseCase
}

func NewPostHandler(useCase usecase.PostUseCase) *PostHandler {
	return &PostHandler{useCase: useCase}
}

func (h *PostHandler) Find(c *gin.Context) {
	filters := parser.ProcessArgs(c.Request.URL.Query())
	posts, err := h.useCase.Find(c.Request.Context(), filters)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, posts)
}

func (h *PostHandler) Create(c *gin.Context) {
	var post domain.Post
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	createdPost, err := h.useCase.Create(c.Request.Context(), post)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, createdPost)
}
