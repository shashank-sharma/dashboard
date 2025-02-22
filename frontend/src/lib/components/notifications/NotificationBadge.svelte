<!-- src/lib/components/notifications/NotificationBadge.svelte -->
<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";
    import type { Notification } from "$lib/stores/notificationStore";

    export let priority: Notification["priority"];

    const variants = {
        urgent: "destructive",
        high: "orange",
        medium: "yellow",
        low: "blue",
    };

    $: variant = variants[priority] || "default";
</script>

<Badge {variant} class="capitalize">
    {priority}
</Badge>
