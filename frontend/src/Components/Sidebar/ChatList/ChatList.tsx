import { useDispatch, useSelector } from "react-redux";
import { openCreateChat } from "../../../store/modalChatSlice";
import {
  selectSearchLoading,
  selectSearchResults,
  selectSearchQuery,
  clearResults,
} from "../../../store/searchSlice";
import style from "./ChatList.module.scss";
import SearchBar from "../SearchBar/SearchBar";
import StoriesSection from "../StoriesSection/StoriesSection";
import { useState, useEffect } from "react";

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  avatar_path: string | boolean;
  last_seen_at: string | null;
  created_at: string;
  updated_at: string;
  is_online?: boolean;
}

interface Chat {
  id: number;
  type: string;
  name: string | null;
  description: string | null;
  avatar_path: string | null;
  created_by: number | null;
  created_at: string;
  updated_at: string;
}

function ChatList() {
  const dispatch = useDispatch();
  const [isSearchFocused, setIsSearchFocused] = useState(false);
  const [activeChatId, setActiveChatId] = useState<number | null>(null);
  const [myChats, setMyChats] = useState<Chat[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentUserId, setCurrentUserId] = useState<number | null>(null);

  const searchResults = useSelector(selectSearchResults);
  const searchLoading = useSelector(selectSearchLoading);
  const searchQuery = useSelector(selectSearchQuery);

  const showSearchResults = isSearchFocused && searchQuery.trim() !== "";

  useEffect(() => {
    const userStr = localStorage.getItem("user");
    if (userStr) {
      const user = JSON.parse(userStr);
      setCurrentUserId(user.id);
    }
    fetchMyChats();
  }, []);

  const fetchMyChats = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/chats/");
      if (response.ok) {
        const data = await response.json();
        const chats = Array.isArray(data) ? data : [];

        const userStr = localStorage.getItem("user");
        let userId = null;
        if (userStr) {
          const user = JSON.parse(userStr);
          userId = user.id;
        }

        const filteredChats = chats.filter(
          (chat: Chat) => chat.created_by === userId,
        );
        setMyChats(filteredChats);
      }
    } catch (error) {
      console.error("Ошибка:", error);
    } finally {
      setLoading(false);
    }
  };
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const msgDate = new Date(
      date.getFullYear(),
      date.getMonth(),
      date.getDate(),
    );

    if (msgDate.getTime() === today.getTime()) {
      return date.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
    }

    const diffDays = Math.floor(
      (today.getTime() - msgDate.getTime()) / (1000 * 60 * 60 * 24),
    );
    if (diffDays === 1) return "Вчера";
    if (diffDays < 7)
      return date.toLocaleDateString("ru-RU", { weekday: "short" });
    return date.toLocaleDateString();
  };

  const getInitials = (name: string | null) => {
    if (!name) return "💬";
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2);
  };

  const getChatName = (chat: Chat) => {
    if (chat.type === "private") {
      return "Приватный чат";
    }
    return chat.name || "Без названия";
  };

  const handleUserClick = (user: User) => {
    setActiveChatId(user.id);
    setIsSearchFocused(false);
    dispatch(clearResults());
    console.log("Выбран пользователь:", user.name);
  };

  const handleChatClick = (chat: Chat) => {
    setActiveChatId(chat.id);
    console.log("Выбран чат:", getChatName(chat));
  };

  const renderChatAvatar = (chat: Chat) => {
    return (
      <div className={style.avatarPlaceholder}>
        {getInitials(getChatName(chat))}
      </div>
    );
  };

  const renderAvatar = (user: User) => {
    const avatarUrl =
      typeof user.avatar_path === "string" && user.avatar_path
        ? user.avatar_path
        : null;

    if (avatarUrl) {
      return <img src={avatarUrl} alt={user.name} />;
    }

    return (
      <div className={style.avatarPlaceholder}>{getInitials(user.name)}</div>
    );
  };

  const formatLastSeen = (dateString: string | null) => {
    if (!dateString) return "Никогда";
    const date = new Date(dateString);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const hours = Math.floor(diff / (1000 * 60 * 60));

    if (hours < 1) return "только что";
    if (hours < 24) return `${hours} ч назад`;
    const days = Math.floor(hours / 24);
    if (days === 1) return "вчера";
    return `${days} д назад`;
  };

  return (
    <div className={style.chatList}>
      <div className={style.header}>
        <div className={style.titleWrapper}>
          <h1>Чаты</h1>
          <button
            className={style.addButton}
            onClick={() => dispatch(openCreateChat())}
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path
                d="M12 5V19M5 12H19"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
          </button>
        </div>
      </div>

      <StoriesSection />
      <SearchBar
        onFocus={() => setIsSearchFocused(true)}
        onBlur={() => setIsSearchFocused(false)}
      />

      <div className={style.chatsContainer}>
        {showSearchResults && (
          <>
            {searchLoading ? (
              <div className={style.loadingState}>
                <div className={style.spinner}></div>
                <span>Поиск пользователей...</span>
              </div>
            ) : searchResults.length === 0 ? (
              <div className={style.emptyState}>
                <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M15.5 15.5L19 19M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                  />
                </svg>
                <p>Ничего не найдено</p>
                <span>Пользователь "{searchQuery}" не найден</span>
              </div>
            ) : (
              <>
                <div className={style.searchHeader}>
                  <span>Результаты поиска</span>
                  <span className={style.resultsCount}>
                    {searchResults.length}
                  </span>
                </div>
                {searchResults.map((user) => (
                  <div
                    key={user.id}
                    className={`${style.chatItem} ${activeChatId === user.id ? style.active : ""}`}
                    onClick={() => handleUserClick(user)}
                  >
                    <div className={style.avatar}>
                      {renderAvatar(user)}
                      {user.is_online && <div className={style.onlineBadge} />}
                    </div>
                    <div className={style.chatInfo}>
                      <div className={style.chatHeader}>
                        <div className={style.chatName}>{user.name}</div>
                        <div className={style.timestamp}>
                          {formatLastSeen(user.last_seen_at)}
                        </div>
                      </div>
                      <div className={style.messagePreview}>
                        <div className={style.lastMessage}>{user.email}</div>
                      </div>
                    </div>
                  </div>
                ))}
              </>
            )}
          </>
        )}

        {!showSearchResults && (
          <>
            {loading ? (
              <div className={style.loadingState}>
                <div className={style.spinner}></div>
                <span>Загрузка чатов...</span>
              </div>
            ) : myChats.length === 0 ? (
              <div className={style.emptyState}>
                <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M21 15C21 15.5304 20.7893 16.0391 20.4142 16.4142C20.0391 16.7893 19.5304 17 19 17H7L3 21V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H19C19.5304 3 20.0391 3.21071 20.4142 3.58579C20.7893 3.96086 21 4.46957 21 5V15Z"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    fill="none"
                  />
                </svg>
                <p>У вас пока нет чатов</p>
                <button
                  className={style.createChatBtn}
                  onClick={() => dispatch(openCreateChat())}
                >
                  Создать чат
                </button>
              </div>
            ) : (
              myChats.map((chat) => (
                <div
                  key={chat.id}
                  className={`${style.chatItem} ${activeChatId === chat.id ? style.active : ""}`}
                  onClick={() => handleChatClick(chat)}
                >
                  <div className={style.avatar}>{renderChatAvatar(chat)}</div>
                  <div className={style.chatInfo}>
                    <div className={style.chatHeader}>
                      <div className={style.chatName}>{getChatName(chat)}</div>
                      <div className={style.timestamp}>
                        {formatDate(chat.updated_at)}
                      </div>
                    </div>
                    <div className={style.messagePreview}>
                      <div className={style.lastMessage}>
                        {chat.description || "Нет описания"}
                      </div>
                    </div>
                  </div>
                </div>
              ))
            )}
          </>
        )}
      </div>
    </div>
  );
}

export default ChatList;
