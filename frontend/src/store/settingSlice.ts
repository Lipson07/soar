import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

export type ThemeType = "dark" | "light";
export type FontSizeType = "small" | "medium" | "large";

interface SettingsState {
  theme: ThemeType;
  fontSize: FontSizeType;
  desktopNotifications: boolean;
  messageSound: boolean;
  soundVolume: number;
  showMessagePreview: boolean;
  showOnlineStatus: boolean;
  readReceipts: boolean;
  typingIndicator: boolean;
}

const initialState: SettingsState = {
  theme: "dark",
  fontSize: "medium",
  desktopNotifications: true,
  messageSound: true,
  soundVolume: 70,
  showMessagePreview: true,
  showOnlineStatus: true,
  readReceipts: true,
  typingIndicator: true,
};

const applyTheme = (theme: ThemeType) => {
  if (typeof document !== "undefined") {
    document.documentElement.setAttribute("data-theme", theme);
    // Сохраняем в localStorage напрямую для надежности
    localStorage.setItem("theme", theme);
  }
};

const applyFontSize = (size: FontSizeType) => {
  if (typeof document !== "undefined") {
    const sizes: Record<FontSizeType, string> = {
      small: "12px",
      medium: "14px",
      large: "16px",
    };
    const root = document.documentElement;
    root.style.setProperty("--font-size-base", sizes[size]);
    root.style.setProperty(
      "--font-size-small",
      `${parseInt(sizes[size]) - 2}px`,
    );
    root.style.setProperty(
      "--font-size-large",
      `${parseInt(sizes[size]) + 2}px`,
    );
    root.style.setProperty("--font-size-xl", `${parseInt(sizes[size]) + 6}px`);
    localStorage.setItem("fontSize", size);
  }
};

const settingsSlice = createSlice({
  name: "settings",
  initialState,
  reducers: {
    setTheme: (state, action: PayloadAction<ThemeType>) => {
      state.theme = action.payload;
      applyTheme(action.payload);
    },
    setFontSize: (state, action: PayloadAction<FontSizeType>) => {
      state.fontSize = action.payload;
      applyFontSize(action.payload);
    },
    setDesktopNotifications: (state, action: PayloadAction<boolean>) => {
      state.desktopNotifications = action.payload;
      if (typeof window !== "undefined" && action.payload) {
        Notification.requestPermission();
      }
    },
    setMessageSound: (state, action: PayloadAction<boolean>) => {
      state.messageSound = action.payload;
    },
    setSoundVolume: (state, action: PayloadAction<number>) => {
      state.soundVolume = action.payload;
    },
    setShowMessagePreview: (state, action: PayloadAction<boolean>) => {
      state.showMessagePreview = action.payload;
    },
    setShowOnlineStatus: (state, action: PayloadAction<boolean>) => {
      state.showOnlineStatus = action.payload;
    },
    setReadReceipts: (state, action: PayloadAction<boolean>) => {
      state.readReceipts = action.payload;
    },
    setTypingIndicator: (state, action: PayloadAction<boolean>) => {
      state.typingIndicator = action.payload;
    },
    resetSettings: (state) => {
      state.theme = "dark";
      state.fontSize = "medium";
      state.desktopNotifications = true;
      state.messageSound = true;
      state.soundVolume = 70;
      state.showMessagePreview = true;
      state.showOnlineStatus = true;
      state.readReceipts = true;
      state.typingIndicator = true;
      applyTheme("dark");
      applyFontSize("medium");
    },
  },
});

export const {
  setTheme,
  setFontSize,
  setDesktopNotifications,
  setMessageSound,
  setSoundVolume,
  setShowMessagePreview,
  setShowOnlineStatus,
  setReadReceipts,
  setTypingIndicator,
  resetSettings,
} = settingsSlice.actions;

export const selectTheme = (state: { settings: SettingsState }) =>
  state.settings.theme;
export const selectFontSize = (state: { settings: SettingsState }) =>
  state.settings.fontSize;
export const selectAllSettings = (state: { settings: SettingsState }) =>
  state.settings;

export default settingsSlice.reducer;
