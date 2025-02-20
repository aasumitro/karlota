import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useEffect, useRef, useState } from "react";
import { useGlobalActionStore } from "@/states/global-action-store";
import { useChatStore } from "@/states/chat-store";
import { useWebSocket } from "@/features/chats/context/websocket";
import { useWebRTC } from "@/features/chats/context/webrtc";
import { CallStage } from "@/types/chat";
import { MediaConnection, Peer } from "peerjs";
import { Button } from "@/components/ui/button";
import { PhoneIcon, PhoneOffIcon } from "lucide-react";
import callingSound from "@/assets/sounds/incoming-call.mp3";

export const IncomingCallModalState = "incoming_modal_state";

export function IncomingCallModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const [callStage, setCallStage] = useState<CallStage>(CallStage.Calling);

  const { states, setState } = useGlobalActionStore();
  const { sendMessage } = useWebSocket();
  const { getMyP2PId, connectPeer } = useWebRTC();
  const { call } = useChatStore();

  // Video and Stream references
  const myVideoRef = useRef<HTMLVideoElement | null>(null);
  const remoteVideoRef = useRef<HTMLVideoElement | null>(null);
  const localStreamRef = useRef<MediaStream | null>(null);
  const callRef = useRef<MediaConnection | null>(null);

  useEffect(() => {
    if (states[IncomingCallModalState]) {
      setDialogOpen(true);
    }
  }, [states]);

  const stopMediaStream = () => {
    if (localStreamRef.current) {
      localStreamRef.current.getTracks().forEach((track) => track.stop());
      localStreamRef.current = null;
    }
  };

  const clearVideoRefs = () => {
    if (myVideoRef.current) myVideoRef.current.srcObject = null;
    if (remoteVideoRef.current) remoteVideoRef.current.srcObject = null;
  };

  const closeCall = () => {
    setState(IncomingCallModalState, false);
    setCallStage(CallStage.Calling);
    setDialogOpen(false);
    stopMediaStream();
    clearVideoRefs();
    callRef.current?.close();
  };

  const acceptCall = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
      if (myVideoRef.current)  myVideoRef.current.srcObject = stream;
      localStreamRef.current = stream;
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
      // Establish Peer connection
      const peerId = call?.peer_id as string;
      connectPeer(peerId, (conn: Peer) => {
        const newCall = conn.call(peerId, stream);
        if (newCall) {
          newCall.on("stream", (remoteStream) => {
            if (remoteVideoRef.current) remoteVideoRef.current.srcObject = remoteStream;
          });
          newCall.on("close", closeCall);
          callRef.current = newCall;
        }
      });
      setCallStage(CallStage.Accept);
    } catch (error) {
      console.error("Error accessing media devices:", error);
    }
  };

  const rejectCall = () => {
    sendMessage({
      action: "answering",
      recipient_id: [call?.vc_caller],
      call: {
        audio: true,
        video: true,
        peer_id: getMyP2PId(),
        action: CallStage.Reject,
      },
    });
    closeCall();
  };

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={closeCall}>
      <AlertDialogContent>
        <AlertDialogHeader className="hidden">
          <AlertDialogTitle />
          <AlertDialogDescription />
        </AlertDialogHeader>

        {callStage === CallStage.Accept ? (
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
            <div className="absolute -bottom-12 left-1/2 transform -translate-x-1/2 flex gap-2">
              <Button variant="destructive" size="icon" className="rounded-full" onClick={closeCall}>
                <PhoneOffIcon />
              </Button>
            </div>
          </div>
        ) : (
          <div className="flex justify-between items-center">
            <span>Incoming call...</span>
            <audio src={callingSound} autoPlay loop />
            <div className="flex gap-2">
              <Button size="icon" className="rounded-full bg-green-500 text-white shadow-sm hover:bg-green-500/80" onClick={acceptCall}>
                <PhoneIcon />
              </Button>
              <Button variant="destructive" size="icon" className="rounded-full" onClick={rejectCall}>
                <PhoneOffIcon />
              </Button>
            </div>
          </div>
        )}
      </AlertDialogContent>
    </AlertDialog>
  );
}

