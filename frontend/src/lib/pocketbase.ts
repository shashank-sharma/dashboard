import PocketBase from 'pocketbase';
import { writable } from 'svelte/store';

export const pb = new PocketBase('http://127.0.0.1:8090'); 

export const currentUser = writable(pb.authStore.model);

pb.authStore.onChange((auth) => {
    console.log('authStore changed', auth);
    currentUser.set(pb.authStore.model);
});

export async function refreshToken() {
  console.log("Is valid: ", !pb.authStore.isValid);
  if (!pb.authStore.isValid) return;
  
  try {
      // Attempt to refresh the token
      const userData = await pb.collection('users').authRefresh();
      console.log('Token refreshed for user:', userData.record.id);
  } catch (err) {
      console.error('Failed to refresh token:', err);
      // If refresh fails, clear the auth store and redirect to login
      pb.authStore.clear();
      window.location.href = '/auth/login';
  }
}

// Set up an interval to check and refresh the token
setInterval(refreshToken, 10 * 60 * 1000);