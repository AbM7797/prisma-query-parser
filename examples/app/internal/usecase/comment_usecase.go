package usecase

import (
    "context"

    "github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
    "github.com/AbM7797/prisma-query-parser/examples/app/internal/service"
    "github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
    "github.com/AbM7797/prisma-query-parser/pkg/parser"
)

type CommentUseCase interface {
    Find(ctx context.Context, filter parser.Filter) ([]db.CommentModel, error)
    Create(ctx context.Context, comment domain.Comment) (*db.CommentModel, error)
}

type commentUseCase struct {
    commentService service.CommentService
}

func NewCommentUseCase(commentService service.CommentService) CommentUseCase {
    return &commentUseCase{
        commentService: commentService,
    }
}

func (uc *commentUseCase) Find(ctx context.Context, filter parser.Filter) ([]db.CommentModel, error) {
    return uc.commentService.Find(ctx, filter, nil, nil)
}

func (uc *commentUseCase) Create(ctx context.Context, comment domain.Comment) (*db.CommentModel, error) {
    return uc.commentService.Create(ctx, comment)
}
