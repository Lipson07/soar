import React from "react";
import style from "./Mainpage.module.scss";
import Sidebar from "../../Sidebar/Sidebar";
function MainPage() {
  return (
    <main className={style.main}>
      <Sidebar />
    </main>
  );
}

export default MainPage;
