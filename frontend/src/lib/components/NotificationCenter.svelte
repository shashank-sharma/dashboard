<!-- src/lib/components/NotificationCenter.svelte -->
<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { fade, slide } from "svelte/transition";
    import { goto } from "$app/navigation";
    import { Bell } from "lucide-svelte";
    import * as Popover from "$lib/components/ui/popover";
    import { Button } from "$lib/components/ui/button";
    import { Card } from "$lib/components/ui/card";
    import { Separator } from "$lib/components/ui/separator";
    import {
        notificationStore,
        groupedNotifications,
        formatNotificationTime,
        getPriorityColor,
        startNotificationPolling,
        stopNotificationPolling,
    } from "$lib/stores/notificationStore";
    import type { Notification } from "$lib/stores/notificationStore";

    let isOpen = false;

    async function handleNotificationClick(notification: Notification) {
        try {
            await notificationStore.markAsRead(notification.id);
            isOpen = false;

            switch (notification.type) {
                case "task":
                    if (notification.metadata?.taskId) {
                        goto(
                            `/dashboard/tasks/${notification.metadata.taskId}`,
                        );
                    }
                    break;
                case "habit":
                    if (notification.metadata?.habitId) {
                        goto(
                            `/dashboard/habits/${notification.metadata.habitId}`,
                        );
                    }
                    break;
                case "calendar":
                    if (notification.metadata?.eventId) {
                        goto(
                            `/dashboard/calendar/event/${notification.metadata.eventId}`,
                        );
                    }
                    break;
                case "focus":
                    if (notification.metadata?.sessionId) {
                        goto(
                            `/dashboard/focus/${notification.metadata.sessionId}`,
                        );
                    }
                    break;
                case "system":
                    if (notification.metadata?.route) {
                        goto(notification.metadata.route);
                    }
                    break;
                case "email":
                    if (notification.metadata?.emailId) {
                        goto(
                            `/dashboard/email/${notification.metadata.emailId}`,
                        );
                    }
                    break;
                default:
                    if (notification.metadata?.route) {
                        goto(notification.metadata.route);
                    }
            }
        } catch (error) {
            console.error("Error handling notification click:", error);
        }
    }

    onMount(() => {
        // startNotificationPolling();
        notificationStore.fetchNotifications();
    });

    // onDestroy(() => {
    //     stopNotificationPolling();
    // });
</script>

<Popover.Root bind:open={isOpen}>
    <Popover.Trigger asChild let:builder>
        <Button variant="ghost" size="icon" builders={[builder]}>
            <Bell class="h-5 w-5" />
            {#if $notificationStore.unreadCount > 0}
                <span class="notification-badge">
                    {$notificationStore.unreadCount > 99
                        ? "99+"
                        : $notificationStore.unreadCount}
                </span>
            {/if}
        </Button>
    </Popover.Trigger>

    <Popover.Content
        class="w-96 p-0 max-h-[32rem] overflow-hidden flex flex-col"
        sideOffset={5}
    >
        <div
            class="p-4 font-semibold border-b flex justify-between items-center"
        >
            <span>Notifications</span>
            {#if $notificationStore.unreadCount > 0}
                <Button
                    variant="ghost"
                    size="sm"
                    class="text-xs"
                    on:click={() => {
                        goto("/dashboard/settings/notifications");
                        isOpen = false;
                    }}
                >
                    View All
                </Button>
            {/if}
        </div>

        {#if $notificationStore.loading}
            <div class="p-4 text-center text-muted-foreground">
                Loading notifications...
            </div>
        {:else if $notificationStore.error}
            <div class="p-4 text-center text-destructive">
                {$notificationStore.error}
            </div>
        {:else if $groupedNotifications.length === 0}
            <div class="p-4 text-center text-muted-foreground">
                No notifications
            </div>
        {:else}
            <div class="overflow-y-auto">
                {#each $groupedNotifications as [date, notifications]}
                    <div class="p-2 bg-muted/50 text-sm text-muted-foreground">
                        {date}
                    </div>

                    {#each notifications as notification (notification.id)}
                        <div
                            class="cursor-pointer transition-colors hover:bg-muted/50 border-b"
                            on:click={() =>
                                handleNotificationClick(notification)}
                            on:keydown={(e) =>
                                e.key === "Enter" &&
                                handleNotificationClick(notification)}
                            tabindex="0"
                            role="button"
                        >
                            <div class="p-4 flex gap-4">
                                <div
                                    class={`w-1 self-stretch rounded-full ${getPriorityColor(notification.priority)}`}
                                />

                                <div class="flex-1 min-w-0">
                                    <div
                                        class="flex items-start justify-between gap-2"
                                    >
                                        <h4 class="font-medium leading-none">
                                            {notification.title}
                                        </h4>
                                        <time
                                            class="text-sm text-muted-foreground whitespace-nowrap"
                                        >
                                            {formatNotificationTime(
                                                notification.created,
                                            )}
                                        </time>
                                    </div>

                                    <p
                                        class="mt-1 text-sm text-muted-foreground line-clamp-2"
                                    >
                                        {notification.content}
                                    </p>

                                    <div class="mt-2 flex items-center gap-2">
                                        <span
                                            class="text-xs font-medium px-2 py-1 rounded-full bg-muted"
                                        >
                                            {notification.type}
                                        </span>
                                        {#if notification.status === "unread"}
                                            <span
                                                class="text-xs font-medium px-2 py-1 rounded-full bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
                                            >
                                                New
                                            </span>
                                        {/if}
                                    </div>
                                </div>
                            </div>
                        </div>
                    {/each}
                {/each}
            </div>
        {/if}
    </Popover.Content>
</Popover.Root>

<style>
    .notification-badge {
        position: absolute;
        top: -8px;
        right: -8px;
        background-color: rgb(239 68 68); /* red-500 */
        color: white;
        border-radius: 9999px;
        padding: 0;
        min-width: 18px;
        height: 18px;
        font-size: 11px;
        font-weight: 500;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 2px solid var(--background);
    }

    /* Dark mode support */
    :global(.dark) .notification-badge {
        border-color: hsl(var(--background));
    }
</style>
