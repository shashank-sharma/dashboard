import { pb } from '$lib/config/pocketbase';
import type { DashboardStats } from '../types';

export class DashboardService {
    static async getStats(): Promise<DashboardStats> {
        try {
            const [tasks, events, projects] = await Promise.all([
                pb.collection('tasks').getList(1, 1, { filter: 'status != "completed"' }),
                pb.collection('events').getList(1, 1, { filter: 'date >= $now' }),
                pb.collection('projects').getList(1, 1, { filter: 'status = "active"' })
            ]);

            return {
                totalTasks: tasks.totalItems,
                completedTasks: await this.getCompletedTasksCount(),
                upcomingEvents: events.totalItems,
                activeProjects: projects.totalItems
            };
        } catch (error) {
            console.error('Error fetching dashboard stats:', error);
            throw error;
        }
    }

    private static async getCompletedTasksCount(): Promise<number> {
        const result = await pb.collection('tasks').getList(1, 1, {
            filter: 'status = "completed"'
        });
        return result.totalItems;
    }
}