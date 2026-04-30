package repository

import (
	"embed"
	"encoding/json"
	"fmt"
	"wingo/internal/models"
)

//go:embed apps.json
var appsFS embed.FS

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
		fmt.Printf("Error reading embedded apps.json: %v\n", err)
		return
	}

	var rawApps map[string]models.AppData
	if err := json.Unmarshal(data, &rawApps); err != nil {
		fmt.Printf("Error parsing apps.json: %v\n", err)
		return
	}

	for id, app := range rawApps {
		app.ID = id
		r.apps[id] = app
	}
}

func (r *AppRepository) GetAll() []models.AppData {
	var result []models.AppData
	for _, app := range r.apps {
		result = append(result, app)
	}
	return result
}

func (r *AppRepository) GetByID(id string) (models.AppData, bool) {
	app, ok := r.apps[id]
	return app, ok
}
