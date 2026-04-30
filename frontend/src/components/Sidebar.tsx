import React from 'react';

interface SidebarProps {
    categories: string[];
    selectedCategory: string;
    onSelectCategory: (category: string) => void;
    onRefresh: () => void;
}

export const Sidebar: React.FC<SidebarProps> = React.memo(({ 
    categories, 
    selectedCategory, 
    onSelectCategory,
    onRefresh
}) => {
    return (
        <aside className="sidebar">
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
                    >
                        {category}
                    </li>
                ))}
            </ul>
        </aside>
    );
});

