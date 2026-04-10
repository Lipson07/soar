import React, { useState } from "react";
import style from "./Sidebar.module.scss";
import SidebarHeader from "./SidebarHeader/SidebarHeader";
import Profile from "./Profile/Profile";
import SecurityPanel from "./SecurityPanel/SecurityPanel";
import FilesPanel from "./FilesPanel/FilesPanel";
import SettingsPanel from "./SettingsPanel/SettingsPanel";
import ChatList from "./ChatList/ChatList";

function Sidebar() {
  const [isProfileOpen, setIsProfileOpen] = useState(false);
  const [isSecurityOpen, setIsSecurityOpen] = useState(false);
  const [isFilesOpen, setIsFilesOpen] = useState(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);

  return (
    <>
      <div className={style.sidebar}>
        <SidebarHeader
          onProfileClick={() => setIsProfileOpen(true)}
          onSecurityClick={() => setIsSecurityOpen(true)}
          onFilesClick={() => setIsFilesOpen(true)}
          onSettingsClick={() => setIsSettingsOpen(true)}
        />
        <div className={style.container}>
          <ChatList />
        </div>
      </div>
      <Profile isOpen={isProfileOpen} onClose={() => setIsProfileOpen(false)} />
      <SecurityPanel
        isOpen={isSecurityOpen}
        onClose={() => setIsSecurityOpen(false)}
      />
      <FilesPanel isOpen={isFilesOpen} onClose={() => setIsFilesOpen(false)} />
      <SettingsPanel
        isOpen={isSettingsOpen}
        onClose={() => setIsSettingsOpen(false)}
      />
    </>
  );
}

export default Sidebar;
