import PocketBase from 'pocketbase';
import { writable } from 'svelte/store';
import { browser } from '$app/environment';

const POCKETBASE_URL = browser ? 
    (localStorage.getItem('pocketbaseUrl') || 'http://127.0.0.1:8090') : 
    'http://127.0.0.1:8090';

export const pb = new PocketBase(POCKETBASE_URL);
export const currentUser = writable(pb.authStore.model);

pb.authStore.onChange((auth) => {
    currentUser.set(pb.authStore.model);
});