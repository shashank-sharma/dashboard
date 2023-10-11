import { COOKIE_OPTIONS } from '$lib/constant';
import { currentUser, pb } from '$lib/pocketbase';

try {
	pb.authStore.loadFromCookie(document.cookie); // load it first
	// get an up-to-date auth store state by veryfing and refreshing the loaded auth model (if any)
	// pb.authStore.isValid && (await pb.collection('users').authRefresh());
    pb.authStore.isValid && (await pb.admins.authRefresh());
} catch (_) {
	// clear the auth store on failed refresh
    console.log("Clearing client authstore");
	pb.authStore.clear();
}
// pb.authStore.loadFromCookie(document.cookie); // load it first
pb.authStore.onChange((_, model) => {
    console.log("Authstore client change");
	if (pb.authStore.isValid) {
        console.log("Authstore client valid");
		currentUser.set(model);
		document.cookie = pb.authStore.exportToCookie(COOKIE_OPTIONS);
	}
});