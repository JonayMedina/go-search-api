package services

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/JonayMedina/go-search-api/internal/models"
)

type MusicService struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
}

func NewMusicService(mongoClient *mongo.Client, redisClient *redis.Client) *MusicService {
	return &MusicService{
		mongoClient: mongoClient,
		redisClient: redisClient,
	}
}

func (s *MusicService) SearchMusic(ctx context.Context, query string) ([]models.Song, error) {
	// Intentar obtener del caché
	cacheKey := "search:" + query
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var songs []models.Song
		if err := json.Unmarshal([]byte(cached), &songs); err == nil {
			return songs, nil
		}
	}

	var wg sync.WaitGroup
	resultChan := make(chan []models.Song, 3)
	errorChan := make(chan error, 3)

	wg.Add(2) // Solo iTunes y ChartLyrics por ahora
	go s.searchITunes(ctx, query, resultChan, errorChan, &wg)
	go s.searchChartLyrics(ctx, query, resultChan, errorChan, &wg)

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	var allSongs []models.Song
	for songs := range resultChan {
		allSongs = append(allSongs, songs...)
	}

	// Verificar errores
	for err := range errorChan {
		if err != nil {
			return nil, err
		}
	}

	// Guardar en caché
	if len(allSongs) > 0 {
		if cached, err := json.Marshal(allSongs); err == nil {
			s.redisClient.Set(ctx, cacheKey, cached, 1*time.Hour)
		}
	}

	// Guardar en MongoDB
	if err := s.saveSongs(ctx, allSongs); err != nil {
		return nil, err
	}

	return allSongs, nil
}

func (s *MusicService) searchITunes(ctx context.Context, query string, resultChan chan<- []models.Song, errorChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("https://itunes.apple.com/search?term=%s", url.QueryEscape(query))
	resp, err := http.Get(url)
	if err != nil {
		errorChan <- err
		return
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			TrackID         int     `json:"trackId"`
			TrackName       string  `json:"trackName"`
			ArtistName      string  `json:"artistName"`
			CollectionName  string  `json:"collectionName"`
			ArtworkUrl100   string  `json:"artworkUrl100"`
			TrackPrice      float64 `json:"trackPrice"`
			TrackTimeMillis int     `json:"trackTimeMillis"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		errorChan <- err
		return
	}

	var songs []models.Song
	for _, item := range result.Results {
		duration := time.Duration(item.TrackTimeMillis) * time.Millisecond
		songs = append(songs, models.Song{
			ID:       fmt.Sprintf("itunes_%d", item.TrackID),
			Name:     item.TrackName,
			Artist:   item.ArtistName,
			Duration: fmt.Sprintf("%02d:%02d", int(duration.Minutes()), int(duration.Seconds())%60),
			Album:    item.CollectionName,
			Artwork:  item.ArtworkUrl100,
			Price:    fmt.Sprintf("USD %.2f", item.TrackPrice),
			Origin:   "itunes",
		})
	}

	resultChan <- songs
}

func (s *MusicService) searchChartLyrics(ctx context.Context, query string, resultChan chan<- []models.Song, errorChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	// ChartLyrics requiere artista y canción separados, así que usaremos el query completo para ambos
	url := fmt.Sprintf("http://api.chartlyrics.com/apiv1.asmx/SearchLyric?artist=%s&song=%s",
		url.QueryEscape(query), url.QueryEscape(query))

	resp, err := http.Get(url)
	if err != nil {
		errorChan <- err
		return
	}
	defer resp.Body.Close()

	var result struct {
		SearchLyricResult []struct {
			LyricId int    `xml:"LyricId"`
			Song    string `xml:"Song"`
			Artist  string `xml:"Artist"`
			SongUrl string `xml:"SongUrl"`
		} `xml:"SearchLyricResult"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		errorChan <- err
		return
	}

	var songs []models.Song
	for _, item := range result.SearchLyricResult {
		songs = append(songs, models.Song{
			ID:       fmt.Sprintf("chartlyrics_%d", item.LyricId),
			Name:     item.Song,
			Artist:   item.Artist,
			Duration: "00:00", // ChartLyrics no proporciona duración
			Album:    "",      // ChartLyrics no proporciona álbum
			Artwork:  "",      // ChartLyrics no proporciona artwork
			Price:    "N/A",   // ChartLyrics no proporciona precio
			Origin:   "chartlyrics",
		})
	}

	resultChan <- songs
}

func (s *MusicService) saveSongs(ctx context.Context, songs []models.Song) error {
	if len(songs) == 0 {
		return nil
	}

	collection := s.mongoClient.Database("musicdb").Collection("songs")

	// Convertir songs a []interface{} para bulk insert
	docs := make([]interface{}, len(songs))
	for i, song := range songs {
		docs[i] = song
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}
