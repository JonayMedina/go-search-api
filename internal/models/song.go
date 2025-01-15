package models

type Song struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	Name       string `json:"name" bson:"name"`
	Artist     string `json:"artist" bson:"artist"`
	Duration   string `json:"duration" bson:"duration,omitempty"`
	Album      string `json:"album" bson:"album,omitempty"`
	Artwork    string `json:"artwork" bson:"artwork,omitempty"`
	Price      string `json:"price" bson:"price,omitempty"`
	Origin     string `json:"origin" bson:"origin,omitempty"`
	Provider   string `json:"provider" bson:"provider,omitempty"`
	PreviewURL string `json:"preview_url" bson:"preview_url,omitempty"`
}
