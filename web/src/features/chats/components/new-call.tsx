import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {useState, useEffect, useRef} from "react";
import {useGlobalActionStore} from "@/states/global-action-store";
import {useWebSocket} from "@/features/chats/context/websocket";
import {useWebRTC} from "@/features/chats/context/webrtc";
import {useAuthStore} from "@/states/auth-store";
import {useChatStore} from "@/states/chat-store";
import {CallStage, Chat} from "@/types/chat";
import {MediaConnection} from "peerjs";
import {cn} from "@/lib/utils";
import {Button} from "@/components/ui/button";
import {PhoneOffIcon} from "lucide-react";
import callingSound from "@/assets/sounds/calling.mp3";

// TODO FIX CAM NOT CLOSED

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

  const remoteVideoRef = useRef<HTMLVideoElement | null>(null);
  const myVideoRef = useRef<HTMLVideoElement | null>(null);
  const localStreamRef = useRef<MediaStream | null>(null);
  const callRef = useRef<MediaConnection | null>(null);

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
    console.log("Closing call...");

    // Stop local video/audio stream
    if (localStreamRef.current) {
      localStreamRef.current.getTracks().forEach(track => {
        console.log("new", track.label)
        track.stop()
      });
      localStreamRef.current = null;
      console.log("Local stream stopped.");
    }

    // Stop my video
    if (myVideoRef.current) {
      myVideoRef.current.remove()
      myVideoRef.current.srcObject = null;
      console.log("My video stopped.");
    }

    // Stop remote video
    if (remoteVideoRef.current) {
      remoteVideoRef.current.remove()
      remoteVideoRef.current.srcObject = null;
      console.log("Remote video stopped.");
    }

    callRef.current?.close()
  }

  const OnAccept = () => {
    onCall(async (call: MediaConnection) => {
      try {
        const stream = await navigator.mediaDevices
          .getUserMedia({ video: true, audio: true });
        if (myVideoRef.current) {
          myVideoRef.current.srcObject = stream;
          call.answer(stream);
          localStreamRef.current = stream
        } else {
          setTimeout(() => {
            if (myVideoRef.current) {
              myVideoRef.current.srcObject = stream;
              call.answer(stream);
              localStreamRef.current = stream
            }
          }, 250);
        }

        call.on('stream', (remoteStream) => {
          if (remoteVideoRef.current) {
            remoteVideoRef.current.srcObject = remoteStream;
          } else {
            setTimeout(() => {
              if (remoteVideoRef.current) {
                remoteVideoRef.current.srcObject = remoteStream;
              }
            }, 250);
          }
        });

        call.on("close", () => {
          console.log("close")
          onClose()
        })

        callRef.current = call
      } catch (err) {
        console.error("Failed to get local stream", err);
      }
    });

    return <div className="relative w-full h-full">
      <video
        ref={remoteVideoRef}
        autoPlay
        playsInline
        className="w-full h-full object-cover rounded-lg"
      />
      <video
        ref={myVideoRef}
        autoPlay
        playsInline
        className="absolute bottom-4 right-4 w-24 h-24 object-cover rounded-md border-2 border-white shadow-lg"
      />
      <div
        className={cn(
          "absolute -bottom-12 left-1/2 transform -translate-x-1/2 ",
          "object-cover rounded-md flex gap-2"
        )}
      >
        <Button
          variant="destructive"
          size="icon"
          className="rounded-full"
          onClick={onClose}
        >
          <PhoneOffIcon />
        </Button>
      </div>
    </div>
  }

  useEffect(() => {
    if (callStage === CallStage.Reject) {
      onClose();
    }
  }, [callStage]);

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader className="hidden">
          <AlertDialogTitle/>
          <AlertDialogDescription />
        </AlertDialogHeader>
        {callStage === CallStage.Calling && <div className="flex justify-between items-center">
          Calling . . .
          <audio src={callingSound} autoPlay loop />
          <div className="object-cover rounded-md flex gap-2">
            <Button
              variant="destructive"
              size="icon"
              className="rounded-full"
              onClick={onClose}
            >
              <PhoneOffIcon />
            </Button>
          </div>
        </div>}
        {callStage === CallStage.Accept && <OnAccept />}
      </AlertDialogContent>
    </AlertDialog>
  );
}