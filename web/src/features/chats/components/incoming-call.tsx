import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {useState, useEffect, useRef} from "react";
import {useGlobalActionStore} from "@/states/global-action-store";
import {useChatStore} from "@/states/chat-store";
import {useWebSocket} from "@/features/chats/context/websocket";
import {useWebRTC} from "@/features/chats/context/webrtc";
import {CallStage} from "@/types/chat";
import {NewCallState} from "@/features/chats/components/new-call";
import {Peer} from "peerjs";

export const IncomingCallModalState = "incoming_modal_state"

export function IncomingCallModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const { states, setState,setStatus } = useGlobalActionStore();
  const { sendMessage } = useWebSocket();
  const { getMyP2PId, connectPeer } = useWebRTC();
  const { call } = useChatStore();
  const [callStage, setCallStage] = useState<CallStage>(CallStage.Calling);

  // Ref for the remote video element
  const myVideoRef = useRef<HTMLVideoElement | null>(null);
  const remoteVideoRef = useRef<HTMLVideoElement | null>(null);

  useEffect(() => {
    if (states[IncomingCallModalState]) {
      setDialogOpen(states[IncomingCallModalState])
    }
  }, [states]);

  const onClose = () => {
    setState(IncomingCallModalState, false)
    setDialogOpen(false)
    setStatus(NewCallState, "")
  }

  const accept = async () => {
    try {
      const stream = await navigator.mediaDevices
        .getUserMedia({ video: true, audio: true });
      if (myVideoRef.current) {
        myVideoRef.current.srcObject = stream;
      } else {
        setTimeout(() => {
          if (myVideoRef.current) {
            myVideoRef.current.srcObject = stream;
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
        <AlertDialogHeader>
          <AlertDialogTitle>Incoming call . .. </AlertDialogTitle>
          <AlertDialogDescription />
        </AlertDialogHeader>
        {callStage === CallStage.Accept && (
          <div>
            <video ref={myVideoRef} autoPlay playsInline />
            <video ref={remoteVideoRef} autoPlay playsInline />
          </div>
        )}
        {callStage === CallStage.Calling && <AlertDialogFooter>
          <button onClick={accept}> accept </button>
          <button onClick={reject}> reject </button>
        </AlertDialogFooter>}
      </AlertDialogContent>
    </AlertDialog>
  );
}