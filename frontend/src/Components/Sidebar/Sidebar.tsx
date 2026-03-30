import React from "react";
import style from "./Sidebar.module.scss";
import SidebarHeader from "./SidebarHeader/SidebarHeader";
import StoriesSection from "./StoriesSection/StoriesSection";
import SearchBar from "./SearchBar/SearchBar";
import ChatList from "./ChatList/ChatList";
function Sidebar() {
  return (
    <div className={style.sidebar}>
      <SidebarHeader />
      <div className={style.container}>
        <ChatList />
      </div>
    </div>
  );
}

export default Sidebar;
