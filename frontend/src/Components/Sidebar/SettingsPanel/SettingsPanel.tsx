import React, { useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  setTheme,
  setFontSize,
  setDesktopNotifications,
  setMessageSound,
  setSoundVolume,
  setShowOnlineStatus,
  setReadReceipts,
  setTypingIndicator,
  resetSettings,
  selectAllSettings,
  type ThemeType,
  type FontSizeType,
} from "../../../store/settingSlice";
import type { AppDispatch } from "../../../store/store";
import style from "./SettingsPanel.module.scss";

interface SettingsPanelProps {
  isOpen: boolean;
  onClose: () => void;
}

function SettingsPanel({ isOpen, onClose }: SettingsPanelProps) {
  const dispatch = useDispatch<AppDispatch>();
  const settings = useSelector(selectAllSettings);
  const [activeTab, setActiveTab] = useState<
    "appearance" | "notifications" | "privacy"
  >("appearance");

  const handleThemeChange = (newTheme: ThemeType) => {
    dispatch(setTheme(newTheme));
  };

  const handleFontSizeChange = (size: FontSizeType) => {
    dispatch(setFontSize(size));
  };

  const handleSettingChange = (key: string, value: boolean) => {
    switch (key) {
      case "desktopNotifications":
        dispatch(setDesktopNotifications(value));
        break;
      case "messageSound":
        dispatch(setMessageSound(value));
        break;
      case "showOnlineStatus":
        dispatch(setShowOnlineStatus(value));
        break;
      case "readReceipts":
        dispatch(setReadReceipts(value));
        break;
      case "typingIndicator":
        dispatch(setTypingIndicator(value));
        break;
    }
  };

  const handleResetSettings = () => {
    dispatch(resetSettings());
  };

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

        <div className={style.tabs}>
          <button
            className={`${style.tab} ${activeTab === "appearance" ? style.active : ""}`}
            onClick={() => setActiveTab("appearance")}
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <circle
                cx="12"
                cy="12"
                r="3"
                stroke="currentColor"
                strokeWidth="2"
              />
              <path
                d="M19.4 15a1.65 1.65 0 0 0 .33-1.82 8 8 0 0 0-14.46 0A1.65 1.65 0 0 0 4.6 15"
                stroke="currentColor"
                strokeWidth="2"
              />
            </svg>
            <span>Оформление</span>
          </button>
          <button
            className={`${style.tab} ${activeTab === "notifications" ? style.active : ""}`}
            onClick={() => setActiveTab("notifications")}
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path
                d="M12 2C10 2 8.3 3.3 8 5.3C5.2 6.1 3 8.7 3 11.7C3 15.3 4.5 17.5 6 18.5V20C6 21.1 6.9 22 8 22H16C17.1 22 18 21.1 18 20V18.5C19.5 17.5 21 15.3 21 11.7C21 8.7 18.8 6.1 16 5.3C15.7 3.3 14 2 12 2Z"
                stroke="currentColor"
                strokeWidth="2"
              />
            </svg>
            <span>Уведомления</span>
          </button>
          <button
            className={`${style.tab} ${activeTab === "privacy" ? style.active : ""}`}
            onClick={() => setActiveTab("privacy")}
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <rect
                x="3"
                y="11"
                width="18"
                height="11"
                rx="2"
                stroke="currentColor"
                strokeWidth="2"
              />
              <path
                d="M7 11V7C7 4.23858 9.23858 2 12 2C14.7614 2 17 4.23858 17 7V11"
                stroke="currentColor"
                strokeWidth="2"
              />
            </svg>
            <span>Приватность</span>
          </button>
        </div>

        <div className={style.content}>
          {activeTab === "appearance" && (
            <div className={style.section}>
              <h3>Тема</h3>
              <div className={style.themeSelector}>
                <button
                  className={`${style.themeOption} ${settings.theme === "dark" ? style.active : ""}`}
                  onClick={() => handleThemeChange("dark")}
                >
                  <div
                    className={style.themePreview}
                    style={{ background: "#1a1a1a" }}
                  >
                    <div
                      style={{
                        background: "#333",
                        width: "100%",
                        height: "8px",
                        borderRadius: "4px",
                      }}
                    />
                    <div
                      style={{
                        background: "#f5a623",
                        width: "70%",
                        height: "8px",
                        borderRadius: "4px",
                        marginTop: "4px",
                      }}
                    />
                  </div>
                  <span>Темная</span>
                </button>
                <button
                  className={`${style.themeOption} ${settings.theme === "light" ? style.active : ""}`}
                  onClick={() => handleThemeChange("light")}
                >
                  <div
                    className={style.themePreview}
                    style={{ background: "#ffffff" }}
                  >
                    <div
                      style={{
                        background: "#e0e0e0",
                        width: "100%",
                        height: "8px",
                        borderRadius: "4px",
                      }}
                    />
                    <div
                      style={{
                        background: "#f5a623",
                        width: "70%",
                        height: "8px",
                        borderRadius: "4px",
                        marginTop: "4px",
                      }}
                    />
                  </div>
                  <span>Светлая</span>
                </button>
              </div>

              <h3>Размер шрифта</h3>
              <div className={style.fontSizeSelector}>
                <button
                  className={`${style.fontOption} ${settings.fontSize === "small" ? style.active : ""}`}
                  onClick={() => handleFontSizeChange("small")}
                >
                  A
                </button>
                <button
                  className={`${style.fontOption} ${settings.fontSize === "medium" ? style.active : ""}`}
                  onClick={() => handleFontSizeChange("medium")}
                >
                  A
                </button>
                <button
                  className={`${style.fontOption} ${settings.fontSize === "large" ? style.active : ""}`}
                  onClick={() => handleFontSizeChange("large")}
                >
                  A
                </button>
              </div>

              <button
                className={style.resetButton}
                onClick={handleResetSettings}
                style={{ marginTop: "32px" }}
              >
                Сбросить настройки
              </button>
            </div>
          )}

          {activeTab === "notifications" && (
            <div className={style.section}>
              <div className={style.settingItem}>
                <div className={style.settingInfo}>
                  <h4>Уведомления на рабочем столе</h4>
                  <p>Показывать уведомления когда приложение свернуто</p>
                </div>
                <label className={style.switch}>
                  <input
                    type="checkbox"
                    checked={settings.desktopNotifications}
                    onChange={(e) =>
                      handleSettingChange(
                        "desktopNotifications",
                        e.target.checked,
                      )
                    }
                  />
                  <span className={style.slider}></span>
                </label>
              </div>

              <div className={style.settingItem}>
                <div className={style.settingInfo}>
                  <h4>Звук сообщений</h4>
                  <p>Воспроизводить звук при получении нового сообщения</p>
                </div>
                <label className={style.switch}>
                  <input
                    type="checkbox"
                    checked={settings.messageSound}
                    onChange={(e) =>
                      handleSettingChange("messageSound", e.target.checked)
                    }
                  />
                  <span className={style.slider}></span>
                </label>
              </div>

              {settings.messageSound && (
                <div className={style.settingItem}>
                  <div className={style.settingInfo}>
                    <h4>Громкость</h4>
                  </div>
                  <input
                    type="range"
                    min="0"
                    max="100"
                    value={settings.soundVolume}
                    onChange={(e) =>
                      dispatch(setSoundVolume(parseInt(e.target.value)))
                    }
                    className={style.sliderInput}
                  />
                  <span
                    style={{ marginLeft: "12px", color: "var(--text-primary)" }}
                  >
                    {settings.soundVolume}%
                  </span>
                </div>
              )}
            </div>
          )}

          {activeTab === "privacy" && (
            <div className={style.section}>
              <div className={style.settingItem}>
                <div className={style.settingInfo}>
                  <h4>Показывать статус онлайн</h4>
                  <p>Другие пользователи будут видеть когда вы в сети</p>
                </div>
                <label className={style.switch}>
                  <input
                    type="checkbox"
                    checked={settings.showOnlineStatus}
                    onChange={(e) =>
                      handleSettingChange("showOnlineStatus", e.target.checked)
                    }
                  />
                  <span className={style.slider}></span>
                </label>
              </div>

              <div className={style.settingItem}>
                <div className={style.settingInfo}>
                  <h4>Отчеты о прочтении</h4>
                  <p>Отправители будут видеть когда вы прочитали сообщение</p>
                </div>
                <label className={style.switch}>
                  <input
                    type="checkbox"
                    checked={settings.readReceipts}
                    onChange={(e) =>
                      handleSettingChange("readReceipts", e.target.checked)
                    }
                  />
                  <span className={style.slider}></span>
                </label>
              </div>

              <div className={style.settingItem}>
                <div className={style.settingInfo}>
                  <h4>Индикатор набора текста</h4>
                  <p>Показывать когда вы печатаете сообщение</p>
                </div>
                <label className={style.switch}>
                  <input
                    type="checkbox"
                    checked={settings.typingIndicator}
                    onChange={(e) =>
                      handleSettingChange("typingIndicator", e.target.checked)
                    }
                  />
                  <span className={style.slider}></span>
                </label>
              </div>
            </div>
          )}
        </div>
      </div>
    </>
  );
}

export default SettingsPanel;
