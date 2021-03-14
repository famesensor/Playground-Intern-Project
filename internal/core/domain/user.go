package domain

import (
	"time"
)

type User struct {
	// UserID is userID's from other platform
	UserId    string    `json:"userId" example:"facebook_0001" firestore:"userId"`
	HgId      string    `firestore:"hgId"`
	Nickname  string    `json:"nickname,omitempty" firestore:"nickname,omitempty"`
	Email     string    `json:"email,omitempty" firestore:"email,omitempty"`
	Platform  string    `json:"platform" example:"facebook" firestore:"platform"`
	Picture   string    `json:"picture" firstore:"picture"`
	CreateAt  time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
	DeletedAt time.Time `firestore:"deletedAt"`
}

type RegisterUser struct {
	ID       string `json:"userID" validate:"required"`
	Nickname string `json:"nickname,omitempty" validate:"required,min=3,max=50"`
	Email    string `json:"email,omitempty" validate:"required_if=Platform social,omitempty,email"`
	Password string `json:"password,omitempty"`
	Platform string `json:"platform" validate:"required,oneof=social otp"`
	Gender   string `json:"gender" validate:"required"`
	Picture  string `json:"picture" validate:"omitempty,dive,url"`
	Otp      string `json:"otp,omitempty" validate:"required_if=Platform otp,omitempty,len=6,numeric"`
}

type LoginUser struct {
	ID       string `json:"userID" validate:"required"`
	Platform string `json:"platform" validate:"required,oneof=social otp"`
}
