<script lang="ts">
    import { getContext } from "svelte";
    import { Card, CardContent, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Heart } from "lucide-svelte";
    import { LineChart, Spline, Tooltip, Svg } from "layerchart";
    import { PeriodType } from "svelte-ux";
    import { format } from "@layerstack/utils";

    const { theme } = getContext("theme");
    
    export let data = [
        { time: "00:00", value: 75 },
        { time: "04:00", value: 68 },
        { time: "08:00", value: 82 },
        { time: "12:00", value: 78 },
        { time: "16:00", value: 85 },
        { time: "20:00", value: 72 },
    ];

    $: processedData = data.map((d) => ({
        date: new Date(`2024-01-01 ${d.time}`),
        value: d.value,
    }));
    $: primaryColor = `hsl(var(--primary))`;
</script>

<Card class="w-full h-[180px]">
    <CardHeader class="p-3">
        <div class="flex items-center space-x-2">
            <Heart class="h-4 w-4 text-red-500" />
            <CardTitle class="text-sm font-semibold">Heart Rate</CardTitle>
        </div>
    </CardHeader>
    <CardContent class="p-3 pt-0">
        <div class="h-[120px]">
            <LineChart
                data={processedData}
                x="date"
                y="value"
                yPadding={[0, 4]}
                padding={{ top: 10, right: 10, bottom: 20, left: 30 }}
                points
                tooltip={{ mode: "voronoi" }}
                props={{
                    spline: { class: "stroke-2" },
                    xAxis: {
                        format: (d) => format(d, PeriodType.Hour),
                        tickLength: 0,
                    },
                    yAxis: {
                        format: (v) => `${v}`,
                        ticks: 3,
                    },
                    grid: { x: false, y: true },
                    highlight: { points: { r: 3 } },
                }}
            >
                <svelte:fragment slot="marks" let:series>
                    {#each series as s}
                        <Spline stroke={primaryColor} class="stroke-1.5" />
                    {/each}
                </svelte:fragment>
                <svelte:fragment slot="tooltip" let:x let:y>
                    <Tooltip.Root let:data>
                        <Tooltip.Header>
                            {format(x(data), PeriodType.Hour)}
                        </Tooltip.Header>
                        <Tooltip.List>
                            <Tooltip.Item label="BPM" value={y(data)} color={primaryColor} />
                        </Tooltip.List>
                    </Tooltip.Root>
                </svelte:fragment>
            </LineChart>
        </div>
    </CardContent>
</Card>

<style>
    :global(.tick text) {
        font-size: 10px;
        fill: hsl(var(--muted-foreground));
    }
    :global(.domain) {
        stroke: hsl(var(--muted-foreground));
    }
</style>