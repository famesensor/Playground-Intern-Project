package domain

import "time"

var AnonymousList = map[int]interface{}{
	0:  "แมวน้ำเงิน",
	1:  "แมวฟ้า",
	2:  "แมวเขียว",
	3:  "แมวส้ม",
	4:  "แมวชมพู",
	5:  "แมวม่วง",
	6:  "แมวแดง",
	7:  "แมวเหลือง",
	8:  "หมาน้ำเงิน",
	9:  "หมาฟ้า",
	10: "หมาเขียว",
	11: "หมาส้ม",
	12: "หมาชมพู",
	13: "หมาแดง",
	14: "หมาหมา",
	15: "หมาเหลือง",
	16: "แพนด้าน้ำเงิน",
	17: "แพนด้าฟ้า",
	18: "แพนด้าเขียว",
	19: "แพนด้าส้ม",
	20: "แพนด้าชมพู",
	21: "แพนด้าม่วง",
	22: "แพนด้าแดง",
	23: "แพนด้าเหลือง",
	24: "หมูน้ำเงิน",
	25: "หมูฟ้า",
	26: "หมูเขียว",
	27: "หมูส้ม",
	28: "หมูชมพู",
	29: "หมูม่วง",
	30: "หมูแดง",
	31: "หมูเหลือง",
	32: "กระต่ายน้ำเงิน",
	33: "กระต่ายฟ้า",
	34: "กระต่ายเขียว",
	35: "กระต่ายส้ม",
	36: "กระต่ายชมพู",
	37: "กระต่ายม่วง",
	38: "กระต่ายแดง",
	39: "กระต่ายเหลือง",
	40: "หนูน้ำเงิน",
	41: "หนูฟ้า",
	42: "หนูเขียว",
	43: "หนูส้ม",
	44: "หนูชมพู",
	45: "หนูม่วง",
	46: "หนูแดง",
	47: "หนูเหลือง",
	48: "แมวน้ำน้ำเงิน",
	49: "แมวน้ำฟ้า",
	50: "แมวน้ำเขียว",
	51: "แมวน้ำส้ม",
	52: "แมวน้ำชมพู",
	53: "แมวน้ำม่วง",
	54: "แมวน้ำแดง",
	55: "แมวน้ำเหลือง",
	56: "เสือน้ำเงิน",
	57: "เสือฟ้า",
	58: "เสือเขียว",
	59: "เสือส้ม",
	60: "เสือชมพู",
	61: "เสือม่วง",
	62: "เสือแดง",
	63: "เสือเหลือง",
}

type HangoComment struct {
	CommentId        string    `json:"commentId" firestore:"commentId"`
	Comment          string    `json:"comment" firestore:"comment"`
	PictureThumbnail string    `json:"pictureThumbnail" firestore:"pictureThumbnail"`
	PictureLarge     string    `json:"pictureLarge" firestore:"pictureLarge"`
	Like             []string  `json:"like" firestore:"like"`
	NestCount        int       `json:"nestCount" firestore:"nestCount"`
	IsNest           bool      `json:"isNest" firestore:"isNest"`
	Hide             bool      `json:"hide" firestore:"hide"`
	HgId             string    `json:"hgId" firestore:"hgId"`
	CreatedAt        time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" firestore:"updatedAt"`
}

type InputComment struct {
	PostId           string `json:"postId" form:"postId" validate:"required"`
	CommentId        string `json:"commentId" form:"commentId" validate:"required_if=Type reply"`
	Comment          string `json:"comment" form:"comment"`
	PictureThumbnail string `json:"pictureThumbnail" firestore:"pictureThumbnail" validate:"omitempty,url"`
	PictureLarge     string `json:"pictureLarge" firestore:"pictureLarge" validate:"omitempty,url"`
	// Type for this comment is root or reply comment
	Type string `json:"type" form:"type" validate:"required,oneof=root reply"`
}

type UpdateComment struct {
	PostId    string `json:"postId" form:"postId" validate:"required"`
	CommentId string `json:"commentId" form:"commentId" validate:"required"`
	ReplyId   string `json:"replyId" form:"replyId" validate:"required_if=Type reply"`
	Comment   string `json:"comment" form:"comment" validate:"required"`
	Picture   string `json:"picture" form:"picture" validate:"omitempty,url"`
	// Type for this comment is root or reply comment
	Type string `json:"type" form:"type" validate:"required,oneof=root reply"`
}

type CommentNotiDocuments struct {
	HangoComment HangoComment
	HgIdNoti     string
}

type PostCommentQuery struct {
	PostId    string `query:"postId" validate:"required"`
	CommentId string `query:"commentId" validate:"required"`
	ReplyId   string `query:"replyId" validate:"required_if=Type reply"`
	Type      string `query:"type" validate:"required,oneof=root reply"`
}

type UploadImageCommentQuery struct {
	Collection string `query:"collection" validate:"required"`
	PostId     string `query:"postId" validate:"required"`
	CommentId  string `query:"commentId" validate:"required_if=Type reply"`
	Type       string `query:"type" validate:"required,oneof=root reply"`
}

type GetCommentQuery struct {
	PostId        string `query:"postId" validate:"required"`
	CommentId     string `query:"commentId" validate:"required_if=Type reply"`
	LastCommentId string `query:"lastCommentId"`
	Limit         int    `query:"limit"`
	Type          string `query:"type" validate:"required,oneof=root reply"`
}
