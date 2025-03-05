<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { foodLogStore } from "../stores/food-log.store";
    import { format, isToday, isYesterday, parseISO } from "date-fns";
    import { Trash2, X, Calendar, Clock, Tag } from "lucide-svelte";
    import { toast } from "svelte-sonner";
    import { Button } from "$lib/components/ui/button";
    import { Badge } from "$lib/components/ui/badge";
    import { pb } from "$lib/config/pocketbase";
    import type { FoodLogEntry } from "../types";
    import * as Dialog from "$lib/components/ui/dialog";
    import type { Writable } from "svelte/store";

    // Accept fileTokenStore as a prop to use shared token
    export let fileTokenStore:
        | Writable<{ token: string; expiresAt: number } | null>
        | undefined = undefined;

    let loading = true;
    let loadingMore = false;
    let observer: IntersectionObserver;
    let loadTriggerElement: HTMLDivElement;
    let selectedEntry: FoodLogEntry | null = null;
    let detailsOpen = false;

    // Store file tokens and their expiry times (local fallback if not provided from parent)
    let fileTokens: { token: string; expiresAt: number } | null = null;
    let tokenRequestInProgress = false;
    let tokenPromise: Promise<string> | null = null;

    // Subscribe to the store
    const unsubscribe = foodLogStore.subscribe((state) => {
        loading = state.isLoading;
    });

    // Subscribe to token store if provided
    let unsubscribeToken: (() => void) | undefined;
    if (fileTokenStore) {
        unsubscribeToken = fileTokenStore.subscribe((value) => {
            fileTokens = value;
        });
    }

    // Format date for grouping
    function formatDateGroup(dateString: string): string {
        if (!dateString || dateString === "unknown") {
            return "Unknown Date";
        }

        try {
            const date = parseISO(dateString);

            if (isNaN(date.getTime())) {
                return "Invalid Date";
            }

            if (isToday(date)) {
                return "Today";
            } else if (isYesterday(date)) {
                return "Yesterday";
            } else {
                return format(date, "EEEE, MMMM d, yyyy");
            }
        } catch (error) {
            console.error(`Error formatting date group: ${dateString}`, error);
            return "Unknown Date";
        }
    }

    // Group entries by date
    function groupEntriesByDate(
        entries: FoodLogEntry[],
    ): Record<string, FoodLogEntry[]> {
        const groups: Record<string, FoodLogEntry[]> = {};

        entries.forEach((entry) => {
            try {
                // Use the custom date field instead of created for grouping
                if (!entry.date) {
                    console.warn(
                        `Entry ${entry.id} has no date field, using created date instead`,
                    );
                    const dateKey = format(
                        parseISO(entry.created),
                        "yyyy-MM-dd",
                    );
                    if (!groups[dateKey]) {
                        groups[dateKey] = [];
                    }
                    groups[dateKey].push(entry);
                    return;
                }

                // Validate the date before parsing
                const parsedDate = parseISO(entry.date);
                if (isNaN(parsedDate.getTime())) {
                    console.warn(
                        `Entry ${entry.id} has invalid date: ${entry.date}, using created date instead`,
                    );
                    const dateKey = format(
                        parseISO(entry.created),
                        "yyyy-MM-dd",
                    );
                    if (!groups[dateKey]) {
                        groups[dateKey] = [];
                    }
                    groups[dateKey].push(entry);
                    return;
                }

                const dateKey = format(parsedDate, "yyyy-MM-dd");
                if (!groups[dateKey]) {
                    groups[dateKey] = [];
                }
                groups[dateKey].push(entry);
            } catch (error) {
                console.error(
                    `Error processing date for entry ${entry.id}:`,
                    error,
                );
                // Fallback to using the created date if the date field has issues
                try {
                    const dateKey = format(
                        parseISO(entry.created),
                        "yyyy-MM-dd",
                    );
                    if (!groups[dateKey]) {
                        groups[dateKey] = [];
                    }
                    groups[dateKey].push(entry);
                } catch (fallbackError) {
                    console.error(
                        `Fallback also failed for entry ${entry.id}:`,
                        fallbackError,
                    );
                    // As a last resort, group under "Unknown Date"
                    const dateKey = "unknown";
                    if (!groups[dateKey]) {
                        groups[dateKey] = [];
                    }
                    groups[dateKey].push(entry);
                }
            }
        });

        return groups;
    }

    // Handle deleting an entry
    async function handleDelete(entryId: string) {
        if (confirm("Are you sure you want to delete this entry?")) {
            const result = await foodLogStore.deleteEntry(entryId);

            if (result.success) {
                toast.success("Entry deleted successfully");
            } else {
                toast.error("Failed to delete entry");
            }
        }
    }

    // Get file token with caching
    async function getFileToken(): Promise<string> {
        const now = Date.now();

        // If we have a valid token, return it
        if (fileTokens && fileTokens.expiresAt > now) {
            return fileTokens.token;
        }

        // If a token request is already in progress, wait for it
        if (tokenRequestInProgress && tokenPromise) {
            return tokenPromise;
        }

        // Start a new token request
        tokenRequestInProgress = true;
        tokenPromise = new Promise<string>(async (resolve, reject) => {
            try {
                console.log("Fetching new file token");
                const token = await pb.files.getToken();

                // Cache the token
                const newToken = {
                    token,
                    expiresAt: now + 110 * 1000, // 110 seconds (slightly less than 2 min)
                };

                // Store in local variable
                fileTokens = newToken;

                // If parent provided a token store, update it too
                if (fileTokenStore) {
                    fileTokenStore.set(newToken);
                }

                tokenRequestInProgress = false;
                resolve(token);
            } catch (error) {
                tokenRequestInProgress = false;
                console.error("Error fetching file token:", error);
                reject("");
            }
        });

        return tokenPromise;
    }

    // Get image URL from PocketBase with token for protected files
    async function getImageUrl(entry: FoodLogEntry): Promise<string> {
        if (!entry.image) return "";

        try {
            const token = await getFileToken();
            if (!token) {
                console.error("Failed to get a valid file token");
                return "";
            }

            // Return URL with token
            return pb.files.getUrl(entry, entry.image, { token });
        } catch (error) {
            console.error("Error getting file URL:", error);
            return "";
        }
    }

    // Open details modal for an entry
    function openDetails(entry: FoodLogEntry) {
        selectedEntry = entry;
        detailsOpen = true;
    }

    // Format time only
    function formatTime(dateString: string): string {
        try {
            const date = parseISO(dateString);
            if (isNaN(date.getTime())) {
                return "Invalid Time";
            }
            return format(date, "h:mm a");
        } catch (error) {
            console.error(`Error formatting time: ${dateString}`, error);
            return "Invalid Time";
        }
    }

    // Format full date and time
    function formatDateTime(dateString: string): string {
        try {
            const date = parseISO(dateString);
            if (isNaN(date.getTime())) {
                return "Invalid Date/Time";
            }
            return format(date, "MMMM d, yyyy 'at' h:mm a");
        } catch (error) {
            console.error(`Error formatting date/time: ${dateString}`, error);
            return "Invalid Date/Time";
        }
    }

    // Set up infinite scrolling
    function setupInfiniteScroll() {
        observer = new IntersectionObserver(
            (entries) => {
                const entry = entries[0];

                if (
                    entry.isIntersecting &&
                    $foodLogStore.hasMore &&
                    !$foodLogStore.isLoading &&
                    $foodLogStore.entries.length > 0
                ) {
                    loadingMore = true;
                    foodLogStore.loadEntries();
                }
            },
            { rootMargin: "100px" },
        );

        if (loadTriggerElement) {
            observer.observe(loadTriggerElement);
        }
    }

    onMount(async () => {
        // Load initial entries
        await foodLogStore.loadEntries(true);
        setupInfiniteScroll();
    });

    onDestroy(() => {
        // Clean up
        unsubscribe();
        if (unsubscribeToken) {
            unsubscribeToken();
        }
        if (observer && loadTriggerElement) {
            observer.unobserve(loadTriggerElement);
        }
    });

    // Derive grouped entries from the store
    $: groupedEntries = $foodLogStore.entries.length
        ? groupEntriesByDate($foodLogStore.entries)
        : {};
    $: dateKeys = Object.keys(groupedEntries).sort((a, b) => {
        // Always put unknown at the end
        if (a === "unknown") return 1;
        if (b === "unknown") return -1;

        // Sort valid dates in reverse chronological order
        try {
            return parseISO(b).getTime() - parseISO(a).getTime();
        } catch (error) {
            console.error("Error sorting date keys:", error);
            return 0;
        }
    });
    $: hasEntries = dateKeys.length > 0;
