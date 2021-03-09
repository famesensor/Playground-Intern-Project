package ports

import (
	"bytes"
	"context"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
)

type AuthFirestoreRepository interface {
	// CreateAuthenDoc -> is the (fn) that will be create the authen document into database
	CreateAuthenDoc(ctx context.Context, authDoc model.AuthDocument) error
	// GetAuthenDocByRefreshToken -> is the (fn) get authen document by refresh token
	GetAuthenDocByRefreshToken(ctx context.Context, refreshTk string) (model.AuthDocument, error)
	// RevokedToken -> is the (fn) that will be revoke token when user logout
	RevokedToken(ctx context.Context, refreshTk string) error
	// CreateOtp -> is the (fn) that will be create the otp document into database
	CreateOtp(ctx context.Context, otpDoc model.OtpDoc) error
	// GetOtpDoc -> is the (fn) get otp document by phone
	GetOtpDoc(ctx context.Context, phone string) (model.OtpDoc, error)
	// DeleteOtp -> is the (fn) that will be delete(Revoke) Otp when user verify otp
	DeleteOtp(ctx context.Context, phone string) error
}

type UserFirestoreRepository interface {
	// CreateUser -> is the (fn) that will be create the user and profile account into database
	CreateUser(ctx context.Context, user model.RegisterUser) (id string, err error)
	// GetUserById -> is the (fn) get user information by ID
	GetUserById(ctx context.Context, id string) (model.User, error)
	// GetUserByEmail -> is the (fn) get user information by Email
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type ProfileFirestoreRepository interface {
	// GetProfileByHgId -> is the (fn) get profile information by hgId
	GetProfileByHgId(ctx context.Context, hgId string) (model.Profile, error)
}

type RestaurantFirestoreRepository interface {
	// CreateRestaurant(ctx context.Context, restDoc model.CreateRestaurant) error
	GetRestaurantById(ctx context.Context, restId string) (model.ResponseRestaurant, error)
	GetRestaurantNearby(ctx context.Context, restaurantQuery model.ListRestaurantQuery) ([]model.ReponseHomeRestaurant, error)
	// GetRestaurantByName(ctx context.Context, restName, lastRestId string) ([]model.ReponseHomeRestaurant, error)

	CreateCheckIn(ctx context.Context, hgId, restId string) error
	GetCheckIn(ctx context.Context, hgId, restId string) (model.CheckInDoc, error)
	GetCheckInByHgId(ctx context.Context, hgId string) (model.CheckInDoc, error)
	CheckOut(ctx context.Context, hgId, restId string) error
	PeakMode(ctx context.Context, hgId string) (string, error)
	GetCheckInAnonymous(ctx context.Context, restId string) (anonymousDocs []model.AnonymousDoc, err error)
	GetCheckInPeakMode(ctx context.Context, restId string) (peakModeDocs []model.PeakModeDoc, err error)
}
type PostFeedsFirestoreRepository interface {
	CreateHangoPost(ctx context.Context, postDoc model.InputPost, interActorDoc model.InterActor) (model.HangoPost, error)
	// EditHangoPost(ctx context.Context, postId string, postDoc model.InputPost) (model.HangoPost, error)
	DeleteHangoPost(ctx context.Context, postId string) error
	GetHangoPostByID(ctx context.Context, postId string) (model.HangoPost, error)
	GetAllHangoPost(ctx context.Context, postQuery model.FeedPostQuery) ([]model.PostFeedReponse, error)
	ReportPost(ctx context.Context, postReport model.ReportPost) error
	UpdateLikePost(ctx context.Context, postId, hgId string) (string, error)
}

type CommentPostFirestoreRepository interface {
	CreateComment(ctx context.Context, postId, commentId string, commentDoc model.HangoComment) (model.HangoComment, error)
	GetCommentById(ctx context.Context, postId, commentId, replyId string) (model.HangoComment, error)
	// EditComment(ctx context.Context, commentDoc model.UpdateComment) (model.HangoComment, error)
	DeleteComment(ctx context.Context, hgId string, commentQuery model.PostCommentQuery) error
	GetRootComment(ctx context.Context, commentQuery model.GetCommentQuery) ([]model.HangoComment, error)
	GetReplyComment(ctx context.Context, commentQuery model.GetCommentQuery) ([]model.HangoComment, error)
	UpdateLikeComment(ctx context.Context, hgId string, commentQuery model.PostCommentQuery) (string, error)

	// Check user interac with post or comment
	CheckInterActor(ctx context.Context, postId, hgId string) error
}

type UploadStorageRepository interface {
	UploadFiletoStorage(ctx context.Context, files []bytes.Buffer, collection, typeFile, id string) ([]string, error)
}
