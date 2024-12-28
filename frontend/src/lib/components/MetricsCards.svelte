<script lang="ts">
    import {
        Card,
        CardContent,
        CardHeader,
        CardTitle,
    } from "$lib/components/ui/card";
    import { Progress } from "$lib/components/ui/progress";
    import {
        LineChart,
        Line,
        XAxis,
        YAxis,
        CartesianGrid,
        Tooltip,
        ResponsiveContainer,
    } from "recharts";
    import {
        Brain,
        Coffee,
        Monitor,
        Utensils,
        Droplets,
        Heart,
        CheckSquare,
        ListTodo,
    } from "lucide-svelte";
    import { Badge } from "$lib/components/ui/badge";
    import { Separator } from "$lib/components/ui/separator";

    // Dummy data - Replace with actual API data later
    const heartRateData = Array.from({ length: 24 }, (_, i) => ({
        hour: i,
        rate: 60 + Math.random() * 20,
    }));

    const appUsageData = [
        { name: "Work", value: 180, color: "bg-blue-500" },
        { name: "Entertainment", value: 60, color: "bg-purple-500" },
        { name: "Social", value: 45, color: "bg-pink-500" },
        { name: "Productivity", value: 90, color: "bg-green-500" },
    ];

    const metrics = [
        {
            title: "Sleep Time",
            value: "7h 30m",
            icon: Brain,
            color: "text-blue-500",
            detail: {
                bedtime: "23:00",
                wakeup: "06:30",
                quality: "Good",
            },
        },
        {
            title: "Screen Time",
            value: "6h 15m",
            icon: Monitor,
            color: "text-purple-500",
            detail: appUsageData,
        },
        {
            title: "Meals",
            value: "1,850 kcal",
            icon: Utensils,
            color: "text-green-500",
            meals: {
                breakfast: { name: "Oatmeal & Fruits", calories: 450 },
                lunch: { name: "Grilled Chicken Salad", calories: 650 },
                dinner: { name: "Salmon with Quinoa", calories: 750 },
            },
        },
        {
            title: "Water Intake",
            value: "1.8L / 2.4L",
            icon: Droplets,
            color: "text-cyan-500",
            progress: 75,
        },
        {
            title: "Heart Rate",
            value: "72 bpm",
            icon: Heart,
            color: "text-red-500",
            chart: heartRateData,
            stats: {
                min: "58 bpm",
                max: "125 bpm",
                avg: "72 bpm",
            },
        },
        {
            title: "Tasks",
            value: "5 pending",
            icon: ListTodo,
            color: "text-yellow-500",
            tasks: [
                { name: "Review PR", priority: "high" },
                { name: "Team Meeting", priority: "medium" },
                { name: "Write Documentation", priority: "low" },
            ],
        },
    ];

    function formatMinutes(mins: number): string {
        const hours = Math.floor(mins / 60);
        const minutes = mins % 60;
        return `${hours}h ${minutes}m`;
    }

    function getPriorityColor(priority: string): string {
        switch (priority) {
            case "high":
                return "bg-red-500";
            case "medium":
                return "bg-yellow-500";
            case "low":
                return "bg-green-500";
            default:
                return "bg-gray-500";
        }
    }
</script>

{#each metrics as metric}
    <Card class="h-full">
        <CardHeader>
            <CardTitle class="flex items-center gap-2 text-lg">
                <svelte:component
                    this={metric.icon}
                    class="w-5 h-5 {metric.color}"
                />
                {metric.title}
            </CardTitle>
        </CardHeader>
        <CardContent>
            <div class="text-2xl font-bold mb-4">{metric.value}</div>

            {#if metric.progress !== undefined}
                <Progress value={metric.progress} class="mt-2" />
            {/if}

            {#if metric.chart}
                <div class="h-32 mt-4">
                    <ResponsiveContainer width="100%" height="100%">
                        <LineChart data={metric.chart}>
                            <Line
                                type="monotone"
                                dataKey="rate"
                                stroke="currentColor"
                                dot={false}
                                strokeWidth={2}
                            />
                            <XAxis
                                dataKey="hour"
                                tick={{ fontSize: 12 }}
                                interval={4}
                            />
                            <YAxis
                                tick={{ fontSize: 12 }}
                                domain={["dataMin - 5", "dataMax + 5"]}
                            />
                            <Tooltip />
                        </LineChart>
                    </ResponsiveContainer>

                    <div
                        class="flex justify-between mt-2 text-sm text-muted-foreground"
                    >
                        <span>Min: {metric.stats?.min}</span>
                        <span>Avg: {metric.stats?.avg}</span>
                        <span>Max: {metric.stats?.max}</span>
                    </div>
                </div>
            {/if}

            {#if metric.detail && !Array.isArray(metric.detail)}
                <div class="space-y-2 mt-4">
                    <div class="flex justify-between">
                        <span class="text-sm text-muted-foreground"
                            >Bedtime</span
                        >
                        <span class="text-sm">{metric.detail.bedtime}</span>
                    </div>
                    <div class="flex justify-between">
                        <span class="text-sm text-muted-foreground"
                            >Wake up</span
                        >
                        <span class="text-sm">{metric.detail.wakeup}</span>
                    </div>
                    <div class="flex justify-between">
                        <span class="text-sm text-muted-foreground"
                            >Quality</span
                        >
                        <span class="text-sm">{metric.detail.quality}</span>
                    </div>
                </div>
            {/if}

            {#if Array.isArray(metric.detail)}
                <div class="space-y-2 mt-4">
                    {#each metric.detail as app}
                        <div class="flex justify-between items-center">
                            <div class="flex items-center gap-2">
                                <div
                                    class="w-2 h-2 rounded-full {app.color}"
                                ></div>
                                <span class="text-sm">{app.name}</span>
                            </div>
                            <span class="text-sm"
                                >{formatMinutes(app.value)}</span
                            >
                        </div>
                    {/each}
                </div>
            {/if}

            {#if metric.meals}
                <div class="space-y-4 mt-4">
                    {#each Object.entries(metric.meals) as [meal, data]}
                        <div>
                            <div class="flex justify-between items-center">
                                <span
                                    class="capitalize text-sm text-muted-foreground"
                                    >{meal}</span
                                >
                                <Badge variant="outline"
                                    >{data.calories} kcal</Badge
                                >
                            </div>
                            <p class="text-sm mt-1">{data.name}</p>
                            <Separator class="mt-2" />
                        </div>
                    {/each}
                </div>
            {/if}

            {#if metric.tasks}
                <div class="space-y-2 mt-4">
                    {#each metric.tasks as task}
                        <div class="flex items-center justify-between">
                            <div class="flex items-center gap-2">
                                <div
                                    class="w-2 h-2 rounded-full {getPriorityColor(
                                        task.priority,
                                    )}"
                                ></div>
                                <span class="text-sm">{task.name}</span>
                            </div>
                            <Badge variant="outline" class="capitalize"
                                >{task.priority}</Badge
                            >
                        </div>
                    {/each}
                </div>
            {/if}
        </CardContent>
    </Card>
{/each}
