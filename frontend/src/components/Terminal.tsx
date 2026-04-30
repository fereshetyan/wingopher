import React, { useEffect, useRef } from 'react';
import { InstallStatus } from '../types';

interface TerminalProps {
    statuses: Record<string, InstallStatus>;
    onClose: () => void;
}

export const Terminal: React.FC<TerminalProps> = ({ statuses, onClose }) => {
    const terminalEndRef = useRef<HTMLDivElement>(null);

    const scrollToBottom = () => {
        terminalEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    useEffect(() => {
        scrollToBottom();
    }, [statuses]);

    const activeLogs = Object.values(statuses).filter(s => s.logs);

    return (
        <div className="terminal-panel">
            <div className="terminal-header">
                <div className="terminal-tabs">
                    <span className="terminal-tab active">System Logs</span>
                </div>
                <button className="terminal-close" onClick={onClose}>×</button>
            </div>
            <div className="terminal-body">
                {activeLogs.length === 0 ? (
                    <div className="terminal-empty">No active logs. Waiting for operations...</div>
                ) : (
                    activeLogs.map(status => (
                        <div key={status.id} className="terminal-group">
                            <div className="terminal-app-name">[{status.id}] - {status.status.toUpperCase()}</div>
                            <pre className="terminal-content">{status.logs}</pre>
                        </div>
                    ))
                )}
                <div ref={terminalEndRef} />
            </div>
        </div>
    );
};
