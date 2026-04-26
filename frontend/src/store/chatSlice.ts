import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

export interface Chat {
  id: string;
  type: string;
  name: string | null;
  creator_id: string;
  avatar_url: string | null;
  created_at: string;
  updated_at: string;
  last_message_at: string | null;
  last_message?: {
    id: string;
    text: string;
    user_id: string;
    created_at: string;
    type?: string;
    file_url?: string;
    file_name?: string;
  } | null;
  unread_count?: number;
  other_user_id?: string | null;
  other_user_name?: string | null;
}

interface ChatState {
  chats: Chat[];
  loading: boolean;
}

const initialState: ChatState = {
  chats: [],
  loading: false,
};

const chatSlice = createSlice({
  name: "chats",
  initialState,
  reducers: {
    setChats: (state, action: PayloadAction<Chat[]>) => {
      state.chats = action.payload;
      state.loading = false;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    updateChatLastMessage: (
      state,
      action: PayloadAction<{ chatId: string; lastMessage: any }>,
    ) => {
      const { chatId, lastMessage } = action.payload;
      const chatIndex = state.chats.findIndex((chat) => chat.id === chatId);

      if (chatIndex !== -1) {
        state.chats[chatIndex].last_message = {
          id: lastMessage.id,
          text: lastMessage.text,
          user_id: lastMessage.user_id,
          created_at: lastMessage.created_at,
          type: lastMessage.type,
          file_url: lastMessage.file_url,
          file_name: lastMessage.file_name,
        };
        state.chats[chatIndex].last_message_at = lastMessage.created_at;
      }
    },
    addChat: (state, action: PayloadAction<Chat>) => {
      // Проверяем, нет ли уже такого чата
      const exists = state.chats.some((chat) => chat.id === action.payload.id);
      if (!exists) {
        state.chats.unshift(action.payload);
      }
    },
    removeChat: (state, action: PayloadAction<string>) => {
      state.chats = state.chats.filter((chat) => chat.id !== action.payload);
    },
    clearChats: (state) => {
      state.chats = [];
      state.loading = false;
    },
  },
});

export const {
  setChats,
  setLoading,
  updateChatLastMessage,
  addChat,
  removeChat,
  clearChats,
} = chatSlice.actions;

// Селекторы с правильными типами для RootState
export const selectChats = (state: { chats: ChatState }) => state.chats.chats;
export const selectChatsLoading = (state: { chats: ChatState }) =>
  state.chats.loading;

export default chatSlice.reducer;
