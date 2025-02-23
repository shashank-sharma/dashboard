<script lang="ts">
    import { dashboardStore } from "../stores/dashboard.store";
    import { Card } from "$lib/components/ui/card";
    import { CheckSquare, Calendar, Briefcase, ListTodo } from "lucide-svelte";

    const statCards = [
        {
            title: "Total Tasks",
            icon: ListTodo,
            getValue: (stats) => stats.totalTasks,
            color: "text-blue-500",
        },
        {
            title: "Completed",
            icon: CheckSquare,
            getValue: (stats) => stats.completedTasks,
            color: "text-green-500",
        },
        {
            title: "Upcoming Events",
            icon: Calendar,
            getValue: (stats) => stats.upcomingEvents,
            color: "text-purple-500",
        },
        {
            title: "Active Projects",
            icon: Briefcase,
            getValue: (stats) => stats.activeProjects,
            color: "text-orange-500",
        },
    ];
</script>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
    {#each statCards as card}
        <Card class="p-6">
            <div class="flex items-center space-x-4">
                <div class={`p-3 rounded-full ${card.color} bg-opacity-10`}>
                    <svelte:component
                        this={card.icon}
                        class={`h-6 w-6 ${card.color}`}
                    />
                </div>
                <div>
                    <p class="text-sm text-muted-foreground">{card.title}</p>
                    <h3 class="text-2xl font-bold">
                        {card.getValue($dashboardStore.stats)}
                    </h3>
                </div>
            </div>
        </Card>
    {/each}
</div>
