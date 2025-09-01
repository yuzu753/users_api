package usecase

import (
	"context"
	"users_api/src/domain"
)

type UserSearchIn struct {
	TenantID string
	UserName string
	Email    string
	Limit    int
	Offset   int
}

type UserSearchUsecase struct {
	repo domain.UserRepository
}

func NewUserSearchUsecase(r domain.UserRepository) *UserSearchUsecase {
	return &UserSearchUsecase{repo: r}
}

func (uc *UserSearchUsecase) Do(ctx context.Context, in UserSearchIn) ([]*domain.User, int, error) {
	return uc.repo.Search(ctx, in.TenantID, in.UserName, in.Email, in.Limit, in.Offset)
}