<script lang="ts">
    import { page } from "$app/stores";
    import NotebookComponent from "$lib/components/sections/NotebookComponent.svelte";

    $: notebookId = $page.params.id;
</script>

<NotebookComponent id={notebookId} />
