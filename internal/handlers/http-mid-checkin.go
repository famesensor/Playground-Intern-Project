package handlers

import (
	"log"

	"github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type middlewareCheckIn struct {
	restUc ports.RestaurantUsecase
}

type MiddlewareCheckIn interface {
	CheckInProdected() fiber.Handler
	FeedCheckInProdected() fiber.Handler
}

func NewMiddlewareCheckIn(restUc ports.RestaurantUsecase) MiddlewareCheckIn {
	return &middlewareCheckIn{
		restUc,
	}
}

func (mid *middlewareCheckIn) CheckInProdected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get hgId from conetext fiber
		tk := c.Locals("user").(*jwt.Token)
		claims := tk.Claims.(jwt.MapClaims)
		hgId := claims["hgId"].(string)

		checkDoc, _ := mid.restUc.GetCheckInByHgId(c.Context(), hgId)
		c.Locals("checkIn", checkDoc)

		return c.Next()
	}
}

func (mid *middlewareCheckIn) FeedCheckInProdected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		postId := c.Locals("bodyData.postId").(string)

		log.Println(postId)
		return c.Next()
	}
}
