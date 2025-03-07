<script lang="ts">
    import { onMount } from "svelte";
    import FoodLogUpload from "$lib/features/food-log/components/FoodLogUpload.svelte";
    import FoodLogEntries from "$lib/features/food-log/components/FoodLogEntries.svelte";
    import { foodLogStore } from "$lib/features/food-log/stores/food-log.store";
    import {
        FOOD_TAGS,
        DEFAULT_FOOD_LOG_FORM,
        type FoodLogFormData,
        type FoodLogEntry,
    } from "$lib/features/food-log/types";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Badge } from "$lib/components/ui/badge";
    import {
        Calendar,
        Plus,
        UtensilsCrossed,
        Search,
        Upload,
        X,
        Utensils,
    } from "lucide-svelte";
    import { format, parseISO } from "date-fns";
    import * as Card from "$lib/components/ui/card";
    import * as Dialog from "$lib/components/ui/dialog";
    import { toast } from "svelte-sonner";
    import { pb } from "$lib/config/pocketbase";
    import { writable, get, type Writable } from "svelte/store";
    import { isMobile } from "$lib/stores/device";
    import {
        getFileToken,
        getAuthenticatedFileUrl,
    } from "$lib/services/file-token";

    let searchTerm = "";
    let selectedDate: string | null = null;
    let selectedTag: string | null = null;
    let showFilters = false;
    let uploadModalOpen = false;
    let currentMealType = "";

    // Meal sections data
    const mealSections = [
        { id: "breakfast", label: "Breakfast", icon: "☕" },
        { id: "lunch", label: "Lunch", icon: "🍲" },
        { id: "dinner", label: "Dinner", icon: "🍽️" },
    ];

    // Store today's entries by meal type
    let todaysMeals: Record<string, FoodLogEntry | null> = {
        breakfast: null,
        lunch: null,
        dinner: null,
    };

    // Cache for meal image URLs
    let mealImages: Record<string, string> = {};

    // Track token request
    let tokenRequestInProgress = false;
    let tokenPromise: Promise<string> | null = null;

    // Apply search filter with debounce
    let searchTimeout: NodeJS.Timeout;
    function handleSearch() {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            applyFilters();
        }, 300);
    }

    // Apply all filters
    function applyFilters() {
        foodLogStore.setFilter({
            searchTerm: searchTerm || undefined,
            date: selectedDate || undefined,
            tag: selectedTag || undefined,
        });
    }

    // Clear all filters
    function clearFilters() {
        searchTerm = "";
        selectedDate = null;
        selectedTag = null;
        foodLogStore.setFilter({});
    }

    // Format date for display
    function formatDate(dateString: string | null): string {
        if (!dateString) return "";
        return format(new Date(dateString), "MMMM d, yyyy");
    }

    // Format time for display
    function formatTime(dateString: string): string {
        try {
            const date = parseISO(dateString);
            if (isNaN(date.getTime())) {
                return "Invalid Time";
            }
            return format(date, "h:mm a");
        } catch (error) {
            return "Invalid Time";
        }
    }

    // Open the upload modal for a specific meal type
    function openUploadModal(mealType: string) {
        currentMealType = mealType;
        uploadModalOpen = true;
    }

    // Improve image loading with more verbose debugging
    async function loadMealImages(entriesToLoad: any[]) {
        console.log(`Loading images for ${entriesToLoad.length} entries...`);

        try {
            for (const entry of entriesToLoad) {
                if (entry && entry.id && entry.image) {
                    try {
                        // Get the base file URL first
                        const baseFileUrl = pb.files.getUrl(entry, entry.image);

                        // Then get the authenticated URL
                        const authenticatedUrl =
                            await getAuthenticatedFileUrl(baseFileUrl);
                        console.log(
                            `Generated URL for ${entry.id} (${entry.tag}):`,
                            authenticatedUrl.substring(0, 50) + "...",
                        );
                        mealImages[entry.id] = authenticatedUrl;
                    } catch (err) {
                        console.error(
                            `Error generating URL for entry ${entry.id}:`,
                            err,
                        );
                    }
                }
            }
        } catch (error) {
            console.error("Error loading meal images:", error);
        }
    }

    // Get image URL for an entry (returns cached URL directly)
    function getImageUrlForEntry(entry: FoodLogEntry | null): string {
        if (!entry) {
            console.log("No entry provided to getImageUrlForEntry");
            return "";
        }
        if (!entry.id) {
            console.log("Entry has no id:", entry);
            return "";
        }
        if (!entry.image) {
            console.log(`Entry ${entry.id} has no image`);
            return "";
        }

        const url = mealImages[entry.id] || "";
        if (!url) {
            console.log(
                `No cached URL found for entry ${entry.id}, meal: ${entry.tag}`,
            );
        }
        return url;
    }

    // Format time for display with safeguards
    function formatTimeForEntry(entry: FoodLogEntry | null): string {
        if (!entry) {
            console.log("No entry provided to formatTimeForEntry");
            return "";
        }
        if (!entry.date) {
            console.log(`Entry ${entry.id} has no date`);
            return "";
        }
        try {
            return formatTime(entry.date);
        } catch (error) {
            console.error(
                `Error formatting time for entry ${entry.id}:`,
                error,
            );
            return "";
        }
    }

    // Find today's meals from entries
    async function findTodaysMeals(entries: FoodLogEntry[]) {
        const today = new Date().toISOString().split("T")[0];
        console.log("Today's date for comparison:", today);

        // Reset meals
        todaysMeals = {
            breakfast: null,
            lunch: null,
            dinner: null,
        };

        // Find the most recent entry for each meal type for today
        entries.forEach((entry) => {
            if (!entry.date) return;

            // Ensure consistent date format parsing with "T" separator
            const entryDate = entry.date.split(" ")[0];
            console.log(
                `Checking entry date: ${entryDate} vs today: ${today}`,
                entry.tag,
                entry.name,
            );

            if (entryDate === today) {
                const mealType = entry.tag;
                console.log(
                    `Found today's ${mealType} entry:`,
                    entry.id,
                    entry.name,
                );

                if (
                    mealType === "breakfast" &&
                    (!todaysMeals.breakfast ||
                        new Date(entry.date) >
                            new Date(todaysMeals.breakfast.date))
                ) {
                    todaysMeals.breakfast = entry;
                }

                if (
                    mealType === "lunch" &&
                    (!todaysMeals.lunch ||
                        new Date(entry.date) > new Date(todaysMeals.lunch.date))
                ) {
                    todaysMeals.lunch = entry;
                }

                if (
                    mealType === "dinner" &&
                    (!todaysMeals.dinner ||
                        new Date(entry.date) >
                            new Date(todaysMeals.dinner.date))
                ) {
                    todaysMeals.dinner = entry;
                }
            }
        });

        console.log(
            "Today's meals found:",
            Object.keys(todaysMeals).filter((key) => todaysMeals[key] !== null),
        );

        // Load images for today's meals
        loadMealImages(
            Object.values(todaysMeals).filter((entry) => entry !== null),
        );
    }

    // Subscribe to entries to find today's meals
    const unsubscribe = foodLogStore.subscribe(async (state) => {
        await findTodaysMeals(state.entries);
    });

    // Handle successful upload
    function handleUploadSuccess() {
        uploadModalOpen = false;
        // Clear image cache after new upload to ensure fresh images
        mealImages = {};
    }

    onMount(() => {
        // Clean up on component unmount
        return () => {
            clearTimeout(searchTimeout);
            unsubscribe();
        };
    });
