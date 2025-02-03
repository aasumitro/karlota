import { useChatStore } from "@/states/chat-store";
import { Message } from "@/types/chat";
import { Fragment } from "react/jsx-runtime";
import { cn } from "@/lib/utils";
import { format } from "date-fns";
import { useAuthStore } from "@/states/auth-store";
import { ChatForm } from "@/features/chats/components/chat-form";
import { ChatToolbar } from "@/features/chats/components/chat-toolbar";
import { useGlobalActionStore } from "@/states/global-action-store";

export function ChatDetail() {
  const sessionID = useAuthStore.getState().auth.user?.id;
  const { states } = useGlobalActionStore();
  const { selectedChat, messages } = useChatStore();

  const currentMessage = messages?.reduce(
    (acc: Record<string, Message[]>, obj) => {
      const key = format(obj.created_at, "d MMM, yyyy");
      if (!acc[key]) {
        acc[key] = [];
      }
      acc[key].push(obj);
      return acc;
    },
    {}
  );

  const TypingIndicator = () => {
    return (
      <div className="flex items-center space-x-1.5 px-3 py-2 w-fit rounded-full bg-muted">
        <span className="w-2 h-2 rounded-full bg-foreground/50 animate-[bounce_1.4s_infinite_.2s] opacity-75" />
        <span className="w-2 h-2 rounded-full bg-foreground/50 animate-[bounce_1.4s_infinite_.4s] opacity-75" />
        <span className="w-2 h-2 rounded-full bg-foreground/50 animate-[bounce_1.4s_infinite_.6s] opacity-75" />
      </div>
    );
  };

  return (
    <>
      <ChatToolbar />

      <div className="flex flex-1 flex-col gap-2 rounded-md px-4 pb-4 pt-0">
        <div className="flex size-full flex-1">
          <div className="chat-text-container relative -mr-4 flex flex-1 flex-col overflow-y-hidden">
            <div className="chat-flex flex h-40 w-full flex-grow flex-col-reverse justify-start gap-4 overflow-y-auto py-2 pb-4 pr-4">
              {states[`typing_${selectedChat?.id}`] && <TypingIndicator />}

              {currentMessage &&
                Object.keys(currentMessage)
                  .reverse() // Reverse the groups so the latest dates are at the bottom
                  .map((key) => (
                    <Fragment key={key}>
                      {currentMessage[key]
                        .reverse() // Reverse the messages within each group
                        .map((msg, index) => (
                          <div
                            key={`${msg.id}-${msg.created_at}-${index}`}
                            className={cn(
                              "chat-box max-w-72 break-words px-3 py-2 shadow-lg",
                              msg.user_id === sessionID
                                ? "self-end rounded-[16px_16px_0_16px] bg-primary/85 text-primary-foreground/75"
                                : "self-start rounded-[16px_16px_16px_0] bg-secondary"
                            )}
                          >
                            {msg.content}{" "}
                            <span
                              className={cn(
                                "mt-1 block text-xs font-light italic text-muted-foreground",
                                msg.user_id === sessionID && "text-right"
                              )}
                            >
                              {format(msg.created_at, "h:mm a")}
                            </span>
                          </div>
                        ))}
                      <div className="text-center text-xs">{key}</div>
                    </Fragment>
                  ))}
            </div>
          </div>
        </div>
        <ChatForm />
      </div>
    </>
  );
}