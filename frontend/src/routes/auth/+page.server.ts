import { redirect, type Actions, fail } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
export const load = (async ({ locals }) => {
	// if user is logged in, send them back to home
	if (locals.pb.authStore.isValid) throw redirect(301, '/');
	//  return { user: structuredClone(locals.pb.authStore.model) };
}) satisfies PageServerLoad;

interface FormResError {
	code: FormFieldKey | 'unknown'; // error can be unknown field
	message: string;
}
enum FormFieldKey {
	EmailOrUsername = 'emailOrUsername',
	Password = 'password'
}
const makeErrObj = (message = '', code: FormFieldKey | 'unknown' = 'unknown'): FormResError => ({
	code,
	message
});

export const actions = {
	login: async ({ locals, request }) => {
		const formData = await request.formData();
		const emailOrUsername = formData.get(FormFieldKey.EmailOrUsername)?.toString();
		const password = formData.get(FormFieldKey.Password)?.toString();
		// validate emailOrUsername
		if (!emailOrUsername || emailOrUsername == '')
			return fail(500, {
				emailOrUsername,
				error: makeErrObj('Empty email address or username', FormFieldKey.EmailOrUsername)
			});
		// validate password
		if (!password || password == '')
			return fail(500, {
				emailOrUsername,
				error: makeErrObj('Empty password', FormFieldKey.Password)
			});

		try {
			// const { token, record } = await locals.pb
			// 	.collection('users')
			// 	.authWithPassword(emailOrUsername.toLowerCase(), password);
            console.log("Logging in")
            const { token, record } = await locals.pb
            .admins
            .authWithPassword(emailOrUsername.toLowerCase(), password);
            console.log("Done", token)
			// return { success: true };
		} catch (err) {
			// in case of user, avoid showing actual error to prevent fishing of information
			return fail(500, {
				emailOrUsername,
				error: makeErrObj('Incorrect credentials!', 'unknown')
			});
		}
		throw redirect(301, '/');
	}
} satisfies Actions;