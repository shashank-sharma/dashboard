import { pb } from '$lib/config/pocketbase';
import type { Task } from '../types';

export class TasksService {
    static async fetchTasks(searchQuery?: string): Promise<Task[]> {
        try {
            const filterConditions = [`user = "${pb.authStore.model?.id}"`];
            if (searchQuery) {
                filterConditions.push(
                    `(title ~ "${searchQuery}" || description ~ "${searchQuery}")`
                );
            }

            const filter = filterConditions.join(" && ");
            const resultList = await pb.collection("tasks").getList(1, 50, {
                sort: "-created",
                filter,
            });

            return resultList.items as Task[];
        } catch (error) {
            throw new Error('Failed to fetch tasks');
        }
    }

    static async updateTask(taskId: string, data: Partial<Task>): Promise<Task> {
        try {
            const updated = await pb.collection("tasks").update(taskId, data);
            return updated as Task;
        } catch (error) {
            throw new Error('Failed to update task');
        }
    }

    static async deleteTask(taskId: string): Promise<void> {
        try {
            await pb.collection("tasks").delete(taskId);
        } catch (error) {
            throw new Error('Failed to delete task');
        }
    }
}