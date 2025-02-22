<!-- src/lib/components/mail/MailDetail.svelte -->
<script lang="ts">
    import { mailMessagesStore } from "$lib/stores/mailMessagesStore";
    import { pb } from "$lib/pocketbase";
    import { Button } from "$lib/components/ui/button";
    import { Star, Trash2, Archive } from "lucide-svelte";
    import {
        Sheet,
        SheetContent,
        SheetHeader,
        SheetTitle,
        SheetDescription,
    } from "$lib/components/ui/sheet";
    import { ScrollArea } from "$lib/components/ui/scroll-area";
    import { toast } from "svelte-sonner";
    import { format } from "date-fns";

    async function toggleStar() {
        if ($mailMessagesStore.selectedMail) {
            try {
                await pb
                    .collection("mail_messages")
                    .update($mailMessagesStore.selectedMail.id, {
                        is_starred: !$mailMessagesStore.selectedMail.is_starred,
                    });
            } catch (error) {
                toast.error("Failed to update star status");
            }
        }
    }

    async function moveToTrash() {
        if ($mailMessagesStore.selectedMail) {
            try {
                await pb
                    .collection("mail_messages")
                    .update($mailMessagesStore.selectedMail.id, {
                        is_trash: true,
                        is_inbox: false,
                    });
                mailMessagesStore.selectMail(null);
                toast.success("Message moved to trash");
            } catch (error) {
                toast.error("Failed to move message to trash");
            }
        }
    }

    async function archiveEmail() {
        if ($mailMessagesStore.selectedMail) {
            try {
                await pb
                    .collection("mail_messages")
                    .update($mailMessagesStore.selectedMail.id, {
                        is_inbox: false,
                    });
                mailMessagesStore.selectMail(null);
                toast.success("Message archived");
            } catch (error) {
                toast.error("Failed to archive message");
            }
        }
    }

    function formatDateTime(dateStr: string) {
        return format(new Date(dateStr), "PPpp");
    }

    $: open = !!$mailMessagesStore.selectedMail;
    $: onOpenChange = (value: boolean) => {
        if (!value) {
            mailMessagesStore.selectMail(null);
        }
    };
</script>

<Sheet {open} {onOpenChange}>
    <SheetContent class="w-[90%] sm:w-[600px] sm:max-w-none">
        {#if $mailMessagesStore.selectedMail}
            <SheetHeader
                class="flex-row items-center justify-between space-y-0 pb-2 border-b"
            >
                <SheetTitle class="flex items-center gap-2">
                    <Button
                        variant="ghost"
                        size="icon"
                        on:click={toggleStar}
                        class={$mailMessagesStore.selectedMail.is_starred
                            ? "text-yellow-400"
                            : ""}
                    >
                        <Star
                            size={20}
                            fill={$mailMessagesStore.selectedMail.is_starred
                                ? "currentColor"
                                : "none"}
                        />
                    </Button>
                    <Button variant="ghost" size="icon" on:click={archiveEmail}>
                        <Archive size={20} />
                    </Button>
                    <Button variant="ghost" size="icon" on:click={moveToTrash}>
                        <Trash2 size={20} />
                    </Button>
                </SheetTitle>
            </SheetHeader>

            <ScrollArea class="h-[calc(100vh-8rem)] mt-6">
                <div class="space-y-6">
                    <div>
                        <h2 class="text-2xl font-bold">
                            {$mailMessagesStore.selectedMail.subject}
                        </h2>
                        <div class="mt-4 space-y-1">
                            <div class="flex items-center justify-between">
                                <div>
                                    <span class="font-medium">From: </span>
                                    <span class="text-muted-foreground"
                                        >{$mailMessagesStore.selectedMail
                                            .from}</span
                                    >
                                </div>
                                <span class="text-sm text-muted-foreground">
                                    {formatDateTime(
                                        $mailMessagesStore.selectedMail
                                            .received_date,
                                    )}
                                </span>
                            </div>
                            <div>
                                <span class="font-medium">To: </span>
                                <span class="text-muted-foreground"
                                    >{$mailMessagesStore.selectedMail.to}</span
                                >
                            </div>
                        </div>
                    </div>

                    <div class="border-t pt-6">
                        <div
                            class="prose prose-sm dark:prose-invert max-w-none"
                        >
                            {@html $mailMessagesStore.selectedMail.body}
                        </div>
                    </div>
                </div>
            </ScrollArea>
        {/if}
    </SheetContent>
</Sheet>
