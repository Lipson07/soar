import React from "react";
import style from "./SecurityPanel.module.scss";

interface SecurityPanelProps {
  isOpen: boolean;
  onClose: () => void;
}

function SecurityPanel({ isOpen, onClose }: SecurityPanelProps) {
  if (!isOpen) return null;

  return (
    <>
      <div className={style.overlay} onClick={onClose} />
      <div className={style.panel}>
        <div className={style.header}>
          <h2>Безопасность</h2>
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
                d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z"
                stroke="currentColor"
                strokeWidth="2"
              />
              <path
                d="M12 8V12M12 16H12.01"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
            <p>Настройки безопасности</p>
            <span>В разработке</span>
          </div>
        </div>
      </div>
    </>
  );
}

export default SecurityPanel;
