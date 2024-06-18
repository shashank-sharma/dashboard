import { redirect } from '@sveltejs/kit';

import type { RequestEvent, RequestHandler } from './$types';

export const GET: RequestHandler = ({ locals }: RequestEvent) => {
	locals.pb?.authStore.clear();
	locals.user = undefined;
	throw redirect(303, '/');
};