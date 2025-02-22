<!-- src/routes/+error.svelte -->
<script lang="ts">
    import { page } from "$app/stores";
    import { Button } from "$lib/components/ui/button";
    import { FileQuestion, Home, AlertCircle, ArrowLeft } from "lucide-svelte";
    import { fade } from "svelte/transition";
    import { goto } from "$app/navigation";

    // Get error details from the page store
    $: status = $page.status;
    $: message = $page.error?.message || "This page could not be found";

    function goBack() {
        window.history.back();
    }

    function goHome() {
        goto("/");
    }
</script>

<div
    class="min-h-screen flex items-center justify-center p-4"
    in:fade={{ duration: 300 }}
>
    <div class="text-center space-y-8 max-w-2xl mx-auto">
        <div class="flex justify-center space-x-4 items-center">
            <div class="text-primary">
                {#if status === 404}
                    <FileQuestion class="h-24 w-24" />
                {:else}
                    <AlertCircle class="h-24 w-24" />
                {/if}
            </div>
            <span class="text-7xl font-bold text-primary">{status}</span>
        </div>

        <div class="space-y-4">
            <h1 class="text-4xl font-bold tracking-tight">
                {#if status === 404}
                    Page Not Found
                {:else}
                    Something Went Wrong
                {/if}
            </h1>
            <p class="text-muted-foreground text-lg">
                {message}
            </p>
            {#if status === 404}
                <p class="text-muted-foreground">
                    The page you're looking for doesn't exist. Please check the
                    URL or try navigating from the home page.
                </p>
            {/if}
        </div>

        <div class="flex justify-center gap-4">
            <Button variant="outline" on:click={goBack} class="gap-2">
                <ArrowLeft class="h-4 w-4" />
                Go Back
            </Button>
            <Button on:click={goHome} class="gap-2">
                <Home class="h-4 w-4" />
                Return Home
            </Button>
        </div>
    </div>
</div>
