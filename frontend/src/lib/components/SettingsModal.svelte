<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Settings } from "lucide-svelte";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { onMount } from "svelte";
    import { Moon, Sun } from "lucide-svelte";
    import { Switch } from "$lib/components/ui/switch";
    import { browser } from "$app/environment";
    import { isDebugModeEnabled } from "$lib/features/pwa/services";

    let isOpen = false;
    let pocketbaseUrl = "http://127.0.0.1:8090";

    let pwaDebugMode = false;
    let pwaDebugModeChanged = false;
    let selectedTheme = "system";

    onMount(() => {
        if (browser) {
            const savedUrl = localStorage.getItem("pocketbase-url");
            if (savedUrl) {
                pocketbaseUrl = savedUrl;
            }

            pwaDebugMode = isDebugModeEnabled();

            const storedTheme = localStorage.getItem("theme");
            if (storedTheme) {
                selectedTheme = storedTheme;
            }
        }
    });

    function updatePocketbaseUrl() {
        if (browser) {
            localStorage.setItem("pocketbase-url", pocketbaseUrl);
            window.location.reload();
        }
    }

    function togglePwaDebugMode(value: boolean) {
        pwaDebugMode = value;
        pwaDebugModeChanged = true;

        if (browser) {
            localStorage.setItem("pwa-debug-mode", pwaDebugMode.toString());
            console.log(
                `[PWA] Debug mode ${pwaDebugMode ? "enabled" : "disabled"}`,
            );
        }
    }

    function handleSave() {
        updatePocketbaseUrl();

        if (pwaDebugModeChanged && browser) {
            window.location.reload();
        } else {
            isOpen = false;
        }
    }
</script>

<Dialog.Root bind:open={isOpen}>
    <Dialog.Trigger asChild let:builder>
        <Button
            variant="outline"
            size="icon"
            builders={[builder]}
            title="Settings"
        >
            <Settings className="h-4 w-4" />
        </Button>
    </Dialog.Trigger>
    <Dialog.Portal>
        <Dialog.Overlay />
        <Dialog.Content class="sm:max-w-[425px]">
            <Dialog.Header>
                <Dialog.Title class="flex items-center gap-2">
                    <Settings class="h-5 w-5" />
                    <span>Settings</span>
                </Dialog.Title>
            </Dialog.Header>

            <div class="space-y-6 py-4">
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right" for="pocketbase-url">
                        Pocketbase URL
                    </Label>
                    <Input
                        id="pocketbase-url"
                        bind:value={pocketbaseUrl}
                        class="col-span-3"
                    />
                </div>

                <!-- Theme Placeholder -->
                <div class="flex items-center justify-between">
                    <div class="space-y-0.5">
                        <Label>Theme</Label>
                        <div class="text-sm text-muted-foreground">
                            Choose your preferred theme
                        </div>
                    </div>
                    <div class="flex gap-2">
                        <Button
                            variant="outline"
                            size="icon"
                            title="Light Mode"
                        >
                            <Sun class="h-4 w-4" />
                        </Button>
                        <Button variant="outline" size="icon" title="Dark Mode">
                            <Moon class="h-4 w-4" />
                        </Button>
                    </div>
                </div>

                <!-- PWA Debug Mode Toggle -->
                <div class="flex items-center justify-between">
                    <div class="space-y-0.5">
                        <Label>PWA Debug Mode</Label>
                        <div class="text-sm text-muted-foreground">
                            Enable detailed PWA installation debugging
                        </div>
                    </div>
                    <Switch
                        checked={pwaDebugMode}
                        onCheckedChange={togglePwaDebugMode}
                    />
                </div>
            </div>

            <Dialog.Footer>
                <Button variant="outline" on:click={() => (isOpen = false)}>
                    Cancel
                </Button>
                <Button on:click={handleSave}>Save changes</Button>
            </Dialog.Footer>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
