package handlers

import (
	model "github.com/HangoKub/Hango-service/internal/core/domain"
	interfaces "github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/HangoKub/Hango-service/pkg/reponseHandler"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type PostFeedsHandler struct {
	postUc interfaces.PostFeedsUsecase
}

func NewPostFeedHandler(postUc interfaces.PostFeedsUsecase) interfaces.PostFeedsHandler {
	return &PostFeedsHandler{
		postUc,
	}
}

// TODO: Handler notify for user like post

func (h *PostFeedsHandler) CreateHangoPost(c *fiber.Ctx) error {
	// Get data from context fiber and Convert type
	postDoc, _ := c.Locals("bodyData").(*model.InputPost)

	// Get hgId from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	res, err := h.postUc.CreateHangoPost(c.Context(), *postDoc, hgId)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", res)
}

// func (h *PostFeedsHandler) EditHangoPost(c *fiber.Ctx) error {
// 	// Get postId form params
// 	postId := c.Params("id")
// 	if postId == "" {
// 		return reponseHandler.ReponseMsg(c, fiber.StatusBadRequest, "failed", "Validation Errors", &fiber.Map{"PostId": "PostId is required"})
// 	}

// 	// Get data from context fiber and Convert type
// 	postDoc, _ := c.Locals("bodyData").(*model.InputPost)

// 	// Get hgId from conetext fiber
// 	tk := c.Locals("user").(*jwt.Token)
// 	claims := tk.Claims.(jwt.MapClaims)
// 	hgId := claims["hgId"].(string)

// 	res, err := h.postUc.EditHangoPost(c.Context(), postId, hgId, *postDoc)
// 	if err != nil {
// 		return reponseError(c, err)
// 	}

// 	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", res)
// }

func (h *PostFeedsHandler) DeleteHangoPost(c *fiber.Ctx) error {
	// Get postId form params
	postId := c.Params("id")
	if postId == "" {
		return reponseHandler.ReponseMsg(c, fiber.StatusBadRequest, "failed", "Validation Errors", &fiber.Map{"PostId": "PostId is required"})
	}

	// Get hgId  from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	if err := h.postUc.DeleteHangoPost(c.Context(), postId, hgId); err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "Post ID : "+postId+" Delete is success", nil)
}

func (h *PostFeedsHandler) GetHangoPostById(c *fiber.Ctx) error {
	postId := c.Params("postId")

	if postId == "" {
		return reponseHandler.ReponseMsg(c, fiber.StatusBadRequest, "failed", "Validation Errors", &fiber.Map{"PostId": "PostId is required"})
	}

	// Get hgId  from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	res, err := h.postUc.GetPostById(c.Context(), hgId, postId)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", res)
}

func (h *PostFeedsHandler) GetAllHangoPost(c *fiber.Ctx) error {
	q := c.Locals("queryData").(*model.FeedPostQuery)

	postDocs, err := h.postUc.GetAllHangoPost(c.Context(), *q)
	if err != nil {
		return reponseError(c, err)
	}

	if len(postDocs) == 0 {
		return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", &fiber.Map{"isMore": false})
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", postDocs)
}

func (h *PostFeedsHandler) ReportPost(c *fiber.Ctx) error {
	reportDoc := c.Locals("bodyData").(*model.ReportPost)

	// Get hgId  from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)
	reportDoc.HgId = hgId

	if err := h.postUc.ReportPost(c.Context(), *reportDoc); err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", &fiber.Map{"type": "report"})
}

func (h *PostFeedsHandler) UploadImagePost(c *fiber.Ctx) error {
	q := c.Locals("queryData").(*model.UploadImagePostQuery)
	picDocs := c.Locals("pictureData").([]model.PictureDoc)

	// Get hgId  from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)
	q.HgId = hgId

	picUrls, err := h.postUc.UploadImagePost(c.Context(), picDocs, *q)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", picUrls)
}

func (h *PostFeedsHandler) UpdateLikePost(c *fiber.Ctx) error {
	postId := c.Query("postId")

	if postId == "" {
		return reponseHandler.ReponseMsg(c, fiber.StatusBadRequest, "failed", "Validation Errors", &fiber.Map{"PostId": "PostId is required"})
	}

	// Get hgId from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	res, err := h.postUc.UpdateLikePost(c.Context(), postId, hgId)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", &fiber.Map{"success": true, "type": res})
}
