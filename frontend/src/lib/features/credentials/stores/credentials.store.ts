import { writable } from "svelte/store";
import { pb } from "$lib/config/pocketbase";
import type { CredentialsState, CredentialsStats } from "../types";

// Initial state
const initialState: CredentialsState = {
    stats: {
        totalTokens: 0,
        totalDeveloperTokens: 0,
        totalApiKeys: 0,
    },
    isLoading: false,
};

// Create store
function createCredentialsStore() {
    const { subscribe, set, update } = writable<CredentialsState>(initialState);

    return {
        subscribe,
        loadStats: async () => {
            update(state => ({ ...state, isLoading: true }));
            try {
                // Fetch counts from all collections
                const tokens = await pb.collection("tokens").getFullList({ fields: "id" });
                const devTokens = await pb.collection("dev_tokens").getFullList({ fields: "id" });
                const apiKeys = await pb.collection("api_keys").getFullList({ fields: "id" });

                // Update stats
                update(state => ({
                    ...state,
                    stats: {
                        totalTokens: tokens.length,
                        totalDeveloperTokens: devTokens.length,
                        totalApiKeys: apiKeys.length,
                    },
                    isLoading: false,
                }));
            } catch (error) {
                console.error("Failed to load credentials stats:", error);
                update(state => ({ ...state, isLoading: false }));
            }
        },
        reset: () => set(initialState),
    };
}

export const credentialsStore = createCredentialsStore(); 