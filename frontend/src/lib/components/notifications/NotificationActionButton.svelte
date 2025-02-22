<!-- src/lib/components/notifications/NotificationActionButton.svelte -->
<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import type { Notification } from "$lib/stores/notificationStore";
    import { goto } from "$app/navigation";

    export let notification: Notification;
    export let onAction: () => void = () => {};

    function handleClick() {
        if (notification.metadata?.action) {
            switch (notification.metadata.action.type) {
                case "navigate":
                    goto(notification.metadata.action.url);
                    break;
                case "external":
                    window.open(notification.metadata.action.url, "_blank");
                    break;
                case "custom":
                    if (
                        typeof notification.metadata.action.handler ===
                        "function"
                    ) {
                        notification.metadata.action.handler();
                    }
                    break;
            }
        }
        onAction();
    }
</script>

<Button
    variant="ghost"
    size="sm"
    class="w-full justify-start"
    on:click={handleClick}
>
    <slot />
</Button>
