import { writable, get } from 'svelte/store';
import { MailService } from '../services';
import type { MailState } from '../types';
import { toast } from 'svelte-sonner';
import { pb } from '$lib/config/pocketbase';

function createMailStore() {
    const { subscribe, set, update } = writable<MailState>({
        isAuthenticated: false,
        isLoading: false,
        isAuthenticating: false,
        syncStatus: null,
        lastChecked: null,
        syncAvailable: false // Add flag for sync availability
    });

    const STATUS_CHECK_INTERVAL = 30000;

    const store = {
        subscribe,
        
        checkStatus: async (force = false) => {
            const currentState = get({ subscribe });
            const now = Date.now();

            if (!force && currentState.lastChecked) {
                if (now - currentState.lastChecked < STATUS_CHECK_INTERVAL) {
                    return currentState.syncStatus;
                }
            }

            const shouldShowLoading = force || !currentState.lastChecked;
            
            if (shouldShowLoading) {
                update(state => ({ ...state, isLoading: true }));
            }

            try {
                const status = await MailService.checkStatus();
                update(state => ({
                    ...state, 
                    isAuthenticated: !!status, 
                    isLoading: false,
                    syncStatus: status,
                    lastChecked: now,
                    syncAvailable: true
                }));
                return status;
            } catch (error) {
                const is404 = error.status === 404;
                update(state => ({ 
                    ...state, 
                    isLoading: false,
                    lastChecked: now,
                    syncAvailable: !is404
                }));
                
                if (!is404) {
                    console.error('Error checking mail status:', error);
                }
                return null;
            }
        },

        startAuth: async () => {
            update(state => ({ ...state, isAuthenticating: true }));
            try {
                const url = await MailService.startAuth();
                if (!url) {
                    toast.error('Failed to start authentication process');
                }
                return url;
            } catch (error) {
                toast.error('Failed to start authentication process');
                return null;
            } finally {
                update(state => ({ ...state, isAuthenticating: false }));
            }
        },

        completeAuth: async (code: string) => {
            try {
                update(state => ({ ...state, isAuthenticating: true }));
                
                const result = await pb.send('/api/mail/auth/callback', {
                    method: 'POST',
                    body: { code }
                });

                // Force a fresh status check after auth
                const status = await store.checkStatus(true);
                
                if (status) {
                    // Start sync subscription if available
                    if (status && store.getState().syncAvailable) {
                        store.subscribeToChanges();
                    }
                    return true;
                }
                return false;
            } catch (error) {
                console.error('Auth completion error:', error);
                return false;
            } finally {
                update(state => ({ ...state, isAuthenticating: false }));
            }
        },

        getState: () => get({ subscribe }),

        subscribeToChanges: () => {
            const currentState = get({ subscribe });
            // Only subscribe if sync is available
            if (currentState.syncAvailable) {
                try {
                    pb.collection('mail_sync').subscribe(currentState.syncStatus.id, (e) => {
                        update(state => ({
                            ...state, 
                            syncStatus: {
                                status: e.record.sync_status,
                                last_synced: e.record.last_synced
                            },
                            lastChecked: Date.now(),
                        }));
                    });
                } catch (error) {
                    console.error('Failed to subscribe to mail changes:', error);
                }
            }
        },

        unsubscribe: () => {
            const currentState = get({ subscribe });
            // Only attempt to unsubscribe if sync was available
            if (currentState.syncAvailable) {
                try {
                    pb.collection('mail_sync').unsubscribe();
                } catch (error) {
                    console.error('Error unsubscribing from mail changes:', error);
                }
            }
        }
    };

    return store;
}

export const mailStore = createMailStore();