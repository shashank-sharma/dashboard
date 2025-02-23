import { writable } from 'svelte/store';
import type { DashboardState } from '../types';

interface DashboardStats {
    totalTasks: number;
    completedTasks: number;
    upcomingEvents: number;
    activeProjects: number;
}

const initialState: DashboardState = {
    activeSection: 'dashboard',
    stats: {
        totalTasks: 0,
        completedTasks: 0,
        upcomingEvents: 0,
        activeProjects: 0
    },
    isLoading: false
};

// Create a single store instance
const store = writable<DashboardState>(initialState);

// Export actions
export const dashboardStore = {
    subscribe: store.subscribe,
    setActiveSection: (section: string) => {
        store.update(state => ({ ...state, activeSection: section }));
    },
    setStats: (stats: DashboardStats) => {
        store.update(state => ({ ...state, stats }));
    },
    setLoading: (isLoading: boolean) => {
        store.update(state => ({ ...state, isLoading }));
    }
};