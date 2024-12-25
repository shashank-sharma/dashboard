<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "$lib/pocketbase";
    import {
        Card,
        CardContent,
        CardHeader,
        CardTitle,
    } from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";

    let servers = [
        {
            name: "Raspberry Pi 4",
            status: "Online",
            ip: "192.168.1.1",
        },
    ];

    onMount(async () => {
        await fetchServers();
    });

    async function fetchServers() {
        try {
            servers = await pb.collection("servers").getFullList({
                sort: "-created",
            });
        } catch (error) {
            console.error("Error fetching servers:", error);
        }
    }
</script>

<div class="space-y-4">
    <h2 class="text-2xl font-bold">Servers</h2>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {#each servers as server}
            <Card>
                <CardHeader>
                    <CardTitle>{server.name}</CardTitle>
                </CardHeader>
                <CardContent>
                    <p>Status: {server.status}</p>
                    <p>IP: {server.ip}</p>
                </CardContent>
            </Card>
        {/each}
    </div>
</div>
