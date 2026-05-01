package models

type AppData struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	Choco       string `json:"choco"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Winget      string `json:"winget"`
	Foss        bool   `json:"foss"`
	IsSystemApp bool   `json:"is_system_app"`
}

type InstallStatus struct {
	ID      string `json:"id"`
	Status  string `json:"status"` // "pending", "installing", "completed", "failed"
	Message string `json:"message"`
	Logs    string `json:"logs"`
}
