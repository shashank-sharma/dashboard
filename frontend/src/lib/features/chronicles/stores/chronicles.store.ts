import { writable } from 'svelte/store';
import { pb } from '$lib/config/pocketbase';
import type { ChronicleState, JournalEntry } from '../types';
import { toast } from 'svelte-sonner';

const initialState: ChronicleState = {
    currentEntry: null,
    isLoading: false,
    error: null
};

function createChroniclesStore() {
    const { subscribe, set, update } = writable<ChronicleState>(initialState);

    return {
        subscribe,
        async loadEntry(date: Date) {
            update(state => ({ ...state, isLoading: true }));
            try {
                const record = await pb.collection('journal_entries').getFirstListItem(
                    `date = "${date.toISOString().split('T')[0]}" && user = "${pb.authStore.model.id}"`
                );
                update(state => ({ 
                    ...state, 
                    currentEntry: record,
                    isLoading: false 
                }));
            } catch (error) {
                update(state => ({ 
                    ...state, 
                    currentEntry: null,
                    isLoading: false 
                }));
            }
        },

        async saveEntry(entry: JournalEntry) {
            update(state => ({ ...state, isLoading: true }));
            try {
                const data = {
                    ...entry,
                    user: pb.authStore.model.id
                };

                let record;
                try {
                    const existingRecord = await pb.collection('journal_entries')
                        .getFirstListItem(`date = "${entry.date}" && user = "${pb.authStore.model.id}"`);
                    record = await pb.collection('journal_entries').update(existingRecord.id, data);
                    toast.success('Updated successfully');
                } catch {
                    record = await pb.collection('journal_entries').create(data);
                    toast.success('Saved successfully');
                }

                update(state => ({ 
                    ...state, 
                    currentEntry: record,
                    isLoading: false 
                }));
            } catch (error) {
                toast.error('Failed to save');
                update(state => ({ 
                    ...state, 
                    isLoading: false,
                    error: 'Failed to save entry' 
                }));
            }
        }
    };
}

export const chroniclesStore = createChroniclesStore();