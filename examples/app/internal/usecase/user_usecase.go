// internal/usecase/user_usecase.go
package usecase

import (
	"context"

	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/service"
	"github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
	"github.com/AbM7797/prisma-query-parser/pkg/parser"
)

type UserUseCase interface {
	Find(ctx context.Context, filter parser.Filter) ([]db.UserModel, error)
	Create(ctx context.Context, user domain.User) (*db.UserModel, error)
}

type userUseCase struct {
	userService service.UserService
}

func NewUserUseCase(userService service.UserService) UserUseCase {
	return &userUseCase{
		userService: userService,
	}
}

func (uc *userUseCase) Find(ctx context.Context, filter parser.Filter) ([]db.UserModel, error) {
	return uc.userService.Find(ctx, filter, nil, nil)
}

func (uc *userUseCase) Create(ctx context.Context, user domain.User) (*db.UserModel, error) {
	return uc.userService.Create(ctx, user)
}
