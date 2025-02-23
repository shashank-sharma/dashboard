import { pb } from '$lib/config/pocketbase';
import { writable } from 'svelte/store';
import type { User } from '$lib/types';

export const currentUser = writable<User | null>(pb.authStore.model);

export function authIsValid(): boolean {
    return pb.authStore.isValid;
}

export async function login(username: string, password: string) {
    const authData = await pb.collection('users').authWithPassword(username, password);
    currentUser.set(authData.record);
    return authData;
}

export async function logout() {
    pb.authStore.clear();
    currentUser.set(null);
}