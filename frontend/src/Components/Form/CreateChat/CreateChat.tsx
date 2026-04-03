import React, { useState, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  closeCreateChat,
  selectModalChat,
} from "../../../store/modalChatSlice";
import style from "./CreateChat.module.scss";

interface User {
  id: number;
  name: string;
  email: string;
  avatar_path: string | boolean;
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
        const response = await fetch(
          `http://localhost:8080/api/users/search?query=${encodeURIComponent(searchQuery)}`,
        );
        const data = await response.json();
        const users = Array.isArray(data) ? data : [];
        const filtered = users.filter(
          (user: User) =>
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
  }, [searchQuery, selectedUsers]);

  const resetForm = () => {
    setChatName("");
    setChatType("private");
    setDescription("");
    setSearchQuery("");
    setSearchResults([]);
    setSelectedUsers([]);
    setIsSubmitting(false);
  };

  const handleOverlayClick = (e: any) => {
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

  const handleRemoveUser = (userId: number) => {
    setSelectedUsers(selectedUsers.filter((u) => u.id !== userId));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (chatType !== "private" && !chatName.trim()) {
      alert("Введите название чата");
      return;
    }

    if (selectedUsers.length === 0) {
      alert("Выберите хотя бы одного участника");
      return;
    }

    setIsSubmitting(true);

    try {
      const userStr = localStorage.getItem("user");
      let currentUserId = null;
      let currentUser = null;
      if (userStr) {
        currentUser = JSON.parse(userStr);
        currentUserId = currentUser.id;
      }

      const requestBody: any = {
        type: chatType,
        description: description || null,
      };

      if (chatType === "group") {
        requestBody.name = chatName;
      }

      const chatResponse = await fetch("http://localhost:8080/api/chats/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(requestBody),
      });

      if (!chatResponse.ok) {
        const error = await chatResponse.json();
        throw new Error(error.error || "Ошибка создания чата");
      }

      const chat = await chatResponse.json();

      // Добавляем всех выбранных пользователей + текущего
      const allUserIds = [...selectedUsers.map((u) => u.id)];
      if (currentUserId && !allUserIds.includes(currentUserId)) {
        allUserIds.push(currentUserId);
      }

      if (allUserIds.length > 0) {
        await fetch(`http://localhost:8080/api/chats/${chat.id}/members/bulk`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            user_ids: allUserIds,
          }),
        });
      }

      alert("Чат успешно создан!");
      dispatch(closeCreateChat());
      window.location.reload();
    } catch (error) {
      console.error("Ошибка:", error);
      alert("Ошибка при создании чата");
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
              />
            </div>
          )}

          <div className={style.inputGroup}>
            <label className={style.label}>Добавить участников</label>
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
                      {typeof user.avatar_path === "string" &&
                      user.avatar_path ? (
                        <img src={user.avatar_path} alt={user.name} />
                      ) : (
                        <div className={style.avatarPlaceholder}>
                          {user.name.charAt(0).toUpperCase()}
                        </div>
                      )}
                    </div>
                    <div className={style.userInfo}>
                      <div className={style.userName}>{user.name}</div>
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
                    <span>{user.name}</span>
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

          <div className={style.inputGroup}>
            <label className={style.label}>Описание</label>
            <textarea
              className={style.textarea}
              placeholder="Введите описание чата..."
              rows={3}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>

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
