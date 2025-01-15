package services

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/JonayMedina/go-search-api/internal/models"
)

type MusicService struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
	providers   []models.Provider
}

func NewMusicService(mongoClient *mongo.Client, redisClient *redis.Client) *MusicService {
	service := &MusicService{
		mongoClient: mongoClient,
		redisClient: redisClient,
		providers:   []models.Provider{},
	}

	// Registrar proveedores
	service.registerProviders()
	return service
}

func (s *MusicService) registerProviders() {
	// iTunes Provider
	s.providers = append(s.providers, models.Provider{
		Name:    "iTunes",
		BaseURL: "https://itunes.apple.com/search",
		Type:    models.ProviderTypeJSON,
		SearchFn: func(ctx context.Context, query string) ([]models.Song, error) {
			return s.searchITunes(ctx, query)
		},
	})

	// ChartLyrics Provider
	s.providers = append(s.providers, models.Provider{
		Name:    "ChartLyrics",
		BaseURL: "http://api.chartlyrics.com/apiv1.asmx/SearchLyric",
		Type:    models.ProviderTypeSOAP,
		SearchFn: func(ctx context.Context, query string) ([]models.Song, error) {
			return s.searchChartLyrics(ctx, query)
		},
	})
}

func (s *MusicService) searchITunes(ctx context.Context, query string) ([]models.Song, error) {
	// Construir URL con parámetros
	url := fmt.Sprintf("%s?term=%s&media=music", s.providers[0].BaseURL, url.QueryEscape(query))

	// Realizar petición HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request para iTunes: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en petición a iTunes: %w", err)
	}
	defer resp.Body.Close()

	log.Println(resp.Body)

	// Estructura para respuesta de iTunes
	var response struct {
		Results []struct {
			TrackName      string `json:"trackName"`
			ArtistName     string `json:"artistName"`
			CollectionName string `json:"collectionName"`
			PreviewURL     string `json:"previewUrl"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de iTunes: %w", err)
	}

	// Convertir resultados al modelo Song
	var songs []models.Song
	for _, item := range response.Results {
		songs = append(songs, models.Song{
			Name:       item.TrackName,
			Artist:     item.ArtistName,
			Album:      item.CollectionName,
			Provider:   "iTunes",
			PreviewURL: item.PreviewURL,
			Origin:     "iTunes",
		})
	}

	return songs, nil
}

func (s *MusicService) searchChartLyrics(ctx context.Context, query string) ([]models.Song, error) {
	// Dividir la consulta en artista y canción (asumiendo formato "artista - canción")
	parts := strings.Split(query, "-")
	artist := strings.TrimSpace(query)
	song := ""
	if len(parts) > 1 {
		artist = strings.TrimSpace(parts[0])
		song = strings.TrimSpace(parts[1])
	}

	// Construir URL con parámetros
	url := fmt.Sprintf("%s?artist=%s&song=%s",
		s.providers[1].BaseURL,
		url.QueryEscape(artist),
		url.QueryEscape(song))

	// Realizar petición HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request para ChartLyrics: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en petición a ChartLyrics: %w", err)
	}
	defer resp.Body.Close()

	log.Println(resp.Body)

	// Estructura para respuesta de ChartLyrics
	var response struct {
		SearchLyricResult struct {
			SearchLyricResponse []struct {
				Song    string `xml:"Song"`
				Artist  string `xml:"Artist"`
				LyricId string `xml:"LyricId"`
				SongUrl string `xml:"SongUrl"`
			} `xml:"SearchLyricResult"`
		} `xml:"SearchLyricResult"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de ChartLyrics: %w", err)
	}

	// Convertir resultados al modelo Song
	var songs []models.Song
	for _, item := range response.SearchLyricResult.SearchLyricResponse {
		songs = append(songs, models.Song{
			Name:       item.Song,
			Artist:     item.Artist,
			Provider:   "ChartLyrics",
			PreviewURL: item.SongUrl,
			Origin:     "ChartLyrics",
		})
	}

	return songs, nil
}

func (s *MusicService) SearchMusic(ctx context.Context, query string) ([]models.Song, error) {
	// Validar entrada
	if query == "" {
		return nil, errors.New("la consulta de búsqueda no puede estar vacía")
	}

	// Intentar obtener del caché
	songs, err := s.getFromCache(ctx, query)
	if err == nil {
		return songs, nil
	}

	// Búsqueda paralela en APIs
	songs, err = s.searchInParallel(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error en búsqueda paralela: %w", err)
	}

	// Guardar resultados
	if err := s.saveResults(ctx, query, songs); err != nil {
		log.Printf("Error al guardar resultados: %v", err)
	}

	return songs, nil
}

func (service *MusicService) searchInParallel(ctx context.Context, query string) ([]models.Song, error) {
	var wg sync.WaitGroup
	results := make(chan []models.Song, 3)
	errors := make(chan error, 3)

	// Ejecutar búsquedas en paralelo
	for _, provider := range service.providers {
		wg.Add(1)
		go func(p models.Provider) {
			defer wg.Done()
			songs, err := p.Search(ctx, query)
			if err != nil {
				errors <- err
				return
			}
			results <- songs
		}(provider)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	var allSongs []models.Song
	for songs := range results {
		allSongs = append(allSongs, songs...)
	}

	// Verificar errores
	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return allSongs, nil
}

func (service *MusicService) getFromCache(ctx context.Context, query string) ([]models.Song, error) {
	// Intentar obtener del caché
	cacheKey := "search:" + query
	if cached, err := service.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var songs []models.Song
		if err := json.Unmarshal([]byte(cached), &songs); err == nil {
			return songs, nil
		}
	}
	return nil, nil
}

func (service *MusicService) saveResults(ctx context.Context, query string, songs []models.Song) error {
	if len(songs) == 0 {
		return nil
	}

	collection := service.mongoClient.Database("musicdb").Collection("searches")

	searchRecord := map[string]interface{}{
		"query":     query,
		"songs":     songs,
		"timestamp": time.Now(),
	}

	_, err := collection.InsertOne(ctx, searchRecord)
	return err
}

type MusicServiceInterface interface {
	SearchMusic(ctx context.Context, query string) ([]models.Song, error)
}
