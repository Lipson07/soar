import { configureStore } from "@reduxjs/toolkit";
import { persistStore, persistReducer } from "redux-persist";
import storage from "./storage";
import user from "./userSlice";
import modal from "./modalChatSlice";
import search from "./searchSlice";
import selectedChat from "./selectedChatSlice";
import settings from "./settingSlice";

const persistUserConfig = {
  key: "user",
  storage,
  whitelist: ["user", "isAuthenticated"],
};

const persistSelectedChatConfig = {
  key: "selectedChat",
  storage,
  whitelist: ["currentChat", "isOpen"],
};

const persistSettingsConfig = {
  key: "settings",
  storage,
  whitelist: [
    "theme",
    "fontSize",
    "desktopNotifications",
    "messageSound",
    "soundVolume",
    "showMessagePreview",
    "showOnlineStatus",
    "readReceipts",
    "typingIndicator",
  ],
};

const persistedUserReducer = persistReducer(persistUserConfig, user);
const persistedSelectedChatReducer = persistReducer(
  persistSelectedChatConfig,
  selectedChat,
);
const persistedSettingsReducer = persistReducer(
  persistSettingsConfig,
  settings,
);

export const store = configureStore({
  reducer: {
    user: persistedUserReducer,
    modal: modal,
    search: search,
    selectedChat: persistedSelectedChatReducer,
    settings: persistedSettingsReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ["persist/PERSIST", "persist/REHYDRATE"],
      },
    }),
});

export const persistor = persistStore(store);

// После rehydrate применяем тему
persistor.subscribe(() => {
  const state = store.getState();
  if (state.settings?.theme) {
    document.documentElement.setAttribute("data-theme", state.settings.theme);
  }
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;
