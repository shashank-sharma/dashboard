<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { foodLogStore } from "../stores/food-log.store";
    import { format, isToday, isYesterday, parseISO } from "date-fns";
    import { Trash2, X, Calendar, Clock, Tag, Image } from "lucide-svelte";
    import { toast } from "svelte-sonner";
    import { Button } from "$lib/components/ui/button";
    import { Badge } from "$lib/components/ui/badge";
    import { pb } from "$lib/config/pocketbase";
    import type { FoodLogEntry } from "../types";
    import * as Dialog from "$lib/components/ui/dialog";
    import {
        getFileToken,
        getAuthenticatedFileUrl,
    } from "$lib/services/file-token";

    let loading = true;
    let loadingMore = false;
    let observer: IntersectionObserver;
    let loadTriggerElement: HTMLDivElement;
    let selectedEntry: FoodLogEntry | null = null;
    let detailsOpen = false;
    let hasTriedLoadingMore = false;
    let allEntriesLoaded = false;

    // Store file tokens and their expiry times (local fallback if not provided from parent)
    let fileTokens: { token: string; expiresAt: number } | null = null;
    let tokenRequestInProgress = false;
    let tokenPromise: Promise<string> | null = null;
    let failedImages: Record<string, number> = {};
    let retryTimeouts: Record<string, number> = {};

    // Add this near the other state variables
    let confirmDeleteOpen = false;
    let entryToDelete: FoodLogEntry | null = null;

    const unsubscribe = foodLogStore.subscribe((state) => {
        loading = state.isLoading;

        if (hasTriedLoadingMore && !state.hasMore && !state.isLoading) {
            allEntriesLoaded = true;
        }

        if (!state.isLoading) {
            loadingMore = false;
        }
    });

    function formatTag(tag: string): string {
        if (!tag) return "Unknown";

        return tag
            .split("_")
            .map(
                (word) =>
                    word.charAt(0).toUpperCase() + word.slice(1).toLowerCase(),
            )
            .join(" ");
    }

    function confirmDelete(entry: FoodLogEntry) {
        entryToDelete = entry;
        confirmDeleteOpen = true;
    }

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
                return format(date, "MMMM d, yyyy");
            }
        } catch (error) {
            console.error(`Error formatting date group: ${dateString}`, error);
            return "Unknown Date";
        }
    }

    function groupEntriesByDate(
        entries: FoodLogEntry[],
    ): Record<string, FoodLogEntry[]> {
        const groups: Record<string, FoodLogEntry[]> = {};

        entries.forEach((entry) => {
            try {
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

    async function handleDelete(entryId: string) {
        const result = await foodLogStore.deleteEntry(entryId);

        if (result.success) {
            toast.success("Entry deleted successfully");
            confirmDeleteOpen = false;
            detailsOpen = false;
        } else {
            toast.error("Failed to delete entry");
        }
    }

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

    // Handle image load failure with retry mechanism
    function handleImageError(event: Event, entryId: string) {
        const img = event.target as HTMLImageElement;

        if (!failedImages[entryId]) {
            failedImages[entryId] = 0;
        }

        // Limit retries to 3 attempts
        if (failedImages[entryId] < 3) {
            failedImages[entryId]++;
            console.log(
                `Retrying image load for entry ${entryId}, attempt ${failedImages[entryId]}`,
            );

            // Clear previous timeout if exists
            if (retryTimeouts[entryId]) {
                window.clearTimeout(retryTimeouts[entryId]);
            }

            // Exponential backoff: 1s, 2s, 4s
            const delay = Math.pow(2, failedImages[entryId] - 1) * 1000;
            retryTimeouts[entryId] = window.setTimeout(() => {
                // Force reload by appending a timestamp
                const currentSrc = img.src;
                const url = new URL(currentSrc);
                url.searchParams.set("_retry", Date.now().toString());
                img.src = url.toString();
            }, delay);
        }
    }

    function openDetails(entry: FoodLogEntry) {
        selectedEntry = entry;
        detailsOpen = true;
    }

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

    function setupInfiniteScroll() {
        console.log("Setting up infinite scroll");
        observer = new IntersectionObserver(
            (entries) => {
                const entry = entries[0];

                if (
                    entry.isIntersecting &&
                    $foodLogStore.hasMore &&
                    !$foodLogStore.isLoading &&
                    $foodLogStore.entries.length > 0
                ) {
                    console.log("Loading more entries...");
                    loadingMore = true;
                    hasTriedLoadingMore = true;
                    foodLogStore.loadEntries();
                } else if (
                    entry.isIntersecting &&
                    !$foodLogStore.hasMore &&
                    !$foodLogStore.isLoading
                ) {
                    console.log("No more entries to load");
                    hasTriedLoadingMore = true;
                }
            },
            { rootMargin: "250px" },
        );

        if (loadTriggerElement) {
            console.log("Observing load trigger element");
            observer.observe(loadTriggerElement);
        }
    }

    function loadMore() {
        console.log(
            "Loading more photos manually",
            "hasMore:",
            $foodLogStore.hasMore,
            "isLoading:",
            $foodLogStore.isLoading,
        );
        if ($foodLogStore.hasMore && !$foodLogStore.isLoading) {
            loadingMore = true;
            hasTriedLoadingMore = true;
            foodLogStore.loadEntries();
        } else if (!$foodLogStore.hasMore) {
            console.log("No more photos to load - showing end message");
            hasTriedLoadingMore = true;
        }
    }

    onMount(async () => {
        console.log("Component mounted, loading initial entries");
        await foodLogStore.loadEntries(true);
        console.log("Initial entries loaded, hasMore:", $foodLogStore.hasMore);

        if (!$foodLogStore.hasMore) {
            console.log("Initial load indicates no more entries available");
            hasTriedLoadingMore = true;
        }

        setupInfiniteScroll();
    });

    onDestroy(() => {
        unsubscribe();

        Object.values(retryTimeouts).forEach((timeoutId) => {
            window.clearTimeout(timeoutId);
        });
    });

    $: groupedEntries = $foodLogStore.entries.length
        ? groupEntriesByDate($foodLogStore.entries)
        : {};
    $: dateKeys = Object.keys(groupedEntries).sort((a, b) => {
        if (a === "unknown") return 1;
        if (b === "unknown") return -1;

        try {
            return parseISO(b).getTime() - parseISO(a).getTime();
        } catch (error) {
            console.error("Error sorting date keys:", error);
            return 0;
        }
    });
    $: hasEntries = dateKeys.length > 0;
    $: hasMoreToLoad = $foodLogStore.hasMore && !loadingMore;
    $: reachedEnd = !$foodLogStore.hasMore && hasTriedLoadingMore;

    $: if (!$foodLogStore.hasMore && hasTriedLoadingMore) {
        console.log("Reactive detection: reached end of food log entries");
    }

    $: if (
        !$foodLogStore.hasMore &&
        $foodLogStore.entries.length > 0 &&
        !loading &&
        !loadingMore
    ) {
        console.log(
            "Store state indicates we may have reached the end:",
            "hasMore:",
            $foodLogStore.hasMore,
            "entries:",
            $foodLogStore.entries.length,
            "loading:",
            loading,
            "loadingMore:",
            loadingMore,
        );
    }
</script>

<div class="food-log-entries">
    {#if loading && !$foodLogStore.entries.length}
        <div class="flex justify-center items-center py-10">
            <div class="loader"></div>
        </div>
    {:else if hasEntries}
        <div class="photos-container">
            {#each dateKeys as dateKey, i}
                <div class="photos-section my-4">
                    <div class="date-marker" id={`date-${dateKey}`}>
                        <h3
                            class="date-header sm:text-sm md:text-md lg:text-lg font-semibold bg-background/95 backdrop-blur-sm py-2 md:px-3 flex items-center gap-2 rounded-lg shadow-sm"
                        >
                            <div
                                class="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0"
                            >
                                <Calendar class="h-4 w-4 text-primary" />
                            </div>
                            <span
                                >{formatDateGroup(
                                    groupedEntries[dateKey][0].date,
                                )}</span
                            >
                        </h3>
                    </div>

                    {#each groupedEntries[dateKey] as entry (entry.id)}
                        <div
                            class="photo-item mx-1 relative overflow-hidden rounded-lg shadow-sm hover:shadow-md transition-all duration-200 cursor-pointer"
                            data-date={dateKey}
                            on:click={() => openDetails(entry)}
                            on:keydown={(e) =>
                                e.key === "Enter" && openDetails(entry)}
                            tabindex="0"
                            role="button"
                            aria-label="View details of {entry.name}"
                        >
                            {#if entry.image}
                                <div class="photo-image">
                                    {#await getImageUrl(entry)}
                                        <div
                                            class="image-loading-placeholder flex items-center justify-center h-full w-full"
                                        >
                                            <div class="w-6 h-6 loader"></div>
                                        </div>
                                    {:then url}
                                        {#if url}
                                            <img
                                                src={url}
                                                alt={entry.name}
                                                class="w-full h-full object-cover"
                                                loading="lazy"
                                                on:error={(e) =>
                                                    handleImageError(
                                                        e,
                                                        entry.id,
                                                    )}
                                            />
                                        {:else}
                                            <div
                                                class="image-error-placeholder flex items-center justify-center h-full w-full bg-muted"
                                            >
                                                <Image
                                                    class="h-8 w-8 text-muted-foreground opacity-50"
                                                />
                                            </div>
                                        {/if}
                                    {:catch}
                                        <div
                                            class="image-error-placeholder flex items-center justify-center h-full w-full bg-muted"
                                        >
                                            <Image
                                                class="h-8 w-8 text-muted-foreground opacity-50"
                                            />
                                        </div>
                                    {/await}
                                </div>
                            {:else}
                                <div
                                    class="photo-image bg-muted flex items-center justify-center"
                                >
                                    <Tag
                                        class="h-8 w-8 text-muted-foreground"
                                    />
                                </div>
                            {/if}

                            <!-- Hover overlay with info -->
                            <div class="photo-overlay">
                                <div class="photo-info">
                                    <span
                                        class="hidden md:block photo-title truncate"
                                        >{entry.name}</span
                                    >
                                    <div class="photo-meta">
                                        <span
                                            class="hidden md:block photo-time flex items-center gap-1"
                                        >
                                            {formatTime(entry.date)}
                                        </span>
                                        <span class="photo-tag">
                                            {formatTag(entry.tag)}
                                        </span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    {/each}
                </div>
            {/each}
        </div>

        <div
            bind:this={loadTriggerElement}
            class="py-4 mt-4 text-center"
            id="scroll-trigger"
        >
            {#if loadingMore}
                <div class="flex flex-col items-center justify-center py-3">
                    <div class="loader mx-auto"></div>
                    <p class="text-sm text-muted-foreground mt-3">
                        Loading more photos...
                    </p>
                </div>
            {:else if reachedEnd || allEntriesLoaded || (!$foodLogStore.hasMore && hasTriedLoadingMore)}
                <div
                    class="end-of-content py-8 text-center border-t border-muted mt-2"
                >
                    <p class="text-muted-foreground">
                        No more photos to display
                    </p>
                    <p class="text-xs text-muted-foreground/70 mt-1">
                        You've reached the end of your food log
                    </p>
                </div>
            {:else if hasMoreToLoad}
                <!-- Invisible element for triggering infinite scroll -->
                <div class="h-20"></div>
                <Button variant="outline" class="mt-2" on:click={loadMore}>
                    Load more photos
                </Button>
            {:else}
                <!-- Fallback in case other conditions don't match -->
                <div class="h-20"></div>
                <Button variant="outline" class="mt-2" on:click={loadMore}>
                    Check for more photos
                </Button>
            {/if}
        </div>
    {:else}
        <div class="py-10 text-center">
            <p class="text-muted-foreground">
                No food entries yet. Add your first meal!
            </p>
        </div>
    {/if}
</div>

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
                        {#await getImageUrl(selectedEntry)}
                            <div
                                class="flex items-center justify-center h-48 bg-muted"
                            >
                                <div class="loader"></div>
                            </div>
                        {:then url}
                            {#if url}
                                <img
                                    src={url}
                                    alt={selectedEntry.name}
                                    class="w-full object-contain max-h-[400px]"
                                    on:error={(e) =>
                                        handleImageError(
                                            e,
                                            selectedEntry?.id || "unknown",
                                        )}
                                />
                            {:else}
                                <div
                                    class="flex items-center justify-center h-48 bg-muted"
                                >
                                    <Image
                                        class="h-10 w-10 text-muted-foreground opacity-50"
                                    />
                                </div>
                            {/if}
                        {:catch}
                            <div
                                class="flex items-center justify-center h-48 bg-muted"
                            >
                                <Image
                                    class="h-10 w-10 text-muted-foreground opacity-50"
                                />
                            </div>
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
                            confirmDelete(selectedEntry);
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

<Dialog.Root bind:open={confirmDeleteOpen}>
    <Dialog.Content class="sm:max-w-[400px]">
        <Dialog.Header>
            <Dialog.Title>Confirm Deletion</Dialog.Title>
            <Dialog.Description>
                {#if entryToDelete}
                    Are you sure you want to delete "{entryToDelete.name}"?
                {:else}
                    Are you sure you want to delete this entry?
                {/if}
                <p class="text-sm text-destructive mt-2">
                    This action cannot be undone.
                </p>
            </Dialog.Description>
        </Dialog.Header>

        <Dialog.Footer>
            <Button
                variant="outline"
                on:click={() => (confirmDeleteOpen = false)}>Cancel</Button
            >
            <Button
                variant="destructive"
                on:click={() => {
                    if (entryToDelete) {
                        handleDelete(entryToDelete.id);
                    }
                }}
            >
                Delete
            </Button>
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

    .image-loading-placeholder,
    .image-error-placeholder {
        background-color: hsl(var(--muted));
    }

    .photos-container {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        width: 100%;
        position: relative;
    }

    .date-marker {
        flex: 0 0 100%;
        margin-bottom: 8px;
        margin-top: 16px;
        width: 50px;
    }

    .date-marker:first-child {
        margin-top: 0;
    }

    .date-header {
        display: inline-flex;
        position: sticky;
        top: 0;
        z-index: 10;
    }

    .photo-item {
        flex: 0 0 auto;
        width: 180px;
        height: 180px;
    }

    .photo-image {
        width: 100%;
        height: 100%;
    }

    .photos-section {
        flex-wrap: wrap;
        display: flex;
    }

    .photo-info {
        flex: 1;
        overflow: hidden;
    }

    .photo-title {
        font-weight: 500;
        font-size: 14px;
        top: 50%;
    }

    .photo-meta {
        display: flex;
        gap: 6px;
        font-size: 12px;
        margin-top: 2px;
        flex-wrap: wrap;
    }

    .photo-time {
        opacity: 0.9;
    }

    .photo-tag {
        font-size: 10px;
        padding: 1px 6px;
        background: rgba(255, 255, 255, 0.2);
        right: 0;
        position: absolute;
        margin: 4px;
        border: 0;
        border-radius: 4px;
        top: 50%;
    }

    .photo-overlay {
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        background: linear-gradient(
            to top,
            rgba(0, 0, 0, 0.7) 0%,
            rgba(0, 0, 0, 0) 100%
        );
        padding: 16px 8px 8px;
        color: white;
        transition: opacity 0.2s ease;
        display: flex;
        justify-content: space-between;
        align-items: flex-end;
    }

    .end-of-content {
        transition: opacity 0.3s ease;
    }

    @media (max-width: 640px) {
        .photos-container {
            gap: 4px;
        }

        .photo-item {
            width: 100px;
            height: 100px;
        }

        .photo-tag {
            position: relative;
            margin: 0;
        }
    }

    @media (max-width: 380px) {
        .photo-item {
            width: 100px;
            height: 100px;
        }
    }
</style>
