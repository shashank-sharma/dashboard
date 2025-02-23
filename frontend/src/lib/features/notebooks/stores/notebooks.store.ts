import { writable } from 'svelte/store';
import type { Notebook } from '../types';
import { notebooksService } from '../services/notebooks.service';

interface NotebooksState {
    notebooks: Notebook[];
    loading: boolean;
    error: string | null;
    currentPage: number;
    itemsPerPage: number;
    totalPages: number;
    searchQuery: string;
}

const initialState: NotebooksState = {
    notebooks: [],
    loading: true,
    error: null,
    currentPage: 1,
    itemsPerPage: 6,
    totalPages: 0,
    searchQuery: ''
};

function createNotebooksStore() {
    const { subscribe, set, update } = writable<NotebooksState>(initialState);

    return {
        subscribe,
        loadNotebooks: async () => {
            update(state => ({ ...state, loading: true }));
            try {
                const notebooks = await notebooksService.listNotebooks();
                update(state => ({ 
                    ...state, 
                    notebooks,
                    totalPages: Math.ceil(notebooks.length / initialState.itemsPerPage)
                }));
            } catch (error) {
                update(state => ({ ...state, error: 'Failed to load notebooks' }));
            } finally {
                update(state => ({ ...state, loading: false }));
            }
        },
        setSearchQuery: (query: string) => {
            update(state => ({ ...state, searchQuery: query, currentPage: 1 }));
        },
        setPage: (page: number) => {
            update(state => ({ ...state, currentPage: page }));
        },
        createNotebook: async (data: CreateNotebookData) => {
            try {
                const notebook = await notebooksService.createNotebook(data);
                update(state => ({
                    ...state,
                    notebooks: [notebook, ...state.notebooks]
                }));
                return notebook;
            } catch (error) {
                throw new Error('Failed to create notebook');
            }
        },
        deleteNotebook: async (id: string) => {
            try {
                await notebooksService.deleteNotebook(id);
                update(state => ({
                    ...state,
                    notebooks: state.notebooks.filter(nb => nb.id !== id)
                }));
            } catch (error) {
                throw new Error('Failed to delete notebook');
            }
        }
    };
}

export const notebooksStore = createNotebooksStore();