package installer

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"wingo/internal/models"
)

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

func (i *WingetInstaller) Install(app models.AppData, onLog func(string)) (string, error) {
	return i.runWingetCommand([]string{"install", "--id", app.Winget, "--exact", "--silent", "--accept-source-agreements", "--accept-package-agreements"}, onLog)
}

func (i *WingetInstaller) Uninstall(app models.AppData, onLog func(string)) (string, error) {
	return i.runWingetCommand([]string{"uninstall", "--id", app.Winget, "--exact", "--silent", "--accept-source-agreements"}, onLog)
}

func (i *WingetInstaller) IsInstalled(wingetID string) bool {
	if wingetID == "" || wingetID == "na" {
		return false
	}
	cmd := exec.Command("winget", "list", "--id", wingetID, "--exact")
	cmd.SysProcAttr = getSysProcAttr()
	err := cmd.Run()
	return err == nil
}

func (i *WingetInstaller) GetInstalledIDs() (map[string]bool, error) {
	// Use a very wide terminal width to minimize wrapping
	cmd := exec.Command("cmd", "/c", "set Microsoft.Winget.Settings_TermWidth=500 && chcp 65001 > nul && winget list --accept-source-agreements")
	cmd.SysProcAttr = getSysProcAttr()
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) == 0 {
			return nil, err
		}
	}

	installed := make(map[string]bool)
	outputStr := string(output)
	
	// Remove carriage returns and split into lines
	outputStr = strings.ReplaceAll(outputStr, "\r", "")
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "-") {
			continue
		}

		// Look for common ID patterns or separators
		// We use Fields to get parts of the line.
		// Usually the ID is the second column, but it might be wrapped.
		fields := strings.Fields(trimmed)
		if len(fields) >= 2 {
			// In many cases, the ID is the second field
			// We'll collect all fields that look like IDs (contain dots or are specific known IDs)
			for _, field := range fields {
				if strings.Contains(field, ".") || strings.Contains(field, "\\") {
					installed[strings.ToLower(field)] = true
				}
			}
		}
	}

	return installed, nil
}

func (i *WingetInstaller) runWingetCommand(args []string, onLog func(string)) (string, error) {
	// Join args and escape if necessary, but here we just need to prepend chcp
	fullArgs := append([]string{"/c", "chcp 65001 > nul && winget"}, args...)
	cmd := exec.Command("cmd", fullArgs...)
	cmd.SysProcAttr = getSysProcAttr()

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	multi := io.MultiReader(stdout, stderr)

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var fullLog strings.Builder
	scanner := bufio.NewScanner(multi)
	for scanner.Scan() {
		line := scanner.Text()

		// If line contains carriage returns, winget is likely updating progress.
		// We want the last version of the line.
		if strings.Contains(line, "\r") {
			parts := strings.Split(line, "\r")
			for j := len(parts) - 1; j >= 0; j-- {
				if trimmed := strings.TrimSpace(parts[j]); trimmed != "" {
					line = parts[j]
					break
				}
			}
		}

		// Clean up the line: remove remaining control characters and trim
		line = strings.ReplaceAll(line, "\r", "")
		trimmed := strings.TrimSpace(line)

		// Filter out winget spinner noise and empty/space-only lines
		isSpinner := trimmed == "-" || trimmed == "\\" || trimmed == "|" || trimmed == "/"
		if isSpinner || trimmed == "" {
			continue
		}

		fullLog.WriteString(line + "\n")
		if onLog != nil {
			onLog(line)
		}
	}

	err := cmd.Wait()
	return fullLog.String(), err
}

func (i *WingetInstaller) GetCleanErrorMessage(err error, output string) string {
	if strings.Contains(output, "0x80070005") {
		return "Access Denied (Needs Admin)"
	}
	if strings.Contains(output, "0x80070643") {
		return "Installation failed (Error 0x80070643)"
	}
	return "Operation failed"
}
