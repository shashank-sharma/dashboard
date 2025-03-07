<script lang="ts">
    import DashboardSidebar from "./DashboardSidebar.svelte";
    import DashboardHeader from "./DashboardHeader.svelte";
    import { fade } from "svelte/transition";
    import type { DashboardSection } from "../types";
    import { onMount, onDestroy } from "svelte";
    import { writable } from "svelte/store";

    export let sections: DashboardSection[];

    // Create a store for tracking if we're on mobile
    const isMobileView = writable(false);

    // Store for tracking if subsections are expanded
    const mobileExpandedSection = writable<string | null>(null);

    // Function to check if we're on mobile based on window width
    function checkMobileView() {
        isMobileView.set(window.innerWidth < 768);
    }

    // Set up event listeners
    onMount(() => {
        checkMobileView();
        window.addEventListener("resize", checkMobileView);
    });

    onDestroy(() => {
        window.removeEventListener("resize", checkMobileView);
    });
</script>

<div class="flex flex-col h-screen bg-background" in:fade={{ duration: 150 }}>
    <div class="flex-1 flex flex-col md:flex-row overflow-hidden">
        <DashboardSidebar
            {sections}
            isMobile={$isMobileView}
            {mobileExpandedSection}
        />

        <div class="flex-1 flex flex-col overflow-hidden">
            <DashboardHeader />

            <main class="flex-1 overflow-auto lg:p-6 md:p4">
                <slot />
                {#if $isMobileView}
                    <!-- Spacer to prevent content from being hidden behind the mobile navigation -->
                    <div
                        class="h-20 w-full mt-6"
                        class:h-40={$mobileExpandedSection !== null}
                    ></div>
                {/if}
            </main>
        </div>
    </div>
</div>
