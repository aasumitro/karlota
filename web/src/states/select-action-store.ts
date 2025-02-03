import {Profile} from "@/types/user";
import {create} from "zustand/index";

interface States {
  user: Profile | null;
}

interface Actions {
  setUser: (user: Profile | null) => void;
}

export const useSelectActionStore = create<States & Actions>((set) => ({
  user: null,
  setUser: (user: Profile | null) =>
    set((state) => {
      return { ...state, user }
    }),
}));