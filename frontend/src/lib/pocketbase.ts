import PocketBase from 'pocketbase';
import { writable } from 'svelte/store';

export const pb = new PocketBase('http://127.0.0.1:8090');

/*
{
  avatar: "",
  collectionId: "_pb_users_auth_",
  collectionName: "users",
  created: "2023-12-08 14:01:39.544Z",
  email: "shashank.sharma98@gmail.com",
  emailVisibility: false,
  id: "k0jhrgadiakg8xb",
  name: "",
  updated: "2023-12-08 14:01:39.544Z",
  username: "users17651",
  verified: true,
}
*/

export const currentUser = writable(pb.authStore.model);
