import {ScrollArea} from "@/components/ui/scroll-area";
import {Fragment} from "react/jsx-runtime";
import {cn} from "@/lib/utils";
import {Avatar, AvatarFallback} from "@/components/ui/avatar";
import {Separator} from "@/components/ui/separator";
import {Chat} from "@/types/chat";
import {useAuthStore} from "@/states/auth-store";
import {ChatListSkeleton} from "@/features/chats/components/chat-list-skeleton";
import {useChatStore} from "@/states/chat-store";
import {useWebSocket} from "@/features/chats/context/websocket";
import {useGlobalActionStore} from "@/states/global-action-store";
import {useEffect} from "react";

interface ChatListProps {
  chats: Chat[]
}

export const ChatListState = "chat_list_state"
export const RefreshChatListState = "refrehs_chat_list_state"

export function ChatList({ chats }: ChatListProps) {
  const sessionID = useAuthStore.getState().auth.user?.id;
  const {states, setState} = useGlobalActionStore();
  const { selectedChat, setSelectedChat } = useChatStore();
  const { sendMessage } = useWebSocket();

  const selectChat = (chat: Chat) => {
    setSelectedChat(chat)
    sendMessage({action: "messages", conversation_id: chat.id})
    if (chat.type === "private") {
      const user = chat.users.find((u) => u.id !== sessionID)
      if (!user) return;
      sendMessage({action: "online_status", target_id: user.id, conversation_id: chat.id})
    }
  }

  useEffect(() => {
    if (states[RefreshChatListState]) {
      sendMessage({action: "chats"});
      setState(RefreshChatListState, false)
    }
  }, [states[RefreshChatListState]]);

  return (
    <ScrollArea className='h-full p-3'>
      {states[ChatListState] && <ChatListSkeleton len={25}/>}

      {!states[ChatListState] && chats.length === 0 && <>create new chat</>}

      {chats.map((chat) => {
        let name = "";
        if (chat.type === "private") {
          const otherUser = chat.users.find((user) => user.id !== sessionID);
          name = otherUser?.display_name || "Unknown";
        } else if (chat.type === "group") {
          name = chat.name;
        }

        return (
          <Fragment key={chat.id}>
            <button
              type='button'
              className={cn(
                `flex w-full rounded-md px-2 py-2 text-left text-sm hover:bg-secondary/75`,
                selectedChat?.id === chat.id && 'sm:bg-muted'
              )}
              onClick={() => selectChat(chat)}
            >
              <div className='flex gap-2'>
                <Avatar>
                  <AvatarFallback
                    className={cn(selectedChat?.id === chat.id && 'bg-white')}
                  >{name[0]?.toUpperCase() ?? "U"}</AvatarFallback>
                </Avatar>
                <div>
                    <span className='col-start-2 row-span-2 font-medium'>
                      {name}
                    </span>
                    {chat.messages?.map((message) => (
                      <span
                        key={message.id}
                        className='col-start-2 row-span-2 row-start-2 line-clamp-2 text-ellipsis text-muted-foreground'
                      >
                         {states[`typing_${chat?.id}`]
                           ? "typing . . ." :
                           (message.user_id === sessionID
                               ? `You: ${message.content}`
                               :  message.content)
                         }
                      </span>
                    ))}
                </div>
              </div>
            </button>
            <Separator className='my-1' />
          </Fragment>
        )})}
    </ScrollArea>
  )
}