</script>

<div
    class="food-log-page container mx-auto py-4 sm:py-6 px-3 sm:px-4 max-w-7xl"
>
    <div class="page-header mb-8">
        <h1 class="text-3xl font-bold">Food Log</h1>
        <p class="text-muted-foreground mt-2">Snapshot for food</p>
    </div>

    <!-- Meal Tracking Sections -->
    <div class="meal-sections mb-10">
        <h2 class="text-xl font-medium mb-4">Today's Meals</h2>

        <!-- Desktop view -->
        <div class="hidden md:grid md:grid-cols-3 gap-4">
            {#each mealSections as meal}
                <Card.Root class="overflow-hidden meal-card">
                    <Card.Header
                        class={todaysMeals[meal.id]
                            ? "bg-green-50 dark:bg-green-900/20"
                            : "bg-orange-50 dark:bg-orange-900/20"}
                    >
                        <div class="flex justify-between items-center">
                            <div class="flex items-center gap-2">
                                <span class="text-2xl" aria-hidden="true"
                                    >{meal.icon}</span
                                >
                                <Card.Title>{meal.label}</Card.Title>
                            </div>
                            <Badge
                                variant={todaysMeals[meal.id]
                                    ? "default"
                                    : "outline"}
                                class={todaysMeals[meal.id]
                                    ? "bg-green-500 badge"
                                    : "badge"}
                            >
                                {todaysMeals[meal.id] ? "Logged" : "Not Logged"}
                            </Badge>
                        </div>
                    </Card.Header>

                    {#if todaysMeals[meal.id]}
                        <!-- Show the meal that's already logged -->
                        <div class="relative">
                            <!-- Display image if it has one -->
                            {#if todaysMeals[meal.id]?.image}
                                <div class="aspect-video overflow-hidden">
                                    {#if getImageUrlForEntry(todaysMeals[meal.id])}
                                        <img
                                            src={getImageUrlForEntry(
                                                todaysMeals[meal.id],
                                            )}
                                            alt={todaysMeals[meal.id]?.name ||
                                                meal.label}
                                            class="w-full h-full object-cover"
                                            on:error={() => {
                                                console.error(
                                                    `Image failed to load for ${meal.id}`,
                                                );
                                                // Clear the cache entry to force a reload - safely handle null
                                                const mealEntry =
                                                    todaysMeals[meal.id];
                                                if (mealEntry && mealEntry.id) {
                                                    delete mealImages[
                                                        mealEntry.id
                                                    ];
                                                    // Try loading again after a delay
                                                    setTimeout(
                                                        () =>
                                                            loadMealImages([
                                                                mealEntry,
                                                            ]),
                                                        1000,
                                                    );
                                                }
                                            }}
                                        />
                                    {:else}
                                        <div
                                            class="w-full h-full flex items-center justify-center bg-muted"
                                        >
                                            <div class="text-center p-4">
                                                <div
                                                    class="loader mx-auto mb-2"
                                                ></div>
                                                <p
                                                    class="text-xs text-muted-foreground"
                                                >
                                                    Loading image...
                                                </p>
                                            </div>
                                        </div>
                                    {/if}
                                </div>
                            {/if}
                            <Card.Content class="p-4 card-content">
                                <h4 class="font-medium mb-1 line-clamp-1">
                                    {todaysMeals[meal.id]?.name || ""}
                                </h4>
                                <div
                                    class="flex justify-between items-center text-sm text-muted-foreground"
                                >
                                    <span>{meal.label}</span>
                                    <span
                                        >{formatTimeForEntry(
                                            todaysMeals[meal.id],
                                        )}</span
                                    >
                                </div>
                                <Button
                                    variant="outline"
                                    class="w-full mt-4"
                                    on:click={() => openUploadModal(meal.id)}
                                >
                                    <Plus class="h-4 w-4 mr-2" />
                                    <span class="whitespace-nowrap"
                                        >Add Another</span
                                    >
                                </Button>
                            </Card.Content>
                        </div>
                    {:else}
                        <!-- Show upload prompt -->
                        <Card.Content class="p-6 card-content">
                            <div
                                class="upload-area flex flex-col items-center justify-center p-6 border-2 border-dashed rounded-lg cursor-pointer hover:bg-muted/50 transition-colors"
                                on:click={() => openUploadModal(meal.id)}
                                on:keydown={(e) =>
                                    e.key === "Enter" &&
                                    openUploadModal(meal.id)}
                                tabindex="0"
                                role="button"
                                aria-label={`Upload your ${meal.label}`}
                            >
                                <Upload
                                    class="h-10 w-10 text-muted-foreground mb-2"
                                    aria-hidden="true"
                                />
                                <p class="text-sm font-medium mb-1">
                                    Upload your {meal.label}
                                </p>
                                <p class="text-xs text-muted-foreground">
                                    Click to add a photo
                                </p>
                            </div>
                        </Card.Content>
                    {/if}
                </Card.Root>
            {/each}
        </div>

        <!-- Mobile compact view -->
        <div class="md:hidden">
            <!-- Replace list with horizontal row of meal cards -->
            <div class="grid grid-cols-3 gap-2">
                {#each mealSections as meal}
                    <div
                        class="flex flex-col items-center p-2 bg-card border rounded-lg shadow-sm cursor-pointer {todaysMeals[
                            meal.id
                        ]
                            ? 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-900'
                            : 'bg-orange-50 dark:bg-orange-900/20 border-orange-200 dark:border-orange-900'}"
                        on:click={() => openUploadModal(meal.id)}
                        on:keydown={(e) =>
                            e.key === "Enter" && openUploadModal(meal.id)}
                        tabindex="0"
                        role="button"
                    >
                        <div class="text-2xl mb-1">{meal.icon}</div>
                        <div class="text-xs font-medium">{meal.label}</div>
                        <Badge
                            variant={todaysMeals[meal.id]
                                ? "default"
                                : "outline"}
                            class="mt-1 text-[9px] px-1 py-0 h-4 {todaysMeals[
                                meal.id
                            ]
                                ? 'bg-green-500'
                                : ''}"
                        >
                            {todaysMeals[meal.id] ? "✓" : "+"}
                        </Badge>
                    </div>
                {/each}
            </div>
        </div>

        <!-- Extra meal button (desktop only) -->
        <div class="mt-4 hidden md:block">
            <Button
                variant="outline"
                class="w-full py-4 sm:py-6 flex items-center justify-center gap-2 border-dashed"
                on:click={() => openUploadModal("snack")}
            >
                <Plus class="h-5 w-5" />
                <span>Add Extra Meal (Snack, Dessert, or Other)</span>
            </Button>
        </div>
    </div>

    <!-- Floating action button for adding extra meals (mobile only) - now square -->

    <div class="hidden md:block fixed bottom-[80px] right-4 z-40">
        <Button
            variant="default"
            size="icon"
            class="h-14 w-14 rounded-md shadow-lg"
            on:click={() => openUploadModal("snack")}
        >
            <Plus class="h-6 w-6" />
        </Button>
    </div>

    <!-- Filters Section -->
    <div class="filters-section mb-8">
        <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-medium">Filter Your Logs</h2>
            <Button
                variant="ghost"
                size="sm"
                on:click={() => (showFilters = !showFilters)}
            >
                {showFilters ? "Hide Filters" : "Show Filters"}
            </Button>
        </div>

        {#if showFilters}
            <div
                class="filter-options grid grid-cols-1 md:grid-cols-3 gap-4 bg-card border rounded-lg p-4"
            >
                <div>
                    <label class="text-sm font-medium mb-2 block">Search</label>
                    <div class="relative">
                        <Search
                            class="absolute left-2 top-3 h-4 w-4 text-muted-foreground"
                        />
                        <Input
                            type="text"
                            placeholder="Search food items..."
                            class="pl-8"
                            bind:value={searchTerm}
                            on:input={handleSearch}
                        />
                    </div>
                </div>

                <div>
                    <label class="text-sm font-medium mb-2 block">Date</label>
                    <div class="flex items-center gap-2">
                        <Input
                            type="date"
                            id="date-filter"
                            bind:value={selectedDate}
                            on:change={applyFilters}
                        />
                        {#if selectedDate}
                            <Button
                                variant="ghost"
                                size="icon"
                                on:click={() => {
                                    selectedDate = null;
                                    applyFilters();
                                }}
                                aria-label="Clear date filter"
                            >
                                <span class="sr-only">Clear</span>
                                <X class="h-4 w-4" />
                            </Button>
                        {/if}
                    </div>
                </div>

                <div>
                    <label class="text-sm font-medium mb-2 block"
                        >Category</label
                    >
                    <div class="food-tags flex flex-wrap gap-2">
                        {#each FOOD_TAGS as tag}
                            <Button
                                variant={selectedTag === tag.value
                                    ? "default"
                                    : "outline"}
                                size="sm"
                                on:click={() => {
                                    selectedTag =
                                        selectedTag === tag.value
                                            ? null
                                            : tag.value;
                                    applyFilters();
                                }}
                            >
                                {tag.label}
                            </Button>
                        {/each}
                    </div>
                </div>

                {#if searchTerm || selectedDate || selectedTag}
                    <div class="md:col-span-3">
                        <Button
                            variant="outline"
                            class="w-full"
                            on:click={clearFilters}
                        >
                            Clear All Filters
                        </Button>
                    </div>
                {/if}
            </div>
        {/if}
    </div>

    <!-- Timeline Section -->
    <div class="mt-12">
        <div class="food-entries-container">
            <FoodLogEntries />
        </div>
    </div>
</div>

<Dialog.Root bind:open={uploadModalOpen}>
    <Dialog.Content class="sm:max-w-[600px] overflow-y-scroll max-h-screen">
        <Dialog.Header>
            <Dialog.Title class="flex items-center gap-2">
                <Utensils class="h-5 w-5" />
                {currentMealType === "breakfast"
                    ? "Add Breakfast"
                    : currentMealType === "lunch"
                      ? "Add Lunch"
                      : currentMealType === "dinner"
                        ? "Add Dinner"
                        : "Add Food"}
            </Dialog.Title>
            <Dialog.Description>
                Upload a photo and details about your {currentMealType}
            </Dialog.Description>
        </Dialog.Header>

        <div class="lg:py-4">
            <FoodLogUpload
                initialTag={currentMealType}
                compact={true}
                onSuccess={handleUploadSuccess}
            />
        </div>

        <Dialog.Footer>
            <Button
                variant="outline"
                on:click={() => (uploadModalOpen = false)}
            >
                Cancel
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<style>
    /* Add loader style for loading images */
    .loader {
        width: 1.5rem;
        height: 1.5rem;
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

    @media (max-width: 640px) {
        .page-header h1 {
            font-size: 1.75rem;
        }

        .page-header {
            margin-bottom: 1.5rem;
        }

        .upload-area {
            padding: 1rem !important;
        }

        .food-log-page {
            padding-left: 1rem !important;
            padding-right: 1rem !important;
        }

        .filter-options {
            gap: 0.75rem;
            padding: 0.75rem;
        }

        .food-tags {
            gap: 0.5rem;
        }

        :global(.meal-sections button) {
            min-height: 2.5rem;
        }
    }
</style>
