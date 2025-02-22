// src/lib/stores/mailMessagesStore.ts
import { writable, get } from 'svelte/store';
import { pb } from '$lib/pocketbase';
import { mailStore } from './mailStore';
import { toast } from 'svelte-sonner';

export interface SyncStatus {
    last_synced: string;
    message_count: number;
    status: string;
}

export interface MailMessage {
    id: string;
    message_id: string;
    thread_id: string;
    from: string;
    to: string;
    subject: string;
    snippet: string;
    body: string;
    is_unread: boolean;
    is_important: boolean;
    is_starred: boolean;
    is_spam: boolean;
    is_inbox: boolean;
    is_trash: boolean;
    is_draft: boolean;
    is_sent: boolean;
    internal_date: string;
    received_date: string;
    created: string;
    updated: string;
}

interface MailMessagesState {
    messages: MailMessage[];
    isLoading: boolean;
    totalItems: number;
    page: number;
    perPage: number;
    selectedMail: MailMessage | null;
}

function createMailMessagesStore() {
    const { subscribe, set, update } = writable<MailMessagesState>({
        messages: [],
        isLoading: false,
        totalItems: 0,
        page: 1,
        perPage: 25,
        selectedMail: null
    });

    let unsubscribe: (() => void) | null = null;

    return {
        subscribe,
        fetchMails: async (page = 1) => {
            const state = get({ subscribe });
            update(s => ({ ...s, isLoading: true }));

            try {
                const result = await pb.collection('mail_messages').getList(page, state.perPage, {
                    sort: '-received_date',
                    // filter: 'is_inbox=true'
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
            update(state => ({ ...state, isLoading: true }));
            try {
                await pb.send('/api/mail/sync', {
                    method: 'POST'
                });
                
                // Update sync status through mailStore
                const status = await mailStore.checkStatus();
                if (status) {
                    toast.success('Mail sync initiated');
                    await store.fetchMails();
                }
            } catch (error) {
                toast.error('Failed to sync emails');
                update(s => ({ ...s, isLoading: false }));
            }
        },

        subscribeToChanges: () => {
            if (unsubscribe) {
                unsubscribe();
            }

            unsubscribe = pb.collection('mail_messages').subscribe('*', 
                async (e) => {
                    if (e.action === 'create') {
                        update(state => ({
                            ...state,
                            messages: [e.record as MailMessage, ...state.messages],
                            totalItems: state.totalItems + 1
                        }));
                    } else if (e.action === 'update') {
                        update(state => ({
                            ...state,
                            messages: state.messages.map(msg =>
                                msg.id === e.record.id ? (e.record as MailMessage) : msg
                            ),
                            selectedMail: state.selectedMail?.id === e.record.id ? 
                                (e.record as MailMessage) : state.selectedMail
                        }));
                    } else if (e.action === 'delete') {
                        update(state => ({
                            ...state,
                            messages: state.messages.filter(msg => msg.id !== e.record.id),
                            totalItems: state.totalItems - 1,
                            selectedMail: state.selectedMail?.id === e.record.id ? 
                                null : state.selectedMail
                        }));
                    }
                }
            );
        },

        unsubscribe: () => {
            if (unsubscribe) {
                unsubscribe();
                unsubscribe = null;
            }
        }
    };
}

export const mailMessagesStore = createMailMessagesStore();