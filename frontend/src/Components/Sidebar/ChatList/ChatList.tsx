import { useDispatch, useSelector } from "react-redux";
import { openCreateChat } from "../../../store/modalChatSlice";
import {
  selectSearchLoading,
  selectSearchResults,
  selectSearchQuery,
  clearResults,
} from "../../../store/searchSlice";
import {
  toggleChat,
  selectCurrentChat,
  selectIsChatOpen,
} from "../../../store/selectedChatSlice";
import style from "./ChatList.module.scss";
import SearchBar from "../SearchBar/SearchBar";
import StoriesSection from "../StoriesSection/StoriesSection";
import { useState, useEffect, useCallback } from "react";

interface User {
  id: string;
  username: string;
  email: string;
  role?: string;
  avatar_url: string | null;
  last_seen?: string | null;
  created_at?: string;
  updated_at?: string;
  status?: string;
  is_online?: boolean;
}

interface Chat {
  id: string;
  type: string;
  name: string | null;
  creator_id: string;
  avatar_url: string | null;
  created_at: string;
  updated_at: string;
  last_message_at: string | null;
  last_message: {
    id: string;
    text: string;
    user_id: string;
    created_at: string;
    file_url?: string;
    file_name?: string;
    type?: string;
  } | null;
  unread_count: number;
  other_user_id?: string | null;
  other_user_name?: string | null;
}

