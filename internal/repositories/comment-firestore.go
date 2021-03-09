package repositories

import (
	"context"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	model "github.com/HangoKub/Hango-service/internal/core/domain"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	commentCollection     string = "hango-comment"
	nestCommentCollection string = "hango-nest-comment"
)

type CommentFirestoreRepository struct {
	firestore *firestore.Client
}

func NewCommentPostFirestore(firestore *firestore.Client) *CommentFirestoreRepository {
	return &CommentFirestoreRepository{
		firestore,
	}
}

func (r *CommentFirestoreRepository) CreateComment(ctx context.Context, postId, commentId string, commentDoc model.HangoComment) (model.HangoComment, error) {
	batch := r.firestore.Batch()
	postRef := r.firestore.Doc(postCollection + "/" + postId)
	ref := postRef.Collection(commentCollection).NewDoc()
	var rootRef *firestore.DocumentRef

	// check if root comment / nest comment
	if !commentDoc.IsNest {
		// update count root comment
		batch.Update(postRef, []firestore.Update{
			{Path: "commentCount", Value: firestore.Increment(1)},
		})
	} else {
		rootRef = r.firestore.Doc(postCollection + "/" + postId).Collection(commentCollection).Doc(commentId)
		ref = rootRef.Collection(nestCommentCollection).NewDoc()
		// update count nest comment
		batch.Update(rootRef, []firestore.Update{
			{Path: "nestCount", Value: firestore.Increment(1)},
		})
	}

	// Create comment or nest comment
	commentDoc.CommentId = ref.ID
	batch.Set(ref, commentDoc)

	if _, err := batch.Commit(ctx); err != nil {
		// delete referrence if batch commit error
		batch.Delete(ref)
		batch.Commit(ctx)
		return model.HangoComment{}, err
	}

	return commentDoc, nil
}

func (r *CommentFirestoreRepository) EditComment(ctx context.Context, commentDoc model.UpdateComment) (model.HangoComment, error) {
	ref := r.firestore.Doc(postCollection + "/" + commentDoc.PostId).Collection(commentCollection).Doc(commentDoc.CommentId)

	// Set ref nested comment
	if commentDoc.ReplyId != "" && commentDoc.Type == "reply" {
		ref = ref.Collection(nestCommentCollection).Doc(commentDoc.ReplyId)
	}

	err := r.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := tx.Get(ref)
		if err != nil {
			return err
		}
		return tx.Set(ref, map[string]interface{}{
			"comment":   commentDoc.Comment,
			"picture":   commentDoc.Picture,
			"updatedAt": time.Now(),
		}, firestore.MergeAll)
	})

	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return model.HangoComment{}, errs.DocumentNotFound
		}
		return model.HangoComment{}, err
	}

	return r.GetCommentById(ctx, commentDoc.PostId, commentDoc.CommentId, commentDoc.ReplyId)
}

// TODO: Handler delete nested comment when user delete comment root
func (r *CommentFirestoreRepository) DeleteComment(ctx context.Context, hgId string, commentQuery model.PostCommentQuery) error {
	postRef := r.firestore.Doc(postCollection + "/" + commentQuery.PostId)
	ref := postRef.Collection(commentCollection).Doc(commentQuery.CommentId)
	var rootRef *firestore.DocumentRef
	// Set ref nested comment
	if commentQuery.Type == "reply" && commentQuery.ReplyId != "" {
		rootRef = postRef.Collection(commentCollection).Doc(commentQuery.CommentId)
		ref = rootRef.Collection(nestCommentCollection).Doc(commentQuery.ReplyId)
	}

	err := r.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		if _, err := tx.Get(ref); err != nil {
			return err
		}
		if err := tx.Delete(ref); err != nil {
			return err
		}
		switch commentQuery.Type {
		case "root":
			if err := tx.Update(postRef, []firestore.Update{{Path: "commentCount", Value: firestore.Increment(-1)}}); err != nil {
				return err
			}
			break
		case "reply":
			if err := tx.Update(rootRef, []firestore.Update{{Path: "nestCount", Value: firestore.Increment(-1)}}); err != nil {
				return err
			}
			break
		}

		return nil
	})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return errs.DocumentNotFound
		}
		return err
	}
	return nil
}

func (r *CommentFirestoreRepository) GetCommentById(ctx context.Context, postId, commentId, replyId string) (model.HangoComment, error) {
	postRef := r.firestore.Doc(postCollection + "/" + postId)
	ref := postRef.Collection(commentCollection).Doc(commentId)

	// Set ref nested comment
	if replyId != "" {
		ref = ref.Collection(nestCommentCollection).Doc(replyId)
	}

	dsnap, err := ref.Get(ctx)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return model.HangoComment{}, errs.DocumentNotFound
		}
		return model.HangoComment{}, err
	}
	comment := new(model.HangoComment)
	dsnap.DataTo(comment)
	comment.CommentId = dsnap.Ref.ID

	return *comment, nil
}

