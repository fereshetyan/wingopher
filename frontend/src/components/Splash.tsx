import React from 'react';
import './Splash.css';
import wingoLogo from '../assets/wingo.svg';

export const Splash: React.FC = () => {
    return (
        <div className="splash-container">
            <div className="gopher-container">
                <img src={wingoLogo} alt="Wingo Logo" className="wingo-logo" />
                <div className="loader-ring"></div>
            </div>
            <h1 className="splash-title">Wingo</h1>
            <p className="splash-text">Scanning your system for installed apps...</p>
        </div>
    );
};
