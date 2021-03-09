package repositories

import (
	"context"
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
	postCollection       string = "hango-feeds"
	interActorCollection string = "inter-actor"
	reportPostCollection string = "report-post"
)

type PostFeedsFirestoreRepository struct {
	firestore *firestore.Client
}

func NewPostFeedFirestore(firestore *firestore.Client) *PostFeedsFirestoreRepository {
	return &PostFeedsFirestoreRepository{
		firestore,
	}
}

func (r *PostFeedsFirestoreRepository) CreateHangoPost(ctx context.Context, postDoc model.InputPost, interActorDoc model.InterActor) (model.HangoPost, error) {
	hangoPost := model.HangoPost{
		PostId:           postDoc.PostId,
		RoomId:           postDoc.RoomId,
		Content:          postDoc.Content,
		PictureThumbnail: postDoc.PictureThumbnail,
		PictureLarge:     postDoc.PictureLarge,
		Type:             postDoc.Type,
		HgId:             interActorDoc.HgId,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	batch := r.firestore.Batch()
	postRef := r.firestore.Collection(postCollection).Doc(postDoc.PostId)
	batch.Set(postRef, hangoPost)
	interActor := postRef.Collection(interActorCollection).NewDoc()
	batch.Set(interActor, model.InterActor{
		HgId:    interActorDoc.HgId,
		Aka:     interActorDoc.Aka,
		Picture: interActorDoc.Picture,
	})
	if _, err := batch.Commit(ctx); err != nil {
		batch.Delete(postRef)
		batch.Delete(interActor)
		batch.Commit(ctx)
		return model.HangoPost{}, err
	}
	hangoPost.PostId = postDoc.PostId
	hangoPost.InterActor = append(hangoPost.InterActor, interActorDoc)

	return hangoPost, nil
}

// func (r *PostFeedsFirestoreRepository) EditHangoPost(ctx context.Context, postId string, postDoc model.InputPost) (model.HangoPost, error) {
// 	ref := r.firestore.Collection(postCollection).Doc(postId)
// 	err := r.firestore.RunTransaction(ctx, func(c context.Context, tx *firestore.Transaction) error {
// 		_, err := tx.Get(ref)
// 		if err != nil {
// 			return err
// 		}
// 		return tx.Set(ref, map[string]interface{}{
// 			"content":   postDoc.Content,
// 			"picture":   postDoc.Picture,
// 			"type":      postDoc.Type,
// 			"updatedAt": time.Now(),
// 		}, firestore.MergeAll)
// 	})
// 	if err != nil {
// 		if grpc.Code(err) == codes.NotFound {
// 			return model.HangoPost{}, errs.DocumentNotFound
// 		}
// 		return model.HangoPost{}, err
// 	}

// 	return r.GetHangoPostByID(ctx, postId)
// }

// TODO: Delete sub collection of post when post deleted
func (r *PostFeedsFirestoreRepository) DeleteHangoPost(ctx context.Context, postId string) error {
	if _, err := r.firestore.Collection(postCollection).Doc(postId).Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PostFeedsFirestoreRepository) GetHangoPostByID(ctx context.Context, postId string) (model.HangoPost, error) {
	dsnap, err := r.firestore.Collection(postCollection).Doc(postId).Get(ctx)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return model.HangoPost{}, errs.DocumentNotFound
		}
		return model.HangoPost{}, err
	}

	hangoPost := new(model.HangoPost)
	if dsnap.Exists() {
		dsnap.DataTo(hangoPost)
		iter := r.firestore.Doc(postCollection + "/" + postId).Collection(interActorCollection).Documents(ctx)
		defer iter.Stop()

		var interDocs []model.InterActor
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return model.HangoPost{}, err
			}
			if doc.Exists() {
				interDoc := new(model.InterActor)
				doc.DataTo(interDoc)
				interDocs = append(interDocs, *interDoc)
			}
		}
		hangoPost.InterActor = interDocs
	}

	return *hangoPost, nil
}

func (r *PostFeedsFirestoreRepository) GetAllHangoPost(ctx context.Context, postQuery model.FeedPostQuery) ([]model.PostFeedReponse, error) {
	state := r.firestore.Collection(postCollection)
	ref := state.OrderBy("createdAt", firestore.Desc).Where("roomId", "==", postQuery.RoomId).Limit(postQuery.Limit)
	var iter *firestore.DocumentIterator

	// Query post
	if postQuery.Type != "" {
		ref = ref.Where("type", "==", postQuery.Type)
	}
	if postQuery.LastPostId != "" {
		dsnap, err := state.Doc(postQuery.LastPostId).Get(ctx)
		if err != nil {
			if grpc.Code(err) == codes.NotFound {
				return []model.PostFeedReponse{}, errs.DocumentNotFound
			}
			return []model.PostFeedReponse{}, err
		}
		iter = ref.StartAfter(dsnap.Data()["createdAt"]).Documents(ctx)
	} else {
		iter = ref.Documents(ctx)
	}
	defer iter.Stop()

	var hangoPosts []model.PostFeedReponse
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []model.PostFeedReponse{}, err
		}

		if doc.Exists() {
			hangoPost := new(model.PostFeedReponse)
			mapstructure.Decode(doc.Data(), hangoPost)
			hangoPosts = append(hangoPosts, *hangoPost)
		}
	}
	// // This is for mock trends posts
	// if postQuery.Popular {
	// 	sort.Slice(hangoPosts, func(i, j int) bool { return hangoPosts[i].LikeCount > hangoPosts[j].LikeCount })
	// }

	return hangoPosts, nil
}

func (r *PostFeedsFirestoreRepository) ReportPost(ctx context.Context, postReport model.ReportPost) error {
	postRef := r.firestore.Collection(postCollection).Doc(postReport.PostId)
	iter := postRef.Collection(reportPostCollection).Where("hgId", "==", postReport.HgId).Where("type", "==", postReport.Type).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		// if user has reported
		if doc.Exists() {
			return nil
		}
	}
	// create report...
	_, err := postRef.Collection(reportPostCollection).NewDoc().Set(ctx, model.ReportPost{
		PostId:    postReport.PostId,
		Type:      postReport.Type,
		Detail:    postReport.Detail,
		HgId:      postReport.HgId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *PostFeedsFirestoreRepository) UpdateLikePost(ctx context.Context, postId, hgId string) (string, error) {
	query := r.firestore.Collection(postCollection).Where("postId", "==", postId).Limit(1).Where("like", "array-contains", hgId)
	var result string

	err := r.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		iter := tx.Documents(query)
		ref := r.firestore.Doc(postCollection + "/" + postId)

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
