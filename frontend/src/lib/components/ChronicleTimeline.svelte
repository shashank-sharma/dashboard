<!-- ChronicleTimeline.svelte -->
<script lang="ts">
    import { onMount } from "svelte";
    import { Button } from "$lib/components/ui/button";
    import {
        ChevronLeft,
        ChevronRight,
        CalendarDays,
        ChevronDown,
    } from "lucide-svelte";
    import { fly } from "svelte/transition";
    import { Root, Trigger, Content } from "$lib/components/ui/tooltip";

    export let onDateSelect = (date: Date) => {};
    export let selectedDate: Date = new Date();

    let scrollContainer: HTMLElement;
    let days: Date[] = [];
    const pastDays = 60; // Days in the past
    const futureDays = 30; // Days in the future
    const totalDays = pastDays + futureDays + 1; // +1 for today

    // Generate array of dates including future
    function generateDays() {
        const today = new Date();
        days = Array.from({ length: totalDays }, (_, i) => {
            const date = new Date();
            date.setDate(today.getDate() - pastDays + i);
            return date;
        });
    }

    // Format date for display
    function formatDate(date: Date) {
        return date.getDate().toString().padStart(2, "0");
    }

    // Format complete date for tooltip
    function formatComplete(date: Date) {
        return date.toLocaleDateString("en-US", {
            weekday: "long",
            year: "numeric",
            month: "long",
            day: "numeric",
        });
    }

    // Format month for display
    function formatMonth(date: Date) {
        return date.toLocaleString("default", { month: "short" });
    }

    // Check if date is today
    function isToday(date: Date) {
        const today = new Date();
        return date.toDateString() === today.toDateString();
    }

    // Check if date is selected
    function isSelected(date: Date) {
        return date.toDateString() === selectedDate.toDateString();
    }

    // Check if date is in future
    function isFuture(date: Date) {
        const today = new Date();
        today.setHours(23, 59, 59, 999); // End of today
        return date > today;
    }

    // Check if it's first day of month
    function isFirstOfMonth(date: Date, index: number) {
        return index === 0 || date.getDate() === 1;
    }

    // Handle scroll buttons
    function scroll(direction: "left" | "right") {
        const scrollAmount = 200;
        if (scrollContainer) {
            scrollContainer.scrollBy({
                left: direction === "left" ? -scrollAmount : scrollAmount,
                behavior: "smooth",
            });
        }
    }

    // Scroll to date and center it
    function scrollToDate(date: Date) {
        const dateIndex = days.findIndex(
            (d) => d.toDateString() === date.toDateString(),
        );
        if (dateIndex !== -1 && scrollContainer) {
            const elements =
                scrollContainer.getElementsByClassName("date-marker");
            if (elements[dateIndex]) {
                elements[dateIndex].scrollIntoView({
                    behavior: "smooth",
                    block: "nearest",
                    inline: "center",
                });
            }
        }
    }

    // Jump to today
    function jumpToPresent() {
        const today = new Date();
        onDateSelect(today);
        scrollToDate(today);
    }

    // Handle date selection
    function handleDateSelect(date: Date) {
        if (!isFuture(date)) {
            onDateSelect(date);
            scrollToDate(date);
        }
    }

    onMount(() => {
        generateDays();
        // Wait for the DOM to update then scroll to selected date
        setTimeout(() => scrollToDate(selectedDate), 0);
    });

    // Watch selectedDate changes
    $: if (selectedDate && scrollContainer) {
        scrollToDate(selectedDate);
    }
</script>

<div class="relative w-full">
    <!-- Center indicator -->
    <div class="absolute top-0 left-1/2 -translate-x-1/2 text-muted-foreground">
        <ChevronDown class="h-4 w-4" />
    </div>

    <!-- Navigation buttons -->
    <Button
        variant="outline"
        size="icon"
        class="absolute left-0 top-1/2 -translate-y-1/2 z-10 bg-background"
        on:click={() => scroll("left")}
    >
        <ChevronLeft class="h-4 w-4" />
    </Button>

    <Button
        variant="outline"
        size="icon"
        class="absolute right-0 top-1/2 -translate-y-1/2 z-10 bg-background"
        on:click={() => scroll("right")}
    >
        <ChevronRight class="h-4 w-4" />
    </Button>

    <!-- Timeline container -->
    <div class="overflow-x-hidden relative mx-12" bind:this={scrollContainer}>
        <div class="flex space-x-2 pt-8">
            {#each days as day, i}
                <Root>
                    <Trigger>
                        <button
                            class="timeline-button flex flex-col items-center min-w-[14px] relative group date-marker"
                            class:opacity-50={isFuture(day)}
                            class:timeline-active={isSelected(day)}
                            on:click={() => handleDateSelect(day)}
                            disabled={isFuture(day)}
                        >
                            {#if isFirstOfMonth(day, i)}
                                <div
                                    class="absolute -top-8 left-1/2 -translate-x-1/2 text-xs text-muted-foreground whitespace-nowrap opacity-80"
                                >
                                    {formatMonth(day)}
                                </div>
                            {/if}

                            <!-- Date marker line -->
                            <div
                                class="timeline-line w-[2px] h-14 bg-border relative transition-colors duration-200"
                                class:bg-primary={isSelected(day)}
                                class:bg-blue-500={isToday(day)}
                                class:hover:bg-primary={isSelected(day)}
                                class:hover:bg-blue-500={isToday(day)}
                                class:opacity-50={isFuture(day)}
                            >
                                <!-- Date number -->
                                <div
                                    class="absolute -top-4 left-1/2 -translate-x-1/2 px-1.5 py-0.5 text-xs font-medium
                                           transition-all duration-200 whitespace-nowrap"
                                    class:text-primary={isSelected(day)}
                                    class:text-blue-500={isToday(day)}
                                    class:opacity-50={isFuture(day)}
                                >
                                    {formatDate(day)}
                                </div>
                            </div>
                        </button>
                    </Trigger>
                    <Content>
                        {formatComplete(day)}
                        {#if isFuture(day)}
                            <span class="text-xs text-muted-foreground">
                                (Future)</span
                            >
                        {/if}
                    </Content>
                </Root>
            {/each}
        </div>
    </div>

    <!-- Jump to present button -->
    {#if !isToday(selectedDate)}
        <div
            class="absolute right-0 mb-2 z-10"
            in:fly={{ y: 20, duration: 300 }}
        >
            <Button
                variant="outline"
                size="sm"
                class="gap-2"
                on:click={jumpToPresent}
            >
                <CalendarDays class="h-4 w-4" />
                Jump to Present
            </Button>
        </div>
    {/if}
</div>

<style>
    .timeline-button:not(:disabled):hover {
        box-shadow: 0 0 8px rgba(var(--primary), 0.8);
    }

    .timeline-line {
        cursor: pointer;
    }

    button:disabled {
        cursor: not-allowed;
    }

    .timeline-active .timeline-line {
        @apply bg-blue-500;
    }

    @keyframes bounce {
        0%,
        100% {
            transform: translateY(0) translateX(-50%);
        }
        50% {
            transform: translateY(4px) translateX(-50%);
        }
    }
</style>
