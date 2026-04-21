import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

interface Message {
  id: string;
  chat_id: string;
  user_id: string;
  type: "text" | "image" | "file";
  text: string;
  file_url?: string;
  file_name?: string;
  file_size?: number;
  mime_type?: string;
  reply_to: string | null;
  is_edited: boolean;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
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
  unread_count: number;
  last_message?: {
    id: string;
    text: string;
    user_id: string;
    created_at: string;
  } | null;
  other_user_id?: string | null;
  other_user_name?: string | null;
}

interface SelectedChatState {
  currentChat: Chat | null;
  messages: Message[];
  isLoading: boolean;
  isOpen: boolean;
}

const initialState: SelectedChatState = {
  currentChat: null,
  messages: [],
  isLoading: false,
  isOpen: false,
};

const selectedChatSlice = createSlice({
  name: "selectedChat",
  initialState,
  reducers: {
    openChat: (state, action: PayloadAction<Chat>) => {
      state.currentChat = action.payload;
      state.isOpen = true;
      state.messages = [];
    },
    closeChat: (state) => {
      state.currentChat = null;
      state.isOpen = false;
      state.messages = [];
    },
    toggleChat: (state, action: PayloadAction<Chat>) => {
      if (state.isOpen && state.currentChat?.id === action.payload.id) {
        state.isOpen = false;
        state.currentChat = null;
        state.messages = [];
      } else {
        state.currentChat = action.payload;
        state.isOpen = true;
        state.messages = [];
      }
    },
    setMessages: (state, action: PayloadAction<Message[]>) => {
      state.messages = action.payload;
    },
    addMessage: (state, action: PayloadAction<Message>) => {
      state.messages.push(action.payload);
    },
    updateMessageText: (
      state,
      action: PayloadAction<{ id: string; text: string; is_edited: boolean }>,
    ) => {
      const index = state.messages.findIndex((m) => m.id === action.payload.id);
      if (index !== -1) {
        state.messages[index].text = action.payload.text;
        state.messages[index].is_edited = action.payload.is_edited;
        state.messages[index].updated_at = new Date().toISOString();
      }
    },
    deleteMessage: (state, action: PayloadAction<string>) => {
      state.messages = state.messages.filter((m) => m.id !== action.payload);
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    clearSelectedChat: (state) => {
      state.currentChat = null;
      state.messages = [];
      state.isOpen = false;
    },
  },
});

export const {
  openChat,
  closeChat,
  toggleChat,
  setMessages,
  addMessage,
  updateMessageText,
  deleteMessage,
  setLoading,
  clearSelectedChat,
} = selectedChatSlice.actions;

export const selectCurrentChat = (state: { selectedChat: SelectedChatState }) =>
  state.selectedChat.currentChat;

export const selectChatMessages = (state: {
  selectedChat: SelectedChatState;
}) => state.selectedChat.messages;

export const selectChatLoading = (state: { selectedChat: SelectedChatState }) =>
  state.selectedChat.isLoading;

export const selectIsChatOpen = (state: { selectedChat: SelectedChatState }) =>
  state.selectedChat.isOpen;

export default selectedChatSlice.reducer;
