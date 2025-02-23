export interface Task {
    id: string;
    title: string;
    description: string;
    category: string;
    due?: string;
    created: string;
    updated: string;
    user: string;
}

export interface Category {
    value: string;
    label: string;
    color: string;
}

export interface TaskState {
    tasks: Task[];
    isLoading: boolean;
    error: string | null;
    selectedTask: Task | null;
    searchQuery: string;
}