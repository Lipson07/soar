import { createSlice } from "@reduxjs/toolkit";

interface ModalState {
  isCreateChatOpen: boolean;
}

const initialState: ModalState = {
  isCreateChatOpen: false,
};

const modalChatSlice = createSlice({
  name: "modalChat",
  initialState,
  reducers: {
    openCreateChat: (state) => {
      state.isCreateChatOpen = true;
    },
    closeCreateChat: (state) => {
      state.isCreateChatOpen = false;
    },
  },
});

export const { openCreateChat, closeCreateChat } = modalChatSlice.actions;
export const selectModalChat = (state: { modal: ModalState }) =>
  state.modal.isCreateChatOpen;
export default modalChatSlice.reducer;
