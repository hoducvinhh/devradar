package app

import (
	"os"
	"time"

	"github.com/devradar/devradar/internal/discovery"
	"github.com/devradar/devradar/internal/github"
	"github.com/devradar/devradar/internal/profile"
	"github.com/devradar/devradar/internal/storage"
)

// App holds all dependencies and shared application services.
type App struct {
	Storage   *storage.Storage
	Cache     *storage.Cache
	Config    *storage.Config
	Client    *github.Client
	Discovery *discovery.Service
	Profile   *profile.Manager
}

// NewApp initializes application storage, configurations, and services.
func NewApp(overrideToken string) (*App, error) {
	store, err := storage.NewStorage()
	if err != nil {
		return nil, err
	}

	cfg, err := store.LoadConfig()
	if err != nil {
		cfg = &storage.Config{}
	}

	// Resolve token: CLI flag > GITHUB_TOKEN env var > stored config token
	token := overrideToken
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		token = cfg.GitHubToken
	}

	cache := storage.NewCache(store.BaseDir(), 30*time.Minute)

	clientOpts := []github.Option{}
	if token != "" {
		clientOpts = append(clientOpts, github.WithToken(token))
	}

	client := github.NewClient(clientOpts...)
	discoverySvc := discovery.NewService(client, store, cache)
	profileMgr := profile.NewManager(store)

	return &App{
		Storage:   store,
		Cache:     cache,
		Config:    cfg,
		Client:    client,
		Discovery: discoverySvc,
		Profile:   profileMgr,
	}, nil
}
