<script lang="ts">
    import { notebooksStore } from "../stores/notebooks.store";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { toast } from "svelte-sonner";

    export let open = false;

    let name = "";
    let loading = false;

    async function handleSubmit() {
        if (!name.trim()) {
            toast.error("Please enter a notebook name");
            return;
        }

        loading = true;
        try {
            await notebooksStore.createNotebook({
                name,
                version: "1.0",
                cells: [],
            });
            toast.success("Notebook created successfully");
            open = false;
            name = "";
        } catch (error) {
            toast.error("Failed to create notebook");
        } finally {
            loading = false;
        }
    }

    function handleClose() {
        name = "";
        open = false;
    }
</script>

<Dialog.Root bind:open on:close={handleClose}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Create New Notebook</Dialog.Title>
            <Dialog.Description>
                Create a new notebook to start organizing your code and notes.
            </Dialog.Description>
        </Dialog.Header>

        <form on:submit|preventDefault={handleSubmit} class="space-y-4">
            <div class="space-y-2">
                <Label for="name">Name</Label>
                <Input
                    id="name"
                    bind:value={name}
                    placeholder="Enter notebook name"
                />
            </div>

            <Dialog.Footer>
                <Button variant="outline" type="button" on:click={handleClose}>
                    Cancel
                </Button>
                <Button type="submit" disabled={loading}>
                    {loading ? "Creating..." : "Create Notebook"}
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>
