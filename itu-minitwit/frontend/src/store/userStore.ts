import { create } from "zustand";

interface UserStore {
	userId?: number;
	setUserId: (userId?: number) => void;
}

const useUserStore = create<UserStore>((set) => ({
	userId: undefined,
	setUserId: (userId) => set({ userId }),
}));

export default useUserStore;
