package usecase

import (
	"context"

	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/service"
	"github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
	"github.com/AbM7797/prisma-query-parser/pkg/parser"
)

type PostUseCase interface {
	Find(ctx context.Context, filter parser.Filter) ([]db.PostModel, error)
	Create(ctx context.Context, post domain.Post) (*db.PostModel, error)
}

type postUseCase struct {
	postService service.PostService
}

func NewPostUseCase(postService service.PostService) PostUseCase {
	return &postUseCase{
		postService: postService,
	}
}

func (uc *postUseCase) Find(ctx context.Context, filter parser.Filter) ([]db.PostModel, error) {
	return uc.postService.Find(ctx, filter, nil, nil)
}

func (uc *postUseCase) Create(ctx context.Context, post domain.Post) (*db.PostModel, error) {
	return uc.postService.Create(ctx, post)
}
