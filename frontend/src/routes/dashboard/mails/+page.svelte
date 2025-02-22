<!-- src/routes/mail/+page.svelte -->
<script lang="ts">
    import { onMount } from "svelte";
    import { mailStore } from "$lib/stores/mailStore";
    import { mailMessagesStore } from "$lib/stores/mailMessagesStore";
    import MailAuth from "$lib/components/mail/MailAuth.svelte";
    import MailTable from "$lib/components/mail/MailTable.svelte";
    import MailDetail from "$lib/components/mail/MailDetail.svelte";

    let isInitializing = true;

    onMount(async () => {
        try {
            await mailStore.checkStatus();
        } finally {
            isInitializing = false;
        }
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
