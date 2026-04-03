import React, { useState, useEffect } from "react";
import { useDispatch } from "react-redux";
import { searchUsers, setSearchQuery } from "../../../store/searchSlice";
import type { AppDispatch } from "../../../store/store";
import style from "./SearchBar.module.scss";

interface SearchBarProps {
  onFocus?: () => void;
  onBlur?: () => void;
}

function SearchBar({ onFocus, onBlur }: SearchBarProps) {
  const [query, setQuery] = useState<string>("");
  const dispatch = useDispatch<AppDispatch>();

  useEffect(() => {
    if (query.trim() === "") {
      dispatch(setSearchQuery(""));
      return;
    }

    const timer = setTimeout(() => {
      dispatch(setSearchQuery(query));
      dispatch(searchUsers(query));
    }, 500);

    return () => clearTimeout(timer);
  }, [query, dispatch]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setQuery(e.target.value);
  };

  const handleFocus = () => {
    onFocus?.();
  };

  const handleBlur = () => {
    setTimeout(() => {
      onBlur?.();
    }, 200);
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
          onFocus={handleFocus}
          onBlur={handleBlur}
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
