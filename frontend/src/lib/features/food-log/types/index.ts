export interface FoodLogEntry {
    id: string;
    user: string;
    name: string;
    image: string;
    tag: string;
    date: string;
    created: string;
    updated: string;
}

export interface FoodLogState {
    entries: FoodLogEntry[];
    isLoading: boolean;
    hasMore: boolean;
    page: number;
    totalItems: number;
    filter: FoodLogFilter;
}

export interface FoodLogFilter {
    date?: string;
    tag?: string;
    searchTerm?: string;
}

export interface FoodLogFormData {
    name: string;
    tag: string;
    image?: File;
    date: string;
}

export const DEFAULT_FOOD_LOG_FORM: FoodLogFormData = {
    name: '',
    tag: 'breakfast',
    image: undefined,
    date: new Date().toISOString()
};

export const FOOD_TAGS = [
    { value: 'breakfast', label: 'Breakfast' },
    { value: 'lunch', label: 'Lunch' },
    { value: 'dinner', label: 'Dinner' },
    { value: 'snack', label: 'Snack' },
    { value: 'dessert', label: 'Dessert' },
    { value: 'drink', label: 'Drink' }
] as const; 