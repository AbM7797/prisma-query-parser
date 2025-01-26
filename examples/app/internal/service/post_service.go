package service

import (
	"context"

	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/repository"
	"github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
	"github.com/AbM7797/prisma-query-parser/pkg/parser"
)

type PostService interface {
	Find(ctx context.Context, filter parser.Filter, customWhereFilters []db.PostWhereParam, customOrderByFilter []db.PostOrderByParam) ([]db.PostModel, error)
	Create(ctx context.Context, post domain.Post) (*db.PostModel, error)
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{
		postRepo: postRepo,
	}
}

func (s *postService) Find(ctx context.Context, filter parser.Filter, customWhereFilters []db.PostWhereParam, customOrderByFilter []db.PostOrderByParam) ([]db.PostModel, error) {
	return s.postRepo.Find(ctx, filter, customWhereFilters, customOrderByFilter)
}

func (s *postService) Create(ctx context.Context, post domain.Post) (*db.PostModel, error) {
	return s.postRepo.Create(ctx, post)
}
