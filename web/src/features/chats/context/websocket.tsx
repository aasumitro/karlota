import React, {createContext, useContext, useEffect, useCallback, useRef} from "react";
import { Chat, Message } from "@/types/chat";
import { useChatStore } from "@/states/chat-store";
import { isJsonString } from "@/lib/utils";
import { API_ENDPOINT } from "@/lib/api";
import {useGlobalActionStore} from "@/states/global-action-store";
import {ChatListState} from "@/features/chats/components/chat-list";
import {toast} from "@/hooks/use-toast";
import {Profile} from "@/types/user.ts";

const WebSocketContext = createContext<any>(null);

export const WebSocketProvider = ({ children }: { children: React.ReactNode }) => {
  const { setChats, setMessages, setNewMessage, setOnlineStatus } = useChatStore();
  const { setState } = useGlobalActionStore();
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimer = useRef<NodeJS.Timeout | null>(null);

  // Connect WebSocket only once
  const connectWebSocket = useCallback(() => {
    if (wsRef.current) return;

    wsRef.current = new WebSocket(API_ENDPOINT.MESSAGING.WEBSOCKET);

    wsRef.current.onopen = () => {
      setState(ChatListState, true)
      wsRef.current?.send(JSON.stringify({action: "chats"}));
    };

    wsRef.current.onmessage = (event) => {
      if (!isJsonString(event.data)) return;
      const callback = JSON.parse(event.data);

      switch (callback.type) {
        case "error":
          toast({variant: 'destructive', title: callback.data})
          break;
        case "chats":
          setChats(callback.data as Chat[]);
          setState(ChatListState, false)
          break;
        case "refresh_chat":
          setState(ChatListState, true)
          setTimeout(() => {wsRef.current?.send(JSON.stringify({ action: "chats" }))}, 500);
          break;
        case "messages":
          setMessages(callback.data as Message[]);
          break;
        case "online_status":
          const data = callback.data
          setOnlineStatus(data.user as Profile, data.conversation_id as number)
          break;
        case "typing":
          setState(`typing_${callback.data.conversation_id}`,
            callback.data.status as boolean)
          break;
        case "new_message":
          setNewMessage(callback.data as Message);
          break;
        default:
          console.warn("Unknown WebSocket event:", callback);
          break;
      }
    };

    wsRef.current.onclose = () => {
      wsRef.current = null;
      if (!reconnectTimer.current) {
        reconnectTimer.current = setTimeout(connectWebSocket, 3000); // Reconnect after 3 seconds
      }
    };
  }, []);

  const sendMessage = useCallback((payload: any) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) return;
    wsRef.current.send(JSON.stringify(payload));
  }, []);

  useEffect(() => {
    connectWebSocket();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
      if (reconnectTimer.current) {
        clearTimeout(reconnectTimer.current);
        reconnectTimer.current = null;
      }
    };
  }, [connectWebSocket]);

  return (
    <WebSocketContext.Provider value={{ sendMessage }}>
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider");
  }
  return context;
};