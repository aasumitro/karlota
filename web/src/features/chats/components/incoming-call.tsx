import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {useState,  useEffect} from "react";
import {useGlobalActionStore} from "@/states/global-action-store";
import {useChatStore} from "@/states/chat-store";
import {useWebSocket} from "@/features/chats/context/websocket";
import {useWebRTC} from "@/features/chats/context/webrtc";
import {CallStage} from "@/types/chat";

// TODO: watch this out

export const IncomingCallModalState = "incoming_modal_state"

export function IncomingCallModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const { states, setState } = useGlobalActionStore();
  const { sendMessage } = useWebSocket();
  const { getMyP2PId, connectPeer } = useWebRTC();
  const { call } = useChatStore();

  useEffect(() => setDialogOpen(states[IncomingCallModalState]), [states[IncomingCallModalState]]);

  const onClose = () => {
    setState(IncomingCallModalState, false)
    setDialogOpen(false)
  }

  // https://glitch.com/edit/#!/peerjs-video?path=public%2Fmain.js%3A6%3A1
  // et messagesEl = document.querySelector('.messages');
  // let peerIdEl = document.querySelector('#connect-to-peer');
  // let videoEl = document.querySelector('.remote-video');
  // let logMessage = (message) => {
  //   let newMessage = document.createElement('div');
  //   newMessage.innerText = message;
  //   messagesEl.appendChild(newMessage);
  // };
  // let renderVideo = (stream) => {
  //   videoEl.srcObject = stream;
  // };

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Incoming call . .. </AlertDialogTitle>
          <AlertDialogDescription />
        </AlertDialogHeader>
        <>
          content {call?.vc_type} {call?.vc_data}
        </>
        <AlertDialogFooter>
          <button onClick={() => {
            sendMessage({
              action:"answering",
              vc_type: "video_audio",
              recipient_id: [call?.vc_caller],
              vc_data: getMyP2PId(),
              vc_action: CallStage.Accepted
            })
            connectPeer(call?.vc_data as string);
          }}> accept </button>
          <button onClick={() => {
            sendMessage({
              action:"answering",
              recipient_id: [call?.vc_caller],
              vc_type: "video_audio",
              vc_data: getMyP2PId(),
              vc_action: CallStage.Reject
            })
          }}> reject </button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}