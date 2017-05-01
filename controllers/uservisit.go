package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"github.com/keyrrae/monimenta_backend/models"
)

// UserVisitTown update visit town
func (c DbController) UserVisitTown(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	townId := params["townid"]
	userId := params["userid"]

	if !bson.IsObjectIdHex(townId) || !bson.IsObjectIdHex(userId) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	townDbID := bson.ObjectIdHex(townId)

	// Stub town
	t := models.Town{}

	// Fetch town
	if err := c.Session.DB(c.DbName).C(c.TownTable).FindId(townDbID).One(&t); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t.VisitedBy = append(t.VisitedBy, userId)

	// Remove user
	if err := c.Session.DB(c.DbName).C(c.TownTable).UpdateId(t.Id, t); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}


