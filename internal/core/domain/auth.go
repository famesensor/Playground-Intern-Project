package domain

import (
	"time"

	"github.com/form3tech-oss/jwt-go"
)

const (
	TOKEN_VALID     string = "VALID"
	TOKEN_INVALID   string = "INVALID"
	TOKEN_EXPIRED   string = "EXPIRED"
	VERIFY_AUTH_OTP string = "AUTH_OTP"
	REGISTER_OTP    string = "REGISTER_OTP"
	LOGIN_OTP       string = "LOGIN_OTP"
)

type AuthClaim struct {
	HgId string `json:"hgId"`
	jwt.StandardClaims
}

type TokenVerificationDocuments struct {
	TokenCard    TokenCard
	AuthDocument AuthDocument
}
type TokenCard struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type AuthDocument struct {
	HgId         string    `firestore:"hgId"`
	RefreshToken string    `firestore:"refreshToken"`
	IsRevoked    bool      `firestore:"isRevoked"`
	TokenType    string    `firestore:"tokenType"`
	ExpiresAt    time.Time `firestore:"expiresAt"`
	CreatedAt    time.Time `firestore:"createdAt"`
	UpdatedAt    time.Time `firestore:"updatedAt"`
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type OtpDoc struct {
	Otp       string    `firestore:"otp"`
	Phone     string    `firestore:"phone"`
	TypeOTP   string    `firestore:"typeOTP"`
	ExpiresAt time.Time `firestore:"expiresAt"`
	CreatedAt time.Time `firestore:"createdAt"`
}

type Otp struct {
	Phone string `json:"phone" validate:"required,len=10,numeric"`
	Otp   string `json:"otp" validate:"required,len=6,numeric"`
}

type Phone struct {
	Phone string `json:"phone" validate:"required,len=10,numeric"`
}
