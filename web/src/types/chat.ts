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