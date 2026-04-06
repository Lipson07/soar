import React, { use } from "react";
import style from "./Mainpage.module.scss";
import Sidebar from "../../Sidebar/Sidebar";
import CreateChat from "../../Form/CreateChat/CreateChat";
import { useSelector } from "react-redux";
import { selectUser } from "../../../store/userSlice";
import ChatZone from "../../ChatZone/ChatZone";
import {
  selectChatLoading,
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
      {isOpen ? <ChatZone /> : ""}
    </main>
  );
}

export default MainPage;
