package models

import "gopkg.in/mgo.v2/bson"

// User represents the structure of our resource
type User struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Cred    string        `json:"cred" bson:"cred"`
	Name    string        `json:"name" bson:"name"`
	Email   string        `json:"email" bson:"email"`
	TownIDs []string      `bson:"towns"`
}

type UserTown struct {
	Id    bson.ObjectId `json:"id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Email string        `json:"email" bson:"email"`
}

type UserInfo struct {
	Id    bson.ObjectId `json:"id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Email string        `json:"email" bson:"email"`
}
