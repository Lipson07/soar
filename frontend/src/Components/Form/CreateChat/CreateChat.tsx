import React, { useState, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  closeCreateChat,
  selectModalChat,
} from "../../../store/modalChatSlice";
import style from "./CreateChat.module.scss";

interface User {
  id: string;
  username: string;
  email: string;
  avatar_url: string | null;
  status: string;
}

function CreateChat() {
  const dispatch = useDispatch();
  const isOpen = useSelector(selectModalChat);

  const [chatName, setChatName] = useState("");
  const [chatType, setChatType] = useState("private");
  const [description, setDescription] = useState("");
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<User[]>([]);
  const [selectedUsers, setSelectedUsers] = useState<User[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [currentUserId, setCurrentUserId] = useState<string | null>(null);

  useEffect(() => {
    const userStr = localStorage.getItem("user");
    if (userStr) {
      const user = JSON.parse(userStr);
      setCurrentUserId(user.id);
    }
  }, []);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "unset";
      resetForm();
    }
    return () => {
      document.body.style.overflow = "unset";
    };
  }, [isOpen]);

  useEffect(() => {
    if (!searchQuery.trim()) {
      setSearchResults([]);
      return;
    }

    const timer = setTimeout(async () => {
      setSearchLoading(true);
      try {
        const token = localStorage.getItem("token");
        const response = await fetch(
          `http://localhost:8080/api/users/search?q=${encodeURIComponent(searchQuery)}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );
        const data = await response.json();
        const users = Array.isArray(data) ? data : [];
        const filtered = users.filter(
          (user: User) =>
            user.id !== currentUserId &&
            !selectedUsers.some((selected) => selected.id === user.id),
        );
        setSearchResults(filtered);
      } catch (error) {
        console.error("Ошибка поиска:", error);
        setSearchResults([]);
      } finally {
        setSearchLoading(false);
      }
    }, 500);

    return () => clearTimeout(timer);
  }, [searchQuery, selectedUsers, currentUserId]);

  const resetForm = () => {
    setChatName("");
    setChatType("private");
    setDescription("");
    setSearchQuery("");
    setSearchResults([]);
    setSelectedUsers([]);
    setIsSubmitting(false);
  };

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      dispatch(closeCreateChat());
    }
  };

  const handleClose = () => {
    dispatch(closeCreateChat());
  };

  const handleAddUser = (user: User) => {
    setSelectedUsers([...selectedUsers, user]);
    setSearchQuery("");
    setSearchResults([]);
  };

  const handleRemoveUser = (userId: string) => {
    setSelectedUsers(selectedUsers.filter((u) => u.id !== userId));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (chatType === "group" && !chatName.trim()) {
      alert("Введите название чата");
      return;
    }

    if (selectedUsers.length === 0 && chatType === "group") {
      alert("Выберите хотя бы одного участника");
      return;
    }

    if (chatType === "private" && selectedUsers.length === 0) {
      alert("Выберите пользователя для личного чата");
      return;
    }

    setIsSubmitting(true);

    try {
      const token = localStorage.getItem("token");

      if (chatType === "private") {
        const response = await fetch(
          "http://localhost:8080/api/chats/private",
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
              user_id: selectedUsers[0].id,
            }),
          },
        );

        if (!response.ok) {
          const error = await response.json();
          throw new Error(error.error || "Ошибка создания чата");
        }

        alert("Чат успешно создан!");
        dispatch(closeCreateChat());
        window.location.reload();
      } else {
        const response = await fetch("http://localhost:8080/api/chats/group", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            name: chatName,
            user_ids: selectedUsers.map((u) => u.id),
            avatar_url: null,
          }),
        });

        if (!response.ok) {
          const error = await response.json();
          throw new Error(error.error || "Ошибка создания группы");
        }

        alert("Группа успешно создана!");
        dispatch(closeCreateChat());
        window.location.reload();
      }
    } catch (error) {
      console.error("Ошибка:", error);
      alert(
        error instanceof Error ? error.message : "Ошибка при создании чата",
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className={style.overlay} onClick={handleOverlayClick}>
      <div className={style.modal}>
        <div className={style.header}>
          <h2>Создать чат</h2>
          <button className={style.closeButton} onClick={handleClose}>
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

        <form onSubmit={handleSubmit} className={style.form}>
          <div className={style.typeSection}>
            <label className={style.sectionLabel}>Тип чата</label>
            <div className={style.typeOptions}>
              <button
                type="button"
                className={`${style.typeBtn} ${chatType === "private" ? style.active : ""}`}
                onClick={() => setChatType("private")}
              >
                <div className={style.typeInfo}>
                  <span className={style.typeTitle}>Приватный</span>
                  <span className={style.typeDesc}>
                    Диалог между двумя людьми
                  </span>
                </div>
              </button>

              <button
                type="button"
                className={`${style.typeBtn} ${chatType === "group" ? style.active : ""}`}
                onClick={() => setChatType("group")}
              >
                <div className={style.typeInfo}>
                  <span className={style.typeTitle}>Группа</span>
                  <span className={style.typeDesc}>
                    Общайтесь с несколькими людьми
                  </span>
                </div>
              </button>
            </div>
          </div>

          {chatType === "group" && (
            <div className={style.inputGroup}>
              <label className={style.label}>Название чата</label>
              <input
                type="text"
                className={style.input}
                placeholder="Введите название"
                value={chatName}
                onChange={(e) => setChatName(e.target.value)}
                required
              />
            </div>
          )}

          <div className={style.inputGroup}>
            <label className={style.label}>
              {chatType === "private"
                ? "Выберите пользователя"
                : "Добавить участников"}
            </label>
            <div className={style.searchContainer}>
              <input
                type="text"
                className={style.searchInput}
                placeholder="Поиск по имени или email..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              {searchLoading && <div className={style.spinner}></div>}
            </div>

            {searchResults.length > 0 && (
              <div className={style.searchResults}>
                {searchResults.map((user) => (
                  <div
                    key={user.id}
                    className={style.searchItem}
                    onClick={() => handleAddUser(user)}
                  >
                    <div className={style.userAvatar}>
                      {user.avatar_url ? (
                        <img src={user.avatar_url} alt={user.username} />
                      ) : (
                        <div className={style.avatarPlaceholder}>
                          {user.username.charAt(0).toUpperCase()}
                        </div>
                      )}
                    </div>
                    <div className={style.userInfo}>
                      <div className={style.userName}>{user.username}</div>
                      <div className={style.userEmail}>{user.email}</div>
                    </div>
                    <button type="button" className={style.addBtn}>
                      +
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          {selectedUsers.length > 0 && (
            <div className={style.selectedSection}>
              <label className={style.label}>
                Выбранные ({selectedUsers.length})
              </label>
              <div className={style.selectedList}>
                {selectedUsers.map((user) => (
                  <div key={user.id} className={style.selectedItem}>
                    <span>{user.username}</span>
                    <button
                      type="button"
                      onClick={() => handleRemoveUser(user.id)}
                    >
                      ×
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}

          <div className={style.actions}>
            <button
              type="button"
              className={style.cancelBtn}
              onClick={handleClose}
            >
              Отмена
            </button>
            <button
              type="submit"
              className={style.submitBtn}
              disabled={isSubmitting}
            >
              {isSubmitting ? "Создание..." : "Создать чат"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default CreateChat;
