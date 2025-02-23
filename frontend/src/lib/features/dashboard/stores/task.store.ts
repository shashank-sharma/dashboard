import { createBaseStore } from './base.store';
import type { Task } from '../types';

export const tasksStore = createBaseStore<Task>('tasks');