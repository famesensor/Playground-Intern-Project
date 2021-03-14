package repositories

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	model "github.com/HangoKub/Hango-service/internal/core/domain"
	"github.com/HangoKub/Hango-service/pkg"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	restaurantCollection  string = "restaurant"
	userCheckInCollection string = "user-checkIn"
)

type RestaurantFirestoreRepo struct {
	firestore *firestore.Client
}

func NewRestaurantFirestore(firestore *firestore.Client) *RestaurantFirestoreRepo {
	return &RestaurantFirestoreRepo{
		firestore,
	}
}

// // mock for create restaurant
// func (r *RestaurantFirestoreRepo) CreateRestaurant(ctx context.Context, restDoc model.CreateRestaurant) error {
// 	batch := r.firestore.Batch()

// 	ref := r.firestore.Collection(restaurantCollection).NewDoc()
// 	batch.Set(ref, model.CreateRestaurant{
// 		RestId:   ref.ID,
// 		RestName: restDoc.RestName,
// 		Location: &latlng.LatLng{
// 			Latitude:  restDoc.Location.Latitude,
// 			Longitude: restDoc.Location.Longitude,
// 		},
// 		GeoHash:   restDoc.GeoHash,
// 		Info:      restDoc.Info,
// 		Tag:       restDoc.Tag,
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	})
// 	if _, err := batch.Commit(ctx); err != nil {
// 		batch.Delete(ref)
// 		batch.Commit(ctx)
// 		return err
// 	}

// 	return nil
// }

func (r *RestaurantFirestoreRepo) GetRestaurantById(ctx context.Context, restId string) (model.ResponseRestaurant, error) {
	restDoc := new(model.ResponseRestaurant)
	dsnap, err := r.firestore.Collection(restaurantCollection).Doc(restId).Get(ctx)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return model.ResponseRestaurant{}, errs.DocumentNotFound
		}
		return model.ResponseRestaurant{}, err
	}
	if dsnap.Exists() {
		dsnap.DataTo(restDoc)
	}

	return *restDoc, nil
}

func (r *RestaurantFirestoreRepo) GetRestaurantNearby(ctx context.Context, restaurantQuery model.ListRestaurantQuery) ([]model.ReponseHomeRestaurant, error) {
	ref := r.firestore.Collection(restaurantCollection)
	// var query firestore.Query
	query := ref.Query

	if restaurantQuery.Lower != "" && restaurantQuery.Upper != "" {
		// log.Println("Query location : ", restaurantQuery.Lower, restaurantQuery.Upper)
		query = query.Where("geohash", ">=", restaurantQuery.Lower).Where("geohash", "<=", restaurantQuery.Upper).OrderBy("geohash", firestore.Desc)
	}
	if len(restaurantQuery.Tag) > 0 {
		// log.Println("Tag : ", restaurantQuery.Tag)
		query = query.Where("tag", "array-contains-any", restaurantQuery.Tag)
	}
	if restaurantQuery.RestId != "" {
		// log.Println("Pagination")
		dsnap, err := ref.Doc(restaurantQuery.RestId).Get(ctx)
		if err != nil {
			if grpc.Code(err) == codes.NotFound {
				return []model.ReponseHomeRestaurant{}, errs.DocumentNotFound
			}
			return []model.ReponseHomeRestaurant{}, err
		}
		query = query.StartAfter(dsnap.Data()["geohash"])
	}
	iter := query.Limit(restaurantQuery.Limit).Documents(ctx)
	defer iter.Stop()

	var restDocs []model.ReponseHomeRestaurant
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []model.ReponseHomeRestaurant{}, err
		}
		if doc.Exists() {
			// Check if time isn't period of time
			if restaurantQuery.TimeNow != "" {
				info, _ := doc.DataAt("info")
				temp := info.(map[string]interface{})
				if restaurantQuery.TimeNow <= temp["open"].(string) && temp["close"].(string) <= restaurantQuery.TimeNow {
					continue
				}
			}

			restDoc := new(model.ReponseHomeRestaurant)
			mapstructure.Decode(doc.Data(), &restDoc)
			restDocs = append(restDocs, *restDoc)
		}
	}

	return restDocs, nil
}

// TODO: Implement performance random anonymous, change url picture anonymous, fix concurrency
func (r *RestaurantFirestoreRepo) CreateCheckIn(ctx context.Context, hgId, restId string) error {
	ref := r.firestore.Doc(restaurantCollection + "/" + restId).Collection(userCheckInCollection)

	// Ramdom number for anonymous after user check-in and peakmode
	uid := pkg.RandomNumber(9999, 0000)
	_ = rand.Intn(64)
	var n int
	for {
		if n == 5 {
			return errors.New("run time error")
		}
		doc, err := ref.Where("uid", "==", uid).Documents(ctx).GetAll()
		if err != nil {
			return err
		}
		if len(doc) == 0 {
			break
		}
		n++
	}

	checkDoc := model.CheckInDoc{
		RestId:    restId,
		UId:       uid,
		NameAnon:  "Hanger" + uid,
		Picture:   "https://firebasestorage.googleapis.com/v0/b/hango-dev-32d20.appspot.com/o/utils%2Fanonymous%2Fhanger.svg?alt=media",
		HgId:      hgId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Overwritten the old document ...
	iter := ref.Where("hgId", "==", hgId).Limit(1).Documents(ctx)
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
			_, err := ref.Doc(doc.Ref.ID).Set(ctx, checkDoc)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// Create new check-in document if doesn't existing... and update count user check-in restaurant
	err := r.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		err := tx.Update(ref.Parent, []firestore.Update{{Path: "info.checkInCount", Value: firestore.Increment(1)}})
		if err != nil {
			return err
		}
		return tx.Set(ref.NewDoc(), checkDoc)
	})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return errs.DocumentNotFound
		}
		return err
	}

	return nil
}

