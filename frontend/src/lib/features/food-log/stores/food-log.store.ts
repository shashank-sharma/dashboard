import { writable } from 'svelte/store';
import { pb } from '$lib/config/pocketbase';
import { toast } from 'svelte-sonner';
import type { FoodLogEntry, FoodLogState, FoodLogFilter, FoodLogFormData } from '../types';
import { format } from 'date-fns';

// Initial state
const initialState: FoodLogState = {
    entries: [],
    isLoading: false,
    hasMore: true,
    page: 1,
    totalItems: 0,
    filter: {}
};

// Create store
function createFoodLogStore() {
    const { subscribe, set, update } = writable<FoodLogState>(initialState);

    // Helper function to check authentication
    const ensureAuthenticated = async (): Promise<boolean> => {
        if (!pb.authStore.isValid) {
            toast.error('Please log in to access your food log');
            return false;
        }
        return true;
    };

    return {
        subscribe,
        
        // Load food log entries with pagination and filters
        loadEntries: async (reset = false) => {
            if (!await ensureAuthenticated()) return;
            
            update(state => ({ ...state, isLoading: true }));
            
            try {
                const currentState = { ...initialState };
                subscribe(s => {
                    currentState.page = s.page;
                    currentState.filter = s.filter;
                    currentState.entries = s.entries;
                })();

                // Reset state if requested
                if (reset) {
                    update(state => ({
                        ...state,
                        entries: [],
                        page: 1,
                        hasMore: true
                    }));
                    currentState.page = 1;
                    currentState.entries = [];
                }

                // Build filter conditions
                const filter = [];
                if (currentState.filter.date) {
                    const date = new Date(currentState.filter.date);
                    const startOfDay = new Date(date);
                    startOfDay.setHours(0, 0, 0, 0);
                    
                    const endOfDay = new Date(date);
                    endOfDay.setHours(23, 59, 59, 999);
                    
                    filter.push(`date >= '${startOfDay.toISOString()}' && date <= '${endOfDay.toISOString()}'`);
                }
                
                if (currentState.filter.tag) {
                    filter.push(`tag = '${currentState.filter.tag}'`);
                }
                
                if (currentState.filter.searchTerm) {
                    filter.push(`name ~ '${currentState.filter.searchTerm}'`);
                }

                const filterString = filter.length > 0 ? filter.join(' && ') : '';

                // Fetch data from PocketBase
                const perPage = 10;
                const resultList = await pb.collection('food_log').getList(currentState.page, perPage, {
                    sort: '-date',
                    filter: filterString,
                    expand: 'user'
                });

                update(state => ({
                    ...state,
                    entries: reset ? resultList.items as unknown as FoodLogEntry[] : [...state.entries, ...resultList.items as unknown as FoodLogEntry[]],
                    hasMore: resultList.items.length === perPage,
                    page: state.page + 1,
                    totalItems: resultList.totalItems,
                    isLoading: false
                }));
            } catch (error) {
                console.error('Failed to load food log entries:', error);
                update(state => ({ ...state, isLoading: false, hasMore: false }));
            }
        },

        // Add a new food log entry
        addEntry: async (data: FoodLogFormData) => {
            if (!await ensureAuthenticated()) return { success: false, error: 'Not authenticated' };
            
            update(state => ({ ...state, isLoading: true }));
            
            try {
                const formData = new FormData();
                formData.append('name', data.name);
                formData.append('tag', data.tag);
                formData.append('date', data.date);
                formData.append('user', pb.authStore.model?.id);
                
                if (data.image) {
                    formData.append('image', data.image);
                }

                const record = await pb.collection('food_log').create(formData);
                
                // Add the new entry to the beginning of the list and update state
                update(state => ({
                    ...state,
                    entries: [record as unknown as FoodLogEntry, ...state.entries],
                    isLoading: false,
                    totalItems: state.totalItems + 1
                }));

                return { success: true, record };
            } catch (error) {
                console.error('Failed to add food log entry:', error);
                update(state => ({ ...state, isLoading: false }));
                return { success: false, error };
            }
        },

        // Delete a food log entry
        deleteEntry: async (id: string) => {
            if (!await ensureAuthenticated()) return { success: false, error: 'Not authenticated' };
            
            update(state => ({ ...state, isLoading: true }));
            
            try {
                await pb.collection('food_log').delete(id);
                
                // Remove the entry from the store
                update(state => ({
                    ...state,
                    entries: state.entries.filter(entry => entry.id !== id),
                    isLoading: false,
                    totalItems: state.totalItems - 1
                }));

                return { success: true };
            } catch (error) {
                console.error('Failed to delete food log entry:', error);
                update(state => ({ ...state, isLoading: false }));
                return { success: false, error };
            }
        },

        // Update filters and reload entries
        setFilter: (filter: Partial<FoodLogFilter>) => {
            update(state => ({
                ...state,
                filter: { ...state.filter, ...filter }
            }));
            
            // Reload entries with the new filter
            foodLogStore.loadEntries(true);
        },

        // Reset the store to initial state
        reset: () => set(initialState)
    };
}

export const foodLogStore = createFoodLogStore(); 