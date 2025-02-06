import {
  AlertDialog,
  // AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {useState, /*FormEvent,*/ useEffect} from "react";
// import { Loader} from "lucide-react";
// import { Button } from "@/components/ui/button";
import {useGlobalActionStore} from "@/states/global-action-store";
// import {useChatStore} from "@/states/chat-store";
import {useWebSocket} from "@/features/chats/context/websocket";
import {useWebRTC} from "@/features/chats/context/webrtc";
import {useAuthStore} from "@/states/auth-store";
import {useChatStore} from "@/states/chat-store";
import {CallStage, Chat} from "@/types/chat";
// import {Chat} from "@/types/chat";
// import {useAuthStore} from "@/states/auth-store";
// import {ChatListState} from "@/features/chats/components/chat-list";

// TODO: watch this out


export const NewCallState = "call_state"
export const NewCallModalState = "call_modal_state"

export function NewCallModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const { states, status, setState } = useGlobalActionStore();
  const { sendMessage } = useWebSocket();
  const { getMyP2PId } = useWebRTC();
  const sessionID = useAuthStore.getState().auth.user?.id;
  const { selectedChat } = useChatStore();
  const [callStage, setCallStage] = useState<CallStage>(CallStage.Calling);

  const getUserIds = (chat: Chat) => chat?.users
    .filter((user) => user.id !== sessionID).map(user => user.id) ?? [];

  useEffect(() => {
    if (states[NewCallModalState]) {
      setDialogOpen(states[NewCallModalState])
      sendMessage({
          action: "calling",
          recipient_id: getUserIds(selectedChat as Chat),
          vc_type: "video_audio",
          vc_data: getMyP2PId()
      })
    }
    if (status[NewCallState]) {
      setCallStage(status[NewCallState] as CallStage)
    }
  }, [states[NewCallModalState], status[NewCallState]]);

  const onClose = () => {
    setState(NewCallModalState, false)
    setDialogOpen(false)
  }

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Calling . .. </AlertDialogTitle>
          <AlertDialogDescription />
        </AlertDialogHeader>
        {callStage === CallStage.Calling && <>calling . . .</>}
        {callStage === CallStage.Accepted && <>accepted</>}
        {callStage === CallStage.Reject && <>rejected</>}
        <AlertDialogFooter>
          test
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}