import { useState, useEffect, useMemo, useCallback } from 'react';
import { GetApps, InstallApps, UninstallApp, IsAdmin, CheckWinget, CheckInstalled, GetInstalledApps } from "../../wailsjs/go/app/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { AppData, InstallStatus } from '../types';

export function useInstallManager() {
    const [apps, setApps] = useState<AppData[]>([]);
    const [selectedApps, setSelectedApps] = useState<Set<string>>(new Set());
    const [search, setSearch] = useState('');
    const [selectedCategory, setSelectedCategory] = useState('All');
    const [statuses, setStatuses] = useState<Record<string, InstallStatus>>({});
    const [isAdmin, setIsAdmin] = useState(true);
    const [hasWinget, setHasWinget] = useState(true);
    const [installedApps, setInstalledApps] = useState<Set<string>>(new Set());
    const [isLoading, setIsLoading] = useState(true);

    const refreshInstalled = useCallback(async () => {
        const installed = await GetInstalledApps();
        setInstalledApps(new Set(installed));
    }, []);

    useEffect(() => {
        const init = async () => {
            setIsLoading(true);
            const admin = await IsAdmin();
            setIsAdmin(admin);
            const winget = await CheckWinget();
            setHasWinget(winget);
            const data = await GetApps();
            setApps(data);
            await refreshInstalled();
            setIsLoading(false);
        };

        init();

        const statusOff = EventsOn("install_status", (status: InstallStatus) => {
            setStatuses(prev => ({ ...prev, [status.id]: status }));
            
            if (status.status === 'completed') {
                if (status.message === 'Done') {
                    // Optimistically add to installed apps
                    setInstalledApps(prev => new Set(prev).add(status.id));
                } else if (status.message === 'Uninstalled') {
                    // Optimistically remove from installed apps
                    setInstalledApps(prev => {
                        const next = new Set(prev);
                        next.delete(status.id);
                        return next;
                    });
                }

                // Verify after a short delay to ensure winget database is updated
                setTimeout(() => {
                    CheckInstalled(status.id).then(isInstalled => {
                        setInstalledApps(prev => {
                            const next = new Set(prev);
                            if (isInstalled) next.add(status.id);
                            else next.delete(status.id);
                            return next;
                        });
                    });
                }, 3000);
            }
        });

        return () => {
            statusOff();
        };
    }, [refreshInstalled]);

    const categories = useMemo(() => {
        return ['All', ...Array.from(new Set(apps.map(a => a.category))).sort()];
    }, [apps]);

    const filteredApps = useMemo(() => {
        return apps.filter(app => {
            const matchesSearch = app.content.toLowerCase().includes(search.toLowerCase()) ||
                app.description.toLowerCase().includes(search.toLowerCase());
            const matchesCategory = selectedCategory === 'All' || app.category === selectedCategory;
            return matchesSearch && matchesCategory;
        }).sort((a, b) => a.content.localeCompare(b.content));
    }, [apps, search, selectedCategory]);

    const toggleApp = useCallback((id: string) => {
        setSelectedApps(prev => {
            const next = new Set(prev);
            if (next.has(id)) next.delete(id);
            else next.add(id);
            return next;
        });
    }, []);

    const install = useCallback(() => {
        InstallApps(Array.from(selectedApps));
    }, [selectedApps]);

    const uninstall = useCallback((id: string) => {
        UninstallApp(id);
    }, []);

    const selectAll = useCallback(() => {
        const allIds = filteredApps.map(app => app.id);
        setSelectedApps(new Set(allIds));
    }, [filteredApps]);

    const clearSelection = useCallback(() => {
        setSelectedApps(new Set());
    }, []);

    return {
        apps,
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
    };
}
