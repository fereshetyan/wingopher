import { useState, useCallback, useEffect } from 'react';
import './App.css';
import { useInstallManager } from './hooks/useInstallManager';
import { Sidebar } from './components/Sidebar';
import { Header } from './components/Header';
import { AppCard } from './components/AppCard';
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
        toggleApp,
        install,
        uninstall,
        selectAll,
        clearSelection,
        refreshInstalled,
        isLoading
    } = useInstallManager();

    const [showTerminal, setShowTerminal] = useState(false);
    const [confirmUninstall, setConfirmUninstall] = useState<{id: string, name: string} | null>(null);

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
        <div className="container">
            <Sidebar 
                categories={categories}
                selectedCategory={selectedCategory}
                onSelectCategory={setSelectedCategory}
                onRefresh={refreshInstalled}
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
                
                <div className="app-grid">
                    {filteredApps.map(app => (
                        <AppCard 
                            key={app.id}
                            app={app}
                            isSelected={selectedApps.has(app.id)}
                            isInstalled={installedApps.has(app.id)}
                            status={statuses[app.id]}
                            onToggle={toggleApp}
                            onUninstall={handleUninstallRequest}
                        />
                    ))}
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
