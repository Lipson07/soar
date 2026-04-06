import React from "react";
import {
  FiPhone,
  FiVideo,
  FiSearch,
  FiMoreVertical,
  FiInfo,
} from "react-icons/fi";
import style from "./Header.module.scss";
import { selectCurrentChat } from "../../../store/selectedChatSlice";
import { useSelector } from "react-redux";

function Header() {
  const chatInfo = useSelector(selectCurrentChat);

  const getDisplayName = () => {
    if (!chatInfo) return "Чат";

    // Для групповых чатов
    if (chatInfo.type !== "private") {
      return chatInfo.name || "Групповой чат";
    }

    // Для приватных чатов - показываем имя собеседника
    if (chatInfo.other_user_name) {
      return chatInfo.other_user_name;
    }

    // Если чат только создан и данных еще нет
    if (chatInfo.name && chatInfo.name !== "Приватный чат") {
      return chatInfo.name;
    }

    return "Пользователь";
  };

  const getAvatarLetter = () => {
    const name = getDisplayName();
    if (name === "Чат" || name === "Групповой чат" || name === "Пользователь")
      return "?";
    return name.charAt(0).toUpperCase();
  };

  const isOnline = true; // Потом можно доработать

  return (
    <header className={style.header}>
      <div className={style.userInfo}>
        <div className={style.avatar}>
          <span>{getAvatarLetter()}</span>
        </div>
        <div className={style.userDetails}>
          <p className={style.userName}>{getDisplayName()}</p>
          <p
            className={`${style.userStatus} ${isOnline ? style.online : style.offline}`}
          >
            {isOnline ? "● Онлайн" : "○ Офлайн"}
          </p>
        </div>
      </div>

      <div className={style.actions}>
        <button className={style.iconButton} title="Поиск">
          <FiSearch size={20} />
        </button>
        <button className={style.iconButton} title="Звонок">
          <FiPhone size={20} />
        </button>
        <button className={style.iconButton} title="Видеозвонок">
          <FiVideo size={20} />
        </button>
        <button className={style.iconButton} title="Информация">
          <FiInfo size={20} />
        </button>
        <button className={style.menuButton} title="Меню">
          <FiMoreVertical size={20} />
        </button>
      </div>
    </header>
  );
}

export default Header;
