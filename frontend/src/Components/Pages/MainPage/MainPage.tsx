import React, { use } from "react";
import style from "./Mainpage.module.scss";
import Sidebar from "../../Sidebar/Sidebar";
import CreateChat from "../../Form/CreateChat/CreateChat";
import { useSelector } from "react-redux";
import { selectUser } from "../../../store/userSlice";
function MainPage() {
  const user = useSelector(selectUser);
  console.log(user);
  return (
    <main className={style.main}>
      <Sidebar />
      <CreateChat />
    </main>
  );
}

export default MainPage;
