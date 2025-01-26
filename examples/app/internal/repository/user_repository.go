package repository

import (
	"context"

	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/pkg/types"
	"github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
	"github.com/AbM7797/prisma-query-parser/pkg/parser"
)

type UserRepository interface {
	Find(ctx context.Context, filter parser.Filter, whereCustomFilters []db.UserWhereParam, orderByCustomFilters []db.UserOrderByParam) ([]db.UserModel, error)
	Create(ctx context.Context, user domain.User) (*db.UserModel, error)
}

type userRepository struct {
	client *db.PrismaClient
	tm     *types.TypeMapper
}

func NewUserRepository(client *db.PrismaClient, tm *types.TypeMapper) UserRepository {
	return &userRepository{
		client: client,
		tm:     tm,
	}
}

func (r *userRepository) Find(ctx context.Context, filters parser.Filter, whereCustomFilters []db.UserWhereParam, orderByCustomFilters []db.UserOrderByParam) ([]db.UserModel, error) {
	whereFilters := parser.BuildWhereFilters(
		filters,
		whereCustomFilters,
		r.tm,
		"user",
	)

	orderByFilters := parser.BuildOrderFilters(filters, orderByCustomFilters, r.tm, "user")

	query := r.client.User.FindMany(
		whereFilters...,
	).OrderBy(orderByFilters...)

	if filters["take"] != nil {
		query.Take(filters["take"].(int))
	}

	if filters["skip"] != nil {
		query.Skip(filters["skip"].(int))
	}

	users, err := query.Exec(ctx)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Create(ctx context.Context, user domain.User) (*db.UserModel, error) {
	createdUser, err := r.client.User.CreateOne(
		db.User.Email.Set(user.Email),
		db.User.Name.Set(user.Name),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
