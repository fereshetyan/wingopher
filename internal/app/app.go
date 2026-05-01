package app

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"wingopher/internal/installer"
	"wingopher/internal/models"
	"wingopher/internal/repository"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	repo      repository.Repository
	installer installer.Installer
}

func NewApp(repo repository.Repository, inst installer.Installer) *App {
	return &App{
		repo:      repo,
		installer: inst,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetApps() []models.AppData {
	return a.repo.GetAll()
}

func (a *App) IsAdmin() bool {
	return a.installer.IsAdmin()
}

func (a *App) CheckWinget() bool {
	return a.installer.CheckWinget()
}

func (a *App) InstallApps(ids []string) {
	if !a.installer.CheckWinget() {
		for _, id := range ids {
			a.emitStatus(id, "failed", "Winget not found", "Please install Winget to use this application.")
		}
		return
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3)

	for _, id := range ids {
		app, ok := a.repo.GetByID(id)
		if !ok || app.Winget == "" || app.Winget == "na" {
			a.emitStatus(id, "failed", "Invalid Winget ID", "")
			continue
		}

		wg.Add(1)
		go func(app models.AppData) {
			defer wg.Done()
			
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-a.ctx.Done():
				a.emitStatus(app.ID, "failed", "Cancelled", "Installation cancelled.")
				return
			}

			var currentLogs strings.Builder
			a.emitStatus(app.ID, "installing", "Installing...", "")

			output, err := a.installer.Install(a.ctx, app, func(line string) {
				currentLogs.WriteString(line + "\n")
				a.emitStatus(app.ID, "installing", "Installing...", currentLogs.String())
			})

			if err != nil {
				msg := a.installer.GetCleanErrorMessage(err, output)
				a.emitStatus(app.ID, "failed", msg, output)
				slog.Error("Failed to install app", "id", app.ID, "error", err)
			} else {
				a.emitStatus(app.ID, "completed", "Done", output)
				slog.Info("Successfully installed app", "id", app.ID)
			}
		}(app)
	}

	go func() {
		wg.Wait()
		runtime.EventsEmit(a.ctx, "all_installations_finished", true)
	}()
}

func (a *App) UninstallApp(id string) {
	app, ok := a.repo.GetByID(id)
	if !ok || app.Winget == "" || app.Winget == "na" {
		a.emitStatus(id, "failed", "Invalid Winget ID", "")
		return
	}

	a.emitStatus(id, "installing", "Uninstalling...", "")

	var currentLogs strings.Builder
	output, err := a.installer.Uninstall(a.ctx, app, func(line string) {
		currentLogs.WriteString(line + "\n")
		a.emitStatus(app.ID, "installing", "Uninstalling...", currentLogs.String())
	})

	if err != nil {
		msg := a.installer.GetCleanErrorMessage(err, output)
		a.emitStatus(id, "failed", msg, output)
		slog.Error("Failed to uninstall app", "id", id, "error", err)
	} else {
		a.emitStatus(id, "completed", "Uninstalled", output)
		slog.Info("Successfully uninstalled app", "id", id)
	}
}

func (a *App) UpgradeApp(id string) {
	app, ok := a.repo.GetByID(id)
	if !ok || app.Winget == "" || app.Winget == "na" {
		a.emitStatus(id, "failed", "Invalid Winget ID", "")
		return
	}

	a.emitStatus(id, "installing", "Updating...", "")

	var currentLogs strings.Builder
	output, err := a.installer.Upgrade(a.ctx, app, func(line string) {
		currentLogs.WriteString(line + "\n")
		a.emitStatus(app.ID, "installing", "Updating...", currentLogs.String())
	})

	if err != nil {
		msg := a.installer.GetCleanErrorMessage(err, output)
		a.emitStatus(id, "failed", msg, output)
		slog.Error("Failed to upgrade app", "id", id, "error", err)
	} else {
		a.emitStatus(id, "completed", "Updated", output)
		slog.Info("Successfully upgraded app", "id", id)
	}
}

func (a *App) CheckInstalled(id string) bool {
	app, ok := a.repo.GetByID(id)
	if !ok {
		return false
	}
	return a.installer.IsInstalled(app.Winget)
}

func (a *App) GetInstalledApps() []string {
	installedMap, err := a.installer.GetInstalledIDs(a.ctx)
	if err != nil {
		slog.Error("Error getting installed IDs", "error", err)
		return []string{}
	}

	apps := a.repo.GetAll()
	result := make([]string, 0)

	for _, app := range apps {
		if app.Winget == "" || app.Winget == "na" {
			continue
		}

		if installedMap[strings.ToLower(app.Winget)] {
			result = append(result, app.ID)
		}
	}

	slog.Info("Found installed apps", "count", len(result))
	return result
}

func (a *App) GetAppsWithUpdates() []string {
	updatesMap, err := a.installer.GetAppsWithUpdates(a.ctx)
	if err != nil {
		slog.Error("Error getting apps with updates", "error", err)
		return []string{}
	}

	apps := a.repo.GetAll()
	result := make([]string, 0)

	for _, app := range apps {
		if app.Winget == "" || app.Winget == "na" {
			continue
		}

		if updatesMap[strings.ToLower(app.Winget)] {
			result = append(result, app.ID)
		}
	}

	slog.Info("Found apps with updates", "count", len(result))
	return result
}

func (a *App) emitStatus(id, status, message, logs string) {
	runtime.EventsEmit(a.ctx, "install_status", models.InstallStatus{
		ID:      id,
		Status:  status,
		Message: message,
		Logs:    logs,
	})
}
