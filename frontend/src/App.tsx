import { useState, useCallback, useEffect } from 'react';
import './App.css';
import { useInstallManager } from './hooks/useInstallManager';
import { Sidebar } from './components/Sidebar';
import { Header } from './components/Header';
import { AppListItem } from './components/AppListItem';
import { Terminal } from './components/Terminal';
import { Splash } from './components/Splash';
import { ConfirmModal } from './components/ConfirmModal';

function App() {
    const {
        filteredApps,
        selectedApps,
        categories,
        search,
        setSearch,
        selectedCategory,
        setSelectedCategory,
        statuses,
        isAdmin,
        hasWinget,
        installedApps,
        appsWithUpdates,
        toggleApp,
        install,
        uninstall,
        upgrade,
        selectAll,
        clearSelection,
        refreshInstalled,
        isLoading
    } = useInstallManager();

    const [showTerminal, setShowTerminal] = useState(false);
    const [confirmUninstall, setConfirmUninstall] = useState<{id: string, name: string} | null>(null);
    const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(() => {
        return localStorage.getItem('sidebar_collapsed') === 'true';
    });

    const toggleSidebar = useCallback(() => {
        setIsSidebarCollapsed(prev => {
            const newValue = !prev;
            localStorage.setItem('sidebar_collapsed', String(newValue));
            return newValue;
        });
    }, []);

    useEffect(() => {
        const handleContextMenu = (e: MouseEvent) => {
            e.preventDefault();
        };
        document.addEventListener('contextmenu', handleContextMenu);
        return () => {
            document.removeEventListener('contextmenu', handleContextMenu);
        };
    }, []);

    const handleInstall = useCallback(() => {
        install();
        setShowTerminal(true);
    }, [install]);

    const handleUpgrade = useCallback((id: string) => {
        upgrade(id);
        setShowTerminal(true);
    }, [upgrade]);

    const handleUninstallRequest = useCallback((id: string) => {
        const app = filteredApps.find(a => a.id === id);
        if (app) {
            setConfirmUninstall({id, name: app.content});
        }
    }, [filteredApps]);

    const confirmAction = useCallback(() => {
        if (confirmUninstall) {
            uninstall(confirmUninstall.id);
            setConfirmUninstall(null);
            setShowTerminal(true);
        }
    }, [confirmUninstall, uninstall]);

    if (isLoading) {
        return <Splash />;
    }

    return (
        <div className={`container ${isSidebarCollapsed ? 'sidebar-collapsed' : ''}`}>
            <Sidebar 
                categories={categories}
                selectedCategory={selectedCategory}
                onSelectCategory={setSelectedCategory}
                onRefresh={refreshInstalled}
                isCollapsed={isSidebarCollapsed}
                onToggleCollapse={toggleSidebar}
            />
            
            <main className="content">
                {!isAdmin && (
                    <div className="alert warning">
                        ⚠️ Running without Admin rights. Some apps may fail to install.
                    </div>
                )}
                {!hasWinget && (
                    <div className="alert error">
                        ❌ Winget not found! Please install App Installer from Microsoft Store.
                    </div>
                )}
                
                <Header 
                    search={search}
                    onSearchChange={setSearch}
                    selectedCount={selectedApps.size}
                    onInstall={handleInstall}
                    onSelectAll={selectAll}
                    onClearSelection={clearSelection}
                    canInstall={hasWinget}
                />
                
                <div className="app-list-header">
                    <div className="col-selection"></div>
                    <div className="col-name">Name</div>
                    <div className="col-description">Description</div>
                    <div className="col-status">Status</div>
                    <div className="col-actions">Actions</div>
                </div>
                
                <div className="app-list">
                    {filteredApps.length > 0 ? (
                        filteredApps.map(app => (
                            <AppListItem 
                                key={app.id}
                                app={app}
                                isSelected={selectedApps.has(app.id)}
                                isInstalled={installedApps.has(app.id)}
                                hasUpdate={appsWithUpdates.has(app.id)}
                                status={statuses[app.id]}
                                onToggle={toggleApp}
                                onUninstall={handleUninstallRequest}
                                onUpgrade={handleUpgrade}
                            />
                        ))
                    ) : (
                        <div className="empty-state">
                            <div className="empty-state-icon">🔍</div>
                            <p>No applications found matching your criteria.</p>
                            <button className="secondary-btn" onClick={() => {setSearch(''); setSelectedCategory('All');}}>
                                Clear all filters
                            </button>
                        </div>
                    )}
                </div>

                <button 
                    className={`terminal-toggle-btn ${showTerminal ? 'active' : ''}`}
                    onClick={() => setShowTerminal(!showTerminal)}
                >
                    {showTerminal ? 'Hide Terminal' : 'Show Terminal'}
                    {Object.values(statuses).some(s => s.status === 'installing') && <span className="pulse-dot"></span>}
                </button>

                {showTerminal && (
                    <Terminal 
                        statuses={statuses} 
                        onClose={() => setShowTerminal(false)} 
                    />
                )}

                {confirmUninstall && (
                    <ConfirmModal 
                        title="Uninstall Application"
                        message={`Are you sure you want to uninstall ${confirmUninstall.name}?`}
                        onConfirm={confirmAction}
                        onCancel={() => setConfirmUninstall(null)}
                    />
                )}
            </main>
        </div>
    );
}

export default App;
