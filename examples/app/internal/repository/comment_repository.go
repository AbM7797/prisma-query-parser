package repository

import (
	"context"

	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/pkg/types"
	"github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
	"github.com/AbM7797/prisma-query-parser/pkg/parser"
)

type CommentRepository interface {
	Find(ctx context.Context, filter parser.Filter, whereCustomFilters []db.CommentWhereParam, orderByCustomFilters []db.CommentOrderByParam) ([]db.CommentModel, error)
	Create(ctx context.Context, comment domain.Comment) (*db.CommentModel, error)
}

type commentRepository struct {
	client *db.PrismaClient
	tm     *types.TypeMapper
}

func NewCommentRepository(client *db.PrismaClient, tm *types.TypeMapper) CommentRepository {
	return &commentRepository{
		client: client,
		tm:     tm,
	}
}

func (r *commentRepository) Find(ctx context.Context, filters parser.Filter, whereCustomFilters []db.CommentWhereParam, orderByCustomFilters []db.CommentOrderByParam) ([]db.CommentModel, error) {
	whereFilters := parser.BuildWhereFilters(
		filters,
		whereCustomFilters,
		r.tm,
		"comment",
	)

	orderbyFilters := parser.BuildOrderFilters(filters, orderByCustomFilters, r.tm, "comment")
	query := r.client.Comment.FindMany(
		whereFilters...,
	).OrderBy(orderbyFilters...).With(
		db.Comment.User.Fetch(),
		db.Comment.Post.Fetch(),
	)

	if filters["take"] != nil {
		query.Take(filters["take"].(int))
	}

	if filters["skip"] != nil {
		query.Skip(filters["skip"].(int))
	}

	comments, err := query.Exec(ctx)

	if err != nil {
		return nil, err
	}
	
	return comments, nil
}

func (r *commentRepository) Create(ctx context.Context, comment domain.Comment) (*db.CommentModel, error) {
	createdComment, err := r.client.Comment.CreateOne(
		db.Comment.Content.Set(comment.Content),
		db.Comment.User.Link(db.User.ID.Equals(comment.UserID)),
		db.Comment.Post.Link(db.Post.ID.Equals(comment.PostID)),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return createdComment, nil
}
