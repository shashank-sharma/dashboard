import { writable } from 'svelte/store'

type Notification = string

// export const notifications = writable<Notification[]>([])
export const notifications = writable<Notification>("")

export function toast(message: string) {
    notifications.set(message);
  // setTimeout(removeToast, 2000)
}
/*
export function toast(message: string) {
  notifications.update((state) => [message, ...state])
  // setTimeout(removeToast, 2000)
}

function removeToast() {
  notifications.update((state) => {
    return [...state.slice(0, state.length - 1)]
  })
}*/
