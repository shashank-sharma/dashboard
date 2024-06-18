import type { ITasks } from '$lib/types';
import type { PageServerLoad } from './$types';
import { pb, currentUser } from '$lib/pocketbase';
import { get } from 'svelte/store';

export const load = (async ({ locals, url }) => {
	try {
		const tasks = await locals.pb?.collection('tasks').getList<ITasks>(1, 10);
		return { tasks: structuredClone(tasks) };
	} catch (err) {
        console.log("Failed to fetch todos");
    }
}) satisfies PageServerLoad;
