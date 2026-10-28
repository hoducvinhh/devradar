package discovery

import (
	"github.com/devradar/devradar/internal/classification"
	"github.com/devradar/devradar/internal/github"
	"github.com/devradar/devradar/internal/ranking"
	"github.com/devradar/devradar/internal/repository"
	"github.com/devradar/devradar/internal/storage"
)

type Service struct {
	client     *github.Client
	storage    *storage.Storage
	cache      *storage.Cache
	classifier *classification.Classifier
	analyzer   *repository.Analyzer
	ranker     *ranking.Ranker
}

func NewService(client *github.Client, store *storage.Storage, cache *storage.Cache) *Service {
	return &Service{
		client:     client,
		storage:    store,
		cache:      cache,
		classifier: classification.NewClassifier(),
		analyzer:   repository.NewAnalyzer(),
		ranker:     ranking.NewRanker(),
	}
}
