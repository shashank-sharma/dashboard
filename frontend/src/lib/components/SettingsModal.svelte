<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Settings } from "lucide-svelte";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { onMount } from "svelte";

    let isOpen = false;
    let pocketbaseUrl = "http://127.0.0.1:8090";

    onMount(() => {
        const savedUrl = localStorage.getItem("pocketbaseUrl");
        if (savedUrl) {
            pocketbaseUrl = savedUrl;
        }
    });

    function updatePocketbaseUrl() {
        localStorage.setItem("pocketbaseUrl", pocketbaseUrl);
        isOpen = false;
        window.location.reload();
    }
</script>

<Button variant="ghost" size="icon" on:click={() => (isOpen = true)}>
    <Settings class="h-[1.2rem] w-[1.2rem]" />
    <span class="sr-only">Open settings</span>
</Button>

<Dialog.Root bind:open={isOpen}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Settings</Dialog.Title>
            <Dialog.Description>
                Configure application settings. Changes will require a reload to
                take effect.
            </Dialog.Description>
        </Dialog.Header>

        <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
                <Label class="text-right" for="pocketbase-url">
                    PocketBase URL
                </Label>
                <Input
                    id="pocketbase-url"
                    class="col-span-3"
                    bind:value={pocketbaseUrl}
                />
            </div>
        </div>

        <Dialog.Footer>
            <Button variant="outline" on:click={() => (isOpen = false)}>
                Cancel
            </Button>
            <Button on:click={updatePocketbaseUrl}>Save changes</Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>
