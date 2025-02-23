<script lang="ts">
    import {
        notificationStore,
        type Notification,
    } from "../stores/notifications.store";
    import { Button } from "$lib/components/ui/button";
    import { Bell } from "lucide-svelte";
    import { Badge } from "$lib/components/ui/badge";
    import { Card } from "$lib/components/ui/card";
    import { onMount, onDestroy } from "svelte";
    import { formatDistanceToNow } from "date-fns";
    import { getPriorityColor } from "$lib/utils";
    import { browser } from "$app/environment";

    let isOpen = false;

    onMount(() => {
        notificationStore.fetchNotifications();
        if (browser) {
            document.addEventListener("click", handleClickOutside);
        }
    });

    onDestroy(() => {
        if (browser) {
            document.removeEventListener("click", handleClickOutside);
        }
    });

    function handleNotificationClick(notification: Notification) {
        if (notification.status === "unread") {
            notificationStore.markAsRead(notification.id);
        }
        isOpen = false;
    }

    function handleClickOutside(event: MouseEvent) {
        if (!browser) return;

        const notificationElement = document.getElementById(
            "notification-center",
        );
        if (
            notificationElement &&
            !notificationElement.contains(event.target as Node)
        ) {
            isOpen = false;
        }
    }
</script>

<div class="relative" id="notification-center">
    <Button variant="ghost" size="icon" on:click={() => (isOpen = !isOpen)}>
        <Bell class="h-5 w-5" />
        {#if $notificationStore.unreadCount > 0}
            <Badge class="absolute -top-1 -right-1 h-5 w-5">
                {$notificationStore.unreadCount}
            </Badge>
        {/if}
    </Button>

    {#if isOpen}
        <Card
            class="absolute right-0 mt-2 w-96 max-h-[500px] z-50 overflow-hidden flex flex-col"
        >
            <div class="p-4 border-b flex items-center justify-between">
                <span class="font-medium">Notifications</span>
                {#if $notificationStore.unreadCount > 0}
                    <Button
                        variant="ghost"
                        size="sm"
                        on:click={() => notificationStore.markAllAsRead()}
                    >
                        Mark all as read
                    </Button>
                {/if}
            </div>

            <div class="overflow-y-auto">
                {#if $notificationStore.loading}
                    <div class="p-4 text-center text-muted-foreground">
                        Loading...
                    </div>
                {:else if $notificationStore.notifications.length === 0}
                    <div class="p-4 text-center text-muted-foreground">
                        No notifications
                    </div>
                {:else}
                    {#each $notificationStore.notifications as notification}
                        <div
                            class="p-4 border-b hover:bg-muted/50 cursor-pointer"
                            class:bg-muted={notification.status === "unread"}
                            role="button"
                            tabindex="0"
                            on:click={() =>
                                handleNotificationClick(notification)}
                            on:keydown={(e) =>
                                e.key === "Enter" &&
                                handleNotificationClick(notification)}
                        >
                            <div class="flex items-start justify-between gap-2">
                                <div>
                                    <p class="text-sm font-medium">
                                        {notification.title}
                                    </p>
                                    <p class="text-sm text-muted-foreground">
                                        {notification.content}
                                    </p>
                                </div>
                                <Badge
                                    variant="outline"
                                    class={getPriorityColor(
                                        notification.priority,
                                    )}
                                >
                                    {notification.priority}
                                </Badge>
                            </div>
                            <p class="text-xs text-muted-foreground mt-1">
                                {formatDistanceToNow(
                                    new Date(notification.created),
                                    { addSuffix: true },
                                )}
                            </p>
                        </div>
                    {/each}
                {/if}
            </div>
        </Card>
    {/if}
</div>
