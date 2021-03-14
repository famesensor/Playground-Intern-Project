package handlers

import (
	model "github.com/HangoKub/Hango-service/internal/core/domain"
	interfaces "github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/HangoKub/Hango-service/pkg/reponseHandler"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// TODO: Handler notify for user comment, like

type CommentHandler struct {
	postUc    interfaces.PostFeedsUsecase
	commentUc interfaces.CommentPostUsecase
}

func NewCommentHandler(postUc interfaces.PostFeedsUsecase, commentUc interfaces.CommentPostUsecase) interfaces.CommentPostHandler {
	return &CommentHandler{
		postUc,
		commentUc,
	}
}

func (h *CommentHandler) CreateCommentPost(c *fiber.Ctx) error {
	// Get comment info from fiber context
	commentDoc := c.Locals("bodyData").(*model.InputComment)

	// Get hgId from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	postDoc, err := h.postUc.GetPostById(c.Context(), hgId, commentDoc.PostId)
	if err != nil {
		return reponseError(c, err)
	}

	commentInfo, err := h.commentUc.CreateComment(c.Context(), hgId, *commentDoc)
	if err != nil {
		return reponseError(c, err)
	}

	if commentInfo.HgIdNoti == "" {
		commentInfo.HgIdNoti = postDoc.HgId
	}

	// TODO: Do Notify to post/comment owner

	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", commentInfo.HangoComment)
}

func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	commentQuery := c.Locals("queryData").(*model.PostCommentQuery)

	// Get hgId from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	if _, err := h.postUc.GetPostById(c.Context(), hgId, commentQuery.PostId); err != nil {
		return reponseError(c, err)
	}

	if err := h.commentUc.DeleteComment(c.Context(), hgId, *commentQuery); err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "Comment delete is success", nil)
}

func (h *CommentHandler) UpdateLikeComment(c *fiber.Ctx) error {
	commentQuery := c.Locals("queryData").(*model.PostCommentQuery)

	// Get hgId from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	if _, err := h.postUc.GetPostById(c.Context(), hgId, commentQuery.PostId); err != nil {
		return reponseError(c, err)
	}

	res, err := h.commentUc.UpdateLikeComment(c.Context(), hgId, *commentQuery)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", &fiber.Map{"success": true, "type": res})

}

func (h *CommentHandler) GetCommentInPost(c *fiber.Ctx) error {
	commentQuery := c.Locals("queryData").(*model.GetCommentQuery)

	// Get hgId from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	if _, err := h.postUc.GetPostById(c.Context(), hgId, commentQuery.PostId); err != nil {
		return reponseError(c, err)
	}

	comments, err := h.commentUc.GetCommentInPost(c.Context(), hgId, *commentQuery)
	if err != nil {
		return reponseError(c, err)
	}

	if len(comments) == 0 {
		return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", &fiber.Map{"isMore": false})
	}
	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", comments)
}

func (h *CommentHandler) UploadImageComment(c *fiber.Ctx) error {
	commentQuery := c.Locals("queryData").(*model.UploadImageCommentQuery)
	picDocs := c.Locals("pictureData").([]model.PictureDoc)

	// Get hgId  from conetext fiber
	tk := c.Locals("user").(*jwt.Token)
	claims := tk.Claims.(jwt.MapClaims)
	hgId := claims["hgId"].(string)

	if _, err := h.postUc.GetPostById(c.Context(), hgId, commentQuery.PostId); err != nil {
		return reponseError(c, err)
	}

	res, err := h.commentUc.UploadImageComment(c.Context(), picDocs, *commentQuery)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", res)
}
