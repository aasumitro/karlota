import {createContext, ReactNode, useCallback, useContext, useEffect, useRef, useState} from "react";
import {MediaConnection, Peer} from "peerjs";

// TODO: need some improvement

interface WebRTCContext {
  getMyP2PId: () => string;
  onCall: (fn?: (call: MediaConnection) => void) => void;
  connectPeer: (peerId: string, fn?: (conn: Peer) => void) => void;
}

const WebRTCContext = createContext<WebRTCContext | null>(null);

export const WebRTCProvider = ({ children }: { children: ReactNode }) => {
  const peerRef = useRef<Peer | null>(null);
  const [p2pid, setP2Pid] = useState<string>("")

  // Function to generate and retrieve a Peer ID
  const genP2PId = (length = 8) => {
    const storedId = localStorage.getItem("p2p_id");
    if (storedId) return storedId;
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let id = "";
    for (let i = 0; i < length; i++) {
      id += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    localStorage.setItem("p2p_id", id);
    return id;
  };

  // Function to initialize PeerJS connection
  const createPeer = useCallback(() => {
    if (!p2pid || peerRef.current) return;
    console.log("Initializing Peer with ID:", p2pid);
    peerRef.current = new Peer(p2pid);
    peerRef.current.on("open", (id) => console.log("Peer connection established with ID:", id));
    peerRef.current.on("connection", (conn) => {
      conn.on('data', (data) => console.log(`received: ${data}`));
      conn.on('open', () => conn.send('hello!'));
    });
    peerRef.current.on('call', (call) => {
      navigator.mediaDevices.getUserMedia({video: true, audio: true})
        .then((stream) => {
          // todo
          console.log(call, stream)
          // call.answer(stream); // Answer the call with an A/V stream.
          // call.on('stream', renderVideo);
        })
        .catch((err) => {
          console.error('Failed to get local stream', err);
        });
    });
    peerRef.current.on("error", (err) => console.error("PeerJS Error:", err));
  }, [p2pid]);

  const getMyP2PId = useCallback(() => p2pid, [p2pid]);

  // // todo
  //  console.log(call, stream)
  //  // call.answer(stream); // Answer the call with an A/V stream.
  //  // call.on('stream', renderVideo);
  // // navigator.mediaDevices.getUserMedia({video: true, audio: true})
  //       //   .then((stream) => {
  //       //     fn(call)
  //       //   })
  //       //   .catch((err) => {
  //       //     console.error('Failed to get local stream', err);
  //       //   });
  const onCall = useCallback((
    fn?: (call: MediaConnection) => void,
  ) => {
    if (!fn) return;
    peerRef.current?.on('call', (call) => fn(call));
  }, []);

  // navigator.mediaDevices.getUserMedia({video: true, audio: true})
  //       .then((stream) => {
  //         // todo
  //         console.log(stream)
  //         // const call = peerRef.current?.call(peerId, stream);
  //         // call?.on('stream', renderVideo);
  //       })
  //       .catch((err) => {
  //         console.log('Failed to get local stream', err);
  //       });
  const connectPeer = useCallback((
    peerId: string, fn?: (conn: Peer) => void,
  ) => {
    if (!peerRef.current || !fn) return;
    const conn =  peerRef.current.connect(peerId)
    conn?.on('data', (data) => console.log(`received: ${data}`));
    conn?.on('open', () => conn.send('hi!'));
    fn(peerRef?.current)
  }, [])

  useEffect(() => {
    // Generate Peer ID and update state
    const id = genP2PId(16);
    setP2Pid(id);
    return () => {
      if (peerRef.current) {
        peerRef.current.destroy();
        peerRef.current = null;
      }
    };
  }, [createPeer]);

  useEffect(() => {
    if (p2pid) createPeer();
  }, [p2pid, createPeer]);

  return (
    <WebRTCContext.Provider value={{ getMyP2PId, onCall, connectPeer }}>
      {children}
    </WebRTCContext.Provider>
  );
}

export const useWebRTC = (): WebRTCContext => {
  const context = useContext(WebRTCContext);
  if (!context) {
    throw new Error('useWebRTC must be used within a WebRTCProvider');
  }
  return context;
};
