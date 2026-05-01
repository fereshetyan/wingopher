import React from 'react';
import { AppData, InstallStatus } from '../types';

interface AppListItemProps {
    app: AppData;
    isSelected: boolean;
    isInstalled: boolean;
    hasUpdate?: boolean;
    status?: InstallStatus;
    onToggle: (id: string) => void;
    onUninstall: (id: string) => void;
    onUpgrade: (id: string) => void;
}

export const AppListItem: React.FC<AppListItemProps> = React.memo(({ 
    app, 
    isSelected, 
    isInstalled,
    hasUpdate,
    status, 
    onToggle, 
    onUninstall,
    onUpgrade,
}) => {
    const handleUninstall = (e: React.MouseEvent) => {
        e.stopPropagation();
        onUninstall(app.id);
    };

    const handleUpgrade = (e: React.MouseEvent) => {
        e.stopPropagation();
        onUpgrade(app.id);
    };

    return (
        <div 
            className={`app-list-item ${isSelected ? 'selected' : ''} ${isInstalled ? 'installed' : ''} ${hasUpdate ? 'has-update' : ''}`} 
            onClick={() => onToggle(app.id)}
        >
            <div className="col-selection">
                <input 
                    type="checkbox" 
                    checked={isSelected}
                    onChange={() => {}} 
                    className="custom-checkbox"
                />
            </div>

            <div className="col-name">
                <div className="app-name-wrapper">
                    <h3>{app.content}</h3>
                </div>
            </div>

            <div className="col-description">
                <p title={app.description}>{app.description}</p>
            </div>

            <div className="col-status">
                <div className="status-and-badge">
                    {status && (
                        <div className={`status-tag ${status.status}`}>
                            <span className="status-dot"></span>
                            {status.message}
                        </div>
                    )}
                    {isInstalled && <span className="installed-badge">Installed</span>}
                    {hasUpdate && <span className="update-badge">Update Available</span>}
                </div>
            </div>

            <div className="col-actions">
                <div className="action-buttons">
                    {hasUpdate && (
                        <button 
                            className="upgrade-btn-small" 
                            onClick={handleUpgrade}
                            disabled={status?.status === 'installing'}
                        >
                            Update
                        </button>
                    )}
                    {isInstalled && !app.is_system_app && (
                        <button 
                            className="uninstall-btn-small" 
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
