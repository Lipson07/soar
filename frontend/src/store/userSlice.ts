import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

interface User {
  id: string;
  username: string;
  email: string;
  password?: string;
  role?: string;
  avatar_url?: string | null;
  last_seen?: string | null;
  created_at: string;
  updated_at: string;
  status?: string;
}

interface UserState {
  user: User | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
  token: string | null;
}

const loadInitialState = (): UserState => {
  try {
    const token = localStorage.getItem("token");
    const userStr = localStorage.getItem("user");

    if (token && userStr) {
      const user = JSON.parse(userStr);
      return {
        user: { ...user, status: "online" },
        token,
        isAuthenticated: true,
        loading: false,
        error: null,
      };
    }
  } catch (error) {
    console.error("Error loading auth state:", error);
  }

  return {
    user: null,
    isAuthenticated: false,
    loading: false,
    error: null,
    token: null,
  };
};

const initialState: UserState = loadInitialState();

const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setUser: (state, action: PayloadAction<{ user: User; token: string }>) => {
      state.user = { ...action.payload.user, status: "online" };
      state.token = action.payload.token;
      state.isAuthenticated = true;
      state.error = null;

      localStorage.setItem("token", action.payload.token);
      localStorage.setItem(
        "user",
        JSON.stringify({ ...action.payload.user, status: "online" }),
      );
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setError: (state, action: PayloadAction<string>) => {
      state.error = action.payload;
      state.loading = false;
    },
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      state.loading = false;
      state.error = null;
      localStorage.removeItem("token");
      localStorage.removeItem("user");
    },
    updateUser: (state, action: PayloadAction<Partial<User>>) => {
      if (state.user) {
        state.user = { ...state.user, ...action.payload };
        localStorage.setItem("user", JSON.stringify(state.user));
      }
    },
    setUserStatus: (state, action: PayloadAction<string>) => {
      if (state.user) {
        state.user.status = action.payload;
        localStorage.setItem("user", JSON.stringify(state.user));
      }
    },
  },
});

export const {
  setUser,
  setLoading,
  setError,
  logout,
  updateUser,
  setUserStatus,
} = userSlice.actions;

export const selectUser = (state: { user: UserState }) => state.user.user;
export const selectIsAuthenticated = (state: { user: UserState }) =>
  state.user.isAuthenticated;
export const selectUserLoading = (state: { user: UserState }) =>
  state.user.loading;
export const selectUserError = (state: { user: UserState }) => state.user.error;
export const selectToken = (state: { user: UserState }) => state.user.token;

export default userSlice.reducer;
