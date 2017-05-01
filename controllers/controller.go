package controllers

import "gopkg.in/mgo.v2"

type (
	// UserController represents the controller for operating on the User resource
	DbController struct {
		DbName    string
		Session   *mgo.Session
		UserTable string
		TownTable string
		LocTable  string
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewDbController(s *mgo.Session, db_name string) *DbController {
	res := &DbController{
		DbName:    db_name,
		Session:   s,
		UserTable: "users",
		TownTable: "towns",
		LocTable:  "location",
	}
	return res
}
