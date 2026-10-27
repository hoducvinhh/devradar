package profile

import (
	"github.com/devradar/devradar/internal/storage"
	"github.com/devradar/devradar/pkg/models"
)

type Manager struct {
	store *storage.Storage
}

func NewManager(store *storage.Storage) *Manager {
	return &Manager{store: store}
}

func (m *Manager) Get() (*models.Profile, error) {
	return m.store.LoadProfile()
}
