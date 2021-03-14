package usecases

import (
	"context"
	"time"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
	interfaces "github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/HangoKub/Hango-service/pkg/resizeImage"
	"github.com/google/uuid"
)

const (
	roomIdPublic string = "Hango-Feed-Public"
)

type PostFeedsUsecase struct {
	RestaurantFirestoreRepo interfaces.RestaurantFirestoreRepository
	PostFirestoreRepo       interfaces.PostFeedsFirestoreRepository
	UploadFileRepo          interfaces.UploadStorageRepository
}

func NewPostUsecase(RestaurantFirestoreRepo interfaces.RestaurantFirestoreRepository, PostFirestoreRepo interfaces.PostFeedsFirestoreRepository, UploadFileRepo interfaces.UploadStorageRepository) *PostFeedsUsecase {
	return &PostFeedsUsecase{
		RestaurantFirestoreRepo,
		PostFirestoreRepo,
		UploadFileRepo,
	}
}

func (u *PostFeedsUsecase) CreateHangoPost(ctx context.Context, postDoc model.InputPost, hgId string) (model.HangoPost, error) {
	if postDoc.RoomId != roomIdPublic {
		if _, err := u.RestaurantFirestoreRepo.GetRestaurantById(ctx, postDoc.RoomId); err != nil {
			return model.HangoPost{}, err
		}
		checkDoc, err := u.RestaurantFirestoreRepo.GetCheckIn(ctx, hgId, postDoc.RoomId)
		if checkDoc.IsRevoked || (checkDoc == model.CheckInDoc{}) {
			return model.HangoPost{}, errs.UserNotAllowed
		}
		if err != nil {
			return model.HangoPost{}, err
		}
	}

	// if user isn't upload image
	if postDoc.PostId == "" {
		postDoc.PostId = "Hango-" + uuid.Must(uuid.NewRandom()).String() + time.Now().Format("2006-01-02-15-04-05")
	}

	interActor := model.InterActor{
		HgId:    hgId,
		Aka:     "Hanger",
		Picture: "https://firebasestorage.googleapis.com/v0/b/hango-dev-32d20.appspot.com/o/utils%2Fanonymous%2Fhanger.svg?alt=media",
	}

	hangPost, err := u.PostFirestoreRepo.CreateHangoPost(ctx, postDoc, interActor)
	if err != nil {
		return model.HangoPost{}, err
	}

	return hangPost, nil
}

func (u *PostFeedsUsecase) DeleteHangoPost(ctx context.Context, postId, hgId string) error {
	postDoc, err := u.GetPostById(ctx, hgId, postId)
	if err != nil {
		return err
	}
	if hgId != postDoc.HgId {
		return errs.UserNotAllowed
	}

	if err := u.PostFirestoreRepo.DeleteHangoPost(ctx, postId); err != nil {
		return err
	}

	return nil
}

func (u *PostFeedsUsecase) GetAllHangoPost(ctx context.Context, postQuery model.FeedPostQuery) ([]model.PostFeedReponse, error) {
	if postQuery.RoomId != roomIdPublic && postQuery.RoomId != "" {
		if _, err := u.RestaurantFirestoreRepo.GetRestaurantById(ctx, postQuery.RoomId); err != nil {
			return []model.PostFeedReponse{}, err
		}
	}
	if postQuery.Limit == 0 {
		postQuery.Limit = 20
	}
	if postQuery.RoomId == "" {
		postQuery.RoomId = roomIdPublic
	}

	return u.PostFirestoreRepo.GetAllHangoPost(ctx, postQuery)
}

func (u *PostFeedsUsecase) GetPostById(ctx context.Context, hgId, postId string) (model.HangoPost, error) {
	postDoc, err := u.PostFirestoreRepo.GetHangoPostByID(ctx, postId)
	if err != nil {
		return model.HangoPost{}, err
	}
	if postDoc.RoomId != roomIdPublic {
		checkDoc, err := u.RestaurantFirestoreRepo.GetCheckIn(ctx, hgId, postDoc.RoomId)
		if checkDoc.IsRevoked || (checkDoc == model.CheckInDoc{}) {
			return model.HangoPost{}, errs.UserNotAllowed
		}
		if err != nil {
			return model.HangoPost{}, err
		}
	}

	return postDoc, nil
}

func (u *PostFeedsUsecase) ReportPost(ctx context.Context, postReport model.ReportPost) error {
	postDoc, err := u.GetPostById(ctx, postReport.HgId, postReport.PostId)
	if err != nil {
		return err
	}
	if postDoc.HgId == postReport.HgId {
		return errs.UserNotAllowed
	}

	return u.PostFirestoreRepo.ReportPost(ctx, postReport)
}

func (u *PostFeedsUsecase) UpdateLikePost(ctx context.Context, postId, hgId string) (string, error) {
	_, err := u.GetPostById(ctx, hgId, postId)
	if err != nil {
		return "", err
	}

	return u.PostFirestoreRepo.UpdateLikePost(ctx, postId, hgId)
}

// TODO: fix storage image thumbnail, raw data
func (u *PostFeedsUsecase) UploadImagePost(ctx context.Context, pictures []model.PictureDoc, imgQuery model.UploadImagePostQuery) (model.PictureDocuments, error) {
	if imgQuery.RoomId != roomIdPublic {
		if _, err := u.RestaurantFirestoreRepo.GetRestaurantById(ctx, imgQuery.RoomId); err != nil {
			return model.PictureDocuments{}, err
		}
		checkDoc, err := u.RestaurantFirestoreRepo.GetCheckIn(ctx, imgQuery.HgId, imgQuery.RoomId)
		if checkDoc.IsRevoked || (checkDoc == model.CheckInDoc{}) {
			return model.PictureDocuments{}, errs.UserNotAllowed
		}
		if err != nil {
			return model.PictureDocuments{}, err
		}
	}

	if imgQuery.PostId == "" && imgQuery.Type != "edit" {
		imgQuery.PostId = "Hango-" + uuid.Must(uuid.NewRandom()).String() + time.Now().Format("2006-01-02-15-04-05")
	}
	if len(pictures) > 6 {
		return model.PictureDocuments{}, errs.PictureExceed
	}

	resizeImage, err := resizeImage.ResizeImage(pictures)
	if err != nil {
		return model.PictureDocuments{}, err
	}

	urlThumbnail, err := u.UploadFileRepo.UploadFiletoStorage(ctx, resizeImage.ImageThumbnail, imgQuery.Collection, "thumbnail", imgQuery.PostId)
	if err != nil {
		return model.PictureDocuments{}, err
	}
	urlLarge, err := u.UploadFileRepo.UploadFiletoStorage(ctx, resizeImage.ImageLarge, imgQuery.Collection, "large", imgQuery.PostId)
	if err != nil {
		return model.PictureDocuments{}, err
	}

	return model.PictureDocuments{UrlThumbnail: urlThumbnail, UrlLarge: urlLarge, PostId: imgQuery.PostId}, nil
}
