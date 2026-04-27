import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { PersistGate } from "redux-persist/integration/react";
import { BrowserRouter } from "react-router-dom";
import App from "./App.tsx";
import store, { persistor } from "./store/store.ts";
import { WebSocketProvider } from "./context/WebSocketContext";
import "./index.scss";

const applyStoredTheme = () => {
  try {
    const persistedData = localStorage.getItem("persist:settings");
    if (persistedData) {
      const parsed = JSON.parse(persistedData);
      const settings = JSON.parse(parsed.settings || "{}");

      if (settings.theme) {
        document.documentElement.setAttribute("data-theme", settings.theme);
      } else {
        document.documentElement.setAttribute("data-theme", "dark");
      }

      if (settings.fontSize) {
        const sizes: Record<string, string> = {
          small: "12px",
          medium: "14px",
          large: "16px",
        };
        const size = sizes[settings.fontSize] || "14px";
        document.documentElement.style.setProperty("--font-size-base", size);
        document.documentElement.style.setProperty(
          "--font-size-small",
          `${parseInt(size) - 2}px`,
        );
        document.documentElement.style.setProperty(
          "--font-size-large",
          `${parseInt(size) + 2}px`,
        );
        document.documentElement.style.setProperty(
          "--font-size-xl",
          `${parseInt(size) + 6}px`,
        );
      }
    } else {
      document.documentElement.setAttribute("data-theme", "dark");
    }
  } catch (e) {
    console.warn("Failed to apply stored theme:", e);
    document.documentElement.setAttribute("data-theme", "dark");
  }
};

applyStoredTheme();

createRoot(document.getElementById("root")!).render(
  <Provider store={store}>
    <PersistGate loading={null} persistor={persistor}>
      <BrowserRouter>
        <WebSocketProvider>
          <App />
        </WebSocketProvider>
      </BrowserRouter>
    </PersistGate>
  </Provider>,
);
