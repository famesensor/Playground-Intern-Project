package domain

import (
	"bytes"
	"image"
	"time"
)

type InputPost struct {
	PostId           string   `json:"postId" firestore:"-"`
	RoomId           string   `json:"roomId" firestore:"roomId" form:"roomId" validate:"required"`
	Content          string   `json:"content" firestore:"content" form:"content"`
	PictureThumbnail []string `json:"pictureThumbnail" firestore:"pictureThumbnail" validate:"omitempty,dive,url"`
	PictureLarge     []string `json:"pictureLarge" firestore:"pictureLarge" validate:"omitempty,dive,url"`
	Type             string   `json:"type" firestore:"type" form:"type" validate:"required"`
}

type HangoPost struct {
	PostId           string       `json:"postId" firestore:"postId,omitempty"`
	RoomId           string       `json:"roomId" firestore:"roomId"`
	Content          string       `json:"content" firestore:"content"`
	PictureThumbnail []string     `json:"pictureThumbnail" firestore:"pictureThumbnail"`
	PictureLarge     []string     `json:"pictureLarge" firestore:"pictureLarge"`
	Type             string       `json:"type" firestore:"type"`
	Like             []string     `json:"like" firestore:"like"`
	LikeCount        int          `json:"likeCount" firestore:"likeCount"`
	CommentCount     int          `json:"commentCount" firestore:"commentCount"`
	InterActor       []InterActor `json:"interActor" firestore:"interActor,omitempty"`
	HgId             string       `json:"hgId" firestore:"hgId"`
	CreatedAt        time.Time    `json:"createdAt" firestore:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt" firestore:"updatedAt"`
}

type PostFeedReponse struct {
	PostId           string    `json:"postId" firestore:"postId"`
	RoomId           string    `json:"roomId" firestore:"roomId"`
	Content          string    `json:"content" firestore:"content"`
	PictureThumbnail []string  `json:"pictureThumbnail" firestore:"pictureThumbnail"`
	Like             []string  `json:"like" firestore:"like"`
	LikeCount        int       `json:"likeCount" firestore:"likeCount"`
	CommentCount     int       `json:"commentCount" firestore:"commentCount"`
	CreatedAt        time.Time `json:"createdAt" firestore:"createdAt"`
}

type InterActor struct {
	HgId    string `json:"hgId" firestore:"hgId"`
	Picture string `json:"picture" firestore:"picture"`
	Aka     string `json:"aka" firestore:"aka"`
	UId     int    `json:"UId" firestore:"UId"`
}

type FeedPostQuery struct {
	RoomId     string `query:"roomId"`
	Type       string `query:"type"`
	Limit      int    `query:"limit"`
	Popular    bool   `query:"popular"`
	LastPostId string `query:"lastPostId"`
}

type UploadImagePostQuery struct {
	Collection string `query:"collection" validate:"required"`
	RoomId     string `query:"roomId" validate:"required"`
	PostId     string `query:"postId"`
	Type       string `query:"type" validate:"required,oneof=post edit"`
	HgId       string
}

type ReportPost struct {
	PostId    string    `json:"postId" form:"postId" firestore:"postId" validate:"required"`
	Type      string    `json:"type" form:"type" firestore:"type" validate:"required"`
	Detail    string    `json:"detail" form:"detail" firestore:"detail" validate:"required"`
	HgId      string    `json:"hgId" firestore:"hgId"`
	CreatedAt time.Time `json:"createdAt" firestore:"cretedAt"`
}

type PictureDoc struct {
	File image.Image
	Type string
}

type PictureResize struct {
	ImageThumbnail []bytes.Buffer
	ImageLarge     []bytes.Buffer
}

type PictureDocuments struct {
	UrlThumbnail []string `json:"urlThumbnail"`
	UrlLarge     []string `json:"urlLarge"`
	PostId       string   `json:"postId"`
}
