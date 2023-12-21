import type { LayoutServerLoad } from './$types';
import { currentUser } from '$lib/pocketbase';

// Define the output type
export type OutputType = { user: object; isLoggedIn: boolean };

// Define the load function
export const load: LayoutServerLoad = async ({ locals }) => {
	if (false) {
        console.log("currentUser =", currentUser);
		// Return the output object
		return { currentUser, isLoggedIn: true };
	}
	// Return the output object
    return {}
	/* return {
		user: undefined,
		isLoggedIn: false
	};*/
};
