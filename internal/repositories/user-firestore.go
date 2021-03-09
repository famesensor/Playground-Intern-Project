package repositories

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	model "github.com/HangoKub/Hango-service/internal/core/domain"
	"github.com/HangoKub/Hango-service/pkg"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

const (
	userCollection string = "users"
)

type UserFirestoreRepo struct {
	firestore *firestore.Client
}

func NewUserFirestore(firestore *firestore.Client) *UserFirestoreRepo {
	return &UserFirestoreRepo{
		firestore,
	}
}

func (r *UserFirestoreRepo) CreateUser(ctx context.Context, user model.RegisterUser) (id string, err error) {
	// Generate HgId ...
	hgId := "#hg." + pkg.RandomNumber(9999999999, 1000000000)
	var n int
	for {
		if n == 5 {
			return "", errors.New("runtime error")
		}
		doc, err := r.firestore.Collection(userCollection).Where("hdId", "==", hgId).Documents(ctx).GetAll()
		if err != nil {
			return "", err
		}
		if len(doc) == 0 {
			break
		}
		n++
	}

	batch := r.firestore.Batch()
	userRef := r.firestore.Collection(userCollection).Doc(hgId)
	batch.Set(userRef, model.User{
		UserId:    user.ID,
		Email:     user.Email,
		HgId:      hgId,
		Nickname:  user.Nickname,
		Platform:  user.Platform,
		CreateAt:  time.Now(),
		UpdatedAt: time.Now(),
	})

	profRef := r.firestore.Collection(profileCollection).NewDoc()
	batch.Set(profRef, model.Profile{
		HgId:      hgId,
		Nickname:  user.Nickname,
		Gender:    user.Gender,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	_, err = batch.Commit(ctx)
	if err != nil {
		batch.Delete(userRef)
		batch.Delete(profRef)
		batch.Commit(ctx)
		return "", err
	}

	return hgId, nil
}

func (r *UserFirestoreRepo) GetUserById(ctx context.Context, id string) (model.User, error) {
	iter := r.firestore.Collection(userCollection).Where("userId", "==", id).Limit(1).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return model.User{}, err
		}
		if doc.Exists() {
			userDoc := new(model.User)
			mapstructure.Decode(doc.Data(), userDoc)
			return *userDoc, nil
		}
	}
	return model.User{}, errs.DocumentNotFound
}

func (r *UserFirestoreRepo) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	iter := r.firestore.Collection(userCollection).Where("email", "==", email).Limit(1).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return model.User{}, err
		}
		if doc.Exists() {
			userDoc := new(model.User)
			mapstructure.Decode(doc.Data(), userDoc)
			return *userDoc, nil
		}
	}
	return model.User{}, errs.DocumentNotFound
}
