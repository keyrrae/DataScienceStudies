package models

import "gopkg.in/mgo.v2/bson"

// Location represents the structure of location
type Location struct {
	Id    bson.ObjectId `json:"id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Lat   float64       `json:"lat" bson:"lat"`
	Lng   float64       `json:"lng" bson:"lng"`
	//Posts []Post        `json:"posts" bson:"posts"`
}
