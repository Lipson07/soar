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
import { useWebSocket } from "../../../context/WebSocketContext";
import CallModal from "../CallModal/CallModal";

interface User {
  id: string;
  username: string;
  avatar_url: string | null;
  status: string;
  last_seen: string | null;
}

interface ChatInfo {
  id: string;
  name: string | null;
  type: string;
  creator_id: string;
  participants_count?: number;
  other_user_id?: string;
  other_user_name?: string;
  avatar_url?: string | null;
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
  const { lastMessage, sendMessage } = useWebSocket();

  const [otherUser, setOtherUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [showMenu, setShowMenu] = useState(false);
  const [showSearch, setShowSearch] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<Message[]>([]);
  const [currentResultIndex, setCurrentResultIndex] = useState(0);
  const [chatDetails, setChatDetails] = useState<ChatInfo | null>(null);
  const [callModal, setCallModal] = useState<{
    isOpen: boolean;
    type: "audio" | "video";
    callId: string;
    roomId: string;
    caller: any;
    isIncoming: boolean;
  } | null>(null);
  const [isOnline, setIsOnline] = useState(false);
  const [lastSeen, setLastSeen] = useState<string | null>(null);
  const [isTyping, setIsTyping] = useState(false);
  const typingTimeoutRef = useRef<NodeJS.Timeout>();

  const menuRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const currentUser = JSON.parse(localStorage.getItem("user") || "{}");
  const currentUserId = currentUser.id;

  useEffect(() => {
    if (chatInfo && chatInfo.type === "private" && chatInfo.other_user_id) {
      fetchUserInfo(chatInfo.other_user_id);
    } else if (chatInfo && chatInfo.type === "group") {
      fetchChatDetails();
      setIsTyping(false);
    } else {
      setOtherUser(null);
      setChatDetails(null);
      setIsTyping(false);
    }
  }, [chatInfo]);

  useEffect(() => {
    if (!lastMessage) return;

    switch (lastMessage.type) {
      case "user-status":
        if (
          chatInfo?.type === "private" &&
          lastMessage.user_id === chatInfo.other_user_id
        ) {
          setIsOnline(lastMessage.status === "online");
          if (lastMessage.last_seen) setLastSeen(lastMessage.last_seen);
        }
        break;
      case "typing":
        if (
          chatInfo?.type === "private" &&
          lastMessage.chat_id === chatInfo.id &&
          lastMessage.user_id === chatInfo.other_user_id
        ) {
          if (lastMessage.is_typing) {
            setIsTyping(true);
            if (typingTimeoutRef.current)
              clearTimeout(typingTimeoutRef.current);
            typingTimeoutRef.current = setTimeout(
              () => setIsTyping(false),
              3000,
            );
          } else {
            setIsTyping(false);
          }
        }
        break;
      case "call-start":
        if (lastMessage.call) {
          setCallModal({
            isOpen: true,
            type: lastMessage.call.type,
            callId: lastMessage.call.id,
            roomId: lastMessage.call.room_id,
            caller: lastMessage.call.caller,
            isIncoming: true,
          });
        }
        break;
      case "call-accept":
        setCallModal((prev) => (prev ? { ...prev, isIncoming: false } : null));
        break;
      case "call-reject":
      case "call-end":
        setCallModal(null);
        break;
    }
  }, [lastMessage, chatInfo]);

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
        setIsOnline(user.status === "online");
        setLastSeen(user.last_seen);
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
        {
          headers: { Authorization: `Bearer ${token}` },
        },
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

    sendMessage({
      type: "call-start",
      call_id: callId,
      room_id: roomId,
      chat_id: chatInfo.id,
      callee_id: otherUser.id,
      call_type: callType,
    });

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
  };

  const getDisplayName = () => {
    if (!chatInfo) return "Чат";
    if (chatInfo.type !== "private") return chatInfo.name || "Групповой чат";
    if (otherUser?.username) return otherUser.username;
    if (chatInfo.other_user_name) return chatInfo.other_user_name;
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
    return <span>{getDisplayName().charAt(0).toUpperCase()}</span>;
  };

  const getStatusText = () => {
    if (isTyping) return "✍️ Печатает...";
    if (chatInfo?.type === "private") {
      if (isOnline) return "● Онлайн";
      if (lastSeen) {
        const date = new Date(lastSeen);
        const diff = Math.floor((Date.now() - date.getTime()) / 1000);
        if (diff < 60) return "Был(а) только что";
        if (diff < 3600) return `Был(а) ${Math.floor(diff / 60)} мин назад`;
        return `Был(а) ${Math.floor(diff / 3600)} ч назад`;
      }
      return "○ Офлайн";
    }
    if (chatInfo?.type === "group" && chatDetails) {
      return `👥 ${chatDetails.participants_count} участников`;
    }
    return "";
  };

  return (
    <header className={style.header}>
      <div className={style.userInfo}>
        <div className={style.avatar}>{getAvatar()}</div>
        <div className={style.userDetails}>
          <p className={style.userName}>{getDisplayName()}</p>
          <p
            className={`${style.userStatus} ${isTyping ? style.typing : isOnline ? style.online : style.offline}`}
          >
            {getStatusText()}
          </p>
        </div>
      </div>
      <div className={style.actions}>{/* ... кнопки действий ... */}</div>
      {callModal && (
        <CallModal
          isOpen={callModal.isOpen}
          type={callModal.type}
          caller={callModal.caller}
          roomId={callModal.roomId}
          callId={callModal.callId}
          isIncoming={callModal.isIncoming}
          ws={null}
          currentUserId={currentUserId}
          onClose={() => setCallModal(null)}
        />
      )}
    </header>
  );
}

export default Header;
