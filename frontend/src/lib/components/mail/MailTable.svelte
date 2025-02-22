<!-- src/lib/components/mail/MailTable.svelte -->
<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import {
        mailMessagesStore,
        type MailMessage,
    } from "$lib/stores/mailMessagesStore";
    import { mailStore } from "$lib/stores/mailStore";
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
    import { format, formatDistanceToNow } from "date-fns";

    let searchQuery = "";

    function formatDate(dateStr: string) {
        const date = new Date(dateStr);
        const now = new Date();
        const isToday = date.toDateString() === now.toDateString();
        return isToday ? format(date, "HH:mm") : format(date, "MMM d");
    }

    function handleRowClick(mail: MailMessage) {
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

    onMount(() => {
        mailMessagesStore.fetchMails();
        mailMessagesStore.subscribeToChanges();
    });

    onDestroy(() => {
        mailMessagesStore.unsubscribe();
    });

    $: lastSyncDisplay = $mailStore.syncStatus?.last_synced
        ? `Last synced ${formatDistanceToNow(new Date($mailStore.syncStatus.last_synced), { addSuffix: true })}`
        : "";
</script>

<div class="space-y-1">
    <div class="flex justify-between items-center mb-2">
        <Input
            type="search"
            placeholder="Search emails..."
            class="w-[200px]"
            bind:value={searchQuery}
        />
        <div class="flex items-center gap-3 text-sm">
            {#if lastSyncDisplay}
                <span class="text-muted-foreground">
                    {lastSyncDisplay}
                </span>
            {/if}
            <Button
                variant="outline"
                size="icon"
                class="h-8 w-8"
                on:click={handleRefresh}
                disabled={$mailMessagesStore.isLoading}
            >
                <div class={$mailMessagesStore.isLoading ? "animate-spin" : ""}>
                    <RefreshCw size={14} />
                </div>
            </Button>
        </div>
    </div>

    <div class="rounded-md border">
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead class="w-8 py-2"></TableHead>
                    <TableHead class="py-2 w-[180px]">From</TableHead>
                    <TableHead class="py-2">Subject & Preview</TableHead>
                    <TableHead class="py-2 w-[60px] text-right">Date</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {#if $mailMessagesStore.isLoading}
                    <TableRow>
                        <TableCell colspan="4" class="text-center py-2">
                            <div class="w-fit mx-auto animate-spin">
                                <RefreshCw size={16} />
                            </div>
                        </TableCell>
                    </TableRow>
                {:else if $mailMessagesStore.messages.length === 0}
                    <TableRow>
                        <TableCell colspan="4" class="text-center py-2">
                            No emails found
                        </TableCell>
                    </TableRow>
                {:else}
                    {#each $mailMessagesStore.messages as mail}
                        <TableRow
                            class={`cursor-pointer hover:bg-muted/50 ${
                                mail.is_unread ? "font-medium" : ""
                            }`}
                            on:click={() => handleRowClick(mail)}
                        >
                            <TableCell class="py-1 w-8">
                                <div class="flex">
                                    <div
                                        class={mail.is_unread
                                            ? "text-primary"
                                            : "text-muted-foreground"}
                                    >
                                        {#if mail.is_unread}
                                            <Mail size={14} />
                                        {:else}
                                            <MailOpen size={14} />
                                        {/if}
                                    </div>
                                    {#if mail.is_starred}
                                        <div class="text-yellow-400 -ml-1">
                                            <Star
                                                size={14}
                                                fill="currentColor"
                                            />
                                        </div>
                                    {/if}
                                </div>
                            </TableCell>
                            <TableCell class="py-1 text-sm">
                                <div class="truncate">
                                    {mail.from}
                                </div>
                            </TableCell>
                            <TableCell class="py-1">
                                <div class="text-sm">
                                    <span>{mail.subject}</span>
                                    <span class="text-muted-foreground">
                                        -
                                    </span>
                                    <span
                                        class="text-muted-foreground truncate"
                                    >
                                        {mail.snippet}
                                    </span>
                                </div>
                            </TableCell>
                            <TableCell class="py-1 text-right">
                                <div class="text-xs text-muted-foreground">
                                    {formatDate(mail.received_date)}
                                </div>
                            </TableCell>
                        </TableRow>
                    {/each}
                {/if}
            </TableBody>
        </Table>
    </div>

    <div class="flex items-center justify-between pt-1 text-sm">
        <div class="text-muted-foreground">
            {($mailMessagesStore.page - 1) * $mailMessagesStore.perPage + 1} -
            {Math.min(
                $mailMessagesStore.page * $mailMessagesStore.perPage,
                $mailMessagesStore.totalItems,
            )} of
            {$mailMessagesStore.totalItems} emails
        </div>
        <div class="flex items-center gap-2">
            <Button
                variant="ghost"
                size="sm"
                class="h-8 px-2"
                disabled={$mailMessagesStore.page <= 1}
                on:click={handlePrevPage}
            >
                Previous
            </Button>
            <span class="text-muted-foreground px-2">
                Page {$mailMessagesStore.page}
            </span>
            <Button
                variant="ghost"
                size="sm"
                class="h-8 px-2"
                disabled={$mailMessagesStore.page *
                    $mailMessagesStore.perPage >=
                    $mailMessagesStore.totalItems}
                on:click={handleNextPage}
            >
                Next
            </Button>
        </div>
    </div>
</div>
