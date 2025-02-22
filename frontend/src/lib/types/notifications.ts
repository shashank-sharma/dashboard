// src/lib/types/notifications.ts
export type NotificationType = 'system' | 'task' | 'habit' | 'focus' | 'calendar' | 'email';
export type NotificationPriority = 'low' | 'medium' | 'high' | 'urgent';
export type NotificationStatus = 'unread' | 'read' | 'archived';

export interface NotificationAction {
    type: 'navigate' | 'external' | 'custom';
    label: string;
    url?: string;
    handler?: () => void;
}

export interface NotificationMetadata {
    action?: NotificationAction;
    taskId?: string;
    habitId?: string;
    eventId?: string;
    [key: string]: any;
}