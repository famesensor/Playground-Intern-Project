package repositories

import (
	"context"

	"cloud.google.com/go/firestore"
	model "github.com/HangoKub/Hango-service/internal/core/domain"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

const (
	authCollection string = "hango-auth"
	otpCollection  string = "hang-otp"
)

type AuthFirestoreRepo struct {
	firestore *firestore.Client
}

func NewAuthFirestore(firestore *firestore.Client) *AuthFirestoreRepo {
	return &AuthFirestoreRepo{
		firestore,
	}
}

func (r *AuthFirestoreRepo) CreateAuthenDoc(ctx context.Context, authDoc model.AuthDocument) error {
	// Overwritten the old document ...
	iter := r.firestore.Collection(authCollection).Where("hgId", "==", authDoc.HgId).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if doc.Exists() {
			_, err := r.firestore.Collection(authCollection).Doc(doc.Ref.ID).Set(ctx, authDoc)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// Create new auth document if doesn't exist ...
	_, _, err := r.firestore.Collection(authCollection).Add(ctx, authDoc)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthFirestoreRepo) GetAuthenDocByRefreshToken(ctx context.Context, refreshTk string) (model.AuthDocument, error) {
	iter := r.firestore.Collection(authCollection).Where("refreshToken", "==", refreshTk).Limit(1).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return model.AuthDocument{}, err
		}
		if doc.Exists() {
			authDoc := new(model.AuthDocument)
			mapstructure.Decode(doc.Data(), &authDoc)
			return *authDoc, nil
		}
	}
	return model.AuthDocument{}, errs.DocumentNotFound
}

func (r *AuthFirestoreRepo) RevokedToken(ctx context.Context, refreshTk string) error {
	state := r.firestore.Collection(authCollection)
	iter := state.Where("refreshToken", "==", refreshTk).Limit(1).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if doc.Exists() {
			// Update Status isRevoked for false -> ture
			_, err := state.Doc(doc.Ref.ID).Update(ctx, []firestore.Update{{Path: "isRevoked", Value: true}})
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errs.InternalServerError
}

func (r *AuthFirestoreRepo) CreateOtp(ctx context.Context, otpDoc model.OtpDoc) error {
	// Overwritten the old document ...
	state := r.firestore.Collection(otpCollection)
	iter := state.Where("phone", "==", otpDoc.Phone).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if doc.Exists() {
			_, err := state.Doc(doc.Ref.ID).Set(ctx, otpDoc)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// Create new opt document if doesn't exist ...
	_, err := state.NewDoc().Set(ctx, otpDoc)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthFirestoreRepo) GetOtpDoc(ctx context.Context, phone string) (model.OtpDoc, error) {
	state := r.firestore.Collection(otpCollection)
	iter := state.Where("phone", "==", phone).Limit(1).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return model.OtpDoc{}, err
		}
		if doc.Exists() {
			otpDoc := new(model.OtpDoc)
			mapstructure.Decode(doc.Data(), &otpDoc)
			return *otpDoc, nil
		}
	}

	return model.OtpDoc{}, errs.DocumentNotFound
}

func (r *AuthFirestoreRepo) DeleteOtp(ctx context.Context, phone string) error {
	state := r.firestore.Collection(otpCollection)
	iter := state.Where("phone", "==", phone).Limit(1).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		if doc.Exists() {
			_, err := state.Doc(doc.Ref.ID).Delete(ctx)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errs.InternalServerError
}
