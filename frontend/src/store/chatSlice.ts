// store/chatSlice.ts
import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

interface LastMessage {
  id: string;
  text: string;
  user_id: string;
  created_at: string;
  type?: string;
  file_url?: string;
  file_name?: string;
}

interface Chat {
  id: string;
  type: string;
  name: string | null;
  creator_id: string;
  avatar_url: string | null;
  created_at: string;
  updated_at: string;
  last_message_at: string | null;
  last_message?: LastMessage | null;
  unread_count?: number;
  other_user_id?: string | null;
  other_user_name?: string | null;
}

interface ChatsState {
  chats: Chat[];
  loading: boolean;
  error: string | null;
}

const initialState: ChatsState = {
  chats: [],
  loading: false,
  error: null,
};

const chatSlice = createSlice({
  name: "chats",
  initialState,
  reducers: {
    setChats: (state, action: PayloadAction<Chat[]>) => {
      state.chats = action.payload;
    },
    addChat: (state, action: PayloadAction<Chat>) => {
      const exists = state.chats.find((c) => c.id === action.payload.id);
      if (!exists) {
        state.chats.unshift(action.payload);
      }
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    updateChatLastMessage: (
      state,
      action: PayloadAction<{
        chatId: string;
        lastMessage: LastMessage;
      }>,
    ) => {
      const { chatId, lastMessage } = action.payload;
      const chatIndex = state.chats.findIndex((c) => c.id === chatId);
      if (chatIndex !== -1) {
        const chat = state.chats[chatIndex];
        chat.last_message = lastMessage;
        chat.last_message_at = lastMessage.created_at;
        // Перемещаем чат вверх списка
        state.chats.splice(chatIndex, 1);
        state.chats.unshift(chat);
      }
    },
    incrementUnread: (state, action: PayloadAction<{ chatId: string }>) => {
      const chat = state.chats.find((c) => c.id === action.payload.chatId);
      if (chat) {
        chat.unread_count = (chat.unread_count || 0) + 1;
        // Перемещаем чат вверх
        const chatIndex = state.chats.findIndex(
          (c) => c.id === action.payload.chatId,
        );
        if (chatIndex > 0) {
          state.chats.splice(chatIndex, 1);
          state.chats.unshift(chat);
        }
      }
    },
    resetUnread: (state, action: PayloadAction<{ chatId: string }>) => {
      const chat = state.chats.find((c) => c.id === action.payload.chatId);
      if (chat) {
        chat.unread_count = 0;
      }
    },
    removeChat: (state, action: PayloadAction<string>) => {
      state.chats = state.chats.filter((c) => c.id !== action.payload);
    },
  },
});

export const {
  setChats,
  addChat,
  setLoading,
  setError,
  updateChatLastMessage,
  incrementUnread,
  resetUnread,
  removeChat,
} = chatSlice.actions;

export const selectChats = (state: { chats: ChatsState }) => state.chats.chats;
export const selectChatsLoading = (state: { chats: ChatsState }) =>
  state.chats.loading;
export const selectChatsError = (state: { chats: ChatsState }) =>
  state.chats.error;

export default chatSlice.reducer;
