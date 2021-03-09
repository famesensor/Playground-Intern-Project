package usecases

import (
	"context"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
	interfaces "github.com/HangoKub/Hango-service/internal/core/ports"
)

type ProfileUsecase struct {
	ProfileFrestoreRepo interfaces.ProfileFirestoreRepository
}

func NewProfileUsecase(ProfileFrestoreRepo interfaces.ProfileFirestoreRepository) *ProfileUsecase {
	return &ProfileUsecase{
		ProfileFrestoreRepo,
	}
}

func (u *ProfileUsecase) GetProfileByHgId(ctx context.Context, hgId string) (model.Profile, error) {
	return u.ProfileFrestoreRepo.GetProfileByHgId(ctx, hgId)
}