func (r *RestaurantFirestoreRepo) GetCheckIn(ctx context.Context, hgId, restId string) (model.CheckInDoc, error) {
	iter := r.firestore.Doc(restaurantCollection+"/"+restId).Collection(userCheckInCollection).Where("hgId", "==", hgId).Limit(1).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return model.CheckInDoc{}, err
		}
		if doc.Exists() {
			checkInDoc := new(model.CheckInDoc)
			doc.DataTo(&checkInDoc)
			return *checkInDoc, nil
		}
	}

	return model.CheckInDoc{}, nil
}

func (r *RestaurantFirestoreRepo) GetCheckInByHgId(ctx context.Context, hgId string) (model.CheckInDoc, error) {
	iter := r.firestore.CollectionGroup(userCheckInCollection).Where("hgId", "==", hgId).Where("isRevoked", "==", false).Limit(1).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return model.CheckInDoc{}, err
		}
		if doc.Exists() {
			checkInDoc := new(model.CheckInDoc)
			mapstructure.Decode(doc.Data(), checkInDoc)
			return *checkInDoc, nil
		}
	}

	return model.CheckInDoc{}, nil
}

func (r *RestaurantFirestoreRepo) PeakMode(ctx context.Context, hgId string) (string, error) {
	iter := r.firestore.CollectionGroup(userCheckInCollection).Where("hgId", "==", hgId).Limit(1).Documents(ctx)
	defer iter.Stop()

	var result string
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}
		if doc.Exists() {
			peakMode := false
			switch doc.Data()["peakMode"] {
			case true:
				result = "Peak-Mode : close"
				peakMode = false
				break
			case false:
				result = "Peak-Mode : open"
				peakMode = true
				break
			}
			_, err := r.firestore.Doc(restaurantCollection+"/"+doc.Ref.Parent.Parent.ID).Collection(userCheckInCollection).Doc(doc.Ref.ID).Update(ctx, []firestore.Update{
				{Path: "peakMode", Value: peakMode},
			})
			if err != nil {
				return "", err
			}
			return result, nil
		}
	}

	return result, errs.DocumentNotFound
}

func (r *RestaurantFirestoreRepo) CheckOut(ctx context.Context, hgId, restId string) error {
	restRef := r.firestore.Doc(restaurantCollection + "/" + restId)
	iter := restRef.Collection(userCheckInCollection).Where("hgId", "==", hgId).Where("isRevoked", "==", false).Limit(1).Documents(ctx)
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
			_, err := restRef.Collection(userCheckInCollection).Doc(doc.Ref.ID).Update(ctx, []firestore.Update{
				{Path: "isRevoked", Value: true},
				{Path: "peakMode", Value: false},
				{Path: "updatedAt", Value: time.Now()},
				{Path: "nameAnon", Value: ""},
				{Path: "picture", Value: ""},
				{Path: "uid", Value: ""},
			})
			if err != nil {
				if grpc.Code(err) == codes.NotFound {
					return errs.DocumentNotFound
				}
				return err
			}
			return nil
		}
	}

	return errs.InternalServerError
}

func (r *RestaurantFirestoreRepo) GetCheckInAnonymous(ctx context.Context, restId string) (anonymousDocs []model.AnonymousDoc, err error) {
	restRef := r.firestore.Doc(restaurantCollection + "/" + restId)
	iter := restRef.Collection(userCheckInCollection).Where("isRevoked", "==", false).Where("peakMode", "==", false).OrderBy("createdAt", firestore.Desc).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return anonymousDocs, err
		}
		if doc.Exists() {
			anonymousDoc := new(model.AnonymousDoc)
			mapstructure.Decode(doc.Data(), anonymousDoc)
			anonymousDocs = append(anonymousDocs, *anonymousDoc)
		}
	}

	return anonymousDocs, nil
}

func (r *RestaurantFirestoreRepo) GetCheckInPeakMode(ctx context.Context, restId string) (peakModeDocs []model.PeakModeDoc, err error) {
	restRef := r.firestore.Doc(restaurantCollection + "/" + restId)
	iter := restRef.Collection(userCheckInCollection).Where("isRevoked", "==", false).Where("peakMode", "==", true).OrderBy("createdAt", firestore.Desc).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return peakModeDocs, err
		}
		if doc.Exists() {
			peakModeDoc := new(model.PeakModeDoc)
			doc.DataTo(peakModeDoc)

			// Get profile user
			userIter, err := r.firestore.Collection(userCollection).Where("hgId", "==", peakModeDoc.HgId).Limit(1).Documents(ctx).GetAll()
			if err != nil {
				return []model.PeakModeDoc{}, err
			}
			profileDoc := new(model.User)
			userIter[0].DataTo(profileDoc)
			peakModeDoc.NickName = profileDoc.Nickname
			peakModeDoc.Picture = profileDoc.Picture

			peakModeDocs = append(peakModeDocs, *peakModeDoc)
		}
	}

	return peakModeDocs, nil
}
