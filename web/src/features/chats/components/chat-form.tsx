import { Button } from "@/components/ui/button";
import { IconPaperclip, IconPhotoPlus, IconPlus, IconSend } from "@tabler/icons-react";
import { toast } from "@/hooks/use-toast";
import { ChangeEvent, KeyboardEvent, useState } from "react";
import { useWebSocket } from "@/features/chats/context/websocket";
import { useChatStore } from "@/states/chat-store";
import { useAuthStore } from "@/states/auth-store";
import {Chat} from "@/types/chat";

export function ChatForm() {
  const sessionID = useAuthStore.getState().auth.user?.id;
  const [message, setMessage] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const [typingTimeout, setTypingTimeout] = useState<NodeJS.Timeout | null>(null);
  const { selectedChat, setNewMessage } = useChatStore();
  const { sendMessage } = useWebSocket();

  // Helper function to get the user IDs for the chat
  const getUserIds = (chat: Chat) => chat?.users
    .filter((user) => user.id !== sessionID).map(user => user.id) ?? [];

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const inputValue = e.target.value;
    setMessage(inputValue);

    if (selectedChat) {
      // Start typing indicator when the user starts typing
      if (inputValue && !isTyping) {
        setIsTyping(true);
        const userIds = getUserIds(selectedChat);
        sendMessage({ action: "typing", conversation_id: selectedChat.id, recipient_id: userIds, typing_status: true });
      }

      // Clear any existing timeout
      if (typingTimeout) {
        clearTimeout(typingTimeout);
      }

      // Set timeout to detect typing end
      const newTimeout = setTimeout(() => {
        setIsTyping(false); // Typing ended
        const userIds = getUserIds(selectedChat);
        sendMessage({ action: "typing", conversation_id: selectedChat.id, recipient_id: userIds, typing_status: false });
      }, 1000); // Typing ended after 1 second of inactivity
      setTypingTimeout(newTimeout);
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && message.trim() && message !== "") {
      e.preventDefault();
      sendMessageToChat(message);
    }
  };

  const handleButtonSent = () => {
    if (message.trim() && message !== "") {
      sendMessageToChat(message);
    }
  };

  const sendMessageToChat = (message: string) => {
    if (selectedChat) {
      const userIds = getUserIds(selectedChat);
      sendMessage({ action: "new_text_message", conversation_id: selectedChat.id, recipient_id: userIds, text_message: message });
      setMessage('');
      setNewMessage({
        id: 1,
        conversation_id: selectedChat.id,
        user_id: sessionID as number,
        type: "text",
        content: message,
        created_at: `${new Date()}`,
        updated_at: `${new Date()}`
      })
      setIsTyping(false)
    }
  };

  const notSupportedToast = () => {
    toast({
      variant: 'destructive',
      title: 'Currently not supported!',
    });
  };

  return (
    <div className='flex w-full flex-none gap-2'>
      <div className='flex flex-1 items-center gap-2 rounded-md border border-input px-2 py-1 focus-within:outline-none focus-within:ring-1 focus-within:ring-ring lg:gap-4'>
        <div className='space-x-1'>
          <Button size='icon' type='button' variant='ghost' className='h-8 rounded-md' onClick={notSupportedToast}>
            <IconPlus size={20} className='stroke-muted-foreground' />
          </Button>
          <Button size='icon' type='button' variant='ghost' className='hidden h-8 rounded-md lg:inline-flex' onClick={notSupportedToast}>
            <IconPhotoPlus size={20} className='stroke-muted-foreground' />
          </Button>
          <Button size='icon' type='button' variant='ghost' className='hidden h-8 rounded-md lg:inline-flex' onClick={notSupportedToast}>
            <IconPaperclip size={20} className='stroke-muted-foreground' />
          </Button>
        </div>
        <label className='flex-1'>
          <span className='sr-only'>Chat Text Box</span>
          <input
            type='text'
            placeholder='Type your messages...'
            className='h-8 w-full bg-inherit focus-visible:outline-none'
            onChange={handleInputChange}
            value={message}
            onKeyDown={handleKeyDown}
          />
        </label>
        <Button variant='ghost' size='icon' className='hidden sm:inline-flex' onClick={handleButtonSent}>
          <IconSend size={20} />
        </Button>
      </div>
      <Button onClick={handleButtonSent} className='h-full sm:hidden'>
        <IconSend size={18} /> Send
      </Button>
    </div>
  );
}