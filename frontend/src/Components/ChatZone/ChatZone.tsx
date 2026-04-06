import React from "react";
import style from "./ChatZone.module.scss";
import Header from "./Header/Header";
import MessageList from "./MessageList/MessageList";
import MessageBar from "./MessageBar/MessageBar";

function ChatZone() {
  return (
    <div className={style.chatZone}>
      <Header />
      <MessageList />
      <MessageBar />
    </div>
  );
}

export default ChatZone;
