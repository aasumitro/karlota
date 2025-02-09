import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {useState, /*FormEvent,*/ useEffect, useRef} from "react";
// import { Loader} from "lucide-react";
// import { Button } from "@/components/ui/button";
import {useGlobalActionStore} from "@/states/global-action-store";
// import {useChatStore} from "@/states/chat-store";
import {useWebSocket} from "@/features/chats/context/websocket";
import {useWebRTC} from "@/features/chats/context/webrtc";
import {useAuthStore} from "@/states/auth-store";
import {useChatStore} from "@/states/chat-store";
import {CallStage, Chat} from "@/types/chat";
import {MediaConnection} from "peerjs";
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
  const { getMyP2PId, onCall } = useWebRTC();
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
          call: {
            "audio": true,
            "video": true,
            "peer_id": getMyP2PId(),
          },
      })
    }
    if (status[NewCallState]) {
      setCallStage(status[NewCallState] as CallStage)
    }
  }, [states, status]);

  const onClose = () => {
    setState(NewCallModalState, false)
    setDialogOpen(false)
  }

  const OnAccept = () => {
    const remoteVideoRef = useRef<HTMLVideoElement | null>(null);
    const myVideoRef = useRef<HTMLVideoElement | null>(null);

    onCall(async (call: MediaConnection) => {
      try {
        const stream = await navigator.mediaDevices
          .getUserMedia({ video: true, audio: true });
        if (myVideoRef.current) {
          myVideoRef.current.srcObject = stream;
          call.answer(stream);
        } else {
          setTimeout(() => {
            if (myVideoRef.current) {
              myVideoRef.current.srcObject = stream;
              call.answer(stream);
            }
          }, 100);
        }

        call.on('stream', (remoteStream) => {
          if (remoteVideoRef.current) {
            remoteVideoRef.current.srcObject = remoteStream;
          } else {
            setTimeout(() => {
              if (remoteVideoRef.current) {
                remoteVideoRef.current.srcObject = remoteStream;
              }
            }, 100);
          }
        });
      } catch (err) {
        console.error("Failed to get local stream", err);
      }
    });

    return <div>
      <video ref={remoteVideoRef} autoPlay playsInline />
      <video ref={myVideoRef} autoPlay playsInline />
    </div>
  }

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Calling . . . </AlertDialogTitle>
          <AlertDialogDescription />
        </AlertDialogHeader>
        {callStage === CallStage.Calling && <>calling . . .</>}
        {callStage === CallStage.Accept && <OnAccept />}
        {callStage === CallStage.Reject && <>rejected</>}
      </AlertDialogContent>
    </AlertDialog>
  );
}