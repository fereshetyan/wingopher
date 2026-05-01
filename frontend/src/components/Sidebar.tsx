import React from 'react';

interface SidebarProps {
    categories: string[];
    selectedCategory: string;
    onSelectCategory: (category: string) => void;
    onRefresh: () => void;
    isCollapsed: boolean;
    onToggleCollapse: () => void;
}

const CATEGORY_ICONS: Record<string, string> = {
    'Utilities': '🛠️',
    'Document': '📄',
    'Pro Tools': '⚙️',
    'AI-Automation': '🤖',
    'Multimedia Tools': '🎬',
    'Development': '💻',
    'Microsoft Tools': '🪟',
    'Browsers': '🌐',
    'Communications': '💬',
    'Games': '🎮',
    'Office': '🏢',
    'System': '🖥️',
    'Security': '🛡️',
    'All': '📦',
    'Installed': '✅',
    'Updates': '🆙'
};

export const Sidebar: React.FC<SidebarProps> = React.memo(({ 
    categories, 
    selectedCategory, 
    onSelectCategory,
    onRefresh,
    isCollapsed,
    onToggleCollapse
}) => {
    return (
        <aside className={`sidebar ${isCollapsed ? 'collapsed' : ''}`}>
            <button 
                className="sidebar-toggle" 
                onClick={onToggleCollapse}
                title={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
            >
                {isCollapsed ? '→' : '←'}
            </button>

            <div className="sidebar-header">
                <h2>Categories</h2>
                <button className="refresh-mini-btn" onClick={onRefresh} title="Refresh installed apps">
                    🔄
                </button>
            </div>
            <ul>
                {categories.map(category => (
                    <li 
                        key={category}
                        className={selectedCategory === category ? 'active' : ''}
                        onClick={() => onSelectCategory(category)}
                        title={isCollapsed ? category : undefined}
                    >
                        <span className="sidebar-icon">
                            {CATEGORY_ICONS[category] || '📁'}
                        </span>
                        <span>{category}</span>
                    </li>
                ))}
            </ul>
        </aside>
    );
});

