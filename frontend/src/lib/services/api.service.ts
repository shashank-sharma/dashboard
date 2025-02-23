import { pb } from '$lib/config/pocketbase';
import type { AuthState } from '$lib/types';

export class ApiService {
    static async refreshToken() {
        try {
            await pb.collection('users').authRefresh();
        } catch (err) {
            pb.authStore.clear();
            window.location.href = '/auth/login';
        }
    }

    static getAuthState(): AuthState {
        return {
            isAuthenticated: pb.authStore.isValid,
            user: pb.authStore.model
        };
    }
}