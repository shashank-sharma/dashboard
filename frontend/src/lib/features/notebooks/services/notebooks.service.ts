import { pb } from "$lib/config/pocketbase";
import type { Notebook, CreateNotebookData } from "../types";

export const notebooksService = {
    async listNotebooks() {
        try {
            const records = await pb.collection('notebooks').getFullList({
                sort: '-created',
                expand: 'user'
            });
            return records as Notebook[];
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
            return record as Notebook;
        } catch (error) {
            console.error('Error fetching notebook:', error);
            throw error;
        }
    },

    async createNotebook(data: CreateNotebookData) {
        try {
            const record = await pb.collection('notebooks').create({
                ...data,
                user: pb.authStore.model?.id,
                cells: data.cells || []
            });
            return record as Notebook;
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
            return record as Notebook;
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