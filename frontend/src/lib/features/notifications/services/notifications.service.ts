import { pb } from '$lib/config/pocketbase';
import type { Notification } from '../types';

export class NotificationsService {
    private static subscription: any = null;

    static async subscribe(userId: string, onNotification: (data: any) => void) {
        // Unsubscribe from any existing subscription
        this.unsubscribe();

        // Subscribe to realtime notifications
        this.subscription = pb.collection('notifications').subscribe('*', async (data) => {
            if (data.record.user === userId) {
                onNotification(data);
            }
        });
    }

    static unsubscribe() {
        if (this.subscription) {
            this.subscription.unsubscribe();
        }
    }

    static async create(notification: Partial<Notification>) {
        return await pb.collection('notifications').create(notification);
    }

    static async markAsRead(id: string) {
        return await pb.collection('notifications').update(id, { read: true });
    }

    static async markAllAsRead(userId: string) {
        const records = await pb.collection('notifications').getFullList({
            filter: `user = "${userId}" && read = false`
        });
        
        return Promise.all(
            records.map(record => 
                pb.collection('notifications').update(record.id, { read: true })
            )
        );
    }
}