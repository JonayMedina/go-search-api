package services

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/JonayMedina/go-search-api/internal/models"
)

type MusicService struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
	providers   []models.ProviderInterface
}

func NewMusicService(mongoClient *mongo.Client, redisClient *redis.Client) *MusicService {
	service := &MusicService{
		mongoClient: mongoClient,
		redisClient: redisClient,
		providers:   []models.ProviderInterface{},
	}

	// Registrar proveedores
	service.registerProviders()
	return service
}

func (service *MusicService) registerProviders() {
	service.providers = append(service.providers, &models.Provider{
		Name:     "iTunes",
		BaseURL:  "https://itunes.apple.com/search",
		Type:     models.ProviderTypeJSON,
		SearchFn: service.searchITunes,
	})

	// ChartLyrics Provider
	service.providers = append(service.providers, models.Provider{
		Name:    "ChartLyrics",
		BaseURL: "http://api.chartlyrics.com/apiv1.asmx/SearchLyric",
		Type:    models.ProviderTypeSOAP,
		SearchFn: func(ctx context.Context, query string) ([]models.Song, error) {
			return service.searchChartLyrics(query)
		},
	})
}

func (service *MusicService) searchITunes(ctx context.Context, query string) ([]models.Song, error) {
	baseURL := "https://itunes.apple.com/search"
	url := fmt.Sprintf("%s?term=%s&media=music", baseURL, url.QueryEscape(query))
	log.Printf("url: %v", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request para iTunes: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error en petición a iTunes: %v", err)
		return nil, fmt.Errorf("error en petición a iTunes: %w", err)
	}
	defer resp.Body.Close()

	// Estructura para respuesta de iTunes
	var response struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			TrackName      string  `json:"trackName"`
			ArtistName     string  `json:"artistName"`
			CollectionName string  `json:"collectionName"`
			PreviewURL     string  `json:"previewUrl"`
			ArtworkUrl100  string  `json:"artworkUrl100"`
			TrackPrice     float64 `json:"trackPrice"`
			ReleaseDate    string  `json:"releaseDate"`
			PrimaryGenre   string  `json:"primaryGenreName"`
		} `json:"results"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error leyendo respuesta: %v", err)
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error decodificando JSON: %v", err)
		return nil, fmt.Errorf("error decodificando JSON: %w", err)
	}

	var songs []models.Song
	for _, item := range response.Results {
		if item.TrackName != "" { // Ignorar resultados sin nombre
			songs = append(songs, models.Song{
				Name:        item.TrackName,
				Artist:      item.ArtistName,
				Album:       item.CollectionName,
				Provider:    "iTunes",
				PreviewURL:  item.PreviewURL,
				ImageURL:    item.ArtworkUrl100,
				Price:       item.TrackPrice,
				ReleaseDate: item.ReleaseDate,
				Genre:       item.PrimaryGenre,
				Origin:      "iTunes",
			})
		}
	}

	log.Printf("iTunes encontró %d canciones", len(songs))
	return songs, nil
}

func (service *MusicService) searchChartLyrics(query string) ([]models.Song, error) {
	baseURL := "http://api.chartlyrics.com/apiv1.asmx/SearchLyric"

	// Dividir la consulta en artista y canción
	parts := strings.Split(query, "-")
	artist := strings.TrimSpace(query)
	song := query
	if len(parts) > 1 {
		artist = strings.TrimSpace(parts[0])
		song = strings.TrimSpace(parts[1])
	}

	url := fmt.Sprintf("%s?artist=%s&song=%s",
		baseURL,
		url.QueryEscape(artist),
		url.QueryEscape(song))

	log.Printf("url: %v", url)

	req, err := http.Get(url)

	if err != nil {
		log.Printf("Error creando request para ChartLyrics: %v", err)
		return nil, fmt.Errorf("error creando request para ChartLyrics: %w", err)
	}
	defer req.Body.Close()

	if req.StatusCode != http.StatusOK {
		log.Printf("Status error: %v", req.StatusCode)
		return nil, fmt.Errorf("status error: %v", req.StatusCode)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error leyendo respuesta: %v", err)
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Estructura para la respuesta XML
	var response models.ChartLyricsResponse

	if err := xml.Unmarshal(body, &response); err != nil {
		// Si hay error al decodificar pero tenemos respuesta, podría ser una respuesta vacía
		if strings.TrimSpace(string(body)) == "" {
			return []models.Song{}, nil
		}
		return nil, fmt.Errorf("error decodificando XML: %w", err)
	}

	log.Printf("response: %v", response.SearchLyricResult)

	var songs []models.Song
	for _, result := range response.SearchLyricResult {
		if result.Artist != "" {
			songs = append(songs, models.Song{
				Name:       result.Song,
				Artist:     result.Artist,
				Provider:   "ChartLyrics",
				PreviewURL: result.SongURL,
				Origin:     "ChartLyrics",
			})
		}
	}

	log.Printf("ChartLyrics encontró %d canciones", len(songs))
	return songs, nil
}

func (service *MusicService) SearchMusic(ctx context.Context, query string) ([]models.Song, error) {
	log.Printf("Iniciando búsqueda de música con query: %s", query)

	// songs, err := service.getFromCache(ctx, query)
	// log.Printf("redis songs: %d ", len(songs))
	// if err == nil && len(songs) > 0 {
	// 	return songs, nil
	// }

	log.Println("Realizando búsqueda en APIs externas")
	songs, err := service.searchInParallel(ctx, query)
	if err != nil {
		log.Printf("Error en búsqueda paralela: %v", err)
		return nil, fmt.Errorf("error en búsqueda paralela: %w", err)
	}

	if err := service.saveToCache(ctx, query, songs); err != nil {
		log.Printf("Error guardando resultados en caché: %v", err)
	}

	// if err := service.saveResults(ctx, query, songs); err != nil {
	// 	log.Printf("Error guardando resultados: %v", err)
	// }

	log.Printf("Búsqueda completada. Encontradas %d canciones", len(songs))
	return songs, nil
}

func (service *MusicService) searchInParallel(ctx context.Context, query string) ([]models.Song, error) {
	var wg sync.WaitGroup
	results := make(chan []models.Song, 3)
	errors := make(chan error, 3)

	for _, provider := range service.providers {
		wg.Add(1)
		go func(p models.ProviderInterface) {
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

func (service *MusicService) saveToCache(ctx context.Context, query string, songs []models.Song) error {
	cacheKey := "search:" + query
	cacheValue, err := json.Marshal(songs)
	if err != nil {
		return fmt.Errorf("error serializando resultados a JSON: %w", err)
	}
	return service.redisClient.Set(ctx, cacheKey, cacheValue, 0).Err()
}

// func (service *MusicService) saveResults(ctx context.Context, query string, songs []models.Song) error {
// 	if len(songs) == 0 {
// 		return nil
// 	}

// 	collection := service.mongoClient.Database("musicdb").Collection("searches")

// 	// Primero eliminamos búsquedas anteriores con el mismo query
// 	_, err := collection.DeleteMany(ctx, map[string]interface{}{
// 		"query": query,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("error eliminando búsquedas anteriores: %w", err)
// 	}

// 	// Guardamos la nueva búsqueda
// 	searchRecord := map[string]interface{}{
// 		"query":     query,
// 		"songs":     songs,
// 		"timestamp": time.Now(),
// 	}

// 	_, err = collection.InsertOne(ctx, searchRecord)
// 	if err != nil {
// 		return fmt.Errorf("error guardando nueva búsqueda: %w", err)
// 	}

// 	return nil
// }

type MusicServiceInterface interface {
	SearchMusic(ctx context.Context, query string) ([]models.Song, error)
}
