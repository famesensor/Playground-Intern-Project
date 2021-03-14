package usecases

import (
	"context"
	"time"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
	interfaces "github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/HangoKub/Hango-service/pkg/geolocation"
)

type RestaurantUsecase struct {
	RestaurantFirestoreRepo interfaces.RestaurantFirestoreRepository
}

func NewRestaurantUsecase(RestaurantFirestoreRepo interfaces.RestaurantFirestoreRepository) *RestaurantUsecase {
	return &RestaurantUsecase{
		RestaurantFirestoreRepo,
	}
}

// func (u *RestaurantUsecase) CreateRestaurant(ctx context.Context, restDoc model.CreateRestaurant) error {
// 	restDoc.GeoHash = geohash.Encode(restDoc.Location.Latitude, restDoc.Location.Longitude)
// 	return u.RestaurantFirestoreRepo.CreateRestaurant(ctx, restDoc)
// }

func (u *RestaurantUsecase) GetRestaurantById(ctx context.Context, restId string) (model.ResponseRestaurant, error) {
	return u.RestaurantFirestoreRepo.GetRestaurantById(ctx, restId)
}

func (u *RestaurantUsecase) GetRestaurantNearby(ctx context.Context, restQuery model.RestaurantQuery) ([]model.ReponseHomeRestaurant, error) {
	var timeNow, lower, upper string
	// Check distance and set default if distabce = 0
	if (restQuery.Distance) == 0 {
		restQuery.Distance = 1
	}

	// Convert km to miles
	mi := restQuery.Distance / 1.609344
	if restQuery.Open {
		timeNow = time.Now().Format("15:04")
	}

	// Hash lat,lng to string for search
	if restQuery.Lat != 0 && restQuery.Lng != 0 {
		lower, upper = geolocation.GetGeoHashRange(restQuery.Lat, restQuery.Lng, mi)
	}

	// Check limit and set default if limit = 0
	if restQuery.Limit == 0 {
		restQuery.Limit = 20
	}

	// Get restaurant nearby
	var restReponse []model.ReponseHomeRestaurant
	restDocs, err := u.RestaurantFirestoreRepo.GetRestaurantNearby(ctx, model.ListRestaurantQuery{Lower: lower, Upper: upper, TimeNow: timeNow, Tag: restQuery.Tag, Limit: restQuery.Limit, RestId: restQuery.RestId})
	if err != nil {
		return nil, err
	}

	if restQuery.Lat != 0 && restQuery.Lng != 0 {
		for _, doc := range restDocs {
			dist := geolocation.DistanceBetween(doc.Location.Latitude, doc.Location.Longitude, restQuery.Lat, restQuery.Lng, "K")
			// log.Printf("Distance : %v", dist)
			if restQuery.Distance >= dist {
				restReponse = append(restReponse, doc)
			}
		}
		return restReponse, nil
	}

	return restDocs, nil
}

// func (u *RestaurantUsecase) GetRestaurantByName(ctx context.Context, restName, lastRestId string) ([]model.ReponseHomeRestaurant, error) {
// 	return u.RestaurantFirestoreRepo.GetRestaurantByName(ctx, restName, lastRestId)
// }

func (u *RestaurantUsecase) CheckInRestaurant(ctx context.Context, hgId string, userDoc model.UserCheckIn) (string, error) {
	// Check if user check-in other restaurant
	checkStay, err := u.RestaurantFirestoreRepo.GetCheckInByHgId(ctx, hgId)
	if (checkStay != model.CheckInDoc{}) {
		if !checkStay.IsRevoked && checkStay.RestId != userDoc.RestId {
			err := u.RestaurantFirestoreRepo.CheckOut(ctx, hgId, checkStay.RestId)
			if err != nil {
				return "", err
			}
		} else {
			return "User checked-in", nil
		}
	}

	// Check resturant existing...
	restDoc, err := u.RestaurantFirestoreRepo.GetRestaurantById(ctx, userDoc.RestId)
	if err != nil {
		return "", err
	}

	// Calculate distance between two location and check reduis user
	dist := geolocation.DistanceBetween(restDoc.Location.Latitude, restDoc.Location.Longitude, userDoc.UserLat, userDoc.Userlng, "K")
	if dist*1000 >= 5000 { // 500 meters or restDoc.info.Reduis
		return "", errs.UserOutofReduis
	}

	// Check if user check-in with out time range
	timeNow := time.Now().Format("15:04")
	if timeNow <= restDoc.Info.Open && restDoc.Info.Close <= timeNow {
		return "", errs.YouCannotCheckInApart
	}

	// Create CheckInDoc
	err = u.RestaurantFirestoreRepo.CreateCheckIn(ctx, hgId, userDoc.RestId)
	if err != nil {
		return "", err
	}
	// Create location temp for check user stay reduis of restaurant

	return "check-in", nil
}

func (u *RestaurantUsecase) GetCheckIn(ctx context.Context, hgId, restId string) (model.CheckInDoc, error) {
	checkDoc, err := u.RestaurantFirestoreRepo.GetCheckIn(ctx, hgId, restId)
	if checkDoc.IsRevoked || (checkDoc == model.CheckInDoc{}) {
		return model.CheckInDoc{}, errs.UserNotAllowed
	}
	if err != nil {
		return model.CheckInDoc{}, err
	}

	return checkDoc, nil
}

func (u *RestaurantUsecase) GetCheckInByHgId(ctx context.Context, hgId string) (model.CheckInDoc, error) {
	checkDoc, err := u.RestaurantFirestoreRepo.GetCheckInByHgId(ctx, hgId)
	if checkDoc.IsRevoked || (checkDoc == model.CheckInDoc{}) {
		return model.CheckInDoc{}, errs.UserNotAllowed
	}
	if err != nil {
		return model.CheckInDoc{}, err
	}

	return checkDoc, nil
}

func (u *RestaurantUsecase) PeakMode(ctx context.Context, hgId string) (string, error) {
	return u.RestaurantFirestoreRepo.PeakMode(ctx, hgId)
}

func (u *RestaurantUsecase) CheckOutRestaurant(ctx context.Context, hgId, restId string) (string, error) {
	if err := u.RestaurantFirestoreRepo.CheckOut(ctx, hgId, restId); err != nil {
		return "", err
	}

	return "check-out", nil
}

func (u *RestaurantUsecase) GetCheckInAnonymous(ctx context.Context, hgId, restId string) (anonymousDocs []model.AnonymousDoc, err error) {
	return u.RestaurantFirestoreRepo.GetCheckInAnonymous(ctx, restId)
}

func (u *RestaurantUsecase) GetCheckInPeakMode(ctx context.Context, hgId, restId string) (peakModeDocs []model.PeakModeDoc, err error) {
	return u.RestaurantFirestoreRepo.GetCheckInPeakMode(ctx, restId)
}
