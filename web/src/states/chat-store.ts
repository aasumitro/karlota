import { create } from 'zustand'
import {Call, Chat, Message} from "@/types/chat";
import {useGlobalActionStore} from "@/states/global-action-store";
import {RefreshChatListState} from "@/features/chats/components/chat-list";
import {Profile} from "@/types/user";

interface States {
  chats: Chat[];
  selectedChat: Chat | null;
  mobileSelectedChat: Chat | null;
  messages: Message[] | null;
  call: Call | null;
}

interface Actions {
  setChats: (chats: Chat[]) => void;
  setSelectedChat: (chat: Chat | null) => void;
  setMessages: (messages: Message[] | null) => void;
  setNewMessage: (message: Message) => void
  setOnlineStatus: (user: Profile, chatId: number) => void
  setCall: (call: Call | null) => void;
}

export const useChatStore = create<States & Actions>((set) => ({
  chats: [],
  selectedChat: null,
  mobileSelectedChat: null,
  messages: [],
  call: null,

  setChats: (chats:  Chat[]) =>
    set((state) => {
      return { ...state, chats }
    }),
  setSelectedChat: (chat: Chat | null) =>
    set((state) => ({
      ...state,
      selectedChat: chat,
      mobileSelectedChat: chat,
    })),

  setMessages: (messages:  Message[] | null) =>
    set((state) => {
      return { ...state, messages }
    }),

  setNewMessage: (message: Message) => {
    set((state) => {
      const updatedChats = state.chats.map((chat) => {
        if (chat.id === message.conversation_id) {
          return {
            ...chat,
            messages: [message],
          };
        }
        return chat;
      });
      const newState = {
        ...state,
        chats: updatedChats,
        messages: state.messages ? [...state.messages, message] : [message],
      };

      const chat = state.chats?.find(
        (chat) => chat.id === message.conversation_id);
      if (!chat) {
        useGlobalActionStore.getState().setState(RefreshChatListState, true)
      }

      if (state.selectedChat) {
        newState.selectedChat = {
          ...state.selectedChat,
          messages: [message],
        };
      }
      return newState;
    });
  },

  setOnlineStatus: (user: Profile, chatId: number) => {
    set((state) => {
      const updatedChats = state.chats.map((chat) => {
        if (chat.id === chatId) {
          return {
            ...chat,
            users: chat.users.map((u) =>
              u.id === user.id ? { ...u, is_online: user.is_online, last_online: user.last_online } : u
            ),
          };
        }
        return chat;
      });

      const updatedState = { ...state, chats: updatedChats };

      if (state.selectedChat && state.selectedChat.id === chatId) {
        updatedState.selectedChat = {
          ...state.selectedChat,
          users: state.selectedChat.users.map((u) =>
            u.id === user.id ? { ...u, is_online: user.is_online, last_online: user.last_online } : u
          ),
        };
      }

      return updatedState;
    });
  },

  setCall: (call:  Call | null) =>
    set((state) => {
      return { ...state, call }
    }),
}));