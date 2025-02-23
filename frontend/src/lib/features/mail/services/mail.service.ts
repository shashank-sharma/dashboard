import { pb } from '$lib/config/pocketbase';
import type { SyncStatus, MailMessage } from '../types';

export class MailService {
    static async checkStatus(): Promise<SyncStatus | null> {
        try {
            return await pb.send<SyncStatus>('/api/mail/sync/status', {
                method: 'GET'
            });
        } catch (error) {
            throw error;
        }
    }

    static async startAuth(): Promise<string | null> {
        try {
            const response = await pb.send<{ url: string }>('/auth/mail/redirect', {
                method: 'GET'
            });
            return response.url;
        } catch (error) {
            return null;
        }
    }

    static async getMessages(page: number, perPage: number): Promise<{
        items: MailMessage[];
        totalItems: number;
    }> {
        const result = await pb.collection('mail_messages').getList(page, perPage, {
            sort: '-received_date'
        });
        
        return {
            items: result.items as MailMessage[],
            totalItems: result.totalItems
        };
    }

    static async updateMessage(id: string, data: Partial<MailMessage>): Promise<void> {
        await pb.collection('mail_messages').update(id, data);
    }

    static async syncMails(): Promise<void> {
        await pb.send('/api/mail/sync', { method: 'POST' });
    }
}