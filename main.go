package main

import (
	"embed"
	"log/slog"
	"os"
	"wingopher/internal/app"
	"wingopher/internal/installer"
	"wingopher/internal/repository"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Initialize dependencies
	repo := repository.NewAppRepository()
	inst := installer.NewWingetInstaller()
	
	// Create an instance of the app structure
	a := app.NewApp(repo, inst)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "WinGopher",
		Width:  1200,
		Height: 800,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 255},
		OnStartup:        a.Startup,
		Bind: []any{
			a,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			Theme:                windows.Dark,
		},
	})

	if err != nil {
		slog.Error("Application failed to run", "error", err)
		os.Exit(1)
	}
}
