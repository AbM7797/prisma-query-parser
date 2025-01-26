// internal/service/user_service.go
package service

import (
	"context"

	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/repository"
	"github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
	"github.com/AbM7797/prisma-query-parser/pkg/parser"
)

type UserService interface {
	Find(ctx context.Context, filter parser.Filter, whereCustomFilters []db.UserWhereParam, orderByCustomFilters []db.UserOrderByParam) ([]db.UserModel, error)
	Create(ctx context.Context, user domain.User) (*db.UserModel, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Find(ctx context.Context, filter parser.Filter, whereCustomFilters []db.UserWhereParam, orderByCustomFilters []db.UserOrderByParam) ([]db.UserModel, error) {
	return s.userRepo.Find(ctx, filter, whereCustomFilters, orderByCustomFilters)
}

func (s *userService) Create(ctx context.Context, user domain.User) (*db.UserModel, error) {
	return s.userRepo.Create(ctx, user)
}
