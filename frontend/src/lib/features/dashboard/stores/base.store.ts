import { writable } from 'svelte/store';
import { pb } from '$lib/config/pocketbase';
import type { BaseItem } from '../types';

export function createBaseStore<T extends BaseItem>(collectionName: string) {
    const { subscribe, set, update } = writable<{
        items: T[];
        loading: boolean;
        error: string | null;
    }>({
        items: [],
        loading: false,
        error: null
    });

    return {
        subscribe,
        async fetchItems() {
            update(s => ({ ...s, loading: true, error: null }));
            try {
                const records = await pb.collection(collectionName).getList(1, 50, {
                    sort: '-created',
                    expand: 'user'
                });
                update(s => ({ ...s, items: records.items as T[] }));
            } catch (error) {
                update(s => ({ ...s, error: (error as Error).message }));
            } finally {
                update(s => ({ ...s, loading: false }));
            }
        },
        async createItem(data: Partial<T>) {
            try {
                const record = await pb.collection(collectionName).create(data);
                update(s => ({ ...s, items: [record as T, ...s.items] }));
            } catch (error) {
                throw error;
            }
        },
        async updateItem(id: string, data: Partial<T>) {
            try {
                const record = await pb.collection(collectionName).update(id, data);
                update(s => ({
                    ...s,
                    items: s.items.map(item => item.id === id ? (record as T) : item)
                }));
            } catch (error) {
                throw error;
            }
        },
        async deleteItem(id: string) {
            try {
                await pb.collection(collectionName).delete(id);
                update(s => ({
                    ...s,
                    items: s.items.filter(item => item.id !== id)
                }));
            } catch (error) {
                throw error;
            }
        }
    };
}