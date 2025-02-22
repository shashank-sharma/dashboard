<!-- src/routes/auth/callback/+page.svelte -->
<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { Loader2 } from "lucide-svelte";
    import { mailStore } from "$lib/stores/mailStore";

    let error: string | null = null;

    async function handleCallback() {
        const code = $page.url.searchParams.get("code");

        if (!code) {
            error = "No authentication code received";
            return;
        }

        try {
            const success = await mailStore.completeAuth(code);
            if (success) {
                window.close();
            } else {
                throw new Error("Authentication failed");
            }
        } catch (e) {
            error = "Failed to complete authentication";
        }
    }

    onMount(() => {
        handleCallback();
    });
</script>

<div class="flex min-h-screen items-center justify-center">
    {#if error}
        <div class="text-destructive text-center">
            <p>{error}</p>
            <p class="text-sm text-muted-foreground mt-2">
                You can close this window
            </p>
        </div>
    {:else}
        <div class="text-center space-y-4">
            <Loader2 class="h-8 w-8 animate-spin mx-auto" />
            <p class="text-muted-foreground">Completing authentication...</p>
        </div>
    {/if}
</div>
