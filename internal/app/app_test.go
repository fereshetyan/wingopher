package app

import (
	"context"
	"testing"
	"wingopher/internal/models"
)

type mockRepository struct {
	apps map[string]models.AppData
}

func (m *mockRepository) GetAll() []models.AppData {
	var res []models.AppData
	for _, a := range m.apps {
		res = append(res, a)
	}
	return res
}

func (m *mockRepository) GetByID(id string) (models.AppData, bool) {
	a, ok := m.apps[id]
	return a, ok
}

type mockInstaller struct {
	admin    bool
	winget   bool
	installed map[string]bool
}

func (m *mockInstaller) IsAdmin() bool { return m.admin }
func (m *mockInstaller) CheckWinget() bool { return m.winget }
func (m *mockInstaller) Install(ctx context.Context, app models.AppData, onLog func(string)) (string, error) { return "installed", nil }
func (m *mockInstaller) Uninstall(ctx context.Context, app models.AppData, onLog func(string)) (string, error) { return "uninstalled", nil }
func (m *mockInstaller) Upgrade(ctx context.Context, app models.AppData, onLog func(string)) (string, error) { return "upgraded", nil }
func (m *mockInstaller) GetAppsWithUpdates(ctx context.Context) (map[string]bool, error) { return nil, nil }
func (m *mockInstaller) IsInstalled(id string) bool { return m.installed[id] }
func (m *mockInstaller) GetInstalledIDs(ctx context.Context) (map[string]bool, error) { return m.installed, nil }
func (m *mockInstaller) GetCleanErrorMessage(err error, output string) string { return "error" }

func TestApp_GetApps(t *testing.T) {
	repo := &mockRepository{
		apps: map[string]models.AppData{
			"test": {ID: "test", Winget: "Test.App"},
		},
	}
	inst := &mockInstaller{}
	a := NewApp(repo, inst)

	apps := a.GetApps()
	if len(apps) != 1 {
		t.Errorf("Expected 1 app, got %d", len(apps))
	}
	if apps[0].ID != "test" {
		t.Errorf("Expected app ID 'test', got '%s'", apps[0].ID)
	}
}

func TestApp_CheckInstalled(t *testing.T) {
	repo := &mockRepository{
		apps: map[string]models.AppData{
			"test": {ID: "test", Winget: "Test.App"},
		},
	}
	inst := &mockInstaller{
		installed: map[string]bool{"Test.App": true},
	}
	a := NewApp(repo, inst)

	if !a.CheckInstalled("test") {
		t.Error("Expected app to be installed")
	}
	if a.CheckInstalled("missing") {
		t.Error("Expected missing app to not be installed")
	}
}
