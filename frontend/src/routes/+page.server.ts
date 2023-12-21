import type { ITasks } from '$lib/types';
import type { PageServerLoad } from './$types';
import { pb, currentUser } from '$lib/pocketbase';

export const load = (async ({ request, cookies }) => {
	try {
		const tasks = await pb.collection('tasks').getList<ITasks>(1, 10);
        console.log("tasks=",tasks);
		return { tasks: structuredClone(tasks) };
	} catch (err) {
        console.log("Failed to fetch todos");
    }
}) satisfies PageServerLoad;
