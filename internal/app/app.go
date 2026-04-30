package app

import (
	"context"
	"strings"
	"sync"
	"wingo/internal/installer"
	"wingo/internal/models"
	"wingo/internal/repository"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	repo      *repository.AppRepository
	installer *installer.WingetInstaller
}

func NewApp() *App {
	return &App{
		repo:      repository.NewAppRepository(),
		installer: installer.NewWingetInstaller(),
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
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			var currentLogs strings.Builder
			a.emitStatus(app.ID, "installing", "Installing...", "")

			output, err := a.installer.Install(app, func(line string) {
				currentLogs.WriteString(line + "\n")
				a.emitStatus(app.ID, "installing", "Installing...", currentLogs.String())
			})
			
			if err != nil {
				msg := a.installer.GetCleanErrorMessage(err, output)
				a.emitStatus(app.ID, "failed", msg, output)
			} else {
				a.emitStatus(app.ID, "completed", "Done", output)
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
	output, err := a.installer.Uninstall(app, func(line string) {
		currentLogs.WriteString(line + "\n")
		a.emitStatus(app.ID, "installing", "Uninstalling...", currentLogs.String())
	})

	if err != nil {
		msg := a.installer.GetCleanErrorMessage(err, output)
		a.emitStatus(id, "failed", msg, output)
	} else {
		a.emitStatus(id, "completed", "Uninstalled", output)
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
	installedMap, err := a.installer.GetInstalledIDs()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error getting installed IDs: %v", err)
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

	runtime.LogInfof(a.ctx, "Found %d installed apps", len(result))
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
