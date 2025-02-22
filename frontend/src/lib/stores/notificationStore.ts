// src/lib/stores/notificationStore.ts
import { writable, derived } from 'svelte/store';
import { toast } from "svelte-sonner";
import { pb } from '$lib/pocketbase';

export interface Notification {
    id: string;
    user: string;
    type: 'system' | 'task' | 'habit' | 'focus' | 'calendar' | 'email';
    title: string;
    content: string;
    priority: 'low' | 'medium' | 'high' | 'urgent';
    status: 'unread' | 'read' | 'archived';
    metadata: Record<string, any>;
    read_at: string | null;
    created: string;
    updated: string;
}

interface NotificationStore {
    notifications: Notification[];
    unreadCount: number;
    loading: boolean;
    error: string | null;
}

const initialState: NotificationStore = {
    notifications: [],
    unreadCount: 0,
    loading: false,
    error: null
};

function createNotificationStore() {
    const { subscribe, set, update } = writable<NotificationStore>(initialState);

    return {
        subscribe,
        
        // Fetch notifications
        async fetchNotifications() {
            update(store => ({ ...store, loading: true, error: null }));
            try {
                const resultList = await pb.collection('notifications').getList(1, 50, {
                    filter: `user = "${pb.authStore.model?.id}"`,
                    sort: '-created',
                    expand: 'user'
                });
                
                // Also update unread count while we're at it
                const unreadCount = resultList.items.filter(
                    (item: any) => item.status === 'unread'
                ).length;

                update(store => ({
                    ...store,
                    notifications: resultList.items as Notification[],
                    unreadCount,
                    loading: false
                }));
            } catch (error) {
                console.error('Error fetching notifications:', error);
                update(store => ({
                    ...store,
                    error: error.message,
                    loading: false
                }));
            }
        },

        // Fetch unread count specifically
        async fetchUnreadCount() {
            try {
                const resultList = await pb.collection('notifications').getList(1, 1, {
                    filter: `user = "${pb.authStore.model?.id}" && status = "unread"`,
                    fields: 'id' // Only fetch IDs to minimize data transfer
                });
                
                update(store => ({
                    ...store,
                    unreadCount: resultList.totalItems
                }));
            } catch (error) {
                console.error('Error fetching unread count:', error);
            }
        },

        // Mark notification as read
        async markAsRead(notificationId: string) {
            try {
                await pb.collection('notifications').update(notificationId, {
                    status: 'read',
                    read_at: new Date().toISOString()
                });
                
                update(store => {
                    const updatedNotifications = store.notifications.map(n => 
                        n.id === notificationId 
                            ? { ...n, status: 'read', read_at: new Date().toISOString() }
                            : n
                    );
                    
                    // Count unread notifications
                    const unreadCount = updatedNotifications.filter(n => n.status === 'unread').length;
                    
                    return {
                        ...store,
                        notifications: updatedNotifications,
                        unreadCount
                    };
                });
            } catch (error) {
                console.error('Error marking notification as read:', error);
                toast.error('Failed to mark notification as read');
            }
        },

        // Mark all as read
        async markAllAsRead() {
            try {
                // Get all unread notifications
                const unreadNotifications = await pb.collection('notifications').getList(1, 50, {
                    filter: `user = "${pb.authStore.model?.id}" && status = "unread"`
                });

                // Update all unread notifications
                await Promise.all(
                    unreadNotifications.items.map(notification =>
                        pb.collection('notifications').update(notification.id, {
                            status: 'read',
                            read_at: new Date().toISOString()
                        })
                    )
                );

                // Update store
                update(store => ({
                    ...store,
                    notifications: store.notifications.map(n => ({
                        ...n,
                        status: 'read',
                        read_at: new Date().toISOString()
                    })),
                    unreadCount: 0
                }));

                toast.success('All notifications marked as read');
            } catch (error) {
                console.error('Error marking all notifications as read:', error);
                toast.error('Failed to mark all notifications as read');
            }
        },

        // Add new notification
        async addNotification(notification: Omit<Notification, 'id' | 'created' | 'updated' | 'user'>) {
            try {
                const newNotification = await pb.collection('notifications').create({
                    ...notification,
                    user: pb.authStore.model?.id,
                    status: 'unread',
                    read_at: null
                });

                update(store => ({
                    ...store,
                    notifications: [newNotification as Notification, ...store.notifications],
                    unreadCount: store.unreadCount + 1
                }));

                // Show toast for high/urgent priority
                if (['high', 'urgent'].includes(notification.priority)) {
                    toast(notification.title, {
                        description: notification.content,
                        duration: 5000,
                    });
                }
            } catch (error) {
                console.error('Error adding notification:', error);
                toast.error('Failed to create notification');
            }
        },

        // Reset store
        reset() {
            set(initialState);
        }
    };
}

export const notificationStore = createNotificationStore();

// Derived stores
export const unreadNotifications = derived(notificationStore, $store => 
    $store.notifications.filter(n => n.status === 'unread')
);

export const groupedNotifications = derived(notificationStore, $store => {
    const groups = $store.notifications.reduce((acc, notification) => {
        const date = new Date(notification.created).toLocaleDateString();
        if (!acc[date]) acc[date] = [];
        acc[date].push(notification);
        return acc;
    }, {} as Record<string, Notification[]>);
    
    return Object.entries(groups).sort((a, b) => 
        new Date(b[0]).getTime() - new Date(a[0]).getTime()
    );
});

// Initialize polling
let pollInterval: NodeJS.Timeout;

export function startNotificationPolling() {
    // Initial fetch
    notificationStore.fetchNotifications();
    
    // Set up polling every 30 seconds
    pollInterval = setInterval(() => {
        notificationStore.fetchNotifications();
    }, 30000);
}

export function stopNotificationPolling() {
    if (pollInterval) clearInterval(pollInterval);
}

// Utility functions
export function formatNotificationTime(timestamp: string): string {
    const date = new Date(timestamp);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (minutes < 60) return `${minutes}m ago`;
    if (hours < 24) return `${hours}h ago`;
    if (days < 7) return `${days}d ago`;
    return date.toLocaleDateString();
}

export function getPriorityColor(priority: Notification['priority']): string {
    switch (priority) {
        case 'urgent': return 'text-red-500 dark:text-red-400';
        case 'high': return 'text-orange-500 dark:text-orange-400';
        case 'medium': return 'text-yellow-500 dark:text-yellow-400';
        case 'low': return 'text-blue-500 dark:text-blue-400';
        default: return 'text-foreground';
    }
}