import { writable, get } from 'svelte/store';
import type { MailMessagesState, MailMessage } from '../types';
import { toast } from 'svelte-sonner';
import { pb } from '$lib/config/pocketbase';
import { mailStore } from './mail.store';

export function createMailMessagesStore() {
    const { subscribe, set, update } = writable<MailMessagesState>({
        messages: [],
        isLoading: false,
        totalItems: 0,
        page: 1,
        perPage: 25,
        selectedMail: null
    });

    const store = {
        subscribe,
        fetchMails: async (page = 1) => {
            update(s => ({ ...s, isLoading: true }));

            try {
                const result = await pb.collection('mail_messages').getList(page, store.getState().perPage, {
                    sort: '-received_date',
                });

                update(s => ({
                    ...s,
                    messages: result.items as MailMessage[],
                    totalItems: result.totalItems,
                    page,
                    isLoading: false
                }));
            } catch (error) {
                toast.error('Failed to fetch emails');
                update(s => ({ ...s, isLoading: false }));
            }
        },

        markAsRead: async (messageId: string) => {
            try {
                await pb.collection('mail_messages').update(messageId, {
                    is_unread: false
                });

                update(state => ({
                    ...state,
                    messages: state.messages.map(msg => 
                        msg.id === messageId ? { ...msg, is_unread: false } : msg
                    )
                }));
            } catch (error) {
                toast.error('Failed to mark email as read');
            }
        },

        selectMail: (mail: MailMessage | null) => {
            update(state => ({ ...state, selectedMail: mail }));
            if (mail?.is_unread) {
                store.markAsRead(mail.id);
            }
        },

        refreshMails: async () => {
            try {
                await pb.send('/api/mail/sync', { method: 'POST' });
                toast.success('Mail sync initiated');
                await store.fetchMails();
                toast.success('Mail sync completed');
            } catch (error) {
                toast.error('Failed to sync emails');
            }
        },

        subscribeToChanges: () => {
            // Check mail store sync availability before subscribing
            const mailState = get(mailStore);
            if (mailState.syncAvailable) {
                try {
                    pb.collection('mail_messages').subscribe('*', () => {
                        store.fetchMails();
                    });
                } catch (error) {
                    console.error('Failed to subscribe to mail messages:', error);
                }
            }
        },

        unsubscribe: () => {
            const mailState = get(mailStore);
            if (mailState.syncAvailable) {
                try {
                    pb.collection('mail_messages').unsubscribe();
                } catch (error) {
                    console.error('Error unsubscribing from mail messages:', error);
                }
            }
        },

        getState: () => get({ subscribe })
    };

    return store;
}

export const mailMessagesStore = createMailMessagesStore();