import type { LayoutServerLoad } from './$types';
import { currentUser } from '$lib/pocketbase';
import { get } from 'svelte/store';

// Define the output type
export type OutputType = { user: object; isLoggedIn: boolean };

// Define the load function
export const load: LayoutServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (user) {
		// Return the output object
		return { user, isLoggedIn: true };
	}
	return {
		user: undefined,
		isLoggedIn: false
	};
};

