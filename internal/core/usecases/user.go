package usecases

import (
	"context"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
	interfaces "github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/HangoKub/Hango-service/pkg/errs"
)

type UserUsecase struct {
	UserFirestoreRepo interfaces.UserFirestoreRepository
}

func NewUserUsecase(UserFirestoreRepo interfaces.UserFirestoreRepository) *UserUsecase {
	return &UserUsecase{
		UserFirestoreRepo,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, user model.RegisterUser) (id string, err error) {
	if userDoc, _ := u.GetUserById(ctx, user.ID); (userDoc != model.User{}) {
		return "", errs.UserAlreadyExsiting
	}

	if userDoc, _ := u.GetUserByEmail(ctx, user.Email); (userDoc != model.User{}) {
		return "", errs.EmailAlreadyExsiting
	}
	return u.UserFirestoreRepo.CreateUser(ctx, user)
}

func (u *UserUsecase) GetUserById(ctx context.Context, id string) (model.User, error) {
	return u.UserFirestoreRepo.GetUserById(ctx, id)
}

func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	return u.UserFirestoreRepo.GetUserByEmail(ctx, email)
}
