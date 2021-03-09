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
	profileCollection string = "profiles"
)

type ProfileFirestoreRepo struct {
	firestore *firestore.Client
}

func NewProfileFirestore(firestore *firestore.Client) *ProfileFirestoreRepo {
	return &ProfileFirestoreRepo{
		firestore,
	}
}

func (r *ProfileFirestoreRepo) GetProfileByHgId(ctx context.Context, hgId string) (model.Profile, error) {
	iter := r.firestore.Collection(profileCollection).Where("hgId", "==", hgId).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return model.Profile{}, err
		}
		if doc.Exists() {
			profileDoc := new(model.Profile)
			mapstructure.Decode(doc.Data(), profileDoc)
			return *profileDoc, nil
		}
	}
	return model.Profile{}, errs.DocumentNotFound
}
