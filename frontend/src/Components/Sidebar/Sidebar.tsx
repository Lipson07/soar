import React from "react";
import style from "./Sidebar.module.scss";
import SidebarHeader from "./SidebarHeader/SidebarHeader";
import StoriesSection from "./StoriesSection/StoriesSection";
import SearchBar from "./SearchBar/SearchBar";
function Sidebar() {
  return (
    <div className={style.sidebar}>
      <SidebarHeader />
      <div className={style.container}>
        <StoriesSection />
        <SearchBar />
      </div>
    </div>
  );
}

export default Sidebar;
