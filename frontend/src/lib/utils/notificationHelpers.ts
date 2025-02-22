// src/lib/utils/notificationHelpers.ts
import type { Notification } from '$lib/stores/notificationStore';
import { goto } from '$app/navigation';

export function handleNotificationAction(notification: Notification): void {
    const metadata = notification.metadata;
    
    if (!metadata) return;

    switch (notification.type) {
        case 'task':
            if (metadata.taskId) {
                goto(`/tasks/${metadata.taskId}`);
            }
            break;
            
        case 'habit':
            if (metadata.habitId) {
                goto(`/habits/${metadata.habitId}`);
            }
            break;
            
        case 'calendar':
            if (metadata.eventId) {
                goto(`/calendar/event/${metadata.eventId}`);
            }
            break;
            
        case 'email':
            if (metadata.emailId) {
                goto(`/email/${metadata.emailId}`);
            }
            break;
            
        case 'focus':
            if (metadata.sessionId) {
                goto(`/focus/${metadata.sessionId}`);
            }
            break;
    }

    // Handle custom actions
    if (metadata.action) {
        switch (metadata.action.type) {
            case 'navigate':
                if (metadata.action.url) {
                    goto(metadata.action.url);
                }
                break;
                
            case 'external':
                if (metadata.action.url) {
                    window.open(metadata.action.url, '_blank');
                }
                break;
                
            case 'custom':
                if (typeof metadata.action.handler === 'function') {
                    metadata.action.handler();
                }
                break;
        }
    }
}

export function getNotificationIcon(type: NotificationType): string {
    switch (type) {
        case 'task':
            return 'check-square';
        case 'habit':
            return 'activity';
        case 'focus':
            return 'clock';
        case 'calendar':
            return 'calendar';
        case 'email':
            return 'mail';
        default:
            return 'bell';
    }
}

export function getPriorityLevel(priority: NotificationPriority): number {
    switch (priority) {
        case 'urgent':
            return 4;
        case 'high':
            return 3;
        case 'medium':
            return 2;
        case 'low':
            return 1;
        default:
            return 0;
    }
}

export function sortNotificationsByPriority(notifications: Notification[]): Notification[] {
    return [...notifications].sort((a, b) => {
        // Sort by priority first
        const priorityDiff = getPriorityLevel(b.priority) - getPriorityLevel(a.priority);
        if (priorityDiff !== 0) return priorityDiff;
        
        // Then by creation date
        return new Date(b.created).getTime() - new Date(a.created).getTime();
    });
}

export function groupNotificationsByDate(notifications: Notification[]): Record<string, Notification[]> {
    return notifications.reduce((groups, notification) => {
        const date = new Date(notification.created).toLocaleDateString();
        if (!groups[date]) {
            groups[date] = [];
        }
        groups[date].push(notification);
        return groups;
    }, {} as Record<string, Notification[]>);
}