package models

import "encoding/xml"

type Song struct {
	ID          string  `json:"id" bson:"_id,omitempty"`
	Name        string  `json:"name" bson:"name"`
	Artist      string  `json:"artist" bson:"artist"`
	Duration    string  `json:"duration" bson:"duration,omitempty"`
	Album       string  `json:"album" bson:"album,omitempty"`
	Artwork     string  `json:"artwork" bson:"artwork,omitempty"`
	Price       float64 `json:"price" bson:"price,omitempty"`
	Origin      string  `json:"origin" bson:"origin,omitempty"`
	Provider    string  `json:"provider" bson:"provider,omitempty"`
	PreviewURL  string  `json:"preview_url" bson:"preview_url,omitempty"`
	ImageURL    string  `json:"image_url,omitempty" bson:"image_url,omitempty"`
	ReleaseDate string  `json:"release_date,omitempty" bson:"release_date,omitempty"`
	Genre       string  `json:"genre,omitempty" bson:"genre,omitempty"`
}

type ChartLyricsResponse struct {
	XMLName           xml.Name `xml:"ArrayOfSearchLyricResult"`
	SearchLyricResult []struct {
		Artist  string `xml:"Artist"`
		Song    string `xml:"Song"`
		SongURL string `xml:"SongUrl"`
		LyricID string `xml:"LyricId"`
		TrackID string `xml:"TrackId"`
	} `xml:"SearchLyricResult"`
}
