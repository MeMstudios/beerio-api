package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Email       string        `bson:"email" json:"email"`
	Password    string        `bson:"password" json:"password"`
	ProfilePic  string        `bson:"cover_image" json:"cover_image"`
	Description string        `bson:"description" json:"description"`
	Favorites   []Favorite    `bson:"favorites" json:"favorites"`
}

type UserCheck struct {
	Email, Password string
}

type Favorite struct {
	UserID, BeerID string
}

// type UserSession struct {
// 	ID      bson.ObjectId `bson:"_id" json:"id"`
// 	UserID  bson.ObjectId `bson:"user_id" json:"user_id"`
// 	LastHit int           `bson:"last" json:"last"`
// }
