import React from "react";
import style from "./SettingsPanel.module.scss";

interface SettingsPanelProps {
  isOpen: boolean;
  onClose: () => void;
}

function SettingsPanel({ isOpen, onClose }: SettingsPanelProps) {
  if (!isOpen) return null;

  return (
    <>
      <div className={style.overlay} onClick={onClose} />
      <div className={style.panel}>
        <div className={style.header}>
          <h2>Настройки</h2>
          <button className={style.closeButton} onClick={onClose}>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path
                d="M18 6L6 18M6 6L18 18"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
          </button>
        </div>
        <div className={style.content}>
          <div className={style.infoMessage}>
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
              <path
                d="M12 15V3M9 12L12 15L15 12"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
              <path
                d="M20 18V20C20 21.1046 19.1046 22 18 22H6C4.89543 22 4 21.1046 4 20V18"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
            <p>Настройки приложения</p>
            <span>В разработке</span>
          </div>
        </div>
      </div>
    </>
  );
}

export default SettingsPanel;
