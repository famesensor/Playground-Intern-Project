package protocol

import (
	model "github.com/HangoKub/Hango-service/internal/core/domain"
	"github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/HangoKub/Hango-service/internal/handlers"
	"github.com/HangoKub/Hango-service/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoute(app *fiber.App, authHandler ports.AuthHanler, restHandler ports.RestaurantHandler, postHandler ports.PostFeedsHandler, CommentHandler ports.CommentPostHandler, authware handlers.MiddlewareAuth) {
	// Docs swagger route
	// app.Get("/docs/*")

	// v1 api
	v1 := app.Group("/api/v1", logger.New())

	// Authication route
	auth := v1.Group("/authentication")
	auth.Post("/registration", middleware.ValidateBody(func() interface{} { return new(model.RegisterUser) }), authHandler.RegisterUser)
	auth.Post("/request-auth-access", middleware.ValidateBody(func() interface{} { return new(model.LoginUser) }), authHandler.ReqToken)
	auth.Put("/refresh-token", middleware.ValidateBody(func() interface{} { return new(model.RefreshToken) }), authHandler.RefreshToken)
	// auth.Patch("/verify-otp", middleware.ValidateBody(func() interface{} { return new(model.Otp) }), authHandler.Otp)         // OTP will return someone, somemonth and someyear
	// auth.Put("/refresh-otp", middleware.ValidateBody(func() interface{} { return new(model.Phone) }), authHandler.RefreshOtp) // OTP will return someone, somemonth and someyear
	auth.Patch("/logout", authware.Prodected(), middleware.ValidateBody(func() interface{} { return new(model.RefreshToken) }), authHandler.Logout)

	// Restaurant route
	rest := v1.Group("/restaurant", authware.Prodected())
	// rest.Post("/create-restaurant", middleware.ValidateBody(func() interface{} { return new(model.CreateRestaurant) }), restHandler.CreateRestaurant)
	rest.Get("/", middleware.ValidateQuery(func() interface{} { return new(model.RestaurantQuery) }), restHandler.GetRestaurantNearby) // ex. ?lat={float}&lng={float}&distance={int(km)}&tag={string}&open={true}
	rest.Post("/check-in", middleware.ValidateBody(func() interface{} { return new(model.UserCheckIn) }), restHandler.CheckInRestaurant)
	rest.Patch("/check-out", restHandler.CheckOutRestaurant)
	rest.Patch("/peakmode", restHandler.PeakMode)
	rest.Get("/checkin-anonymous-list", restHandler.GetCheckInAnonymous) // ex. ?restId={string}
	rest.Get("/checkin-peakmode-list", restHandler.GetCheckInPeakMode)   // ex. ?restId={string}
	rest.Get("/:id", restHandler.GetRestaurantById)

	// Post feed route
	feed := v1.Group("/feeds", authware.Prodected())
	post := feed.Group("/post")
	post.Get("/", middleware.ValidateQuery(func() interface{} { return new(model.FeedPostQuery) }), postHandler.GetAllHangoPost) // ex. ?roomId={string}&type={[]string}&lastPostId={string}&limit={int}&popular={true}
	post.Post("/", middleware.ValidateBody(func() interface{} { return new(model.InputPost) }), postHandler.CreateHangoPost)
	// post.Patch("/edit/:id", middleware.ValidateBody(func() interface{} { return new(model.InputPost) }), postHandler.EditHangoPost) // Edit post will return someone, somemonth and someyear
	post.Delete("/delete/:id", postHandler.DeleteHangoPost)
	post.Patch("/like-unlike/", postHandler.UpdateLikePost)                                                                                                                               // ex. ?postId={string}
	post.Post("/upload-picture-post", middleware.ValidateQuery(func() interface{} { return new(model.UploadImagePostQuery) }), middleware.ValidatePicture(), postHandler.UploadImagePost) // ex. ?collection={string}&postId={string}&type={string}
	post.Post("/report-post", middleware.ValidateBody(func() interface{} { return new(model.ReportPost) }), postHandler.ReportPost)
	post.Get("/:postId", postHandler.GetHangoPostById)

	// Comment post route
	comment := feed.Group("/comment")
	comment.Post("/", middleware.ValidateBody(func() interface{} { return new(model.InputComment) }), CommentHandler.CreateCommentPost)
	comment.Get("/", middleware.ValidateQuery(func() interface{} { return new(model.GetCommentQuery) }), CommentHandler.GetCommentInPost) // ex. ?postId={string}&type={string}&commentId={string}&lastCommentId={string}
	// comment.Patch("/edit", middleware.ValidateBody(func() interface{} { return new(model.UpdateComment) }), CommentHandler.EditCommentPost)      // Edit comment will return someone, somemonth and someyear
	comment.Delete("/delete", middleware.ValidateQuery(func() interface{} { return new(model.PostCommentQuery) }), CommentHandler.DeleteComment)                                                         // ex. ?postId={string}&commentId={string}&replyId={string}&type={string}
	comment.Post("/upload-picture-comment", middleware.ValidateQuery(func() interface{} { return new(model.UploadImageCommentQuery) }), middleware.ValidatePicture(), CommentHandler.UploadImageComment) // ex. ?collection={string}&postId={string}&commentId={string&type={string}
	comment.Patch("/like-unlike", middleware.ValidateQuery(func() interface{} { return new(model.PostCommentQuery) }), CommentHandler.UpdateLikeComment)                                                 // ex. ?postId={string}&commentId={string}&replyId={string}&type={string}
}
