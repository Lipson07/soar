import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  avatar_path: string | boolean;
  last_seen_at: string | null;
  created_at: string;
  updated_at: string;
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
  async (query: string) => {
    if (!query.trim()) return [];

    const response = await fetch(
      `http://localhost:8080/api/users/search?query=${encodeURIComponent(query)}`,
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data as User[];
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
