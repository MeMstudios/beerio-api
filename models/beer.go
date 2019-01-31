package models

import "gopkg.in/mgo.v2/bson"

type Beer struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	BreweryDbId string        `bson:"bdb_id" json:"bdb_id"`
	Description string        `bson:"description" json:"description"`
	Labels      Label         `bson:"labels" json:"labels"`
	ABV         string        `bson:"abv" json:"abv"`
	Status      string        `bson:"status" json:"status"`
	Available   string        `bson:"available" json:"available"`
	Style       string        `bson:"style" json:"style"`
	Rating      float32       `bson:"rating" json:"rating"`
}

type Label struct {
	Icon   string `bson:"icon" json:"icon"`
	Medium string `bson:"medium" json:"medium"`
	Large  string `bson:"large" json:"large"`
}
