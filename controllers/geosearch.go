package controllers

import (
	"net/http"
	"github.com/keyrrae/monimenta_backend/models"
	"encoding/json"
	"math"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func (c DbController) GeoSearch(w http.ResponseWriter, r *http.Request) {
	gs := models.GeoSearchRequest{}

	err := json.NewDecoder(r.Body).Decode(&gs)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//geoID := genGeoId(gs.Lat, gs.Lng)
	center := models.CoOrd{
		Lat: int(gs.Lat * 10) ,
		Lng: int(gs.Lng * 10),
	}

	radius := int(math.Ceil(gs.Radius))

	step := 1
	numStep := int(2 * radius) + 1

	ul := models.CoOrd{
		Lat: center.Lat - radius,
		Lng: center.Lng - radius,
	}

	//townChan := make(chan *models.Town, numStep * numStep)

	townChan := make(chan []models.Town, numStep * numStep)

	for i := 0; i < numStep; i++ {
		for j := 0; j < numStep; j++ {
			geoId := genIntGeoId(ul.Lat + i * step, ul.Lng + j * step)
			fmt.Println(geoId)
			go func(){
				// Remove user

				var towns []models.Town
				collection := c.Session.DB(c.DbName).C(c.TownTable)
				if err := collection.Find(bson.M{"geoid": geoId}).All(&towns); err != nil {
					townChan <- nil
					return
				}
				fmt.Println(towns)

				townChan <- towns

			}()
		}
	}

	//var res []models.Town

	for i := 0; i < numStep * numStep; i++ {
		town := <- townChan
		if town != nil {
			//res = append(res, *town)
			fmt.Println(town)
		}
	}
  // TODO: priority queue
}