import React from 'react';
import { AppData, InstallStatus } from '../types';

interface AppCardProps {
    app: AppData;
    isSelected: boolean;
    isInstalled: boolean;
    status?: InstallStatus;
    onToggle: (id: string) => void;
    onUninstall: (id: string) => void;
}

export const AppCard: React.FC<AppCardProps> = React.memo(({ 
    app, 
    isSelected, 
    isInstalled,
    status, 
    onToggle, 
    onUninstall,
}) => {
    const getIconUrl = (app: AppData) => {
        if (!app.link || app.link === "na") return null;
        try {
            const url = new URL(app.link);
            return `https://www.google.com/s2/favicons?domain=${url.hostname}&sz=64`;
        } catch (e) {
            return null;
        }
    };

    const handleUninstall = (e: React.MouseEvent) => {
        e.stopPropagation();
        onUninstall(app.id);
    };

    return (
        <div 
            className={`app-card ${isSelected ? 'selected' : ''} ${isInstalled ? 'installed' : ''}`} 
            onClick={() => onToggle(app.id)}
        >
            <div className="app-card-selection">
                <input 
                    type="checkbox" 
                    checked={isSelected}
                    onChange={() => {}} 
                    className="custom-checkbox"
                />
            </div>

            <div className="app-header">
                <div className="app-icon-container">
                    {getIconUrl(app) ? (
                        <img 
                            src={getIconUrl(app) || ''} 
                            alt="" 
                            className="app-icon"
                            onError={(e) => (e.currentTarget.style.display = 'none')}
                        />
                    ) : (
                        <div className="app-icon-placeholder">
                            {app.content[0].toUpperCase()}
                        </div>
                    )}
                </div>
                <div className="app-title-container">
                    <h3>{app.content}</h3>
                    {isInstalled && <span className="installed-badge">Installed</span>}
                </div>
            </div>
            
            <div className="app-body">
                <p title={app.description}>{app.description}</p>
            </div>
            
            <div className="card-footer">
                <div className="status-container">
                    {status && (
                        <div className={`status-tag ${status.status}`}>
                            <span className="status-dot"></span>
                            {status.message}
                        </div>
                    )}
                </div>
                <div className="card-actions">
                    {isInstalled && (
                        <button 
                            className="uninstall-btn" 
                            onClick={handleUninstall}
                            disabled={status?.status === 'installing'}
                        >
                            Uninstall
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
});
