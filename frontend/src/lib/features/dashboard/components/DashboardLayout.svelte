<script lang="ts">
    import DashboardSidebar from "./DashboardSidebar.svelte";
    import DashboardHeader from "./DashboardHeader.svelte";
    import { fade } from "svelte/transition";
    import type { DashboardSection } from "../types";

    export let sections: DashboardSection[];
</script>

<div class="flex h-screen bg-background" in:fade={{ duration: 150 }}>
    <DashboardSidebar {sections} />

    <div class="flex-1 flex flex-col overflow-hidden">
        <DashboardHeader />

        <main class="flex-1 overflow-auto p-6">
            <slot />
        </main>
    </div>
</div>
