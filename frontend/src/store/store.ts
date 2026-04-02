import { configureStore } from "@reduxjs/toolkit";
import { persistStore, persistReducer } from "redux-persist";
import storage from "./storage";
import user from "./userSlice";
import modal from "./modalChatSlice";
import search from "./searchSlice";

const persistUserConfig = {
  key: "user",
  storage,
  whitelist: ["user", "isAuthenticated"],
};

const persistedUserReducer = persistReducer(persistUserConfig, user);

export const store = configureStore({
  reducer: {
    user: persistedUserReducer,
    modal: modal,
    search: search,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ["persist/PERSIST", "persist/REHYDRATE"],
      },
    }),
});

export const persistor = persistStore(store);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;
