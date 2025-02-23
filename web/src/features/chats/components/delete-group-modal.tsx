import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {useState, FormEvent, useEffect} from "react";
import { Loader} from "lucide-react";
import { Button } from "@/components/ui/button";
import {useGlobalActionStore} from "@/states/global-action-store";
import {useChatStore} from "@/states/chat-store";
import {useWebSocket} from "@/features/chats/context/websocket";
import {Chat} from "@/types/chat";
import {useAuthStore} from "@/states/auth-store";
import {ChatListState} from "@/features/chats/components/chat-list";

export const DeleteGroupModalState = "delete_group_modal_state"

export function DeleteGroupModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const [isSubmitted, setSubmitted] = useState(false);
  const sessionID = useAuthStore.getState().auth.user?.id;
  const { states, setState } = useGlobalActionStore();
  const { selectedChat, setSelectedChat } = useChatStore();
  const { sendMessage } = useWebSocket();

  useEffect(() =>{
    if (states[DeleteGroupModalState]) {
      setDialogOpen(states[DeleteGroupModalState])
    }
  }, [states]);

  const getUserIds = (chat: Chat) => chat?.users
    .filter((user) => user.id !== sessionID).map(user => user.id) ?? [];

  const onSubmit = (event: FormEvent) => {
    event.preventDefault();
    setSubmitted(true);
    if (selectedChat) {
      const userIds = getUserIds(selectedChat);
      sendMessage({
        action: "delete_group",
        conversation_id: selectedChat.id,
        recipient_id: userIds,
      });
      setTimeout(() => {
        setSelectedChat(null);
        setState(ChatListState, true)
        sendMessage({action: "chats"});
        onClose();
      }, 1000)
    }
  };

  const onClose = () => {
    setState(DeleteGroupModalState, false)
    setDialogOpen(false)
    setSubmitted(false);
  }

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete Group</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to delete {selectedChat?.name}?
            You will not be able to access this group and its data again.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={onClose}>
            Cancel
          </AlertDialogCancel>
          <Button
            className="bg-red-500 hover:bg-red-600 text-white"
            type="submit"
            onClick={onSubmit}
            disabled={isSubmitted}
          >
            <Loader
              className={`${
                isSubmitted ? "block animate-spin mr-2" : "hidden"
              }`}
            />
            Confirm
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}