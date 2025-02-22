<!-- src/lib/components/notifications/EmptyNotifications.svelte -->
<script lang="ts">
    import { Bell } from "lucide-svelte";
</script>

<div
    class="flex flex-col items-center justify-center p-8 text-center text-muted-foreground"
>
    <Bell class="h-12 w-12 mb-4 opacity-50" />
    <h3 class="font-medium mb-1">No notifications</h3>
    <p class="text-sm">
        You're all caught up! Check back later for new notifications.
    </p>
</div>
