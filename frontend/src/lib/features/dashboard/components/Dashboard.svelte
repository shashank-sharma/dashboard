<script lang="ts">
    import { onMount } from "svelte";
    import { dashboardStore } from "../stores/dashboard.store";
    import { DashboardService } from "../services/dashboard.service";
    import { toast } from "svelte-sonner";
    import StatsGrid from "./StatsGrid.svelte";
    import RecentActivity from "./RecentActivity.svelte";

    onMount(async () => {
        try {
            dashboardStore.setLoading(true);
            const stats = await DashboardService.getStats();
            dashboardStore.setStats(stats);
        } catch (error) {
            toast.error("Failed to load dashboard data");
        } finally {
            dashboardStore.setLoading(false);
        }
    });
</script>

<div class="space-y-6">
    <StatsGrid />
    <RecentActivity />
</div>
