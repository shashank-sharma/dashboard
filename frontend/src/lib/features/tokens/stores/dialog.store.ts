import { writable } from 'svelte/store';
import type { TokenDialogState } from '../types';

function createDialogStore() {
    const { subscribe, set, update } = writable<TokenDialogState>({
        isOpen: false,
        tokenToDelete: null
    });

    return {
        subscribe,
        openDialog: () => update(state => ({ ...state, isOpen: true })),
        closeDialog: () => update(state => ({ ...state, isOpen: false })),
        setTokenToDelete: (id: string | null) => update(state => ({ ...state, tokenToDelete: id }))
    };
}

export const dialogStore = createDialogStore();