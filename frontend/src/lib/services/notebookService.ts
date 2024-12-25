import { pb } from '$lib/pocketbase';

export interface Cell {
    id: string;
    content: string;
    output: string;
    type: 'code' | 'markdown';
    language: string;
}

export interface Notebook {
    id: string;
    name: string;
    version: string;
    cells: Cell[];
    created: string;
    updated: string;
    user: string;
}

export const notebookService = {
    async listNotebooks() {
        try {
            const records = await pb.collection('notebooks').getFullList({
                sort: '-created'
            });
            return records;
        } catch (error) {
            console.error('Error fetching notebooks:', error);
            throw error;
        }
    },

    async getNotebook(id: string) {
        try {
            const record = await pb.collection('notebooks').getOne(id, {
                expand: 'user'
            });
            return record;
        } catch (error) {
            console.error('Error fetching notebook:', error);
            throw error;
        }
    },

    async createNotebook(data: Partial<Notebook>) {
        try {
            const record = await pb.collection('notebooks').create({
                ...data,
                user: pb.authStore.model.id,
                cells: data.cells || []
            });
            return record;
        } catch (error) {
            console.error('Error creating notebook:', error);
            throw error;
        }
    },

    async updateNotebook(id: string, data: Partial<Notebook>) {
        try {
            const record = await pb.collection('notebooks').update(id, {
                ...data,
                updated: new Date().toISOString()
            });
            return record;
        } catch (error) {
            console.error('Error updating notebook:', error);
            throw error;
        }
    },

    async deleteNotebook(id: string) {
        try {
            await pb.collection('notebooks').delete(id);
        } catch (error) {
            console.error('Error deleting notebook:', error);
            throw error;
        }
    }
};