import {Profile} from "@/types/user";

export type Chat = {
  id: number;
  type: string;
  name: string;
  users: Profile[];
  messages?: Message[];
  created_at: string;
  updated_at: string;
}

export type Message = {
  id: number;
  conversation_id: number;
  user_id: number;
  type: string;
  content: string;
  created_at: string;
  updated_at: string;
}

// TODO: FIX
export type Call = {
  vc_type: string;
  vc_data: string
  vc_action: string;
  vc_caller: number;
}

export enum CallStage {
  Calling = "calling",
  Accepted = "accepted",
  Reject = "rejected",
  Terminate = "terminated"
}