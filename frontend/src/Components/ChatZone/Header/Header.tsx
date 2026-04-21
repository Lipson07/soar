import React, { useState, useEffect, useRef } from "react";
import {
  FiPhone,
  FiVideo,
  FiSearch,
  FiMoreVertical,
  FiInfo,
  FiUserPlus,
  FiLogOut,
  FiTrash2,
  FiX,
  FiChevronUp,
  FiChevronDown,
} from "react-icons/fi";
import style from "./Header.module.scss";
import {
  selectCurrentChat,
  selectChatMessages,
} from "../../../store/selectedChatSlice";
import { useSelector } from "react-redux";
import CallModal from "../CallModal/CallModal";

interface User {
  id: string;
  username: string;
  avatar_url: string | null;
  status: string;
  last_seen: string | null;
  is_online?: boolean;
}

interface ChatInfo {
  id: string;
  name: string | null;
  type: string;
  creator_id: string;
  participants_count?: number;
}

interface Message {
  id: string;
  chat_id: string;
  user_id: string;
  type: "text" | "image" | "file";
  text: string;
  file_url?: string;
  file_name?: string;
  created_at: string;
}

function Header() {
  const chatInfo = useSelector(selectCurrentChat);
  const messages = useSelector(selectChatMessages);
  const [otherUser, setOtherUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [showMenu, setShowMenu] = useState(false);
  const [showSearch, setShowSearch] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<Message[]>([]);
  const [currentResultIndex, setCurrentResultIndex] = useState(0);
  const [chatDetails, setChatDetails] = useState<ChatInfo | null>(null);
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [callModal, setCallModal] = useState<{
    isOpen: boolean;
    type: "audio" | "video";
    callId: string;
    roomId: string;
    caller: any;
    isIncoming: boolean;
  } | null>(null);

  const menuRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const currentUser = JSON.parse(localStorage.getItem("user") || "{}");
  const currentUserId = currentUser.id;

  useEffect(() => {
    if (chatInfo && chatInfo.type === "private" && chatInfo.other_user_id) {
      fetchUserInfo(chatInfo.other_user_id);
    } else if (chatInfo && chatInfo.type === "group") {
      fetchChatDetails();
    } else {
      setOtherUser(null);
      setChatDetails(null);
    }
  }, [chatInfo]);

  useEffect(() => {
    const token = localStorage.getItem("token");
    const socket = new WebSocket(
      `ws://localhost:8080/ws?token=${token}&user_id=${currentUserId}`,
    );

    socket.onopen = () => console.log("WebSocket connected");

    socket.onmessage = (event) => {
      const signal = JSON.parse(event.data);

      if (signal.type === "call-start" && signal.call) {
        setCallModal({
          isOpen: true,
          type: signal.call.type,
          callId: signal.call.id,
          roomId: signal.call.room_id,
          caller: signal.call.caller,
          isIncoming: true,
        });
      }

      if (signal.type === "call-accept") {
        setCallModal((prev) => (prev ? { ...prev, isIncoming: false } : null));
      }

      if (signal.type === "call-reject" || signal.type === "call-end") {
        setCallModal(null);
      }
    };

    socket.onclose = () => console.log("WebSocket disconnected");
    socket.onerror = (error) => console.error("WebSocket error:", error);

    setWs(socket);

    return () => socket.close();
  }, [currentUserId]);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setShowMenu(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  useEffect(() => {
    if (showSearch && searchInputRef.current) {
      searchInputRef.current.focus();
    }
  }, [showSearch]);

  // Поиск по сообщениям
  useEffect(() => {
    if (!searchQuery.trim()) {
      setSearchResults([]);
      return;
    }

    const query = searchQuery.toLowerCase();
    const results = messages.filter(
      (msg) => msg.type === "text" && msg.text.toLowerCase().includes(query),
    );
    setSearchResults(results);
    setCurrentResultIndex(results.length > 0 ? 0 : -1);
  }, [searchQuery, messages]);

  const scrollToMessage = (messageId: string) => {
    const element = document.getElementById(`message-${messageId}`);
    if (element) {
      element.scrollIntoView({ behavior: "smooth", block: "center" });
      element.classList.add(style.highlighted);
      setTimeout(() => element.classList.remove(style.highlighted), 2000);
    }
  };

  const handleNextResult = () => {
    if (searchResults.length === 0) return;
    const nextIndex = (currentResultIndex + 1) % searchResults.length;
    setCurrentResultIndex(nextIndex);
    scrollToMessage(searchResults[nextIndex].id);
  };

  const handlePrevResult = () => {
    if (searchResults.length === 0) return;
    const prevIndex =
      (currentResultIndex - 1 + searchResults.length) % searchResults.length;
    setCurrentResultIndex(prevIndex);
    scrollToMessage(searchResults[prevIndex].id);
  };

  const handleSearchSubmit = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && searchResults.length > 0) {
      scrollToMessage(searchResults[currentResultIndex].id);
    }
  };

  const fetchUserInfo = async (userId: string) => {
    setLoading(true);
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(
        `http://localhost:8080/api/users/${userId}`,
        {
          headers: { Authorization: `Bearer ${token}` },
        },
      );
      if (response.ok) {
        const user = await response.json();
        setOtherUser(user);
      }
    } catch (error) {
      console.error("Failed to fetch user info:", error);
    } finally {
      setLoading(false);
    }
  };

  const fetchChatDetails = async () => {
    if (!chatInfo) return;
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(
        `http://localhost:8080/api/participants?chat_id=${chatInfo.id}`,
        { headers: { Authorization: `Bearer ${token}` } },
      );
      if (response.ok) {
        const participants = await response.json();
        setChatDetails({
          ...chatInfo,
          participants_count: participants.length,
        });
      }
    } catch (error) {
      console.error("Failed to fetch chat details:", error);
    }
  };

  const handleCall = (callType: "audio" | "video") => {
    if (!chatInfo || !otherUser) return;

    const callId = crypto.randomUUID();
    const roomId = crypto.randomUUID();

    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          type: "call-start",
          call_id: callId,
          room_id: roomId,
          chat_id: chatInfo.id,
          callee_id: otherUser.id,
          call_type: callType,
        }),
      );

      setCallModal({
        isOpen: true,
        type: callType,
        callId,
        roomId,
        caller: {
          id: otherUser.id,
          username: otherUser.username,
          avatar_url: otherUser.avatar_url,
        },
        isIncoming: false,
      });
    }
  };

  const handleSearch = () => {
    setShowSearch(!showSearch);
    setSearchQuery("");
    setSearchResults([]);
  };

  const handleInfo = () => {
    if (chatInfo) {
      const info =
        chatInfo.type === "private"
          ? `Чат с ${getDisplayName()}\nID: ${chatInfo.id}`
          : `Группа: ${getDisplayName()}\nУчастников: ${chatDetails?.participants_count || "..."}\nID: ${chatInfo.id}`;
      alert(info);
    }
  };

  const handleAddParticipant = () => {
    if (chatInfo?.type === "group") {
      alert("Добавление участников в разработке");
    }
  };

  const handleLeaveChat = async () => {
    if (!chatInfo) return;
    if (!confirm(`Вы уверены, что хотите покинуть чат "${getDisplayName()}"?`))
      return;

    try {
      const token = localStorage.getItem("token");
      const response = await fetch(
        `http://localhost:8080/api/participants/leave`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ chat_id: chatInfo.id }),
        },
      );

      if (response.ok) {
        alert("Вы покинули чат");
        window.location.reload();
      } else {
        const error = await response.json();
        alert(error.error || "Ошибка при выходе из чата");
      }
    } catch (error) {
      console.error("Failed to leave chat:", error);
    }
  };

  const handleDeleteChat = async () => {
    if (!chatInfo) return;
    if (!confirm(`Вы уверены, что хотите удалить чат "${getDisplayName()}"?`))
      return;

    try {
      const token = localStorage.getItem("token");
      const response = await fetch(
        `http://localhost:8080/api/chats/${chatInfo.id}`,
        {
          method: "DELETE",
          headers: { Authorization: `Bearer ${token}` },
        },
      );

      if (response.ok) {
        alert("Чат удален");
        window.location.reload();
      } else {
        const error = await response.json();
        alert(error.error || "Ошибка при удалении чата");
      }
    } catch (error) {
      console.error("Failed to delete chat:", error);
    }
  };

  const getDisplayName = () => {
    if (!chatInfo) return "Чат";
    if (chatInfo.type !== "private") {
      return chatInfo.name || "Групповой чат";
    }
    if (chatInfo.other_user_name) return chatInfo.other_user_name;
    if (chatInfo.name && chatInfo.name !== "Приватный чат")
      return chatInfo.name;
    return "Пользователь";
  };

  const getAvatar = () => {
    if (chatInfo?.type === "private" && otherUser?.avatar_url) {
      return (
        <img
          src={otherUser.avatar_url}
          alt={getDisplayName()}
          className={style.avatarImage}
        />
      );
    }
    if (chatInfo?.avatar_url) {
      return (
        <img
          src={chatInfo.avatar_url}
          alt={getDisplayName()}
          className={style.avatarImage}
        />
      );
    }
    const name = getDisplayName();
    const letter =
      name === "Чат" || name === "Групповой чат" || name === "Пользователь"
        ? "?"
        : name.charAt(0).toUpperCase();
    return <span>{letter}</span>;
  };

  const getStatus = () => {
    if (chatInfo?.type === "private" && otherUser) {
      const isOnline = otherUser.status === "online";
      return {
        online: isOnline,
        text: isOnline ? "● Онлайн" : "○ Офлайн",
        lastSeen: otherUser.last_seen,
      };
    }
    if (chatInfo?.type === "group" && chatDetails) {
      return {
        online: true,
        text: `👥 ${chatDetails.participants_count} участников`,
        lastSeen: null,
      };
    }
    return { online: true, text: "● Онлайн", lastSeen: null };
  };

  const formatLastSeen = (lastSeen: string | null) => {
    if (!lastSeen) return "";
    const date = new Date(lastSeen);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) return "Был(а) только что";
    if (minutes < 60) return `Был(а) ${minutes} мин назад`;
    if (hours < 24) return `Был(а) ${hours} ч назад`;
    if (days === 1) return "Был(а) вчера";
    return `Был(а) ${days} дн назад`;
  };

  const isCreator = chatInfo?.creator_id === currentUserId;
  const status = getStatus();
  const isPrivateChat = chatInfo?.type === "private";

  return (
    <>
      <header className={style.header}>
        <div className={style.userInfo}>
          <div className={style.avatar}>{getAvatar()}</div>
          <div className={style.userDetails}>
            <p className={style.userName}>{getDisplayName()}</p>
            <p
              className={`${style.userStatus} ${status.online ? style.online : style.offline}`}
            >
              {status.online
                ? status.text
                : formatLastSeen(status.lastSeen) || status.text}
            </p>
          </div>
        </div>

        <div className={style.actions}>
          {showSearch ? (
            <div className={style.searchWrapper}>
              <div className={style.searchInputWrapper}>
                <FiSearch size={16} className={style.searchIconInside} />
                <input
                  ref={searchInputRef}
                  type="text"
                  className={style.searchInput}
                  placeholder="Поиск в чате..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  onKeyPress={handleSearchSubmit}
                />
                {searchResults.length > 0 && (
                  <div className={style.searchCounter}>
                    <span>
                      {currentResultIndex + 1} / {searchResults.length}
                    </span>
                    <button onClick={handlePrevResult} title="Предыдущее">
                      <FiChevronUp size={14} />
                    </button>
                    <button onClick={handleNextResult} title="Следующее">
                      <FiChevronDown size={14} />
                    </button>
                  </div>
                )}
              </div>
              <button
                className={style.closeSearch}
                onClick={() => setShowSearch(false)}
              >
                <FiX size={16} />
              </button>
            </div>
          ) : (
            <>
              <button
                className={style.iconButton}
                onClick={handleSearch}
                title="Поиск"
              >
                <FiSearch size={20} />
              </button>

              {isPrivateChat && (
                <>
                  <button
                    className={style.iconButton}
                    onClick={() => handleCall("audio")}
                    title="Звонок"
                  >
                    <FiPhone size={20} />
                  </button>
                  <button
                    className={style.iconButton}
                    onClick={() => handleCall("video")}
                    title="Видеозвонок"
                  >
                    <FiVideo size={20} />
                  </button>
                </>
              )}

              <button
                className={style.iconButton}
                onClick={handleInfo}
                title="Информация"
              >
                <FiInfo size={20} />
              </button>

              <div className={style.menuWrapper} ref={menuRef}>
                <button
                  className={style.menuButton}
                  onClick={() => setShowMenu(!showMenu)}
                  title="Меню"
                >
                  <FiMoreVertical size={20} />
                </button>

                {showMenu && chatInfo && (
                  <div className={style.dropdownMenu}>
                    {chatInfo.type === "group" && isCreator && (
                      <button onClick={handleAddParticipant}>
                        <FiUserPlus size={16} />
                        Добавить участника
                      </button>
                    )}
                    {chatInfo.type === "group" && !isCreator && (
                      <button
                        onClick={handleLeaveChat}
                        className={style.danger}
                      >
                        <FiLogOut size={16} />
                        Покинуть группу
                      </button>
                    )}
                    {isCreator && (
                      <button
                        onClick={handleDeleteChat}
                        className={style.danger}
                      >
                        <FiTrash2 size={16} />
                        Удалить чат
                      </button>
                    )}
                  </div>
                )}
              </div>
            </>
          )}
        </div>
      </header>

      {callModal && (
        <CallModal
          isOpen={callModal.isOpen}
          type={callModal.type}
          caller={callModal.caller}
          roomId={callModal.roomId}
          callId={callModal.callId}
          isIncoming={callModal.isIncoming}
          ws={ws}
          currentUserId={currentUserId}
          onClose={() => setCallModal(null)}
        />
      )}
    </>
  );
}

export default Header;
