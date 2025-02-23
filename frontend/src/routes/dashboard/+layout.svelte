<script lang="ts">
    import DashboardLayout from "$lib/features/dashboard/components/DashboardLayout.svelte";
    import { DASHBOARD_SECTIONS } from "$lib/features/dashboard/constants";
    import { onMount, onDestroy } from "svelte";
    import { goto } from "$app/navigation";
    import { authIsValid } from "$lib/services/authService";
    import { notificationStore } from "$lib/features/notifications/stores/notifications.store";

    onMount(() => {
        if (!authIsValid()) {
            goto("/auth/login");
        }
    });
</script>

<DashboardLayout sections={DASHBOARD_SECTIONS}>
    <slot />
</DashboardLayout>
