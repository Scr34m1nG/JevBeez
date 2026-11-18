package search

import (
	"context"
	"jevbeez/api/internal/model"
)

type SearchProvider interface {
	Name() string
	Search(ctx context.Context, query string) ([]model.Evidence, error)
}
