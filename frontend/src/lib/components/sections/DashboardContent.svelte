<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "$lib/pocketbase";
    import { fade } from "svelte/transition";
    import {
        Card,
        CardContent,
        CardHeader,
        CardTitle,
    } from "$lib/components/ui/card";
    import { Calendar } from "$lib/components/ui/calendar";
    import { Badge } from "$lib/components/ui/badge";
    import { ScrollArea } from "$lib/components/ui/scroll-area";
    import {
        Clock,
        Trophy,
        Activity,
        Brain,
        Calendar as CalendarIcon,
    } from "lucide-svelte";
    import { Progress } from "$lib/components/ui/progress";

    let selectedDate = new Date();
    let tasksData = [];
    let isLoading = true;

    const mockTasksData = [
        {
            id: 1,
            title: "Review Project Proposal",
            priority: "high",
            category: "focus",
            due: "2024-12-27",
        },
        {
            id: 2,
            title: "Update Documentation",
            priority: "medium",
            category: "goals",
            due: "2024-12-28",
        },
        {
            id: 3,
            title: "Team Meeting",
            priority: "low",
            category: "fitin",
            due: "2024-12-26",
        },
    ];

    function getPriorityColor(priority: string): string {
        const colors = {
            high: "text-red-500 dark:text-red-400",
            medium: "text-yellow-500 dark:text-yellow-400",
            low: "text-green-500 dark:text-green-400",
        };
        return colors[priority] || colors.low;
    }

    async function fetchDashboardData() {
        try {
            const today = new Date().toISOString().split("T")[0];

            const tasksRecords = await pb.collection("tasks").getList(1, 5, {
                sort: "-created",
            });

            tasksData = mockTasksData;
        } catch (error) {
            console.error("Error fetching dashboard data:", error);
        } finally {
            isLoading = false;
        }
    }

    onMount(() => {
        fetchDashboardData();
    });
</script>

<div
    class="grid gap-4 md:grid-cols-2 lg:grid-cols-3"
    in:fade={{ duration: 300 }}
>
    <!-- Quick Stats -->
    <Card class="col-span-1 md:col-span-2 lg:col-span-2">
        <CardHeader>
            <CardTitle class="flex items-center gap-2">
                <Activity class="h-5 w-5" />
                Overview
            </CardTitle>
        </CardHeader>
        <CardContent class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            <div class="space-y-2">
                <p class="text-sm font-medium text-muted-foreground">
                    Tasks Completed
                </p>
                <p class="text-2xl font-bold">12/15</p>
                <Progress value={80} />
            </div>
            <div class="space-y-2">
                <p class="text-sm font-medium text-muted-foreground">
                    Habits Streak
                </p>
                <p class="text-2xl font-bold">7 days</p>
                <Progress value={70} />
            </div>
            <div class="space-y-2">
                <p class="text-sm font-medium text-muted-foreground">
                    Daily Score
                </p>
                <p class="text-2xl font-bold">4.5/5</p>
                <Progress value={90} />
            </div>
        </CardContent>
    </Card>

    <!-- Calendar Card -->
    <Card class="col-span-1">
        <CardHeader>
            <CardTitle class="flex items-center gap-2">
                <CalendarIcon class="h-5 w-5" />
                Calendar
            </CardTitle>
        </CardHeader>
        <CardContent>
            <Calendar
                type="single"
                bind:selectedDate
                class="rounded-md border"
            />
        </CardContent>
    </Card>

    <!-- Tasks Card -->
    <Card class="col-span-1 md:col-span-2">
        <CardHeader>
            <CardTitle class="flex items-center gap-2">
                <Activity class="h-5 w-5" />
                Recent Tasks
            </CardTitle>
        </CardHeader>
        <CardContent>
            <ScrollArea class="h-[300px] pr-4">
                {#each tasksData as task (task.id)}
                    <div
                        class="mb-4 rounded-lg border p-3 transition-all hover:shadow-md"
                        in:fade={{ duration: 200 }}
                    >
                        <div class="flex items-center justify-between">
                            <h3 class="font-medium">{task.title}</h3>
                            <Badge
                                variant="outline"
                                class={getPriorityColor(task.priority)}
                            >
                                {task.priority}
                            </Badge>
                        </div>
                        <div
                            class="mt-2 flex items-center text-sm text-muted-foreground"
                        >
                            <Clock class="mr-2 h-4 w-4" />
                            {task.due}
                        </div>
                    </div>
                {/each}
            </ScrollArea>
        </CardContent>
    </Card>

    <!-- Habits Summary -->
    <Card class="col-span-1">
        <CardHeader>
            <CardTitle class="flex items-center gap-2">
                <Trophy class="h-5 w-5" />
                Habits Summary
            </CardTitle>
        </CardHeader>
        <CardContent>
            <div class="space-y-4">
                <div class="space-y-2">
                    <div class="flex items-center justify-between">
                        <p class="text-sm font-medium">Morning Routine</p>
                        <Badge variant="outline">7 days</Badge>
                    </div>
                    <Progress value={100} />
                </div>
                <div class="space-y-2">
                    <div class="flex items-center justify-between">
                        <p class="text-sm font-medium">Exercise</p>
                        <Badge variant="outline">5 days</Badge>
                    </div>
                    <Progress value={71} />
                </div>
                <div class="space-y-2">
                    <div class="flex items-center justify-between">
                        <p class="text-sm font-medium">Reading</p>
                        <Badge variant="outline">3 days</Badge>
                    </div>
                    <Progress value={43} />
                </div>
            </div>
        </CardContent>
    </Card>
</div>
