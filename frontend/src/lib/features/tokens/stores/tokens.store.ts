import { writable } from 'svelte/store';
import type { Token, TokenState } from '../types';
import { pb } from '$lib/config/pocketbase';
import { toast } from 'svelte-sonner';

function createTokensStore() {
    const { subscribe, set, update } = writable<TokenState>({
        tokens: [],
        isLoading: false,
        error: null
    });

    return {
        subscribe,
        
        async fetchTokens() {
            update(state => ({ ...state, isLoading: true }));
            try {
                const records = await pb.collection("tokens").getFullList({
                    sort: "-created",
                    expand: "user"
                });
                update(state => ({ ...state, tokens: records }));
            } catch (error) {
                toast.error("Failed to load tokens");
            } finally {
                update(state => ({ ...state, isLoading: false }));
            }
        },

        async toggleStatus(id: string, currentStatus: boolean) {
            try {
                await pb.collection("tokens").update(id, {
                    is_active: !currentStatus
                });
                toast.success("Token status updated");
                this.fetchTokens();
            } catch (error) {
                toast.error("Failed to update token status");
            }
        },

        async deleteToken(id: string) {
            try {
                await pb.collection("tokens").delete(id);
                toast.success("Token deleted successfully");
                this.fetchTokens();
            } catch (error) {
                toast.error("Failed to delete token");
            }
        }
    };
}

export const tokensStore = createTokensStore();