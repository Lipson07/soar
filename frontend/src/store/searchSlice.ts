import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";

interface User {
  id: string;
  username: string;
  email: string;
  role?: string;
  avatar_url?: string | null;
  last_seen?: string | null;
  created_at?: string;
  updated_at?: string;
  status?: string;
  is_online?: boolean;
}

interface SearchState {
  query: string;
  results: User[];
  loading: boolean;
  error: string | null;
}

const initialState: SearchState = {
  query: "",
  results: [],
  loading: false,
  error: null,
};

export const searchUsers = createAsyncThunk(
  "search/searchUsers",
  async (query: string, { getState }) => {
    if (!query.trim()) return [];

    const token = localStorage.getItem("token");

    const response = await fetch(
      `http://localhost:8080/api/users/search?q=${encodeURIComponent(query)}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      },
    );

    if (response.status === 401) {
      throw new Error("Не авторизован. Пожалуйста, войдите снова.");
    }

    if (response.status === 400) {
      throw new Error("Некорректный поисковый запрос");
    }

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    const users = Array.isArray(data) ? data : data.users || [];

    return users as User[];
  },
);

const searchSlice = createSlice({
  name: "search",
  initialState,
  reducers: {
    setSearchQuery: (state, action) => {
      state.query = action.payload;
    },
    clearResults: (state) => {
      state.results = [];
      state.query = "";
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(searchUsers.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(searchUsers.fulfilled, (state, action) => {
        state.loading = false;
        state.results = action.payload;
      })
      .addCase(searchUsers.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || "Ошибка поиска";
      });
  },
});

export const { setSearchQuery, clearResults } = searchSlice.actions;

export const selectSearchQuery = (state: { search: SearchState }) =>
  state.search.query;
export const selectSearchResults = (state: { search: SearchState }) =>
  state.search.results || [];
export const selectSearchLoading = (state: { search: SearchState }) =>
  state.search.loading;
export const selectSearchError = (state: { search: SearchState }) =>
  state.search.error;

export default searchSlice.reducer;
