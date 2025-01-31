package models

import (
	"context"
)

type ProviderType string

const (
	ProviderTypeJSON ProviderType = "json"
	ProviderTypeSOAP ProviderType = "soap"
)

type Provider struct {
	Name     string
	BaseURL  string
	Type     ProviderType
	SearchFn func(ctx context.Context, query string) ([]Song, error)
}

type SearchResponse struct {
	Songs []Song `json:"songs"`
}

type GetUsersResponse struct {
	Users []User `json:"users"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterResponse struct {
	Message string            `json:"message"`
	User    User              `json:"user"`
	Details map[string]string `json:"details,omitempty"`
}

func (p Provider) Search(ctx context.Context, query string) ([]Song, error) {
	return p.SearchFn(ctx, query)
}

type ProviderInterface interface {
	Search(ctx context.Context, query string) ([]Song, error)
}
