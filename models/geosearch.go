package models

// Location represents the structure of location
type GeoSearchRequest struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Radius float64 `json:"radius"`
	Count  int `json:"count"`
}
