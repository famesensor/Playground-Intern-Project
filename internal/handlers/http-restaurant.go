package handlers

import (
	model "github.com/HangoKub/Hango-service/internal/core/domain"
	interfaces "github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/HangoKub/Hango-service/pkg/reponseHandler"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type RestaurantHandler struct {
	restUc interfaces.RestaurantUsecase
}

func NewRestaurantHandler(restUc interfaces.RestaurantUsecase) interfaces.RestaurantHandler {
	return &RestaurantHandler{
		restUc,
	}
}

// // TODO: input keyword for search restaurant by name
// func (h *RestaurantHandler) CreateRestaurant(c *fiber.Ctx) error {
// 	rest := new(model.CreateRestaurant)

// 	if err := c.BodyParser(rest); err != nil {
// 		return reponseError(c, errs.CannotParseData)
// 	}

// 	if err := h.restUc.CreateRestaurant(c.Context(), *rest); err != nil {
// 		return reponseError(c, err)
// 	}
// 	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", nil)
// }

func (h *RestaurantHandler) GetRestaurantById(c *fiber.Ctx) error {
	restId := c.Params("id")

	if restId == "" {
		return reponseHandler.ReponseMsg(c, fiber.StatusBadRequest, "failed", "Validation Errors", &fiber.Map{"RestId": "RestaurantId is required"})
	}

	restDoc, err := h.restUc.GetRestaurantById(c.Context(), restId)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", restDoc)
}

func (h *RestaurantHandler) GetRestaurantNearby(c *fiber.Ctx) error {
	q := c.Locals("queryData").(*model.RestaurantQuery)

	restDoc, err := h.restUc.GetRestaurantNearby(c.Context(), *q)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", restDoc)
}

// func (h *RestaurantHandler) GetRestaurantByName(c *fiber.Ctx) error {
// 	restName := c.Query("restName")
// 	lastRestId := c.Query("lastRestId")

// 	restDocs, err := h.restUc.GetRestaurantByName(c.Context(), restName, lastRestId)
// 	if err != nil {
// 		return reponseError(c, err)
// 	}

// 	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", restDocs)
// }

func (h *RestaurantHandler) CheckInRestaurant(c *fiber.Ctx) error {
	// Get query from context fiber
	query := c.Locals("bodyData").(*model.UserCheckIn)

	// Get hgId from context fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	// Create ckeckin doc
	checkIn, err := h.restUc.CheckInRestaurant(c.Context(), hgId, *query)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", &fiber.Map{"type": checkIn})
}

// TODO: Handle check distance between user and restaurant
func (h *RestaurantHandler) PeakMode(c *fiber.Ctx) error {
	// Get hgId from context fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	checkDoc, err := h.restUc.GetCheckInByHgId(c.Context(), hgId)
	if err != nil {
		return reponseError(c, err)
	}
	if _, err := h.restUc.GetRestaurantById(c.Context(), checkDoc.RestId); err != nil {
		return reponseError(c, err)
	}

	peak, err := h.restUc.PeakMode(c.Context(), hgId)
	if err != nil {
		return reponseError(c, err)
	}
	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", &fiber.Map{"type": peak})
}

func (h *RestaurantHandler) CheckOutRestaurant(c *fiber.Ctx) error {
	// Get hgId from context fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	checkDoc, err := h.restUc.GetCheckInByHgId(c.Context(), hgId)
	if err != nil {
		return reponseError(c, err)
	}
	if _, err := h.restUc.GetRestaurantById(c.Context(), checkDoc.RestId); err != nil {
		return reponseError(c, err)
	}

	checkOut, err := h.restUc.CheckOutRestaurant(c.Context(), hgId, checkDoc.RestId)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", &fiber.Map{"type": checkOut})
}

func (h *RestaurantHandler) GetCheckInAnonymous(c *fiber.Ctx) error {
	restId := c.Query("restId")

	if restId == "" {
		return reponseHandler.ReponseMsg(c, fiber.StatusBadRequest, "failed", "Validation Errors", &fiber.Map{"restId": "RestaurantId is required"})
	}

	// Get hgId from context fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	if _, err := h.restUc.GetRestaurantById(c.Context(), restId); err != nil {
		return reponseError(c, err)
	}
	if _, err := h.restUc.GetCheckIn(c.Context(), hgId, restId); err != nil {
		return reponseError(c, err)
	}

	res, err := h.restUc.GetCheckInAnonymous(c.Context(), hgId, restId)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", res)
}

func (h *RestaurantHandler) GetCheckInPeakMode(c *fiber.Ctx) error {
	restId := c.Query("restId")

	if restId == "" {
		return reponseHandler.ReponseMsg(c, fiber.StatusBadRequest, "failed", "Validation Errors", &fiber.Map{"restId": "RestaurantId is required"})
	}

	// Get hgId from context fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	if _, err := h.restUc.GetRestaurantById(c.Context(), restId); err != nil {
		return reponseError(c, err)
	}
	checkDoc, err := h.restUc.GetCheckInByHgId(c.Context(), hgId)
	if !checkDoc.Peak {
		return reponseError(c, errs.UserNotAllowed)
	}
	if err != nil {
		return reponseError(c, err)
	}

	res, err := h.restUc.GetCheckInPeakMode(c.Context(), hgId, restId)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", res)
}
