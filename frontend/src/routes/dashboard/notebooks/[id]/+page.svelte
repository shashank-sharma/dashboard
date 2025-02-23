<script lang="ts">
    import { page } from "$app/stores";
    import { NotebookEditor } from "$lib/features/notebooks/components";

    $: notebookId = $page.params.id;
</script>

<NotebookEditor id={notebookId} />
