import { writable, get } from 'svelte/store';
import { pb } from '$lib/config/pocketbase';

export interface FileToken {
    token: string;
    expiresAt: number;
}

// Create a writable store for the file token
export const fileTokenStore = writable<FileToken | null>(null);

// For tracking token requests
let tokenRequestInProgress = false;
let tokenPromise: Promise<string> | null = null;

/**
 * Get a valid file token for accessing protected files
 * - Checks for an existing valid token in the store
 * - If no valid token exists, fetches a new one
 * - Handles concurrent requests by returning the same promise
 */
export async function getFileToken(): Promise<string> {
    const now = Date.now();

    const currentToken = get(fileTokenStore);

    if (currentToken && currentToken.expiresAt > now) {
        return currentToken.token;
    }

    if (tokenRequestInProgress && tokenPromise) {
        return tokenPromise;
    }

    tokenRequestInProgress = true;
    tokenPromise = new Promise<string>(async (resolve, reject) => {
        try {
            const token = await pb.files.getToken();
            const tokenData: FileToken = {
                token,
                expiresAt: now + 110 * 1000, // 110 seconds (slightly less than the default 2 min)
            };
            fileTokenStore.set(tokenData);

            tokenRequestInProgress = false;
            resolve(token);
        } catch (error) {
            tokenRequestInProgress = false;
            console.error("Error fetching file token:", error);
            reject("");
        }
    });

    return tokenPromise;
}

/**
 * Get the URL for a file with authentication token
 * @param fileUrl The base file URL
 * @returns Full URL with authentication token
 */
export async function getAuthenticatedFileUrl(fileUrl: string): Promise<string> {
    try {
        const token = await getFileToken();
        const separator = fileUrl.includes('?') ? '&' : '?';
        return `${fileUrl}${separator}token=${token}`;
    } catch {
        return fileUrl;
    }
} 