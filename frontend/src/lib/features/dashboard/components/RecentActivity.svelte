<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "$lib/config/pocketbase";
    import { Card } from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";
    import { ScrollArea } from "$lib/components/ui/scroll-area";
    import { formatDistanceToNow } from "date-fns";

    interface Activity {
        id: string;
        type: string;
        description: string;
        created: string;
        metadata: Record<string, any>;
    }

    let activities: Activity[] = [];
    let loading = true;

    async function fetchRecentActivity() {
        try {
            const records = await pb.collection("activities").getList(1, 5, {
                sort: "-created",
                expand: "user",
            });
            activities = records.items;
        } catch (error) {
            console.error("Error fetching activities:", error);
        } finally {
            loading = false;
        }
    }

    function getActivityIcon(type: string) {
        const icons = {
            task: "✓",
            note: "📝",
            login: "🔑",
            project: "📊",
        };
        return icons[type] || "📌";
    }

    onMount(fetchRecentActivity);
</script>

<Card class="p-6">
    <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-semibold">Recent Activity</h2>
        <Button variant="ghost" size="sm">View all</Button>
    </div>

    <ScrollArea class="h-[400px] pr-4">
        {#if loading}
            <div class="flex items-center justify-center h-32">
                <div
                    class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"
                />
            </div>
        {:else if activities.length === 0}
            <div class="text-center text-muted-foreground py-8">
                No recent activity
            </div>
        {:else}
            <div class="space-y-4">
                {#each activities as activity}
                    <div class="flex items-start gap-4">
                        <div
                            class="w-8 h-8 flex items-center justify-center rounded-full bg-muted"
                        >
                            <span>{getActivityIcon(activity.type)}</span>
                        </div>
                        <div class="flex-1 space-y-1">
                            <p class="text-sm">
                                {activity.description}
                            </p>
                            <p class="text-xs text-muted-foreground">
                                {formatDistanceToNow(
                                    new Date(activity.created),
                                    { addSuffix: true },
                                )}
                            </p>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </ScrollArea>
</Card>
