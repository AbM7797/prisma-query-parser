package service

import (
    "context"

    "github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
    "github.com/AbM7797/prisma-query-parser/examples/app/internal/repository"
    "github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
    "github.com/AbM7797/prisma-query-parser/pkg/parser"
)

type CommentService interface {
    Find(ctx context.Context, filter parser.Filter, whereCustomFilters []db.CommentWhereParam, orderByCustomFilters []db.CommentOrderByParam) ([]db.CommentModel, error)
    Create(ctx context.Context, comment domain.Comment) (*db.CommentModel, error)
}

type commentService struct {
    commentRepo repository.CommentRepository
}

func NewCommentService(commentRepo repository.CommentRepository) CommentService {
    return &commentService{
        commentRepo: commentRepo,
    }
}

func (s *commentService) Find(ctx context.Context, filter parser.Filter, whereCustomFilters []db.CommentWhereParam, orderByCustomFilters []db.CommentOrderByParam) ([]db.CommentModel, error) {
    return s.commentRepo.Find(ctx, filter, whereCustomFilters, orderByCustomFilters)
}

func (s *commentService) Create(ctx context.Context, comment domain.Comment) (*db.CommentModel, error) {
    return s.commentRepo.Create(ctx, comment)
}
