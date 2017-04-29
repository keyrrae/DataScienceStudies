package models

import "gopkg.in/mgo.v2/bson"

// User represents the structure of our resource
type Post struct {
	Id    bson.ObjectId `json:"id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Story string        `json:"story" bson:"story"`
	Lat float64        `json:"lat" bson:"lat"`
	Lng float64          `json:"lng" bson:"lng"`
}
