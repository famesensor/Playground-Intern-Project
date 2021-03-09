package handlers

import interfaces "github.com/HangoKub/Hango-service/internal/core/ports"

type ProfileHandler struct {
	profileUc interfaces.ProfileUsecase
}

func NewProfileHandler(profileUc interfaces.ProfileUsecase) interfaces.ProfileHandler {
	return &ProfileHandler{
		profileUc,
	}
}
