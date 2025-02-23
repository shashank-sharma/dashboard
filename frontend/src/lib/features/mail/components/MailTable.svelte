<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { mailMessagesStore, mailStore } from "../stores";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import {
        Table,
        TableBody,
        TableCell,
        TableHead,
        TableHeader,
        TableRow,
    } from "$lib/components/ui/table";
    import { RefreshCw, Star, Mail, MailOpen } from "lucide-svelte";
    import { format } from "date-fns";

    let searchQuery = "";

    function formatDate(dateStr: string) {
        const date = new Date(dateStr);
        const now = new Date();
        const isToday = date.toDateString() === now.toDateString();
        return isToday ? format(date, "HH:mm") : format(date, "MMM d");
    }

    function handleRowClick(mail) {
        mailMessagesStore.selectMail(mail);
    }

    function handleRefresh() {
        mailMessagesStore.refreshMails();
    }

    function handlePrevPage(e: Event) {
        e.preventDefault();
        if ($mailMessagesStore.page > 1) {
            mailMessagesStore.fetchMails($mailMessagesStore.page - 1);
        }
    }

    function handleNextPage(e: Event) {
        e.preventDefault();
        if (
            $mailMessagesStore.page * $mailMessagesStore.perPage <
            $mailMessagesStore.totalItems
        ) {
            mailMessagesStore.fetchMails($mailMessagesStore.page + 1);
        }
    }

    // Check if sync is in progress
    $: isSyncing = $mailStore.syncStatus?.status === "in-progress";

    onMount(() => {
        mailMessagesStore.fetchMails();
        mailMessagesStore.subscribeToChanges();
    });

    onDestroy(() => {
        mailMessagesStore.unsubscribe();
    });
</script>

<div class="space-y-4">
    <div class="flex items-center justify-between">
        <Input
            type="search"
            placeholder="Search emails..."
            class="w-[300px]"
            bind:value={searchQuery}
        />
        <div class="flex items-center gap-4">
            {#if $mailStore.syncStatus}
                <div class="text-sm text-muted-foreground flex flex-col">
                    <span
                        >Last synced: {new Date(
                            $mailStore.syncStatus.last_synced,
                        ).toLocaleString()}</span
                    >
                    <span class="text-right"
                        >Status: {$mailStore.syncStatus.status}</span
                    >
                </div>
            {/if}
            <Button
                variant="outline"
                size="icon"
                on:click={handleRefresh}
                disabled={isSyncing}
                class={isSyncing ? "opacity-50" : ""}
            >
                <RefreshCw class="h-4 w-4 {isSyncing ? 'animate-spin' : ''}" />
            </Button>
        </div>
    </div>

    <div class="border rounded-md">
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead class="w-12" />
                    <TableHead>From</TableHead>
                    <TableHead>Subject</TableHead>
                    <TableHead class="text-right">Date</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {#if $mailMessagesStore.isLoading}
                    <TableRow>
                        <TableCell colspan="4" class="text-center py-8">
                            <RefreshCw class="h-5 w-5 animate-spin mx-auto" />
                        </TableCell>
                    </TableRow>
                {:else if $mailMessagesStore.messages.length === 0}
                    <TableRow>
                        <TableCell colspan="4" class="text-center py-8">
                            No emails found
                        </TableCell>
                    </TableRow>
                {:else}
                    {#each $mailMessagesStore.messages as mail}
                        <TableRow
                            class="cursor-pointer"
                            on:click={() => handleRowClick(mail)}
                        >
                            <TableCell>
                                <div class="flex items-center gap-2">
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        class={mail.is_starred
                                            ? "text-yellow-400"
                                            : ""}
                                    >
                                        <Star
                                            size={16}
                                            fill={mail.is_starred
                                                ? "currentColor"
                                                : "none"}
                                        />
                                    </Button>
                                </div>
                            </TableCell>
                            <TableCell class="font-medium">
                                <div class="flex items-center gap-2">
                                    {#if mail.is_unread}
                                        <Mail class="h-4 w-4" />
                                    {:else}
                                        <MailOpen
                                            class="h-4 w-4 text-muted-foreground"
                                        />
                                    {/if}
                                    <span
                                        class={mail.is_unread
                                            ? "font-semibold"
                                            : ""}
                                    >
                                        {mail.from}
                                    </span>
                                </div>
                            </TableCell>
                            <TableCell
                                class={mail.is_unread ? "font-semibold" : ""}
                            >
                                {mail.subject}
                            </TableCell>
                            <TableCell class="text-right">
                                {formatDate(mail.received_date)}
                            </TableCell>
                        </TableRow>
                    {/each}
                {/if}
            </TableBody>
        </Table>
    </div>

    {#if $mailMessagesStore.totalItems > $mailMessagesStore.perPage}
        <div class="flex justify-end gap-2">
            <Button
                variant="outline"
                size="sm"
                on:click={handlePrevPage}
                disabled={$mailMessagesStore.page === 1}
            >
                Previous
            </Button>
            <Button
                variant="outline"
                size="sm"
                on:click={handleNextPage}
                disabled={$mailMessagesStore.page *
                    $mailMessagesStore.perPage >=
                    $mailMessagesStore.totalItems}
            >
                Next
            </Button>
        </div>
    {/if}
</div>
