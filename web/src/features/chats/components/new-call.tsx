import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useState, useEffect, useRef } from "react";
import { useGlobalActionStore } from "@/states/global-action-store";
import { useWebSocket } from "@/features/chats/context/websocket";
import { useWebRTC } from "@/features/chats/context/webrtc";
import { useAuthStore } from "@/states/auth-store";
import { useChatStore } from "@/states/chat-store";
import { CallStage, Chat } from "@/types/chat";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { PhoneOffIcon } from "lucide-react";
import callingSound from "@/assets/sounds/calling.mp3";
import {MediaConnection} from "peerjs";

export const NewCallState = "call_state";
export const NewCallModalState = "call_modal_state";

export function NewCallModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const { states, status, setState } = useGlobalActionStore();
  const { sendMessage } = useWebSocket();
  const { getMyP2PId, onCall } = useWebRTC();
  const sessionID = useAuthStore.getState().auth.user?.id;
  const { selectedChat } = useChatStore();
  const [callStage, setCallStage] = useState<CallStage>(CallStage.Calling);

  const remoteVideoRef = useRef<HTMLVideoElement | null>(null);
  const myVideoRef = useRef<HTMLVideoElement | null>(null);
  const localStreamRef = useRef<MediaStream | null>(null);
  const callRef = useRef<MediaConnection | null>(null);

  const getUserIds = (chat: Chat) => chat?.users
    .filter((user) => user.id !== sessionID).map((user) => user.id) ?? [];

  useEffect(() => {
    if (states[NewCallModalState]) {
      setDialogOpen(true);
      sendMessage({
        action: "calling",
        recipient_id: getUserIds(selectedChat as Chat),
        call: { audio: true, video: true, peer_id: getMyP2PId() },
      });
    }
    if (status[NewCallState]) {
      setCallStage(status[NewCallState] as CallStage);
    }
  }, [states, status]);

  const stopMediaStream = () => {
    if (localStreamRef.current) {
      localStreamRef.current.getTracks().forEach((track) => track.stop());
      localStreamRef.current = null;
    }
  };

  const clearVideoRefs = () => {
    [myVideoRef, remoteVideoRef].forEach((ref) => {
      if (ref.current) {
        ref.current.srcObject = null;
      }
    });
  };

  const onClose = () => {
    setState(NewCallModalState, false);
    setDialogOpen(false);
    stopMediaStream();
    clearVideoRefs();
    callRef.current?.close();
  };

  const onAcceptCall = () => {
    onCall(async (call) => {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
        if (myVideoRef.current) myVideoRef.current.srcObject = stream;
        call.answer(stream);
        localStreamRef.current = stream;
        call.on("stream", (remoteStream) => {
          if (remoteVideoRef.current) remoteVideoRef.current.srcObject = remoteStream;
        });
        call.on("close", onClose);
        callRef.current = call;
      } catch (err) {
        console.error("Failed to get local stream", err);
      }
    });
  };

  useEffect(() => {
    if (callStage === CallStage.Accept) onAcceptCall();
    if (callStage === CallStage.Reject) onClose();
  }, [callStage]);

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader className="hidden">
          <AlertDialogTitle />
          <AlertDialogDescription />
        </AlertDialogHeader>
        {callStage === CallStage.Calling && (
          <div className="flex justify-between items-center">
            Calling . . .
            <audio src={callingSound} autoPlay loop />
            <Button variant="destructive" size="icon" className="rounded-full" onClick={onClose}>
              <PhoneOffIcon />
            </Button>
          </div>
        )}
        {callStage === CallStage.Accept && (
          <div className="relative w-full h-full">
            <video ref={remoteVideoRef} autoPlay playsInline className="w-full h-full object-cover rounded-lg" />
            <video ref={myVideoRef} autoPlay playsInline className="absolute bottom-4 right-4 w-24 h-24 object-cover rounded-md border-2 border-white shadow-lg" />
            <div className={cn("absolute -bottom-12 left-1/2 transform -translate-x-1/2 flex gap-2")}>
              <Button variant="destructive" size="icon" className="rounded-full" onClick={onClose}>
                <PhoneOffIcon />
              </Button>
            </div>
          </div>
        )}
      </AlertDialogContent>
    </AlertDialog>
  );
}
