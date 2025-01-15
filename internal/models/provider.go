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

func (p Provider) Search(ctx context.Context, query string) ([]Song, error) {
	return p.SearchFn(ctx, query)
}
