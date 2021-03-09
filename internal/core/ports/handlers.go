package ports

import (
	"github.com/gofiber/fiber/v2"
)

type AuthHanler interface {
	RegisterUser(c *fiber.Ctx) error
	ReqToken(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	Otp(c *fiber.Ctx) error
	RefreshOtp(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type UserHandler interface {
	// CheckInRestaurant(c *fiber.Ctx) error
	// CheckOutRestaurant(c *fiber.Ctx) error
}

type ProfileHandler interface {
	// GetProfileByHgId(c *fiber.Ctx) error
	// UpdateProfile(c *fiber.Ctx) error
	// PrivacyManagement(c *fiber.Ctx) error

}

type RestaurantHandler interface {
	// CreateRestaurant(c *fiber.Ctx) error
	GetRestaurantById(c *fiber.Ctx) error
	GetRestaurantNearby(c *fiber.Ctx) error
	// GetRestaurantByName(c *fiber.Ctx) error

	CheckInRestaurant(c *fiber.Ctx) error
	CheckOutRestaurant(c *fiber.Ctx) error
	PeakMode(c *fiber.Ctx) error
	GetCheckInAnonymous(c *fiber.Ctx) error
	GetCheckInPeakMode(c *fiber.Ctx) error
}

type PostFeedsHandler interface {
	CreateHangoPost(c *fiber.Ctx) error
	// EditHangoPost(c *fiber.Ctx) error
	DeleteHangoPost(c *fiber.Ctx) error
	GetAllHangoPost(c *fiber.Ctx) error
	GetHangoPostById(c *fiber.Ctx) error
	ReportPost(c *fiber.Ctx) error
	UpdateLikePost(c *fiber.Ctx) error

	UploadImagePost(c *fiber.Ctx) error
}

type CommentPostHandler interface {
	CreateCommentPost(c *fiber.Ctx) error
	// EditCommentPost(c *fiber.Ctx) error
	DeleteComment(c *fiber.Ctx) error
	UpdateLikeComment(c *fiber.Ctx) error
	GetCommentInPost(c *fiber.Ctx) error

	UploadImageComment(c *fiber.Ctx) error
}