func (r *CommentFirestoreRepository) GetRootComment(ctx context.Context, commentQuery model.GetCommentQuery) (comments []model.HangoComment, err error) {
	ref := r.firestore.Doc(postCollection + "/" + commentQuery.PostId).Collection(commentCollection)
	commetRef := ref.OrderBy("createdAt", firestore.Desc).Limit(commentQuery.Limit)

	// Pagination
	if commentQuery.LastCommentId != "" && commentQuery.Type == "root" {
		dsnap, err := ref.Doc(commentQuery.LastCommentId).Get(ctx)
		if err != nil {
			if grpc.Code(err) == codes.NotFound {
				return []model.HangoComment{}, errs.DocumentNotFound
			}
			return []model.HangoComment{}, err
		}
		commetRef = commetRef.StartAfter(dsnap.Data()["createdAt"])
	}
	iter := commetRef.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []model.HangoComment{}, errs.DocumentNotFound
		}
		if doc.Exists() {
			comment := new(model.HangoComment)
			mapstructure.Decode(doc.Data(), comment)
			comments = append(comments, *comment)
		}
	}
	sort.Slice(comments, func(i, j int) bool { return comments[i].CreatedAt.Before(comments[j].CreatedAt) })

	return comments, nil
}

func (r *CommentFirestoreRepository) GetReplyComment(ctx context.Context, commentQuery model.GetCommentQuery) (comments []model.HangoComment, err error) {
	ref := r.firestore.Doc(postCollection + "/" + commentQuery.PostId).Collection(commentCollection).Doc(commentQuery.CommentId).Collection(nestCommentCollection)
	commetRef := ref.OrderBy("createdAt", firestore.Desc).Limit(commentQuery.Limit)

	// Pagination
	if commentQuery.LastCommentId != "" && commentQuery.Type == "reply" {
		dsnap, err := ref.Doc(commentQuery.LastCommentId).Get(ctx)
		if err != nil {
			if grpc.Code(err) == codes.NotFound {
				return []model.HangoComment{}, errs.DocumentNotFound
			}
			return []model.HangoComment{}, err
		}
		commetRef = commetRef.StartAfter(dsnap.Data()["createdAt"])
	}
	iter := commetRef.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []model.HangoComment{}, errs.DocumentNotFound
		}
		if doc.Exists() {
			comment := new(model.HangoComment)
			mapstructure.Decode(doc.Data(), comment)
			comments = append(comments, *comment)
		}
	}

	return comments, nil
}

func (r *CommentFirestoreRepository) UpdateLikeComment(ctx context.Context, hgId string, commentQuery model.PostCommentQuery) (string, error) {
	postRef := r.firestore.Doc(postCollection + "/" + commentQuery.PostId)
	var query firestore.Query
	var result string
	ref := postRef.Collection(commentCollection).Doc(commentQuery.CommentId)

	if commentQuery.Type == "root" && commentQuery.ReplyId == "" {
		query = postRef.Collection(commentCollection).Where("commentId", "==", commentQuery.CommentId).Limit(1).Where("like", "array-contains", hgId)
	} else {
		query = postRef.Collection(commentCollection).Doc(commentQuery.CommentId).Collection(nestCommentCollection).Where("commentId", "==", commentQuery.ReplyId).Limit(1).Where("like", "array-contains", hgId)
		ref = ref.Collection(nestCommentCollection).Doc(commentQuery.ReplyId)
	}
	err := r.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		iter := tx.Documents(query)
		defer iter.Stop()

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				result = "like"
				return tx.Update(ref, []firestore.Update{
					{Path: "like", Value: firestore.ArrayUnion(hgId)},
					{Path: "likeCount", Value: firestore.Increment(1)},
				})
			}
			if err != nil {
				return err
			}
			if doc.Exists() {
				result = "unlike"
				return tx.Update(ref, []firestore.Update{
					{Path: "like", Value: firestore.ArrayRemove(hgId)},
					{Path: "likeCount", Value: firestore.Increment(-1)},
				})
			}
		}
	})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return "", errs.DocumentNotFound
		}
		return "", err
	}

	return result, nil
}

// TODO: Do something with anonymous random
func (r *CommentFirestoreRepository) CheckInterActor(ctx context.Context, postId, hgId string) error {
	doc, err := r.firestore.Doc(postCollection+"/"+postId).Collection(interActorCollection).Where("hgId", "==", hgId).Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	if len(doc) == 0 {
		doc, _ := r.firestore.Doc(postCollection + "/" + postId).Collection(interActorCollection).Documents(ctx).GetAll()
		docLength := len(doc)
		number := rand.Intn(64)
		animal := model.AnonymousList[number].(string)
		aka := animal + strconv.Itoa(docLength)

		_, err := r.firestore.Collection(postCollection).Doc(postId).Collection(interActorCollection).NewDoc().Set(ctx, model.InterActor{
			HgId:    hgId,
			Aka:     aka,
			Picture: "https://firebasestorage.googleapis.com/v0/b/hango-dev-32d20.appspot.com/o/utils%2Fanonymous%2F" + strconv.Itoa(number) + ".svg?alt=media",
			UId:     docLength,
		})
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
