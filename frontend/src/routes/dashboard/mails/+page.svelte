<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { MailAuth, MailTable, MailDetail } from "$lib/features/mail";
    import { mailStore, mailMessagesStore } from "$lib/features/mail";

    let isInitializing = true;

    onMount(async () => {
        try {
            await mailStore.checkStatus(true); // Force initial check
            mailStore.subscribeToChanges();
        } catch (error) {
            console.error("Failed to initialize mail:", error);
        } finally {
            isInitializing = false;
        }
    });

    onDestroy(() => {
        mailStore.unsubscribe();
    });
</script>

<div class="container py-8 space-y-8">
    {#if isInitializing}
        <div class="flex items-center justify-center min-h-[200px]">
            <div class="animate-pulse">Loading...</div>
        </div>
    {:else if !$mailStore.isAuthenticated}
        <MailAuth />
    {:else}
        <div class="space-y-8">
            <h1 class="text-3xl font-bold">Mail</h1>
            <div class="relative">
                <div class:pr-[40%]={$mailMessagesStore.selectedMail}>
                    <MailTable />
                </div>
                <MailDetail />
            </div>
        </div>
    {/if}
</div>
