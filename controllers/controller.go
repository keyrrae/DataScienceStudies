package controllers

import (
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
)

type (
	// UserController represents the controller for operating on the User resource
	DbController struct {
		DbName      string
		Session     *mgo.Session
		RedisClient *redis.Client
		UserTable   string
		TownTable   string
		LocTable    string
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewDbController(s *mgo.Session, r *redis.Client, db_name string) *DbController {
	res := &DbController{
		DbName:      db_name,
		Session:     s,
		RedisClient: r,
		UserTable:   "users",
		TownTable:   "towns",
		LocTable:    "location",
	}
	return res
}
