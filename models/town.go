package models

import "gopkg.in/mgo.v2/bson"

type Town struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Title       string        `json:"title" bson:"title"`
	UserId      string        `json:"user_id" bson:"user_id"`
	Category    string        `json:"category" bson:"category"`
	Address     string        `json:"address" bson:"address"`
	Description string        `json:"description" bson:"description"`
	Lat         float64       `json:"lat" bson:"lat"`
	Lng         float64       `json:"lng" bson:"lng"`
	ImageUrls   []string      `json:"image_urls" bson:"image_urls"`
	VisitedBy   []string      `json:"visited_by" bson:"visited_by"`
	Sketch      string        `json:"sketch" bson:"sketch"`
}

type TownGeoRequest struct {
	Lat    float64 `json:"lat" bson:"lat"`
	Lng    float64 `json:"lng" bson:"lng"`
	Radius float64 `json:"radius" bson:"radius"`
	Count  int     `json:"display_count" bson:"display_count"`
}

type Coord struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
