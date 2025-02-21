import {createContext, ReactNode, useCallback, useContext, useEffect, useRef, useState} from "react";
import {MediaConnection, Peer} from "peerjs";

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
    peerRef.current = new Peer(p2pid);
    peerRef.current.on("error", (err) => console.error("PeerJS Error:", err));
  }, [p2pid]);

  const getMyP2PId = useCallback(() => p2pid, [p2pid]);

  const onCall = useCallback((
    fn?: (call: MediaConnection) => void,
  ) => {
    if (!fn) return;
    const conn =  peerRef.current
    conn?.on('call', (call) => fn(call));
  }, []);

  const connectPeer = useCallback((
    peerId: string, fn?: (conn: Peer) => void,
  ) => {
    if (!peerRef.current || !fn) return;
    peerRef.current.connect(peerId)
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

// eslint-disable-next-line react-refresh/only-export-components
export const useWebRTC = (): WebRTCContext => {
  const context = useContext(WebRTCContext);
  if (!context) {
    throw new Error('useWebRTC must be used within a WebRTCProvider');
  }
  return context;
};
