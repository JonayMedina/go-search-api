package models

type Song struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Artist   string `json:"artist" bson:"artist"`
	Duration string `json:"duration" bson:"duration"`
	Album    string `json:"album" bson:"album"`
	Artwork  string `json:"artwork" bson:"artwork"`
	Price    string `json:"price" bson:"price"`
	Origin   string `json:"origin" bson:"origin"`
}
