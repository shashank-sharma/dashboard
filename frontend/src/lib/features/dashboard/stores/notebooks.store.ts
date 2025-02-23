import { createBaseStore } from './base.store';
import type { Notebook } from '../types';

export const notebooksStore = createBaseStore<Notebook>('notebooks');