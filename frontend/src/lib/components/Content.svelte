<script lang="ts">
    import DashboardContent from "./sections/DashboardContent.svelte";
    import TasksContent from "./sections/TasksContent.svelte";
    // Import other content components as needed

    export let activeSection: string;

    const contentComponents = {
        dashboard: DashboardContent,
        tasks: TasksContent,
        // Add other sections here
    };
</script>

<main class="flex-1 p-6 overflow-auto">
    <h1 class="text-3xl font-bold mb-6 capitalize">{activeSection}</h1>
    <svelte:component this={contentComponents[activeSection]} active={true} />
</main>
