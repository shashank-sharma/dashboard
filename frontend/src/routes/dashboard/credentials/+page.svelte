<script lang="ts">
    import { onMount } from "svelte";
    import { credentialsStore } from "$lib/features/credentials/stores/credentials.store";
    import { DASHBOARD_SECTIONS } from "$lib/features/dashboard/constants";
    import { goto } from "$app/navigation";
    import {
        RefreshCcw,
        KeySquare,
        Terminal,
        Key,
        KeyRound,
    } from "lucide-svelte";
    import * as Card from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";

    onMount(() => {
        credentialsStore.loadStats();
    });

    function navigateToSection(path: string) {
        goto(path);
    }
</script>

<div class="p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold">Credentials Overview</h2>
        <Button variant="outline" on:click={() => credentialsStore.loadStats()}>
            <RefreshCcw class="w-4 h-4 mr-2" />
            Refresh
        </Button>
    </div>

    {#if $credentialsStore.isLoading}
        <div class="flex justify-center items-center h-64">
            <RefreshCcw class="w-8 h-8 animate-spin" />
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <!-- API Tokens Card -->
            <Card.Root
                class="cursor-pointer hover:bg-accent/50 transition-colors"
                on:click={() =>
                    navigateToSection("/dashboard/credentials/tokens")}
            >
                <Card.Header
                    class="flex flex-row items-center justify-between space-y-0 pb-2"
                >
                    <Card.Title class="text-xl font-medium"
                        >API Tokens</Card.Title
                    >
                    <KeySquare class="h-5 w-5 text-muted-foreground" />
                </Card.Header>
                <Card.Content>
                    <div class="text-3xl font-bold">
                        {$credentialsStore.stats.totalTokens}
                    </div>
                    <p class="text-xs text-muted-foreground">
                        OAuth and API tokens for external services
                    </p>
                </Card.Content>
            </Card.Root>

            <!-- Developer Tokens Card -->
            <Card.Root
                class="cursor-pointer hover:bg-accent/50 transition-colors"
                on:click={() =>
                    navigateToSection("/dashboard/credentials/developer")}
            >
                <Card.Header
                    class="flex flex-row items-center justify-between space-y-0 pb-2"
                >
                    <Card.Title class="text-xl font-medium"
                        >Developer Tokens</Card.Title
                    >
                    <Terminal class="h-5 w-5 text-muted-foreground" />
                </Card.Header>
                <Card.Content>
                    <div class="text-3xl font-bold">
                        {$credentialsStore.stats.totalDeveloperTokens}
                    </div>
                    <p class="text-xs text-muted-foreground">
                        Personal access tokens for development
                    </p>
                </Card.Content>
            </Card.Root>

            <!-- API Keys Card -->
            <Card.Root
                class="cursor-pointer hover:bg-accent/50 transition-colors"
                on:click={() =>
                    navigateToSection("/dashboard/credentials/api-keys")}
            >
                <Card.Header
                    class="flex flex-row items-center justify-between space-y-0 pb-2"
                >
                    <Card.Title class="text-xl font-medium">API Keys</Card.Title
                    >
                    <Key class="h-5 w-5 text-muted-foreground" />
                </Card.Header>
                <Card.Content>
                    <div class="text-3xl font-bold">
                        {$credentialsStore.stats.totalApiKeys}
                    </div>
                    <p class="text-xs text-muted-foreground">
                        Service-specific API keys and secrets
                    </p>
                </Card.Content>
            </Card.Root>

            <!-- Security Keys Card -->
            <Card.Root
                class="cursor-pointer hover:bg-accent/50 transition-colors"
                on:click={() =>
                    navigateToSection("/dashboard/credentials/security-keys")}
            >
                <Card.Header
                    class="flex flex-row items-center justify-between space-y-0 pb-2"
                >
                    <Card.Title class="text-xl font-medium"
                        >Security Keys</Card.Title
                    >
                    <KeyRound class="h-5 w-5 text-muted-foreground" />
                </Card.Header>
                <Card.Content>
                    <div class="text-3xl font-bold">
                        {$credentialsStore.stats.totalSecurityKeys}
                    </div>
                    <p class="text-xs text-muted-foreground">
                        SSH key pairs for secure server connections
                    </p>
                </Card.Content>
            </Card.Root>
        </div>
    {/if}
</div>
