<script lang="ts">
    import { Card } from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";
    import { Progress } from "$lib/components/ui/progress";
    import { Badge } from "$lib/components/ui/badge";
    import { Separator } from "$lib/components/ui/separator";
    import { toast } from "svelte-sonner";
    import HeartRateChart from "./HeartRateChart.svelte";
    import { Clock, Laptop, Droplet, CheckCircle2, Activity } from "lucide-svelte";

    let behaviorRating = 3;
    let waterIntake = 5;
    const waterGoal = 8;

    const appUsageData = {
        Work: 180,
        Entertainment: 120,
        Social: 90,
        Development: 240,
    };

    const tasks = [
        { title: "Review report", due: "5 PM", priority: "high" },
        { title: "Team meeting", due: "3:30", priority: "medium" },
    ];

    function updateBehaviorRating(rating: number) {
        behaviorRating = rating;
        toast.success(`Rating updated: ${rating}`);
    }

    function addWater() {
        if (waterIntake < waterGoal) {
            waterIntake++;
            toast.success("Water intake logged");
        }
    }
</script>

<div class="grid grid-cols-2 gap-2 h-[calc(100vh-150px)] overflow-y-auto pr-2">
    <!-- Sleep -->
    <Card class="p-3 col-span-1">
        <div class="flex items-center justify-between">
            <div class="flex items-center gap-1.5">
                <Clock class="w-4 h-4 text-blue-500" />
                <span class="text-sm font-medium">Sleep</span>
            </div>
            <Badge variant="outline" class="text-xs">7h 30m</Badge>
        </div>
        <div class="mt-1 text-xs text-muted-foreground">
            10:30 PM - 6:00 AM
        </div>
    </Card>

    <!-- Water -->
    <Card class="p-3 col-span-1">
        <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-1.5">
                <Droplet class="w-4 h-4 text-blue-500" />
                <span class="text-sm font-medium">Water</span>
            </div>
            <Button variant="ghost" size="sm" class="h-7" on:click={addWater}>+</Button>
        </div>
        <Progress value={(waterIntake / waterGoal) * 100} class="h-2" />
        <div class="mt-1 text-xs text-muted-foreground">
            {waterIntake}/{waterGoal} glasses
        </div>
    </Card>

    <!-- Screen Time -->
    <Card class="p-3 col-span-2">
        <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-1.5">
                <Laptop class="w-4 h-4 text-purple-500" />
                <span class="text-sm font-medium">Screen Time</span>
            </div>
            <Badge variant="outline" class="text-xs">11h 30m</Badge>
        </div>
        <div class="grid grid-cols-2 gap-2">
            {#each Object.entries(appUsageData) as [category, minutes]}
                <div>
                    <div class="flex justify-between text-xs mb-1">
                        <span>{category}</span>
                        <span class="text-muted-foreground">{Math.floor(minutes / 60)}h {minutes % 60}m</span>
                    </div>
                    <Progress value={(minutes / 300) * 100} class="h-1.5" />
                </div>
            {/each}
        </div>
    </Card>

    <!-- Heart Rate -->
    <div class="col-span-2">
        <HeartRateChart />
    </div>

    <!-- Tasks -->
    <Card class="p-3 col-span-2">
        <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-1.5">
                <CheckCircle2 class="w-4 h-4 text-green-500" />
                <span class="text-sm font-medium">Tasks</span>
            </div>
            <Badge variant="outline" class="text-xs">{tasks.length}</Badge>
        </div>
        <div class="space-y-1.5">
            {#each tasks as task}
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-xs font-medium">{task.title}</p>
                        <p class="text-xs text-muted-foreground">Due: {task.due}</p>
                    </div>
                    <Badge
                        variant={task.priority === "high" ? "destructive" : "outline"}
                        class="text-xs"
                    >
                        {task.priority}
                    </Badge>
                </div>
            {/each}
        </div>
    </Card>

    <!-- Day Rating -->
    <Card class="p-3 col-span-2">
        <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-1.5">
                <Activity class="w-4 h-4 text-yellow-500" />
                <span class="text-sm font-medium">Day Rating</span>
            </div>
        </div>
        <div class="flex justify-between gap-1">
            {#each Array(5) as _, i}
                <Button
                    variant={behaviorRating === i + 1 ? "default" : "outline"}
                    size="sm"
                    class="flex-1 h-8"
                    on:click={() => updateBehaviorRating(i + 1)}
                >
                    {i + 1}
                </Button>
            {/each}
        </div>
    </Card>
</div>

<style>
    :global(.metadata-scroll::-webkit-scrollbar) {
        width: 4px;
    }
    :global(.metadata-scroll::-webkit-scrollbar-track) {
        background: transparent;
    }
    :global(.metadata-scroll::-webkit-scrollbar-thumb) {
        background-color: hsl(var(--muted));
        border-radius: 2px;
    }
</style>