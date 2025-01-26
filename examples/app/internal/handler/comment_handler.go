package handler

import (
    "github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
    "github.com/AbM7797/prisma-query-parser/examples/app/internal/usecase"
    "github.com/AbM7797/prisma-query-parser/pkg/parser"
    "github.com/gin-gonic/gin"
)

type CommentHandler struct {
    useCase usecase.CommentUseCase
}

func NewCommentHandler(useCase usecase.CommentUseCase) *CommentHandler {
    return &CommentHandler{useCase: useCase}
}

func (h *CommentHandler) Find(c *gin.Context) {
    filters := parser.ProcessArgs(c.Request.URL.Query())
    comments, err := h.useCase.Find(c.Request.Context(), filters)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, comments)
}

func (h *CommentHandler) Create(c *gin.Context) {
    var comment domain.Comment
    err := c.BindJSON(&comment)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    createdComment, err := h.useCase.Create(c.Request.Context(), comment)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(201, createdComment)
}
