import React, {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
  ReactNode,
  useCallback,
} from "react";

interface WebSocketContextType {
  ws: WebSocket | null;
  sendMessage: (data: any) => void;
  lastMessage: any;
}

const WebSocketContext = createContext<WebSocketContextType>({
  ws: null,
  sendMessage: () => {},
  lastMessage: null,
});

export const useWebSocket = () => useContext(WebSocketContext);

export function WebSocketProvider({ children }: { children: ReactNode }) {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [lastMessage, setLastMessage] = useState<any>(null);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const token = localStorage.getItem("token");
    const userStr = localStorage.getItem("user");

    if (!token || !userStr) return;

    const user = JSON.parse(userStr);
    const userId = user.id;

    const socket = new WebSocket(
      `ws://localhost:8080/ws?token=${token}&user_id=${userId}`,
    );

    socket.onopen = () => {
      console.log("Global WebSocket connected");
      socket.send(
        JSON.stringify({
          type: "user-status",
          user_id: userId,
          status: "online",
        }),
      );
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        setLastMessage(data);
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    socket.onclose = () => {
      console.log("Global WebSocket disconnected");
    };

    socket.onerror = (error) => {
      console.error("Global WebSocket error:", error);
    };

    wsRef.current = socket;
    setWs(socket);

    window.addEventListener("beforeunload", () => {
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(
          JSON.stringify({
            type: "user-status",
            user_id: userId,
            status: "offline",
          }),
        );
      }
    });

    return () => {
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(
          JSON.stringify({
            type: "user-status",
            user_id: userId,
            status: "offline",
          }),
        );
        socket.close();
      }
    };
  }, []);

  const sendMessage = useCallback((data: any) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data));
    }
  }, []);

  return (
    <WebSocketContext.Provider value={{ ws, sendMessage, lastMessage }}>
      {children}
    </WebSocketContext.Provider>
  );
}
