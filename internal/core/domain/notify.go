package domain

import "time"

type NotiPostDoc struct {
	NotiId        string    `firestore:"notiId"`
	PostId        string    `json:"postId" firestore:"postId"`
	TitleContent  string    `json:"titleContent" firestore:"titleContent"`
	DetaliContent string    `json:"detailContent" firestore:"detailContent"`
	ReponseRef    string    `json:"reponseRef" firestore:"reponseRef"`
	OwnerRef      string    `json:"ownerRef" firestore:"ownerRef"`
	IsRead        bool      `json:"isRead" firestore:"isRead"`
	CreatedAt     time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" firestore:"updatedAt"`
}

type NotiMatchDoc struct {
	NotiId string `firestore:"notiId"`

	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`
}
