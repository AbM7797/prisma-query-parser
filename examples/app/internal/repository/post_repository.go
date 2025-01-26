package repository

import (
	"context"
	"log"

	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/pkg/types"
	"github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
	"github.com/AbM7797/prisma-query-parser/pkg/parser"
)

type PostRepository interface {
	Find(ctx context.Context, filter parser.Filter, customWhereFilters []db.PostWhereParam, customOrderByFilter []db.PostOrderByParam) ([]db.PostModel, error)
	Create(ctx context.Context, post domain.Post) (*db.PostModel, error)
}

type postRepository struct {
	client *db.PrismaClient
	tm     *types.TypeMapper
}

func NewPostRepository(client *db.PrismaClient, tm *types.TypeMapper) PostRepository {
	return &postRepository{
		client: client,
		tm:     tm,
	}
}

func (r *postRepository) Find(ctx context.Context, filters parser.Filter, customWhereFilters []db.PostWhereParam, customOrderByFilter []db.PostOrderByParam) ([]db.PostModel, error) {
	log.Println("filters", filters)
	whereFilters := parser.BuildWhereFilters(
		filters,
		customWhereFilters,
		r.tm,
		"post",
	)
	orderByFilters := parser.BuildOrderFilters(filters, customOrderByFilter, r.tm, "post")
	query := r.client.Post.FindMany(whereFilters...).
		OrderBy(orderByFilters...).
		With(db.Post.User.Fetch())

	if filters["take"] != nil {
		query = query.Take(filters["take"].(int))
	}

	if filters["skip"] != nil {
		query = query.Skip(filters["skip"].(int))
	}

	posts, err := query.Exec(ctx)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) Create(ctx context.Context, post domain.Post) (*db.PostModel, error) {
	createdPost, err := r.client.Post.CreateOne(
		db.Post.Title.Set(post.Title),
		db.Post.Content.Set(post.Content),
		db.Post.User.Link(db.User.ID.Equals(post.UserID)),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return createdPost, nil
}
