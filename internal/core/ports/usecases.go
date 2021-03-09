package ports

import (
	"context"
	"time"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
)

type AuthenUsecase interface {
	// ReqTokenDocument -> is the (fn) will be create authen document and generate token
	ReqTokenDocument(ctx context.Context, acT time.Duration, rfT time.Duration, tokenType string, claims *model.AuthClaim) (model.TokenVerificationDocuments, error)
	// CreateAuthDoc -> is the (fn) will be create authen document into repository
	CreateAuthDoc(ctx context.Context, authDoc model.AuthDocument) error
	// VerifyToken -> is the (fn) for verify stauts token
	VerifyToken(token string) (string, error)
	// RefreshToken -> is the (fn) will be check refresh token and create access toekn
	RefreshToken(ctx context.Context, refreshTk string) (model.TokenCard, error)
	// ReqOtpDocument -> is the (fn) will be create otp document and send otp to user
	ReqOtpDocument(ctx context.Context, phone string, otpT time.Duration, typeOTP string) error
	// CreateOtpDoc -> is the (fn) will be create otp document into reposity
	RefreshOtp(ctx context.Context, phone, typeOTP string) error
	// VerifyOtp -> is a (fn) that uses for verify otp user
	VerifyOtp(ctx context.Context, opt, phone, option, typeOTP string) error
	// Logout
	Logout(ctx context.Context, refreshTk string) error
}

type UserUsecase interface {
	// CreateUser -> is the (fn) will be create user into repository
	CreateUser(ctx context.Context, user model.RegisterUser) (id string, err error)
	// GetUserById -> get user information by ID
	GetUserById(ctx context.Context, id string) (model.User, error)
	// GetUserByEmail -> get user information by Email
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type ProfileUsecase interface {
	// GetProfileByHgId -> get profile information by hgId
	GetProfileByHgId(ctx context.Context, hgId string) (model.Profile, error)
}

type RestaurantUsecase interface {
	// CreateRestaurant(ctx context.Context, restDoc model.CreateRestaurant) error
	GetRestaurantById(ctx context.Context, restId string) (model.ResponseRestaurant, error)
	GetRestaurantNearby(ctx context.Context, restQuery model.RestaurantQuery) ([]model.ReponseHomeRestaurant, error)
	// GetRestaurantByName(ctx context.Context, restName, lastRestId string) ([]model.ReponseHomeRestaurant, error)

	CheckInRestaurant(ctx context.Context, hgId string, checkInDoc model.UserCheckIn) (string, error)
	GetCheckIn(ctx context.Context, hgId, restId string) (model.CheckInDoc, error)
	GetCheckInByHgId(ctx context.Context, hgId string) (model.CheckInDoc, error)
	CheckOutRestaurant(ctx context.Context, hgId, restId string) (string, error)
	PeakMode(ctx context.Context, hgId string) (string, error)
	GetCheckInAnonymous(ctx context.Context, restId, hgId string) (anonymousDocs []model.AnonymousDoc, err error)
	GetCheckInPeakMode(ctx context.Context, hgId, restId string) (peakModeDocs []model.PeakModeDoc, err error)
}

type PostFeedsUsecase interface {
	CreateHangoPost(ctx context.Context, postDoc model.InputPost, userId string) (model.HangoPost, error)
	// EditHangoPost(ctx context.Context, postId, userId string, postDoc model.InputPost) (model.HangoPost, error)
	DeleteHangoPost(ctx context.Context, postId, userId string) error
	GetPostById(ctx context.Context, hgId, postId string) (model.HangoPost, error)
	GetAllHangoPost(ctx context.Context, postQuery model.FeedPostQuery) ([]model.PostFeedReponse, error)
	ReportPost(ctx context.Context, postReport model.ReportPost) error
	UpdateLikePost(ctx context.Context, postId, userId string) (string, error)

	UploadImagePost(ctx context.Context, pictures []model.PictureDoc, imgQuery model.UploadImagePostQuery) (model.PictureDocuments, error)
}

type CommentPostUsecase interface {
	CreateComment(ctx context.Context, hgId string, commentInput model.InputComment) (model.CommentNotiDocuments, error)
	// UpdsateComment(ctx context.Context, hgId string, commentUpdate model.UpdateComment) (model.HangoComment, error)
	DeleteComment(ctx context.Context, hgId string, commentQuery model.PostCommentQuery) error
	GetCommentInPost(ctx context.Context, hgId string, commentQuery model.GetCommentQuery) ([]model.HangoComment, error)
	UpdateLikeComment(ctx context.Context, hgId string, commentQuery model.PostCommentQuery) (string, error)

	UploadImageComment(ctx context.Context, files []model.PictureDoc, uploadImageComment model.UploadImageCommentQuery) (model.PictureDocuments, error)
}
