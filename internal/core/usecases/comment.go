package usecases

import (
	"context"
	"time"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
	interfaces "github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/HangoKub/Hango-service/pkg/resizeImage"
)

type CommentUsecase struct {
	CommentFirestoreRepo interfaces.CommentPostFirestoreRepository
	UploadFileRepo       interfaces.UploadStorageRepository
}

func NewCommentUsecase(CommentFirestoreRepo interfaces.CommentPostFirestoreRepository, UploadFileRepo interfaces.UploadStorageRepository) *CommentUsecase {
	return &CommentUsecase{
		CommentFirestoreRepo,
		UploadFileRepo,
	}
}

func (u *CommentUsecase) CreateComment(ctx context.Context, hgId string, commentInput model.InputComment) (model.CommentNotiDocuments, error) {
	// Check type comment
	if commentInput.CommentId != "" && commentInput.Type == "root" {
		return model.CommentNotiDocuments{}, errs.OperationIsnotAllowed
	}

	isNest := false
	var owner string
	if commentInput.Type == "reply" {
		isNest = true
		// Check comment root existing
		commentDoc, err := u.CommentFirestoreRepo.GetCommentById(ctx, commentInput.PostId, commentInput.CommentId, "")
		if err != nil {
			return model.CommentNotiDocuments{}, err
		}
		owner = commentDoc.HgId
	}

	// Check user existing in this post
	if err := u.CommentFirestoreRepo.CheckInterActor(ctx, commentInput.PostId, hgId); err != nil {
		return model.CommentNotiDocuments{}, err
	}

	commentDoc := model.HangoComment{
		Comment:          commentInput.Comment,
		PictureThumbnail: commentInput.PictureThumbnail,
		PictureLarge:     commentInput.PictureLarge,
		IsNest:           isNest,
		Hide:             false,
		HgId:             hgId,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	result, err := u.CommentFirestoreRepo.CreateComment(ctx, commentInput.PostId, commentInput.CommentId, commentDoc)
	if err != nil {
		// TODO: Delete inter-actor if create comment errror
		return model.CommentNotiDocuments{}, err
	}

	return model.CommentNotiDocuments{HangoComment: result, HgIdNoti: owner}, nil
}

// func (u *CommentUsecase) UpdateComment(ctx context.Context, hgId string, commentUpdate model.UpdateComment) (model.HangoComment, error) {
// 	if commentUpdate.Type != "reply" && commentUpdate.ReplyId != "" {
// 		return model.HangoComment{}, errs.OperationIsnotAllowed
// 	}
// 	doc, err := u.CommentFirestoreRepo.GetCommentById(ctx, commentUpdate.PostId, commentUpdate.CommentId, commentUpdate.ReplyId)
// 	if err != nil {
// 		return model.HangoComment{}, err
// 	}
// 	if doc.HgId != hgId {
// 		return model.HangoComment{}, errs.OperationIsnotAllowed
// 	}

// 	commentDoc, err := u.CommentFirestoreRepo.EditComment(ctx, commentUpdate)
// 	if err != nil {
// 		return model.HangoComment{}, err
// 	}

// 	return commentDoc, nil
// }

func (u *CommentUsecase) DeleteComment(ctx context.Context, hgId string, commentQuery model.PostCommentQuery) error {
	if commentQuery.ReplyId != "" && commentQuery.Type == "root" {
		return errs.OperationIsnotAllowed
	}

	doc, err := u.CommentFirestoreRepo.GetCommentById(ctx, commentQuery.PostId, commentQuery.CommentId, commentQuery.ReplyId)
	if err != nil {
		return err
	}
	if doc.HgId != hgId {
		return errs.OperationIsnotAllowed
	}

	return u.CommentFirestoreRepo.DeleteComment(ctx, hgId, commentQuery)
}

func (u *CommentUsecase) GetCommentInPost(ctx context.Context, hgId string, commentQuery model.GetCommentQuery) (comments []model.HangoComment, err error) {
	if commentQuery.CommentId != "" && commentQuery.Type == "root" {
		return comments, errs.OperationIsnotAllowed
	}

	if commentQuery.Limit == 0 {
		commentQuery.Limit = 10
	}

	switch commentQuery.Type {
	case "root":
		comments, err = u.CommentFirestoreRepo.GetRootComment(ctx, commentQuery)
		break
	case "reply":
		comments, err = u.CommentFirestoreRepo.GetReplyComment(ctx, commentQuery)
		break
	}
	if err != nil {
		return []model.HangoComment{}, err
	}

	return comments, err
}

func (u *CommentUsecase) UpdateLikeComment(ctx context.Context, hgId string, commentQuery model.PostCommentQuery) (string, error) {
	return u.CommentFirestoreRepo.UpdateLikeComment(ctx, hgId, commentQuery)
}

func (u *CommentUsecase) UploadImageComment(ctx context.Context, files []model.PictureDoc, uploadImageComment model.UploadImageCommentQuery) (model.PictureDocuments, error) {
	var err error
	if uploadImageComment.Type == "reply" {
		_, err = u.CommentFirestoreRepo.GetCommentById(ctx, uploadImageComment.PostId, uploadImageComment.CommentId, "")
	}
	if err != nil {
		return model.PictureDocuments{}, err
	}

	resizeImage, err := resizeImage.ResizeImage(files)
	if err != nil {
		return model.PictureDocuments{}, err
	}
	urlThumbnail, err := u.UploadFileRepo.UploadFiletoStorage(ctx, resizeImage.ImageThumbnail, uploadImageComment.Collection, "thumbnail", uploadImageComment.PostId)
	if err != nil {
		return model.PictureDocuments{}, err
	}
	urlLarge, err := u.UploadFileRepo.UploadFiletoStorage(ctx, resizeImage.ImageLarge, uploadImageComment.Collection, "large", uploadImageComment.PostId)
	if err != nil {
		return model.PictureDocuments{}, err
	}

	return model.PictureDocuments{UrlThumbnail: urlThumbnail, UrlLarge: urlLarge}, nil
}
