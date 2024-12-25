import { get } from 'svelte/store';
import { currentUser, pb } from '$lib/pocketbase';

export async function createWithUser(collectionName: string, data: any) {
    const userId = get(currentUser)?.id;
    if (!userId) {
        throw new Error('User not authenticated');
    }
    const dataWithUser = { ...data, user: userId };
    return pb.collection(collectionName).create(dataWithUser);
}

export async function updateWithUser(collectionName: string, id: string, data: any) {
    const userId = get(currentUser)?.id;
    if (!userId) {
        throw new Error('User not authenticated');
    }
    const dataWithUser = { ...data, user: userId };
    return pb.collection(collectionName).update(id, dataWithUser);
}