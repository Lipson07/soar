import React from "react";
import style from "./MainPage.module.scss";
import Sidebar from "../../Sidebar/Sidebar";
import CreateChat from "../../Form/CreateChat/CreateChat";
import { useSelector } from "react-redux";
import { selectUser } from "../../../store/userSlice";
import ChatZone from "../../ChatZone/ChatZone";
import {
  selectCurrentChat,
  selectIsChatOpen,
} from "../../../store/selectedChatSlice";

function MainPage() {
  const user = useSelector(selectUser);
  const chat = useSelector(selectCurrentChat);
  const isOpen = useSelector(selectIsChatOpen);

  console.log(user);
  console.log(chat);
  console.log(isOpen);

  return (
    <main className={style.main}>
      <Sidebar />
      <CreateChat />
      {isOpen && chat ? (
        <ChatZone />
      ) : (
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
          <h3>Добро пожаловать в Soar</h3>
          <p>Выберите чат для начала общения</p>
        </div>
      )}
    </main>
  );
}

export default MainPage;
