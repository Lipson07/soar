import React, { useState } from "react";
import style from "./Sidebar.module.scss";
import SidebarHeader from "./SidebarHeader/SidebarHeader";
import Profile from "../Pages/Profile/Profile";
import ChatList from "./ChatList/ChatList";

function Sidebar() {
  const [isProfileOpen, setIsProfileOpen] = useState(false);

  return (
    <>
      <div className={style.sidebar}>
        <SidebarHeader onProfileClick={() => setIsProfileOpen(true)} />
        <div className={style.container}>
          <ChatList />
        </div>
      </div>
      <Profile isOpen={isProfileOpen} onClose={() => setIsProfileOpen(false)} />
    </>
  );
}

export default Sidebar;