function ChatList() {
  const dispatch = useDispatch();
  const [isSearchFocused, setIsSearchFocused] = useState(false);
  const [myChats, setMyChats] = useState<Chat[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentUserId, setCurrentUserId] = useState<string | null>(null);
  const [usersCache, setUsersCache] = useState<Map<string, User>>(new Map());

  const searchResults = useSelector(selectSearchResults);
  const searchLoading = useSelector(selectSearchLoading);
  const searchQuery = useSelector(selectSearchQuery);
  const currentChat = useSelector(selectCurrentChat);
  const isChatOpen = useSelector(selectIsChatOpen);

  const showSearchResults = isSearchFocused && searchQuery.trim() !== "";

  const fetchMyChats = useCallback(async (userId: string) => {
    try {
      const token = localStorage.getItem("token");

      if (!token) {
        console.error("Токен не найден");
        setLoading(false);
        return;
      }

      const response = await fetch("http://localhost:8080/api/chats", {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json();
        const chats: Chat[] = Array.isArray(data) ? data : [];

        const enrichedChats = await enrichPrivateChatNames(
          chats,
          token,
          userId,
        );
        setMyChats(enrichedChats);
      } else {
        console.error("Ошибка при получении чатов:", response.status);
      }
    } catch (error) {
      console.error("Ошибка:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const userStr = localStorage.getItem("user");
    if (userStr) {
      const user = JSON.parse(userStr);
      setCurrentUserId(user.id);
      fetchMyChats(user.id);
    }
  }, [fetchMyChats]);

  useEffect(() => {
    if (!currentUserId) return;

    const interval = setInterval(() => {
      fetchMyChats(currentUserId);
    }, 5000);

    return () => clearInterval(interval);
  }, [currentUserId, fetchMyChats]);

  useEffect(() => {
    const handleMessageSent = () => {
      if (currentUserId) {
        fetchMyChats(currentUserId);
      }
    };

    window.addEventListener("messageSent", handleMessageSent);
    return () => window.removeEventListener("messageSent", handleMessageSent);
  }, [currentUserId, fetchMyChats]);

  const enrichPrivateChatNames = async (
    chats: Chat[],
    token: string,
    userId: string,
  ) => {
    const enriched = [...chats];

    for (let i = 0; i < enriched.length; i++) {
      const chat = enriched[i];

      if (chat.type === "private") {
        try {
          const participantsResponse = await fetch(
            `http://localhost:8080/api/participants?chat_id=${chat.id}`,
            {
              method: "GET",
              headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json",
              },
            },
          );

          if (participantsResponse.ok) {
            const participants = await participantsResponse.json();

            const otherParticipant = participants.find(
              (p: any) => p.user_id !== userId,
            );

            if (otherParticipant) {
              const otherUserId = otherParticipant.user_id;
              let user = usersCache.get(otherUserId);

              if (!user) {
                const userResponse = await fetch(
                  `http://localhost:8080/api/users/${otherUserId}`,
                  {
                    method: "GET",
                    headers: {
                      Authorization: `Bearer ${token}`,
                      "Content-Type": "application/json",
                    },
                  },
                );

                if (userResponse.ok) {
                  user = await userResponse.json();
                  setUsersCache((prev) => new Map(prev).set(user.id, user));
                }
              }

              if (user) {
                enriched[i].name = user.username;
                enriched[i].other_user_id = user.id;
                enriched[i].other_user_name = user.username;
                enriched[i].avatar_url = user.avatar_url;
              }
            }
          }
        } catch (error) {
          console.error("Ошибка получения имени для чата:", chat.id, error);
        }
      }
    }

    return enriched;
  };

  const formatDate = (dateString: string | null | undefined) => {
    if (!dateString) return "";
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
      return chat.name || chat.other_user_name || "Приватный чат";
    }
    return chat.name || "Без названия";
  };

  const getLastMessageText = (chat: Chat): string => {
    if (!chat.last_message || !chat.last_message.id) {
      return "Нет сообщений";
    }

    const msg = chat.last_message;

    if (msg.type === "image") {
      return "📷 Изображение";
    }
    if (msg.type === "file") {
      return `📎 ${msg.file_name || "Файл"}`;
    }

    if (msg.text) {
      if (msg.user_id === currentUserId) {
        return `Вы: ${msg.text}`;
      }
      return msg.text;
    }

    return "Нет сообщений";
  };

  const handleUserClick = async (user: User) => {
    const token = localStorage.getItem("token");
    if (!token || !currentUserId) return;

    try {
      const response = await fetch("http://localhost:8080/api/chats/private", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          user_id: user.id,
        }),
      });

      if (response.ok) {
        const chat = await response.json();

        const enrichedChat: Chat = {
          ...chat,
          name: user.username,
          other_user_id: user.id,
          other_user_name: user.username,
          avatar_url: user.avatar_url,
          unread_count: 0,
          last_message: null,
        };

        dispatch(toggleChat(enrichedChat));
        setIsSearchFocused(false);
        dispatch(clearResults());

        fetchMyChats(currentUserId);
      } else {
        console.error("Ошибка создания чата");
      }
    } catch (error) {
      console.error("Ошибка:", error);
    }
  };

  const handleChatClick = (chat: Chat) => {
    dispatch(toggleChat(chat));
  };

  const isChatActive = (chatId: string) => {
    return isChatOpen && currentChat?.id === chatId;
  };

  const renderChatAvatar = (chat: Chat) => {
    if (chat.avatar_url) {
      return <img src={chat.avatar_url} alt={getChatName(chat)} />;
    }

    return (
      <div className={style.avatarPlaceholder}>
        {getInitials(getChatName(chat))}
      </div>
    );
  };

  const renderAvatar = (user: User) => {
    const avatarUrl = user.avatar_url;
    const displayName = user.username;

    if (avatarUrl) {
      return <img src={avatarUrl} alt={displayName} />;
    }

    return (
      <div className={style.avatarPlaceholder}>{getInitials(displayName)}</div>
    );
  };

  const formatLastSeen = (dateString: string | null | undefined) => {
    if (!dateString) return "";
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
                    className={`${style.chatItem} ${
                      isChatActive(user.id) ? style.active : ""
                    }`}
                    onClick={() => handleUserClick(user)}
                  >
                    <div className={style.avatar}>
                      {renderAvatar(user)}
                      {user.status === "online" && (
                        <div className={style.onlineBadge} />
                      )}
                    </div>
                    <div className={style.chatInfo}>
                      <div className={style.chatHeader}>
                        <div className={style.chatName}>{user.username}</div>
                        <div className={style.timestamp}>
                          {formatLastSeen(user.last_seen)}
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
                  className={`${style.chatItem} ${
                    isChatActive(chat.id) ? style.active : ""
                  } ${chat.unread_count > 0 ? style.unread : ""}`}
                  onClick={() => handleChatClick(chat)}
                >
                  <div className={style.avatar}>{renderChatAvatar(chat)}</div>
                  <div className={style.chatInfo}>
                    <div className={style.chatHeader}>
                      <div className={style.chatName}>{getChatName(chat)}</div>
                      <div className={style.timestamp}>
                        {formatDate(chat.last_message_at || chat.updated_at)}
                      </div>
                    </div>
                    <div className={style.messagePreview}>
                      {chat.last_message && chat.last_message.id ? (
                        <div className={style.lastMessage}>
                          {getLastMessageText(chat)}
                        </div>
                      ) : (
                        <div
                          className={style.lastMessage}
                          style={{ opacity: 0.5 }}
                        >
                          Нет сообщений
                        </div>
                      )}
                      {chat.unread_count > 0 && (
                        <div className={style.unreadBadge}>
                          {chat.unread_count > 99 ? "99+" : chat.unread_count}
                        </div>
                      )}
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
