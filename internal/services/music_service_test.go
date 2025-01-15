package services

import (
	"context"
	"testing"

	"github.com/JonayMedina/go-search-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProvider struct {
	mock.Mock
}

func (m *MockProvider) Search(ctx context.Context, query string) ([]models.Song, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]models.Song), args.Error(1)
}

func TestMusicService_SearchMusic(t *testing.T) {
	ctx := context.Background()
	mockProvider := new(MockProvider)

	expectedSongs := []models.Song{
		{
			ID:       "1",
			Name:     "Test Song",
			Artist:   "Test Artist",
			Duration: "3:30",
		},
	}

	mockProvider.On("Search", ctx, "test query").Return(expectedSongs, nil)

	service := &MusicService{
		providers: []models.ProviderInterface{mockProvider},
	}

	songs, err := service.SearchMusic(ctx, "test query")

	assert.NoError(t, err)
	assert.Equal(t, expectedSongs, songs)
	mockProvider.AssertExpectations(t)
}
