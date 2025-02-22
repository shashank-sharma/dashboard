// src/lib/stores/mailStore.ts
import { writable } from 'svelte/store';
import { pb } from '$lib/pocketbase';
import { toast } from 'svelte-sonner';
import type { SyncStatus } from './mailMessagesStore';

interface MailState {
    isAuthenticated: boolean;
    isLoading: boolean;
    isAuthenticating: boolean;
    syncStatus: SyncStatus | null;
}

function createMailStore() {
    const { subscribe, set, update } = writable<MailState>({
        isAuthenticated: false,
        isLoading: true,
        isAuthenticating: false,
        syncStatus: null
    });

    return {
        subscribe,
        checkStatus: async () => {
            update(state => ({ ...state, isLoading: true }));
            try {
                const status = await pb.send<SyncStatus>('/api/mail/sync/status', {
                    method: 'GET'
                });
                update(state => ({ 
                    ...state, 
                    isAuthenticated: true, 
                    isLoading: false,
                    syncStatus: status
                }));
                return status;
            } catch (error) {
                update(state => ({ 
                    ...state, 
                    isAuthenticated: false, 
                    isLoading: false,
                    syncStatus: null
                }));
                return null;
            }
        },
        startAuth: async () => {
            update(state => ({ ...state, isAuthenticating: true }));
            try {
                const response = await pb.send<{ url: string }>('/auth/mail/redirect', {
                    method: 'GET'
                });
                return response.url;
            } catch (error) {
                toast.error('Failed to start authentication process');
                update(state => ({ ...state, isAuthenticating: false }));
                return null;
            }
        },
        completeAuth: async (code: string) => {
            try {
                await pb.send('/auth/mail/callback', {
                    method: 'POST',
                    body: JSON.stringify({
                        code,
                        provider: 'gmail'
                    })
                });
                const status = await pb.send<SyncStatus>('/api/mail/sync/status', {
                    method: 'GET'
                });
                update(state => ({ 
                    ...state, 
                    isAuthenticated: true, 
                    isAuthenticating: false,
                    syncStatus: status
                }));
                return true;
            } catch (error) {
                update(state => ({ ...state, isAuthenticating: false }));
                return false;
            }
        },
        reset: () => {
            set({
                isAuthenticated: false,
                isLoading: false,
                isAuthenticating: false,
                syncStatus: null
            });
        }
    };
}

export const mailStore = createMailStore();