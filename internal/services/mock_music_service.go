package services

import (
	"context"

	"github.com/JonayMedina/go-search-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockMusicService struct {
	mock.Mock
}

func (m *MockMusicService) SearchMusic(ctx context.Context, query string) ([]models.Song, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]models.Song), args.Error(1)
}
