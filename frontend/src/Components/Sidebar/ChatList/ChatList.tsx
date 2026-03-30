import React from "react";
import style from "./ChatList.module.scss";
import SearchBar from "../SearchBar/SearchBar";
import StoriesSection from "../StoriesSection/StoriesSection";

function ChatList() {
  return (
    <div className={style.chatList}>
      <div className={style.header}>
        <div className={style.titleWrapper}>
          <h1>Чаты</h1>
          <button className={style.addButton}>
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
          <button className={style.createChatBtn}>Создать чат</button>
        </div>
      </div>
    </div>
  );
}

export default ChatList;
