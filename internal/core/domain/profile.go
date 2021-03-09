package domain

import "time"

type Profile struct {
	HgId         string    `json:"userDocId" firestore:"hgId"`
	Nickname     string    `json:"nickname,omitempty" firestore:"nickname,omitempty"`
	Gender       string    `firestore:"gender,omitempty"`
	FeelingLevel int       `json:"feelinglevel" firestore:"feelinglevel"`
	IsPublic     bool      `json:"ispublic" firestore:"ispublic"`
	Info         Info      `json:"info" firestore:"info"`
	Picture      Picture   `json:"picture" firestore:"picture"`
	CreatedAt    time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" firestore:"updatedAt"`
}

type Info struct {
	Status        string `json:"status,omitempty" firestore:"status,omitempty"`
	Age           string `json:"age,omitempty" firestore:"age,omitempty"`
	Interest      string `json:"interest,omitempty" firestore:"interest,omitempty"`
	Position      string `json:"position,omitempty" firestore:"position,omitempty"`
	Company       string `json:"company,omitempty" firestore:"company,omitempty"`
	University    string `json:"university,omitempty" firestore:"university,omitempty"`
	Address       string `json:"address,omitempty" firestore:"address,omitempty"`
	FavoriteDrink string `json:"favoritedrink,omitempty" firestore:"favoritedrink,omitempty"`
	FavoriteSong  string `json:"favoritesong,omitempty" firestore:"favoritesong,omitempty"`
	FavoriteFood  string `json:"favoritefood,omitempty" firestore:"favoritefood,omitempty"`
}

type Picture struct {
	Img1 string `json:"img1,omitempty" firestore:"img1,omitempty"`
	Img2 string `json:"img2,omitempty" firestore:"img2,omitempty"`
	Img3 string `json:"img3,omitempty" firestore:"img3,omitempty"`
	Img4 string `json:"img4,omitempty" firestore:"img4,omitempty"`
	Img5 string `json:"img5,omitempty" firestore:"img5,omitempty"`
	Img6 string `json:"img6,omitempty" firestore:"img6,omitempty"`
	Img7 string `json:"img7,omitempty" firestore:"img7,omitempty"`
	Img8 string `json:"img8,omitempty" firestore:"img8,omitempty"`
	Img9 string `json:"img9,omitempty" firestore:"img9,omitempty"`
}
