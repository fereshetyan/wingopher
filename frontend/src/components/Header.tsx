import React from 'react';

interface HeaderProps {
    search: string;
    onSearchChange: (value: string) => void;
    selectedCount: number;
    onInstall: () => void;
    onSelectAll: () => void;
    onClearSelection: () => void;
    canInstall: boolean;
}

export const Header: React.FC<HeaderProps> = React.memo(({
    search,
    onSearchChange,
    selectedCount,
    onInstall,
    onSelectAll,
    onClearSelection,
    canInstall
}) => {
    return (
        <header>
            <div className="search-container">
                <input
                    type="text"
                    placeholder="Search apps by name or description..."
                    value={search}
                    onChange={e => onSearchChange(e.target.value)}
                />
                <div className="header-actions">
                    <button className="secondary-btn" onClick={onSelectAll}>Select All</button>
                    <button className="secondary-btn" onClick={onClearSelection}>Clear</button>
                </div>
            </div>
            <button
                disabled={selectedCount === 0 || !canInstall}
                onClick={onInstall}
                className="install-btn"
            >
                Install Selected ({selectedCount})
            </button>
        </header>
    );
});
