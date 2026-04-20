import React, { useState, useEffect } from "react";
import style from "./SecurityPanel.module.scss";

interface SecuritySettings {
  two_factor_enabled: boolean;
  biometric_enabled: boolean;
  end_to_end_encryption: boolean;
  screen_security: boolean;
  login_alerts: boolean;
}

interface UserSession {
  id: number;
  device_info: string;
  device_type: string;
  location: string;
  last_active: string;
  is_active: boolean;
  is_current: boolean;
}

interface TwoFactorSetup {
  secret: string;
  qr_code: string;
  backup_codes: string[];
}

interface SecurityPanelProps {
  isOpen: boolean;
  onClose: () => void;
}

const API_BASE = "http://localhost:8080/api/security";

function SecurityPanel({ isOpen, onClose }: SecurityPanelProps) {
  const [settings, setSettings] = useState<SecuritySettings>({
    two_factor_enabled: false,
    biometric_enabled: false,
    end_to_end_encryption: true,
    screen_security: true,
    login_alerts: true,
  });
  const [sessions, setSessions] = useState<UserSession[]>([]);
  const [loading, setLoading] = useState(false);
  const [twoFactorSetup, setTwoFactorSetup] = useState<TwoFactorSetup | null>(
    null,
  );
  const [verificationCode, setVerificationCode] = useState("");
  const [showTwoFactorModal, setShowTwoFactorModal] = useState(false);
  const [showClearDataModal, setShowClearDataModal] = useState(false);
  const [showLogoutModal, setShowLogoutModal] = useState(false);

  useEffect(() => {
    if (isOpen) {
      fetchSettings();
      fetchSessions();
    }
  }, [isOpen]);

  const fetchSettings = async () => {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${API_BASE}/settings`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setSettings(data);
      }
    } catch (error) {
      console.error("Failed to fetch settings:", error);
    }
  };

  const fetchSessions = async () => {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${API_BASE}/sessions`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setSessions(Array.isArray(data) ? data : []);
      } else {
        setSessions([]);
      }
    } catch (error) {
      console.error("Failed to fetch sessions:", error);
      setSessions([]);
    }
  };

  const updateSetting = async (key: keyof SecuritySettings, value: boolean) => {
    setSettings((prev) => ({ ...prev, [key]: value }));
    try {
      const token = localStorage.getItem("token");
      await fetch(`${API_BASE}/settings`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ [key]: value }),
      });
    } catch (error) {
      console.error("Failed to update setting:", error);
      setSettings((prev) => ({ ...prev, [key]: !value }));
    }
  };

  const setupTwoFactor = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${API_BASE}/2fa/setup`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setTwoFactorSetup(data);
        setShowTwoFactorModal(true);
      }
    } catch (error) {
      console.error("Failed to setup 2FA:", error);
    } finally {
      setLoading(false);
    }
  };

  const verifyTwoFactor = async () => {
    if (!twoFactorSetup || !verificationCode) return;

    setLoading(true);
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${API_BASE}/2fa/verify`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          code: verificationCode,
          secret: twoFactorSetup.secret,
        }),
      });
      if (response.ok) {
        setSettings((prev) => ({ ...prev, two_factor_enabled: true }));
        setShowTwoFactorModal(false);
        setTwoFactorSetup(null);
        setVerificationCode("");
      } else {
        alert("Неверный код подтверждения");
      }
    } catch (error) {
      console.error("Failed to verify 2FA:", error);
    } finally {
      setLoading(false);
    }
  };

  const disableTwoFactor = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("token");
      await fetch(`${API_BASE}/2fa`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      setSettings((prev) => ({ ...prev, two_factor_enabled: false }));
    } catch (error) {
      console.error("Failed to disable 2FA:", error);
    } finally {
      setLoading(false);
    }
  };

  const terminateSession = async (sessionId: number) => {
    try {
      const token = localStorage.getItem("token");
      await fetch(`${API_BASE}/sessions/${sessionId}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      setSessions((prev) => prev.filter((s) => s.id !== sessionId));
    } catch (error) {
      console.error("Failed to terminate session:", error);
    }
  };

  const terminateAllOtherSessions = async () => {
    try {
      const token = localStorage.getItem("token");
      await fetch(`${API_BASE}/sessions`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      fetchSessions();
    } catch (error) {
      console.error("Failed to terminate sessions:", error);
    }
  };

  const exportSecurityReport = async () => {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${API_BASE}/report`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = `security-report-${Date.now()}.json`;
        a.click();
        URL.revokeObjectURL(url);
      }
    } catch (error) {
      console.error("Failed to export report:", error);
    }
  };

  const clearAllData = () => {
    localStorage.clear();
    sessionStorage.clear();
    window.location.href = "/login";
  };

  const getDeviceIcon = (deviceType: string) => {
    switch (deviceType) {
      case "desktop":
        return (
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
            <rect
              x="2"
              y="3"
              width="20"
              height="14"
              rx="2"
              stroke="currentColor"
              strokeWidth="2"
            />
            <path
              d="M8 21H16M12 17V21"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
            />
          </svg>
        );
      case "mobile":
        return (
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
            <rect
              x="7"
              y="2"
              width="10"
              height="20"
              rx="2"
              stroke="currentColor"
              strokeWidth="2"
            />
            <path
              d="M12 18H12.01"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
            />
          </svg>
        );
      default:
        return (
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
            <circle
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="2"
            />
            <circle cx="12" cy="12" r="2" fill="currentColor" />
          </svg>
        );
    }
  };

  const formatLastActive = (lastActive: string) => {
    const date = new Date(lastActive);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) return "Сейчас";
    if (minutes < 60) return `${minutes} мин назад`;
    if (hours < 24) return `${hours} ч назад`;
    return `${days} дн назад`;
  };

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
          <div className={style.section}>
            <h3>Основные настройки</h3>

            <div className={style.settingItem}>
              <div className={style.settingInfo}>
                <div className={style.settingIcon}>
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
                    <circle cx="12" cy="16" r="2" fill="currentColor" />
                  </svg>
                </div>
                <div className={style.settingText}>
                  <h4>Двухфакторная аутентификация</h4>
                  <p>Дополнительная защита аккаунта</p>
                </div>
              </div>
              <label className={style.switch}>
                <input
                  type="checkbox"
                  checked={settings.two_factor_enabled}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setupTwoFactor();
                    } else {
                      disableTwoFactor();
                    }
                  }}
                />
                <span className={style.slider}></span>
              </label>
            </div>

            <div className={style.settingItem}>
              <div className={style.settingInfo}>
                <div className={style.settingIcon}>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path
                      d="M12 2L2 7L12 12L22 7L12 2Z"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinejoin="round"
                    />
                    <path
                      d="M2 17L12 22L22 17"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinejoin="round"
                    />
                    <path
                      d="M2 12L12 17L22 12"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinejoin="round"
                    />
                  </svg>
                </div>
                <div className={style.settingText}>
                  <h4>Сквозное шифрование</h4>
                  <p>Защита сообщений E2E</p>
                </div>
              </div>
              <label className={style.switch}>
                <input
                  type="checkbox"
                  checked={settings.end_to_end_encryption}
                  onChange={(e) =>
                    updateSetting("end_to_end_encryption", e.target.checked)
                  }
                />
                <span className={style.slider}></span>
              </label>
            </div>
          </div>

          <div className={style.section}>
            <h3>Конфиденциальность</h3>

            <div className={style.settingItem}>
              <div className={style.settingInfo}>
                <div className={style.settingIcon}>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path
                      d="M2 12H22M12 2V22"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                    />
                    <rect
                      x="4"
                      y="4"
                      width="16"
                      height="16"
                      rx="2"
                      stroke="currentColor"
                      strokeWidth="2"
                    />
                  </svg>
                </div>
                <div className={style.settingText}>
                  <h4>Защита скриншотов</h4>
                  <p>Блокировка снимков экрана</p>
                </div>
              </div>
              <label className={style.switch}>
                <input
                  type="checkbox"
                  checked={settings.screen_security}
                  onChange={(e) =>
                    updateSetting("screen_security", e.target.checked)
                  }
                />
                <span className={style.slider}></span>
              </label>
            </div>

            <div className={style.settingItem}>
              <div className={style.settingInfo}>
                <div className={style.settingIcon}>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path
                      d="M18 8C18 4.68629 15.3137 2 12 2C8.68629 2 6 4.68629 6 8V11.1"
                      stroke="currentColor"
                      strokeWidth="2"
                    />
                    <path
                      d="M22 12L2 12V18C2 19.1046 2.89543 20 4 20H20C21.1046 20 22 19.1046 22 18V12Z"
                      stroke="currentColor"
                      strokeWidth="2"
                    />
                    <circle cx="12" cy="16" r="1" fill="currentColor" />
                  </svg>
                </div>
                <div className={style.settingText}>
                  <h4>Уведомления о входе</h4>
                  <p>Оповещения о новых сессиях</p>
                </div>
              </div>
              <label className={style.switch}>
                <input
                  type="checkbox"
                  checked={settings.login_alerts}
                  onChange={(e) =>
                    updateSetting("login_alerts", e.target.checked)
                  }
                />
                <span className={style.slider}></span>
              </label>
            </div>
          </div>

          <div className={style.section}>
            <h3>Активные сессии</h3>

            <div className={style.activeSessions}>
              {sessions.length > 0 ? (
                sessions.map((session) => (
                  <div
                    key={session.id}
                    className={`${style.sessionItem} ${session.is_current ? style.current : ""}`}
                  >
                    <div className={style.sessionIcon}>
                      {getDeviceIcon(session.device_type)}
                    </div>
                    <div className={style.sessionInfo}>
                      <div className={style.sessionName}>
                        {session.device_info || "Неизвестное устройство"}
                      </div>
                      <div className={style.sessionMeta}>
                        {session.location || "Неизвестно"} •{" "}
                        {session.last_active
                          ? formatLastActive(session.last_active)
                          : "Недавно"}
                      </div>
                    </div>
                    {session.is_current ? (
                      <div className={style.sessionBadge}>Текущий</div>
                    ) : (
                      <button
                        className={style.button}
                        onClick={() => terminateSession(session.id)}
                      >
                        Завершить
                      </button>
                    )}
                  </div>
                ))
              ) : (
                <div className={style.emptySessions}>
                  <p>Нет активных сессий</p>
                </div>
              )}
            </div>

            {sessions.length > 0 && (
              <button
                className={style.button}
                onClick={terminateAllOtherSessions}
                style={{ marginTop: "12px", width: "100%" }}
              >
                Завершить все остальные сессии
              </button>
            )}
          </div>

          <div className={style.section}>
            <h3>Управление данными</h3>

            <div
              style={{ display: "flex", flexDirection: "column", gap: "8px" }}
            >
              <button className={style.button} onClick={exportSecurityReport}>
                Экспорт отчета безопасности
              </button>

              <button
                className={style.button}
                onClick={() => setShowClearDataModal(true)}
              >
                Очистить локальные данные
              </button>

              <button
                className={`${style.button} ${style.danger}`}
                onClick={() => setShowLogoutModal(true)}
              >
                Выйти на всех устройствах
              </button>
            </div>
          </div>
        </div>
      </div>

      {showTwoFactorModal && twoFactorSetup && (
        <>
          <div
            className={style.overlay}
            onClick={() => setShowTwoFactorModal(false)}
          />
          <div className={style.modal}>
            <h3>Настройка 2FA</h3>
            <p>Отсканируйте QR-код в приложении Google Authenticator</p>
            <div className={style.qrCode}>
              <img
                src={`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(twoFactorSetup.qr_code)}`}
                alt="QR Code"
              />
            </div>
            <p>
              Или введите секрет вручную:{" "}
              <strong>{twoFactorSetup.secret}</strong>
            </p>
            <input
              type="text"
              className={style.input}
              placeholder="Введите 6-значный код"
              value={verificationCode}
              onChange={(e) => setVerificationCode(e.target.value)}
              maxLength={6}
            />
            <div className={style.modalActions}>
              <button
                className={style.button}
                onClick={() => setShowTwoFactorModal(false)}
              >
                Отмена
              </button>
              <button
                className={style.button}
                onClick={verifyTwoFactor}
                disabled={loading}
              >
                {loading ? "Проверка..." : "Подтвердить"}
              </button>
            </div>
          </div>
        </>
      )}

      {showClearDataModal && (
        <>
          <div
            className={style.overlay}
            onClick={() => setShowClearDataModal(false)}
          />
          <div className={style.modal}>
            <h3>Очистить данные?</h3>
            <p>
              Все локальные данные будут удалены. Вам потребуется войти заново.
            </p>
            <div className={style.modalActions}>
              <button
                className={style.button}
                onClick={() => setShowClearDataModal(false)}
              >
                Отмена
              </button>
              <button
                className={`${style.button} ${style.danger}`}
                onClick={clearAllData}
              >
                Очистить
              </button>
            </div>
          </div>
        </>
      )}

      {showLogoutModal && (
        <>
          <div
            className={style.overlay}
            onClick={() => setShowLogoutModal(false)}
          />
          <div className={style.modal}>
            <h3>Выйти на всех устройствах?</h3>
            <p>Вы будете выведены из аккаунта на всех устройствах.</p>
            <div className={style.modalActions}>
              <button
                className={style.button}
                onClick={() => setShowLogoutModal(false)}
              >
                Отмена
              </button>
              <button
                className={`${style.button} ${style.danger}`}
                onClick={terminateAllOtherSessions}
              >
                Выйти
              </button>
            </div>
          </div>
        </>
      )}
    </>
  );
}

export default SecurityPanel;
