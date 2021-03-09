package domain

import (
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"
)

type ResponseRestaurant struct {
	RestId    string         `json:"restId" fierstore:"restId"`
	RestName  string         `json:"restName" firestore:"restName"`
	Location  *latlng.LatLng `json:"location" firestore:"location"`
	Picture   []string       `json:"picture" firestore:"picture"`
	Info      RestInfo       `json:"info" firestore:"info"`
	Tag       []string       `json:"tag" firestore:"tag"`
	CreatedAt time.Time      `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" firestore:"updatedAt"`
}

type ReponseHomeRestaurant struct {
	RestId   string         `json:"restId" fierstore:"restId"`
	RestName string         `json:"restName" firestore:"restName"`
	Location *latlng.LatLng `json:"location" firestore:"location"`
	Tag      []string       `json:"tag" firestore:"tag"`
	Picture  []string       `json:"picture" firestore:"picture"`
}

type RestInfo struct {
	CarPark      bool     `json:"carPark" firestore:"carPark"`
	RestPhone    []string `json:"restPhone" firestore:"restPhone"`
	Open         string   `json:"open" firestore:"open"`
	Close        string   `json:"close" firestore:"close"`
	Payment      []string `json:"payment" firestore:"payment"`
	CheckInCount int      `json:"checkInCount" firestore:"checkInCount"`
	Reduis       float64  `firestore:"reduis"`
}

type Promotion struct {
	PromotionId string    `json:"omitempty,promotionId" firestore:"omitempty,promotionId"`
	Title       string    `json:"title" firestore:"title"`
	Description string    `json:"description" firestore:"description"`
	Picture     []string  `json:"picture" firestore:"picture"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updatedAt"`
}

// For create mock restaurant
type CreateRestaurant struct {
	RestId    string         `firestore:"restId"`
	RestName  string         `json:"restName" firestore:"restName" form:"restName" validate:"required"`
	Location  *latlng.LatLng `json:"location" firestore:"location" form:"location" validate:"required"`
	GeoHash   string         `json:"geohash" firestore:"geohash"`
	Picture   []string       `json:"picture" firestore:"picture"`
	Info      RestInfo       `json:"info" firestore:"info" form:"info" validate:"required"`
	Tag       []string       `json:"tag" firestore:"tag" form:"tag" validate:"required"`
	CreatedAt time.Time      `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" firestore:"updatedAt"`
}

type RestaurantQuery struct {
	Lat      float64  `query:"lat"`
	Lng      float64  `query:"lng"`
	Distance float64  `query:"distance"`
	Tag      []string `query:"tag"`
	Open     bool     `query:"open"`
	Limit    int      `query:"limit"`
	RestId   string   `query:"restId"`
}

type ListRestaurantQuery struct {
	Lower   string
	Upper   string
	Tag     []string
	TimeNow string
	Limit   int
	RestId  string
}

type UserCheckIn struct {
	UserLat float64 `json:"userLat" form:"userLat" validate:"required"`
	Userlng float64 `json:"userLng" form:"userLng" validate:"required"`
	RestId  string  `json:"restId" form:"restId" validate:"required"`
}

type CheckInDoc struct {
	RestId    string    `firestore:"restId"`
	HgId      string    `firestore:"hgId"`
	UId       string    `firestore:"uid"`
	IsRevoked bool      `firestore:"isRevoked"`
	Peak      bool      `firestore:"peakMode"`
	NameAnon  string    `firestore:"nameAnon"`
	Picture   string    `firestore:"picture"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}

type AnonymousDoc struct {
	UId       string `firestore:"uid"`
	IsRevoked bool   `firestore:"isRevoked"`
	Peak      bool   `firestore:"peakMode"`
	NameAnon  string `firestore:"nameAnon"`
	Picture   string `firestore:"picture"`
}

type PeakModeDoc struct {
	HgId      string `firestore:"hgId"`
	NickName  string `firestore:"nickName"`
	Picture   string `firestore:"picture"`
	Status    string `firestore:"status"`
	IsRevoked bool   `firestore:"isRevoked"`
	Peak      bool   `firestore:"peakMode"`
}
