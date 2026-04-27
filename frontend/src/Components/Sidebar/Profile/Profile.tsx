import React, { useState, useRef } from "react";
import { useSelector, useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { selectUser, updateUser, logout } from "../../../store/userSlice";
import style from "./Profile.module.scss";

interface ProfileProps {
  isOpen: boolean;
  onClose: () => void;
}

function Profile({ isOpen, onClose }: ProfileProps) {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const user = useSelector(selectUser);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState({
    username: "",
    email: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [avatarLoading, setAvatarLoading] = useState(false);

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2);
  };

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return "Не указано";
    return new Date(dateString).toLocaleDateString("ru-RU", {
      day: "numeric",
      month: "long",
      year: "numeric",
    });
  };

  const handleLogout = () => {
    dispatch(logout());
    onClose();
    navigate("/login");
  };

  const handleEditClick = () => {
    if (user) {
      setEditForm({
        username: user.username,
        email: user.email,
      });
    }
    setIsEditing(true);
    setError(null);
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
    setError(null);
  };

  const handleSaveEdit = async () => {
    if (!user) return;

    setLoading(true);
    setError(null);

    try {
      const token = localStorage.getItem("token");

      const response = await fetch(
        `http://localhost:8080/api/users/${user.id}`,
        {
          method: "PUT",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            username: editForm.username,
            email: editForm.email,
          }),
        },
      );

      if (response.ok) {
        const updatedUser = await response.json();
        dispatch(updateUser(updatedUser));
        setIsEditing(false);
      } else {
        const data = await response.json();
        setError(data.error || "Ошибка при обновлении профиля");
      }
    } catch (err) {
      setError("Ошибка соединения с сервером");
    } finally {
      setLoading(false);
    }
  };

  const handleAvatarClick = () => {
    fileInputRef.current?.click();
  };

  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !user) return;

    if (file.size > 5 * 1024 * 1024) {
      setError("Файл слишком большой. Максимальный размер 5MB");
      return;
    }

    setAvatarLoading(true);
    setError(null);

    try {
      const token = localStorage.getItem("token");
      const formData = new FormData();
      formData.append("avatar", file);

      const response = await fetch(
        `http://localhost:8080/api/users/${user.id}/avatar`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
          body: formData,
        },
      );

      if (response.ok) {
        const updatedUser = await response.json();
        dispatch(updateUser(updatedUser));
      } else {
        const data = await response.json();
        setError(data.error || "Ошибка при загрузке аватара");
      }
    } catch (err) {
      setError("Ошибка соединения с сервером");
    } finally {
      setAvatarLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <>
      <div className={style.overlay} onClick={onClose} />
      <div className={style.profilePanel}>
        <div className={style.header}>
          <h2>Профиль</h2>
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
          {user ? (
            <>
              {isEditing ? (
                <div className={style.editForm}>
                  <div className={style.formGroup}>
                    <label htmlFor="username">Имя пользователя</label>
                    <input
                      id="username"
                      type="text"
                      value={editForm.username}
                      onChange={(e) =>
                        setEditForm({ ...editForm, username: e.target.value })
                      }
                      placeholder="Введите имя"
                      className={style.input}
                    />
                  </div>
                  <div className={style.formGroup}>
                    <label htmlFor="email">Email</label>
                    <input
                      id="email"
                      type="email"
                      value={editForm.email}
                      onChange={(e) =>
                        setEditForm({ ...editForm, email: e.target.value })
                      }
                      placeholder="Введите email"
                      className={style.input}
                    />
                  </div>
                  {error && <div className={style.errorMessage}>{error}</div>}
                  <div className={style.formActions}>
                    <button
                      className={style.cancelButton}
                      onClick={handleCancelEdit}
                      disabled={loading}
                    >
                      Отмена
                    </button>
                    <button
                      className={style.saveButton}
                      onClick={handleSaveEdit}
                      disabled={loading}
                    >
                      {loading ? "Сохранение..." : "Сохранить"}
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  <div className={style.avatarSection}>
                    <input
                      ref={fileInputRef}
                      type="file"
                      accept="image/*"
                      onChange={handleAvatarChange}
                      style={{ display: "none" }}
                    />
                    <div
                      className={style.avatarWrapper}
                      onClick={handleAvatarClick}
                    >
                      {avatarLoading ? (
                        <div className={style.avatarLoading}>
                          <div className={style.spinner} />
                        </div>
                      ) : (
                        <>
                          {user.avatar_url ? (
                            <img
                              src={user.avatar_url}
                              alt={user.username}
                              className={style.avatar}
                            />
                          ) : (
                            <div className={style.avatarPlaceholder}>
                              {getInitials(user.username)}
                            </div>
                          )}
                          <div className={style.avatarOverlay}>
                            <svg
                              width="24"
                              height="24"
                              viewBox="0 0 24 24"
                              fill="none"
                            >
                              <path
                                d="M12 5V19M5 12H19"
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                              />
                            </svg>
                          </div>
                        </>
                      )}
                    </div>
                    <div className={style.userInfo}>
                      <h3>{user.username}</h3>
                      <span className={style.email}>{user.email}</span>
                      {user.status && (
                        <span className={style.status}>
                          <span
                            className={`${style.statusDot} ${
                              user.status === "online" ? style.online : ""
                            }`}
                          />
                          {user.status === "online" ? "В сети" : "Не в сети"}
                        </span>
                      )}
                    </div>
                  </div>

                  <div className={style.infoSection}>
                    <div className={style.infoItem}>
                      <span className={style.label}>ID пользователя</span>
                      <span className={style.value}>{user.id}</span>
                    </div>
                    <div className={style.infoItem}>
                      <span className={style.label}>Дата регистрации</span>
                      <span className={style.value}>
                        {formatDate(user.created_at)}
                      </span>
                    </div>
                    <div className={style.infoItem}>
                      <span className={style.label}>Последняя активность</span>
                      <span className={style.value}>
                        {user.last_seen
                          ? new Date(user.last_seen).toLocaleString("ru-RU")
                          : "Неизвестно"}
                      </span>
                    </div>
                    {user.role && (
                      <div className={style.infoItem}>
                        <span className={style.label}>Роль</span>
                        <span className={style.value}>{user.role}</span>
                      </div>
                    )}
                  </div>

                  {error && <div className={style.errorMessage}>{error}</div>}

                  <div className={style.actions}>
                    <button
                      className={style.actionButton}
                      onClick={handleAvatarClick}
                      disabled={avatarLoading}
                    >
                      <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                      >
                        <path
                          d="M12 15V3M9 12L12 15L15 12M5 21H19"
                          stroke="currentColor"
                          strokeWidth="2"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                        />
                      </svg>
                      {avatarLoading ? "Загрузка..." : "Изменить аватар"}
                    </button>
                    <button
                      className={style.actionButton}
                      onClick={handleEditClick}
                    >
                      <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                      >
                        <path
                          d="M20 14.66V20C20 21.1 19.1 22 18 22H6C4.9 22 4 21.1 4 20V14.66M12 2V15M12 2L8 6M12 2L16 6"
                          stroke="currentColor"
                          strokeWidth="2"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                        />
                      </svg>
                      Редактировать профиль
                    </button>
                    <button
                      className={`${style.actionButton} ${style.logout}`}
                      onClick={handleLogout}
                    >
                      <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                      >
                        <path
                          d="M9 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H9M16 17L21 12M21 12L16 7M21 12H9"
                          stroke="currentColor"
                          strokeWidth="2"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                        />
                      </svg>
                      Выйти
                    </button>
                  </div>
                </>
              )}
            </>
          ) : (
            <div className={style.error}>
              <span>Не удалось загрузить профиль</span>
            </div>
          )}
        </div>
      </div>
    </>
  );
}

export default Profile;
