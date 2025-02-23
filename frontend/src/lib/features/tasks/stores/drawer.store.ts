import { writable } from 'svelte/store';
import type { Task } from '../types';

type DrawerState = {
    isOpen: boolean;
    task: Task | null;
};

function createDrawerStore() {
    const { subscribe, set, update } = writable<DrawerState>({
        isOpen: false,
        task: null
    });

    return {
        subscribe,
        open: (task: Task) => update(state => ({ isOpen: true, task: { ...task } })),
        close: () => set({ isOpen: false, task: null }),
        updateTask: (updates: Partial<Task>) => 
            update(state => ({
                ...state,
                task: state.task ? { ...state.task, ...updates } : null
            }))
    };
}

export const drawerStore = createDrawerStore();