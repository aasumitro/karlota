export type Tokens = {
  access_token: Token;
  refresh_token: Token;
  user: Profile;
}

type Token = {
  str: string;
  dtm: string;
}

export type Profile = {
  id: number;
  display_name: string;
  email: string;
  is_online: boolean;
  last_online: number;
  role?: string;
}