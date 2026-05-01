import React from 'react';
import './Splash.css';
import wingopherLogo from '../assets/wingopher.svg';

export const Splash: React.FC = () => {
    return (
        <div className="splash-container">
            <div className="gopher-container">
                <img src={wingopherLogo} alt="WinGopher Logo" className="wingopher-logo" />
                <div className="loader-ring"></div>
            </div>
            <h1 className="splash-title">WinGopher</h1>
            <p className="splash-text">Scanning your system for installed apps...</p>
        </div>
    );
};
