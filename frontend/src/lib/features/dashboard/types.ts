import type { ComponentType } from "svelte";

export interface DashboardSection {
    id: string;
    label: string;
    icon: ComponentType;
    path: string;
    collapsible?: boolean;
    children?: DashboardSection[];
}

export interface DashboardStats {
    totalTasks: number;
    completedTasks: number;
    upcomingEvents: number;
    activeProjects: number;
}

export interface DashboardState {
    activeSection: string;
    stats: DashboardStats;
    isLoading: boolean;
}