</script>

<div class="food-log-entries">
    {#if loading && !$foodLogStore.entries.length}
        <div class="flex justify-center items-center py-10">
            <div class="loader"></div>
        </div>
    {:else if hasEntries}
        <div class="space-y-8">
            {#each dateKeys as dateKey}
                <div class="date-group">
                    <h3
                        class="text-xl font-semibold mb-4 sticky top-0 bg-background/95 backdrop-blur-sm py-2 z-10"
                    >
                        <!-- Use the first entry's date field for the group header -->
                        {formatDateGroup(groupedEntries[dateKey][0].date)}
                    </h3>

                    <div
                        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
                    >
                        {#each groupedEntries[dateKey] as entry (entry.id)}
                            <div
                                class="food-entry bg-card rounded-lg overflow-hidden shadow-sm border transition-all hover:shadow-md cursor-pointer"
                                on:click={() => openDetails(entry)}
                                on:keydown={(e) =>
                                    e.key === "Enter" && openDetails(entry)}
                                tabindex="0"
                                role="button"
                                aria-label="View details of {entry.name}"
                            >
                                {#if entry.image}
                                    <div
                                        class="aspect-video overflow-hidden relative"
                                    >
                                        {#await getImageUrl(entry) then url}
                                            <img
                                                src={url}
                                                alt={entry.name}
                                                class="w-full h-full object-cover"
                                                loading="lazy"
                                            />
                                        {/await}
                                    </div>
                                {/if}

                                <div class="p-4">
                                    <div
                                        class="flex justify-between items-start mb-2"
                                    >
                                        <h4 class="text-lg font-medium">
                                            {entry.name}
                                        </h4>
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            class="h-8 w-8 text-destructive/70 hover:text-destructive -mr-2 -mt-1 z-10"
                                            on:click={(e) => {
                                                e.stopPropagation();
                                                handleDelete(entry.id);
                                            }}
                                            aria-label="Delete entry"
                                        >
                                            <Trash2 class="h-4 w-4" />
                                        </Button>
                                    </div>

                                    <div
                                        class="flex justify-between items-center"
                                    >
                                        <Badge variant="outline"
                                            >{entry.tag}</Badge
                                        >
                                        <span
                                            class="text-xs text-muted-foreground"
                                        >
                                            {formatTime(entry.date)}
                                        </span>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/each}
        </div>

        {#if $foodLogStore.hasMore}
            <div bind:this={loadTriggerElement} class="py-4 mt-4 text-center">
                {#if loadingMore}
                    <div class="loader mx-auto"></div>
                {/if}
            </div>
        {/if}
    {:else}
        <div class="py-10 text-center">
            <p class="text-muted-foreground">
                No food entries yet. Add your first meal!
            </p>
        </div>
    {/if}
</div>

<!-- Details Modal -->
<Dialog.Root bind:open={detailsOpen}>
    <Dialog.Content class="sm:max-w-[600px]">
        <Dialog.Header>
            <Dialog.Title>
                {selectedEntry?.name || "Food Details"}
            </Dialog.Title>
            <Dialog.Description>
                Detailed view of your food log entry
            </Dialog.Description>
        </Dialog.Header>

        {#if selectedEntry}
            <div class="food-details space-y-4 my-4">
                {#if selectedEntry.image}
                    <div
                        class="image-container max-h-[400px] overflow-hidden rounded-md"
                    >
                        {#await getImageUrl(selectedEntry) then url}
                            <img
                                src={url}
                                alt={selectedEntry.name}
                                class="w-full object-contain max-h-[400px]"
                            />
                        {/await}
                    </div>
                {/if}

                <div class="metadata space-y-3 mt-4">
                    <div class="flex items-center gap-2">
                        <Tag class="h-4 w-4 text-primary" />
                        <span class="text-sm font-medium">Category:</span>
                        <Badge>{selectedEntry.tag}</Badge>
                    </div>

                    <div class="flex items-center gap-2">
                        <Calendar class="h-4 w-4 text-primary" />
                        <span class="text-sm font-medium">Date:</span>
                        <span class="text-sm"
                            >{formatDateTime(selectedEntry.date)}</span
                        >
                    </div>

                    <div class="flex items-center gap-2">
                        <Clock class="h-4 w-4 text-primary" />
                        <span class="text-sm font-medium">Added:</span>
                        <span class="text-sm"
                            >{formatDateTime(selectedEntry.created)}</span
                        >
                    </div>
                </div>
            </div>
        {/if}

        <Dialog.Footer>
            <Button variant="outline" on:click={() => (detailsOpen = false)}
                >Close</Button
            >
            {#if selectedEntry}
                <Button
                    variant="destructive"
                    on:click={() => {
                        if (selectedEntry) {
                            handleDelete(selectedEntry.id);
                            detailsOpen = false;
                        }
                    }}
                >
                    Delete
                </Button>
            {/if}
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<style>
    .loader {
        width: 2rem;
        height: 2rem;
        border: 3px solid rgba(0, 0, 0, 0.1);
        border-radius: 50%;
        border-top-color: var(--primary);
        animation: spin 1s ease-in-out infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
