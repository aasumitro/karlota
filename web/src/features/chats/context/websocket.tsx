/* eslint-disable @typescript-eslint/no-explicit-any,react-hooks/exhaustive-deps */
import {createContext, useContext, useEffect, useCallback, useRef, ReactNode} from "react";
import {CallStage, Chat, Message} from "@/types/chat";
import { useChatStore } from "@/states/chat-store";
import { isJsonString } from "@/lib/utils";
import { API_ENDPOINT } from "@/lib/api";
import {useGlobalActionStore} from "@/states/global-action-store";
import {ChatListState} from "@/features/chats/components/chat-list";
import {Profile} from "@/types/user";
import {IncomingCallModalState} from "@/features/chats/components/incoming-call";
import {NewCallState} from "@/features/chats/components/new-call";
// import {toast} from "@/hooks/use-toast";

interface WebSocketContext {
  sendMessage: (payload: any) => void;
}

const WebSocketContext = createContext<WebSocketContext | null>(null);

export const WebSocketProvider = ({ children }: { children: ReactNode }) => {
  const { setChats, setMessages, setNewMessage, setOnlineStatus, setCall } = useChatStore();
  const { setState,setStatus } = useGlobalActionStore();
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
          // toast({variant: 'destructive', title: callback.data})
          console.log("error callback", callback.data)
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
          setOnlineStatus(callback.data.user as Profile, callback.data.conversation_id as number)
          break;
        case "typing_state":
          setState(`typing_${callback.data.conversation_id}`, callback.data.status as boolean)
          break;
        case "new_message":
          setNewMessage(callback.data as Message);
          break;
        case "incoming_call":
          // TODO: fix
          setCall({vc_type: callback.data.type, vc_data: callback.data.payload,
            vc_action: callback.data.action, vc_caller: callback.data.recipient})
          setState(IncomingCallModalState, true)
          break;
        case "answer_call":
          // TODO: FIx
          setStatus(NewCallState, callback.data.action as CallStage)
          break;
        default:
          break;
      }
    };

    wsRef.current.onclose = () => {
      if (wsRef.current?.readyState !== WebSocket.CLOSED) return;
      wsRef.current = null;
      if (reconnectTimer.current) clearTimeout(reconnectTimer.current);
      if (!reconnectTimer.current) reconnectTimer.current = setTimeout(connectWebSocket, 3000);
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

// eslint-disable-next-line react-refresh/only-export-components
export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider");
  }
  return context;
};