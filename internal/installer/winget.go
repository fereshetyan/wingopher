package installer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"wingopher/internal/models"
)

type Installer interface {
	IsAdmin() bool
	CheckWinget() bool
	Install(ctx context.Context, app models.AppData, onLog func(string)) (string, error)
	Uninstall(ctx context.Context, app models.AppData, onLog func(string)) (string, error)
	Upgrade(ctx context.Context, app models.AppData, onLog func(string)) (string, error)
	GetAppsWithUpdates(ctx context.Context) (map[string]bool, error)
	IsInstalled(ctx context.Context, wingetID string) bool
	GetInstalledIDs(ctx context.Context) (map[string]bool, error)
	GetCleanErrorMessage(err error, output string) string
}

type WingetInstaller struct{}

func NewWingetInstaller() *WingetInstaller {
	return &WingetInstaller{}
}

func (i *WingetInstaller) IsAdmin() bool {
	cmd := exec.Command("net", "session")
	cmd.SysProcAttr = getSysProcAttr()
	err := cmd.Run()
	return err == nil
}

func (i *WingetInstaller) CheckWinget() bool {
	cmd := exec.Command("winget", "--version")
	cmd.SysProcAttr = getSysProcAttr()
	err := cmd.Run()
	return err == nil
}

func (i *WingetInstaller) Install(ctx context.Context, app models.AppData, onLog func(string)) (string, error) {
	args := []string{"install", "--id", app.Winget, "--exact", "--silent", "--accept-source-agreements", "--accept-package-agreements"}
	if source := i.determineSource(app.Winget); source != "" {
		args = append(args, "--source", source)
	}
	return i.runWingetCommand(ctx, args, onLog)
}

func (i *WingetInstaller) Uninstall(ctx context.Context, app models.AppData, onLog func(string)) (string, error) {
	args := []string{"uninstall", "--id", app.Winget, "--exact", "--silent", "--accept-source-agreements"}
	if source := i.determineSource(app.Winget); source != "" {
		args = append(args, "--source", source)
	}
	return i.runWingetCommand(ctx, args, onLog)
}

func (i *WingetInstaller) Upgrade(ctx context.Context, app models.AppData, onLog func(string)) (string, error) {
	args := []string{"upgrade", "--id", app.Winget, "--exact", "--silent", "--accept-source-agreements"}
	if source := i.determineSource(app.Winget); source != "" {
		args = append(args, "--source", source)
	}
	return i.runWingetCommand(ctx, args, onLog)
}

func (i *WingetInstaller) GetAppsWithUpdates(ctx context.Context) (map[string]bool, error) {
	cmd := exec.CommandContext(ctx, "cmd", "/c", "set Microsoft.Winget.Settings_TermWidth=500 && chcp 65001 > nul && winget upgrade --accept-source-agreements")
	cmd.SysProcAttr = getSysProcAttr()

	output, err := cmd.CombinedOutput()
	if err != nil && len(output) == 0 {
		return nil, fmt.Errorf("winget upgrade: %w", err)
	}

	return i.parseWingetTable(string(output)), nil
}

func (i *WingetInstaller) determineSource(id string) string {
	if id == "" || id == "na" {
		return ""
	}
	// MS Store IDs are alphanumeric and don't contain dots (e.g. 9nt1r1c2hh7j)
	// Winget IDs usually follow Publisher.App format (contain dots)
	if strings.Contains(id, ".") {
		return "winget"
	}
	// If it doesn't contain a dot, it's likely an MS Store ID
	return "msstore"
}

func (i *WingetInstaller) IsInstalled(ctx context.Context, wingetID string) bool {
	if wingetID == "" || wingetID == "na" {
		return false
	}
	cmd := exec.CommandContext(ctx, "winget", "list", "--id", wingetID, "--exact")
	cmd.SysProcAttr = getSysProcAttr()
	err := cmd.Run()
	return err == nil
}

func (i *WingetInstaller) GetInstalledIDs(ctx context.Context) (map[string]bool, error) {
	cmd := exec.CommandContext(ctx, "cmd", "/c", "set Microsoft.Winget.Settings_TermWidth=500 && chcp 65001 > nul && winget list --accept-source-agreements")
	cmd.SysProcAttr = getSysProcAttr()

	output, err := cmd.CombinedOutput()
	if err != nil && len(output) == 0 {
		return nil, fmt.Errorf("winget list: %w", err)
	}

	return i.parseWingetTable(string(output)), nil
}

func (i *WingetInstaller) parseWingetTable(output string) map[string]bool {
	res := make(map[string]bool)

	// Remove carriage returns and split into lines
	output = strings.ReplaceAll(output, "\r", "")
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines, separators, and headers
		if trimmed == "" || strings.HasPrefix(trimmed, "-") || strings.Contains(trimmed, "Id") && strings.Contains(trimmed, "Version") {
			continue
		}

		// Look for common ID patterns or separators
		// We collect all fields and let the App layer match them against the repository
		fields := strings.Fields(trimmed)
		if len(fields) >= 2 {
			for _, field := range fields {
				// We lowercase for consistency, matching the App layer's check
				res[strings.ToLower(field)] = true
			}
		}
	}

	return res
}

func (i *WingetInstaller) runWingetCommand(ctx context.Context, args []string, onLog func(string)) (string, error) {
	fullArgs := append([]string{"/c", "chcp 65001 > nul && winget"}, args...)
	cmd := exec.CommandContext(ctx, "cmd", fullArgs...)
	cmd.SysProcAttr = getSysProcAttr()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("creating stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("creating stderr pipe: %w", err)
	}
	multi := io.MultiReader(stdout, stderr)

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var fullLog strings.Builder
	scanner := bufio.NewScanner(multi)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "\r") {
			parts := strings.Split(line, "\r")
			for j := len(parts) - 1; j >= 0; j-- {
				if trimmed := strings.TrimSpace(parts[j]); trimmed != "" {
					line = parts[j]
					break
				}
			}
		}

		line = strings.ReplaceAll(line, "\r", "")
		trimmed := strings.TrimSpace(line)

		isSpinner := trimmed == "-" || trimmed == "\\" || trimmed == "|" || trimmed == "/"
		if isSpinner || trimmed == "" {
			continue
		}

		fullLog.WriteString(line + "\n")
		if onLog != nil {
			onLog(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fullLog.String(), fmt.Errorf("reading output: %w", err)
	}

	err = cmd.Wait()
	return fullLog.String(), err
}

func (i *WingetInstaller) GetCleanErrorMessage(err error, output string) string {
	if strings.Contains(output, "0x80070005") {
		return "Access Denied (Needs Admin)"
	}
	if strings.Contains(output, "0x80070643") {
		return "Installation failed (Error 0x80070643)"
	}
	if strings.Contains(output, "install technology is different") {
		return "Installer mismatch (Manual update required)"
	}
	if strings.Contains(output, "A newer version was found, but the install technology is different") {
		return "Installer mismatch"
	}
	if strings.Contains(output, "0x80041010") {
		return "Application not found or mismatch"
	}
	if strings.Contains(output, "Another installation is already in progress") {
		return "Busy (Another install running)"
	}
	return "Operation failed"
}
