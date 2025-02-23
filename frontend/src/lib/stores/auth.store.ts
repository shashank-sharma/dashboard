import { writable } from 'svelte/store';
import type { AuthState } from '$lib/types';
import { ApiService } from '$lib/services/api.service';

function createAuthStore() {
    const { subscribe, set, update } = writable<AuthState>(ApiService.getAuthState());

    return {
        subscribe,
        login: (user) => update(state => ({ ...state, isAuthenticated: true, user })),
        logout: () => set({ isAuthenticated: false, user: null })
    };
}

export const auth = createAuthStore();