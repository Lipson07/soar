import React, { useState } from "react";
import style from "./SearchBar.module.scss";

function SearchBar() {
  const [query, setQuery] = useState<string>("");

  const handleChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setQuery(value);

    if (value.trim() === "") {
      console.log("Пустой запрос, поиск не выполняется");
      return;
    }

    try {
      const response = await fetch(
        `http://localhost:8080/api/users/search?query=${encodeURIComponent(value)}`,
      );
      const data = await response.json();
      console.log("Результаты поиска:", data);
    } catch (error) {
      console.error("Ошибка при поиске:", error);
    }
  };

  return (
    <div className={style.searchContainer}>
      <div className={style.searchWrapper}>
        <input
          type="text"
          placeholder="Поиск людей по имени"
          className={style.searchInput}
          value={query}
          onChange={handleChange}
        />
        <svg
          className={style.searchIcon}
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M15.5 15.5L19 19M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
          />
        </svg>
      </div>
    </div>
  );
}

export default SearchBar;
