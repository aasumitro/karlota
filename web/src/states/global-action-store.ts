import { create } from "zustand";

interface LoadingState {
  [key: string]: boolean;
}

interface StatusState {
  [key: string]: string;
}

interface States {
  states: LoadingState;
  status: StatusState;
}

interface Actions {
  setState: (key: string, value: boolean) => void;
  setStatus: (key: string, value: string) => void;
  resetState: () => void;
}

export const useGlobalActionStore = create<States & Actions>((set) => ({
  states: {},
  status: {},
  setState: (key: string, value: boolean) =>
    set((state) => {
      const newStates = { ...state.states };
      if (value) {
        newStates[key] = value;
      } else {
        delete newStates[key];
      }
      return { states: newStates };
    }),
  setStatus: (key: string, value: string) =>
    set((state) => {
      const newStatus = { ...state.status };
      if (value) {
        newStatus[key] = value;
      } else {
        delete newStatus[key];
      }
      return { status: newStatus };
    }),
  resetState: () => set({ states: {}, status: {} }),
}));