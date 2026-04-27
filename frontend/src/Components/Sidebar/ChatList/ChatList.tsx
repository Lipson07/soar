import { useDispatch, useSelector } from "react-redux";
import { openCreateChat } from "../../../store/modalChatSlice";
import {
  selectSearchLoading,
  selectSearchResults,
  clearResults,
} from "../../../store/searchSlice";
import {
  toggleChat,
  selectCurrentChat,
  selectIsChatOpen,
} from "../../../store/selectedChatSlice";
import { selectChats, setChats, setLoading } from "../../../store/chatSlice";
import { useWebSocket } from "../../../context/WebSocketContext";
import style from "./ChatList.module.scss";
import SearchBar from "../SearchBar/SearchBar";
import StoriesSection from "../StoriesSection/StoriesSection";
import { useState, useEffect, useCallback } from "react";

interface User {
  id: string;
  username: string;
  email: string;
  avatar_url: string | null;
  last_seen?: string | null;
  status?: string;
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
  last_message?: {
    id: string;
    text: string;
    user_id: string;
    created_at: string;
    type?: string;
    file_url?: string;
    file_name?: string;
  } | null;
  unread_count?: number;
  other_user_id?: string | null;
  other_user_name?: string | null;
}

function ChatList() {
  const dispatch = useDispatch();
  const [currentUserId, setCurrentUserId] = useState<string | null>(null);
  const [usersCache, setUsersCache] = useState<Map<string, User>>(new Map());
  const { lastMessage } = useWebSocket();

  const searchResults = useSelector(selectSearchResults);
  const searchLoading = useSelector(selectSearchLoading);
  const currentChat = useSelector(selectCurrentChat);
  const isChatOpen = useSelector(selectIsChatOpen);
  const myChats = useSelector(selectChats);
  const loading = useSelector((state: any) => state.chats.loading);

  const fetchMyChats = useCallback(
    async (userId: string) => {
      try {
        const token = localStorage.getItem("token");
        if (!token) return;
        dispatch(setLoading(true));
        const response = await fetch("http://localhost:8080/api/chats", {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (response.ok) {
          const data = await response.json();
          const chats: Chat[] = Array.isArray(data) ? data : [];
          const enrichedChats = await enrichPrivateChatNames(
            chats,
            token,
            userId,
          );
          dispatch(setChats(enrichedChats));
        }
      } catch (error) {
        console.error("Ошибка:", error);
      } finally {
        dispatch(setLoading(false));
      }
    },
    [dispatch],
  );

  useEffect(() => {
    const userStr = localStorage.getItem("user");
    if (userStr) {
      const user = JSON.parse(userStr);
      setCurrentUserId(user.id);
      fetchMyChats(user.id);
    }
  }, [fetchMyChats]);

  // Обработка входящих сообщений через WebSocket для обновления последнего сообщения и счетчика
  useEffect(() => {
    if (!lastMessage || !currentUserId) return;

    if (lastMessage.type === "new_message" && lastMessage.message) {
      const msg = lastMessage.message;
      const chatId = msg.chat_id;

      dispatch({
        type: "chats/updateLastMessage",
        payload: {
          chatId,
          lastMessage: {
            id: msg.id,
            text: msg.text,
            user_id: msg.user_id,
            created_at: msg.created_at,
            type: msg.type,
            file_url: msg.file_url,
            file_name: msg.file_name,
          },
          lastMessageAt: msg.created_at,
        },
      });

      // Увеличиваем счетчик непрочитанных, если сообщение не от текущего пользователя
      // и чат не открыт
      if (
        msg.user_id !== currentUserId &&
        !(isChatOpen && currentChat?.id === chatId)
      ) {
        dispatch({
          type: "chats/incrementUnread",
          payload: { chatId },
        });
      }
    }
  }, [lastMessage, currentUserId, isChatOpen, currentChat, dispatch]);

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
            { headers: { Authorization: `Bearer ${token}` } },
          );
          if (participantsResponse.ok) {
            const participants = await participantsResponse.json();
            const otherParticipant = participants.find(
              (p: any) => p.user_id !== userId,
            );
            if (otherParticipant) {
              let user = usersCache.get(otherParticipant.user_id);
              if (!user) {
                const userResponse = await fetch(
                  `http://localhost:8080/api/users/${otherParticipant.user_id}`,
                  { headers: { Authorization: `Bearer ${token}` } },
                );
                if (userResponse.ok) {
                  user = await userResponse.json();
                  setUsersCache((prev) => new Map(prev).set(user!.id, user!));
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

  const handleChatClick = (chat: Chat) => {
    // Сбрасываем счетчик непрочитанных при открытии чата
    if (chat.unread_count && chat.unread_count > 0) {
      dispatch({
        type: "chats/resetUnread",
        payload: { chatId: chat.id },
      });
    }
    dispatch(toggleChat(chat));
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
    if (!name) return "?";
    return name.charAt(0).toUpperCase();
  };

  const getChatName = (chat: Chat) => {
    if (chat.type === "private")
      return chat.name || chat.other_user_name || "Пользователь";
    return chat.name || "Без названия";
  };

  const getLastMessageText = (chat: Chat) => {
    if (chat.last_message) {
      const msg = chat.last_message;
      if (msg.type === "image") return "📷 Изображение";
      if (msg.type === "file") return `📎 ${msg.file_name || "Файл"}`;
      if (msg.text) {
        return msg.user_id === currentUserId ? `Вы: ${msg.text}` : msg.text;
      }
    }
    return "Нет сообщений";
  };

  const handleUserClick = (user: User) => {
    const existingChat = myChats.find(
      (chat) => chat.type === "private" && chat.other_user_id === user.id,
    );
    if (existingChat) {
      dispatch(toggleChat(existingChat));
    } else {
      const tempChat: Chat = {
        id: `temp_${user.id}`,
        type: "private",
        name: user.username,
        creator_id: currentUserId || "",
        avatar_url: user.avatar_url,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        last_message_at: null,
        last_message: null,
        unread_count: 0,
        other_user_id: user.id,
        other_user_name: user.username,
      };
      dispatch(toggleChat(tempChat));
    }
    dispatch(clearResults());
  };

  const isChatActive = (chatId: string) =>
    isChatOpen && currentChat?.id === chatId;

  const renderAvatar = (item: Chat | User) => {
    const avatarUrl = "avatar_url" in item ? item.avatar_url : null;
    const name = "username" in item ? item.username : getChatName(item as Chat);

    if (avatarUrl) {
      return <img src={avatarUrl} alt={name} />;
    }
    return <div className={style.avatarPlaceholder}>{getInitials(name)}</div>;
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
      <SearchBar />

      <div className={style.chatsContainer}>
        {loading ? (
          <div className={style.loadingState}>
            <div className={style.spinner}></div>
            <span>Загрузка чатов...</span>
          </div>
        ) : myChats.length === 0 ? (
          <div className={style.emptyState}>
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
              className={`${style.chatItem} ${isChatActive(chat.id) ? style.active : ""} ${chat.unread_count && chat.unread_count > 0 ? style.unread : ""}`}
              onClick={() => handleChatClick(chat)}
            >
              <div className={style.avatar}>{renderAvatar(chat)}</div>
              <div className={style.chatInfo}>
                <div className={style.chatHeader}>
                  <div className={style.chatName}>{getChatName(chat)}</div>
                  <div className={style.timestamp}>
                    {formatDate(chat.last_message_at || chat.updated_at)}
                  </div>
                </div>
                <div className={style.messagePreview}>
                  <div className={style.lastMessage}>
                    {getLastMessageText(chat)}
                  </div>
                  {chat.unread_count && chat.unread_count > 0 ? (
                    <div className={style.unreadBadge}>{chat.unread_count}</div>
                  ) : null}
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

export default ChatList;
