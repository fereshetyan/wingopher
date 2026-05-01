package repository

import (
	"embed"
	"encoding/json"
	"log/slog"
	"sort"
	"wingopher/internal/models"
)

//go:embed apps.json
var appsFS embed.FS

type Repository interface {
	GetAll() []models.AppData
	GetByID(id string) (models.AppData, bool)
}

type AppRepository struct {
	apps map[string]models.AppData
}

func NewAppRepository() *AppRepository {
	repo := &AppRepository{
		apps: make(map[string]models.AppData),
	}
	repo.loadApps()
	return repo
}

func (r *AppRepository) loadApps() {
	data, err := appsFS.ReadFile("apps.json")
	if err != nil {
		slog.Error("Failed to read embedded apps.json", "error", err)
		return
	}

	var rawApps map[string]models.AppData
	if err := json.Unmarshal(data, &rawApps); err != nil {
		slog.Error("Failed to parse apps.json", "error", err)
		return
	}

	for id, app := range rawApps {
		app.ID = id
		r.apps[id] = app
	}
}

func (r *AppRepository) GetAll() []models.AppData {
	result := make([]models.AppData, 0, len(r.apps))
	for _, app := range r.apps {
		result = append(result, app)
	}
	
	// Sort by ID for deterministic output
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	
	return result
}

func (r *AppRepository) GetByID(id string) (models.AppData, bool) {
	app, ok := r.apps[id]
	return app, ok
}
