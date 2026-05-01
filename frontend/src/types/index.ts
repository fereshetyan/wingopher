export interface AppData {
    id: string;
    category: string;
    content: string;
    description: string;
    winget: string;
    link: string;
    is_system_app?: boolean;
}

export interface InstallStatus {
    id: string;
    status: string;
    message: string;
    logs: string;
}
