<script lang="ts">
    import {
        Card,
        CardContent,
        CardHeader,
        CardTitle,
    } from "$lib/components/ui/card";
</script>

<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
    <Card>
        <CardHeader
            class="flex flex-row items-center justify-between space-y-0 pb-2"
        >
            <CardTitle class="text-sm font-medium">Total Users</CardTitle>
            <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                class="h-4 w-4 text-muted-foreground"
            >
                <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                <circle cx="9" cy="7" r="4" />
                <path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
        </CardHeader>
        <CardContent>
            <div class="text-2xl font-bold">10,245</div>
            <p class="text-xs text-muted-foreground">+10% from last month</p>
        </CardContent>
    </Card>
</div>

<Card class="mt-6">
    <CardHeader>
        <CardTitle>Monthly Overview</CardTitle>
    </CardHeader>
    <CardContent>
        <p>Loading chart...</p>
    </CardContent>
</Card>
