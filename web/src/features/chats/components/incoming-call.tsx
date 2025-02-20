import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {useEffect, useRef, useState} from "react";
import {useGlobalActionStore} from "@/states/global-action-store";
import {useChatStore} from "@/states/chat-store";
import {useWebSocket} from "@/features/chats/context/websocket";
import {useWebRTC} from "@/features/chats/context/webrtc";
import {CallStage} from "@/types/chat";
import {MediaConnection, Peer} from "peerjs";
import {cn} from "@/lib/utils";
import {Button} from "@/components/ui/button";
import {PhoneIcon, PhoneOffIcon} from "lucide-react";
import callingSound from "@/assets/sounds/incoming-call.mp3";

export const IncomingCallModalState = "incoming_modal_state"

export function IncomingCallModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const { states, setState } = useGlobalActionStore();
  const { sendMessage } = useWebSocket();
  const { getMyP2PId, connectPeer } = useWebRTC();
  const { call } = useChatStore();
  const [callStage, setCallStage] = useState<CallStage>(CallStage.Calling);

  // Ref for the remote video element
  const myVideoRef = useRef<HTMLVideoElement | null>(null);
  const remoteVideoRef = useRef<HTMLVideoElement | null>(null);
  const localStreamRef = useRef<MediaStream | null>(null);
  const callRef = useRef<MediaConnection | null>(null);

  useEffect(() => {
    if (states[IncomingCallModalState]) {
      setDialogOpen(states[IncomingCallModalState])
    }
  }, [states]);

  const onClose = () => {
    setState(IncomingCallModalState, false)
    setCallStage(CallStage.Calling)
    setDialogOpen(false)
    // Stop local video/audio stream
    if (localStreamRef.current) {
      localStreamRef.current.getTracks().forEach(track => {
        console.log("income", track.label)
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

  const accept = async () => {
    try {
      const stream = await navigator.mediaDevices
        .getUserMedia({ video: true, audio: true });
      if (myVideoRef.current) {
        myVideoRef.current.srcObject = stream;
        localStreamRef.current = stream
      } else {
        setTimeout(() => {
          if (myVideoRef.current) {
            myVideoRef.current.srcObject = stream;
            localStreamRef.current = stream
          }
        }, 100);
      }

      sendMessage({
        recipient_id: [call?.vc_caller],
        call: {
          audio: true,
          video: true,
          peer_id: getMyP2PId(),
          action: CallStage.Accept,
        },
        action: "answering",
      });

      // Establish peer connection
      const peerId = call?.peer_id as string
      connectPeer(peerId, (conn: Peer) => {
        const onCall = conn.call(peerId, stream);
        if (onCall) {
          onCall.on("stream", (remoteStream) => {
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
        }
        onCall.on("close", () => onClose())
        callRef.current = onCall
      });
      setCallStage(CallStage.Accept);
    } catch (err) {
      console.error("Failed to get local stream", err);
    }
  }

  const reject =() => {
    sendMessage({
      action:"answering",
      recipient_id: [call?.vc_caller],
      call: {
        "audio": true,
        "video": true,
        "peer_id": getMyP2PId(),
        "action": CallStage.Reject
      },
    })
    onClose();
  }

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader className="hidden">
          <AlertDialogTitle/>
          <AlertDialogDescription />
        </AlertDialogHeader>
        {callStage === CallStage.Accept && (
          <div className="relative w-full h-full">
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
        )}
        {callStage === CallStage.Calling && <div className="flex justify-between items-center">
          Incoming call. . .
          <audio src={callingSound} autoPlay loop />
          <div className="object-cover rounded-md flex gap-2">
            <Button
              size="icon"
              className="rounded-full bg-green-500 text-white shadow-sm hover:bg-green-500/80"
              onClick={accept}
            >
              <PhoneIcon />
            </Button>
            <Button
              variant="destructive"
              size="icon"
              className="rounded-full"
              onClick={reject}
            >
              <PhoneOffIcon />
            </Button>
          </div>
        </div>}
      </AlertDialogContent>
    </AlertDialog>
  );
}