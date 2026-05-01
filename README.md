<div align="center">

<img src="assets/images/wingopher.png" width="400" alt="WinGopher Logo">

# WinGopher - Parallel Windows App Installer

A modern, high-performance application installer for Windows that provides a GUI for the winget package manager, built with Wails (Go + React/TS).

</div>

## Features

### Core Functionality
- **Parallel Installations**: Install multiple apps simultaneously with configurable concurrency (default: 3).
- **Winget Integration**: Full interaction with Windows Package Manager CLI.
- **Batch Operations**: Select and install multiple apps at once.
- **Uninstall Support**: Remove installed applications directly from the interface.
- **Real-time Log Streaming**: Live CLI output display during installations and uninstalls.

### User Interface
- **Modern Clean Design**: Professional and focused UI for efficient application management.
- **Dark Theme**: Consistent dark color scheme across all components.
- **200+ Applications**: Curated app database across 10+ categories (Browsers, Development, Games, etc.).
- **Category Filtering**: Browse apps by category with dynamic sidebar.
- **Real-time Search**: Filter apps by name or description instantly.
- **App Icons**: Automatic favicon retrieval with fallback placeholders.
- **Status Tracking**: Visual indicators for installing/completed/failed states.
- **Installed Detection**: Automatic detection and badging of already-installed apps.
- **Terminal Panel**: Toggleable log viewer with auto-scroll functionality.
- **Splash Screen**: Animated loading screen during system scanning.
- **Confirmation Modals**: Prevents accidental uninstalls.

### Technical Features
- **Single Binary**: Fully portable executable with all assets embedded via Go embed.
- **Admin Detection**: Warning prompts when running without elevated privileges.
- **Winget Verification**: Checks winget availability on startup.
- **Event-Driven Architecture**: Real-time communication between Go backend and React frontend.
- **Optimistic UI Updates**: Immediate feedback with background verification.
- **UTF-8 Support**: Proper encoding for international application names.
- **Silent Execution**: Hidden command windows during winget operations.
- **User-Friendly Errors**: Translated error codes with helpful messages.

## Project Structure

The project follows a scalable, modular architecture:

### Backend (Go)

- `internal/app`: Wails application handlers, lifecycle events, and method bindings.
- `internal/installer`: Winget CLI interaction, command execution, and log parsing.
- `internal/models`: Domain data structures (`AppData`, `InstallStatus`).
- `internal/repository`: Data management, `apps.json` loading, and app queries.

### Frontend (React/TS)

- `src/components`: Atomic UI components (AppCard, Sidebar, Header, Terminal, etc.).
- `src/hooks`: Custom React hooks for business logic (`useInstallManager`).
- `src/types`: Centralized TypeScript type definitions.
- `src/assets`: Static assets including fonts and SVG logos.

### Data

- `internal/repository/apps.json`: Curated database of 200+ applications with winget IDs, categories, and metadata.

## Prerequisites

- **Windows 10/11** with [Winget](https://learn.microsoft.com/en-us/windows/package-manager/) installed.
- [Go 1.23+](https://go.dev/dl/) (for development).
- [Node.js](https://nodejs.org/) (for development).
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) (for development/build).

## Development

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Run with hot-reload
wails dev

# Build production executable
wails build
```

## License

MIT
