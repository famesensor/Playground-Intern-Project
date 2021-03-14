package usecases

import (
	"context"
	"crypto/rsa"
	"time"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
	interfaces "github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/HangoKub/Hango-service/pkg"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/HangoKub/Hango-service/pkg/notify"
	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
)

type AuthUsecase struct {
	firestoreRepo  interfaces.AuthFirestoreRepository
	prirKey        *rsa.PrivateKey
	publicKey      *rsa.PublicKey
	messageBirdKey string
}

func NewAuthUsecase(firestoreRepo interfaces.AuthFirestoreRepository, prirKey *rsa.PrivateKey, publicKey *rsa.PublicKey, messageBirdKey string) *AuthUsecase {
	return &AuthUsecase{
		firestoreRepo,
		prirKey,
		publicKey,
		messageBirdKey,
	}
}

func (u *AuthUsecase) ReqTokenDocument(ctx context.Context, tokenType string, claims *model.AuthClaim) (model.TokenVerificationDocuments, error) {
	iat := time.Now()
	exp := time.Now().Add(time.Minute * 10)
	rfExp := time.Now().Add(time.Hour * 168)
	stdClaim := jwt.StandardClaims{
		IssuedAt:  iat.Unix(),
		ExpiresAt: exp.Unix(),
	}

	claims.StandardClaims = stdClaim
	inf := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	acTk, err := inf.SignedString(u.prirKey)
	if err != nil {
		return model.TokenVerificationDocuments{}, err
	}
	refreshTk := uuid.Must(uuid.NewRandom()).String() + pkg.RandUUID(9999999999, 1000000000, 32)
	tkVerifyDoc := generateTokenVerifyDocument(acTk, refreshTk, tokenType, rfExp, claims)

	return tkVerifyDoc, nil
}

func (u *AuthUsecase) CreateAuthDoc(ctx context.Context, authDoc model.AuthDocument) error {
	return u.firestoreRepo.CreateAuthenDoc(ctx, authDoc)
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshTk string) (model.TokenCard, error) {
	rfToken, err := u.firestoreRepo.GetAuthenDocByRefreshToken(ctx, refreshTk)
	if err != nil {
		return model.TokenCard{}, errs.OperationIsnotAllowed
	}

	if rfToken.IsRevoked || rfToken.ExpiresAt.Before(time.Now()) {
		return model.TokenCard{}, errs.InvalidToken
	}

	clamis := &model.AuthClaim{
		HgId: rfToken.HgId,
	}

	TkDoc, err := u.ReqTokenDocument(ctx, "bearer", clamis)
	if err != nil {
		return model.TokenCard{}, err
	}

	if -(time.Now().Sub(rfToken.ExpiresAt).Hours()) > 24 {
		TkDoc.AuthDocument.RefreshToken = rfToken.RefreshToken
		TkDoc.TokenCard.RefreshToken = rfToken.RefreshToken
		TkDoc.AuthDocument.ExpiresAt = rfToken.ExpiresAt
	}
	if err := u.CreateAuthDoc(ctx, TkDoc.AuthDocument); err != nil {
		return model.TokenCard{}, err
	}

	return TkDoc.TokenCard, nil
}

func (u *AuthUsecase) Logout(ctx context.Context, refreshTk string) error {
	return u.firestoreRepo.RevokedToken(ctx, refreshTk)
}

func (u *AuthUsecase) VerifyToken(token string) (string, error) {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return u.publicKey, nil
	})
	if err != nil {
		v, _ := err.(jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return model.TOKEN_EXPIRED, err
		} else {
			return model.TOKEN_INVALID, err
		}
	}
	return model.TOKEN_VALID, nil
}

func (u *AuthUsecase) ReqOtpDocument(ctx context.Context, phone string, otpT time.Duration, typeOTP string) error {
	otpExp := time.Now().Add(otpT)

	otp := pkg.RandomNumber(999999, 000000)
	if err := notify.SendOTP(phone, otp, u.messageBirdKey); err != nil {
		return err
	}

	if err := u.firestoreRepo.CreateOtp(ctx, model.OtpDoc{Otp: otp, Phone: phone, TypeOTP: typeOTP, ExpiresAt: otpExp, CreatedAt: time.Now()}); err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) RefreshOtp(ctx context.Context, phone, typeOTP string) error {
	otpDoc, err := u.firestoreRepo.GetOtpDoc(ctx, phone)
	if err != nil {
		return errs.OperationIsnotAllowed
	}

	if -((time.Now().Sub(otpDoc.ExpiresAt)).Minutes()) > 1 {
		return errs.FailedReqOTP
	}

	otpExp := time.Minute * 5
	if err := u.ReqOtpDocument(ctx, phone, otpExp, typeOTP); err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) VerifyOtp(ctx context.Context, otp, phone, option, typeOTP string) error {
	otpDoc, err := u.firestoreRepo.GetOtpDoc(ctx, phone)
	if err != nil {
		if err == errs.DocumentNotFound {
			return err
		}
		return errs.OperationIsnotAllowed
	}

	if otpDoc.Otp != otp || otpDoc.TypeOTP != typeOTP {
		return errs.InvalidOTP
	}

	switch option {
	case model.LOGIN_OTP:
		if otpDoc.ExpiresAt.Before(time.Now()) {
			return errs.OtpExpires
		}
		if err = u.firestoreRepo.DeleteOtp(ctx, phone); err != nil {
			return err
		}
		return nil
	case model.REGISTER_OTP:
		if otpDoc.ExpiresAt.Before(time.Now()) {
			return errs.OtpExpires
		}
		return nil
	default:
		if err = u.firestoreRepo.DeleteOtp(ctx, phone); err != nil {
			return err
		}
		return nil
	}
}

func generateTokenVerifyDocument(accessTk, refreshTk, tokenType string, rfTk time.Time, authClaims *model.AuthClaim) model.TokenVerificationDocuments {
	return model.TokenVerificationDocuments{
		TokenCard: model.TokenCard{
			AccessToken:  accessTk,
			RefreshToken: refreshTk,
		},
		AuthDocument: model.AuthDocument{
			HgId:         authClaims.HgId,
			RefreshToken: refreshTk,
			IsRevoked:    false,
			TokenType:    tokenType,
			ExpiresAt:    rfTk,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}
}
