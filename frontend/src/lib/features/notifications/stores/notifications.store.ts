import { writable, derived } from 'svelte/store';
import { toast } from "svelte-sonner";
import { pb } from '$lib/config/pocketbase';

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

function createNotificationStore() {
    const { subscribe, set, update } = writable<{
        notifications: Notification[];
        unreadCount: number;
        loading: boolean;
        error: string | null;
    }>({
        notifications: [],
        unreadCount: 0,
        loading: false,
        error: null
    });

    return {
        subscribe,
        
        async fetchNotifications() {
            update(store => ({ ...store, loading: true, error: null }));
            try {
                const resultList = await pb.collection('notifications').getList(1, 50, {
                    filter: `user = "${pb.authStore.model?.id}"`,
                    sort: '-created',
                    expand: 'user'
                });
                
                update(store => ({
                    ...store,
                    notifications: resultList.items as Notification[],
                    unreadCount: resultList.items.filter((n: Notification) => n.status === 'unread').length,
                    loading: false
                }));
            } catch (error: any) {
                console.error('Error fetching notifications:', error);
                update(store => ({
                    ...store,
                    error: error.message,
                    loading: false
                }));
            }
        },

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
                    
                    return {
                        ...store,
                        notifications: updatedNotifications,
                        unreadCount: updatedNotifications.filter(n => n.status === 'unread').length
                    };
                });
            } catch (error: any) {
                console.error('Error marking notification as read:', error);
                toast.error('Failed to mark notification as read');
            }
        },

        async markAllAsRead() {
            try {
                const unreadNotifications = await pb.collection('notifications').getList(1, 50, {
                    filter: `user = "${pb.authStore.model?.id}" && status = "unread"`
                });

                await Promise.all(
                    unreadNotifications.items.map(notification =>
                        pb.collection('notifications').update(notification.id, {
                            status: 'read',
                            read_at: new Date().toISOString()
                        })
                    )
                );

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
            } catch (error: any) {
                console.error('Error marking all notifications as read:', error);
                toast.error('Failed to mark all notifications as read');
            }
        },

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

                if (['high', 'urgent'].includes(notification.priority)) {
                    toast(notification.title, {
                        description: notification.content,
                        duration: 5000,
                    });
                }
            } catch (error: any) {
                console.error('Error adding notification:', error);
                toast.error('Failed to create notification');
            }
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