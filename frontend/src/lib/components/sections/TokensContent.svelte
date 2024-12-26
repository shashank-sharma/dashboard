<script lang="ts">
    import { pb } from "$lib/pocketbase";
    import { onMount } from "svelte";
    import { toast } from "svelte-sonner";
    import {
        Card,
        CardContent,
        CardDescription,
        CardFooter,
        CardHeader,
        CardTitle,
    } from "$lib/components/ui/card";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Button } from "$lib/components/ui/button";
    import * as Dialog from "$lib/components/ui/dialog";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
    import {
        Select,
        SelectContent,
        SelectItem,
        SelectTrigger,
        SelectValue,
    } from "$lib/components/ui/select";
    import { Badge } from "$lib/components/ui/badge";
    import { Switch } from "$lib/components/ui/switch";
    import {
        Eye,
        EyeOff,
        Plus,
        RefreshCcw,
        Trash2,
        Key,
        Calendar as CalendarIcon,
    } from "lucide-svelte";
    import { Calendar } from "$lib/components/ui/calendar";
    import * as Popover from "$lib/components/ui/popover";
    import {
        DateFormatter,
        type DateValue,
        getLocalTimeZone,
    } from "@internationalized/date";
    import { cn } from "$lib/utils";
    import { buttonVariants } from "$lib/components/ui/button";

    const df = new DateFormatter("en-US", {
        dateStyle: "long",
    });

    let tokens: any[] = [];
    let loading = true;
    let showTokenDialog = false;
    let showTokenValue = false;
    let selectedProvider = "";
    let expiryDate: DateValue | undefined = undefined;
    let contentRef: HTMLElement | null = null;
    let tokenToDelete: string | null = null;
    let showDeleteDialog = false;

    // Form data
    let formData = {
        provider: "",
        account: "",
        access_token: "",
        token_type: "Bearer",
        refresh_token: "",
        expiry: "",
        scope: "",
    };

    const providers = [
        { value: "google", label: "Google" },
        { value: "coolify", label: "Coolify" },
        { value: "github", label: "GitHub" },
        { value: "gitlab", label: "GitLab" },
    ];

    function handleProviderChange(value: string) {
        formData.provider = value;
    }

    async function loadTokens() {
        try {
            const records = await pb.collection("tokens").getFullList({
                sort: "-created",
                expand: "user",
            });
            tokens = records;
        } catch (error) {
            toast.error("Failed to load tokens");
        } finally {
            loading = false;
        }
    }

    async function handleSubmit() {
        try {
            const data = {
                ...formData,
                user: pb.authStore.model.id,
                is_active: true,
                expiry: expiryDate
                    ? expiryDate.toDate(getLocalTimeZone()).toISOString()
                    : null,
            };
            await pb.collection("tokens").create(data);
            toast.success("Token created successfully");
            showTokenDialog = false;
            formData = {
                provider: "",
                account: "",
                access_token: "",
                token_type: "Bearer",
                refresh_token: "",
                expiry: "",
                scope: "",
            };
            expiryDate = undefined;
            loadTokens();
        } catch (error) {
            toast.error("Failed to create token");
        }
    }

    async function toggleTokenStatus(id: string, currentStatus: boolean) {
        try {
            await pb.collection("tokens").update(id, {
                is_active: !currentStatus,
            });
            toast.success("Token status updated");
            loadTokens();
        } catch (error) {
            toast.error("Failed to update token status");
        }
    }

    function handleDeleteClick(id: string) {
        tokenToDelete = id;
        showDeleteDialog = true;
    }

    async function confirmDelete() {
        if (!tokenToDelete) return;

        try {
            await pb.collection("tokens").delete(tokenToDelete);
            toast.success("Token deleted successfully");
            loadTokens();
        } catch (error) {
            toast.error("Failed to delete token");
        } finally {
            showDeleteDialog = false;
            tokenToDelete = null;
        }
    }

    function formatDate(date: string) {
        return new Date(date).toLocaleString();
    }

    onMount(() => {
        loadTokens();
    });
</script>

<AlertDialog.Root bind:open={showDeleteDialog}>
    <AlertDialog.Content>
        <AlertDialog.Header>
            <AlertDialog.Title>Are you sure?</AlertDialog.Title>
            <AlertDialog.Description>
                This action cannot be undone. This will permanently delete the
                token and remove it from our servers.
            </AlertDialog.Description>
        </AlertDialog.Header>
        <AlertDialog.Footer>
            <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
            <AlertDialog.Action on:click={confirmDelete}>
                Delete
            </AlertDialog.Action>
        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>

<div class="container mx-auto p-4">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold">API Tokens</h2>
        <Dialog.Root bind:open={showTokenDialog}>
            <Dialog.Trigger>
                <Button class="flex items-center gap-2">
                    <Plus class="w-4 h-4" />
                    New Token
                </Button>
            </Dialog.Trigger>
            <Dialog.Content class="sm:max-w-[425px]">
                <Dialog.Header>
                    <Dialog.Title>Create New Token</Dialog.Title>
                </Dialog.Header>

                <form class="space-y-4" on:submit|preventDefault={handleSubmit}>
                    <div class="space-y-2">
                        <Label for="provider">Provider</Label>
                        <Select value={formData.provider}>
                            <SelectTrigger class="w-full">
                                <SelectValue placeholder="Select provider" />
                            </SelectTrigger>
                            <SelectContent>
                                {#each providers as provider}
                                    <SelectItem
                                        value={provider.value}
                                        on:click={() =>
                                            handleProviderChange(
                                                provider.value,
                                            )}
                                    >
                                        {provider.label}
                                    </SelectItem>
                                {/each}
                            </SelectContent>
                        </Select>
                    </div>

                    <div class="space-y-2">
                        <Label for="account">Account</Label>
                        <Input
                            id="account"
                            bind:value={formData.account}
                            placeholder="Account name or email"
                            required
                        />
                    </div>

                    <div class="space-y-2">
                        <Label for="access_token">Access Token</Label>
                        <div class="relative">
                            <Input
                                id="access_token"
                                type={showTokenValue ? "text" : "password"}
                                bind:value={formData.access_token}
                                placeholder="Access token"
                                required
                            />
                            <button
                                type="button"
                                class="absolute right-2 top-1/2 -translate-y-1/2"
                                on:click={() =>
                                    (showTokenValue = !showTokenValue)}
                            >
                                {#if showTokenValue}
                                    <EyeOff class="w-4 h-4" />
                                {:else}
                                    <Eye class="w-4 h-4" />
                                {/if}
                            </button>
                        </div>
                    </div>

                    <div class="space-y-2">
                        <Label for="token_type">Token Type</Label>
                        <Input
                            id="token_type"
                            bind:value={formData.token_type}
                            placeholder="Bearer"
                            required
                        />
                    </div>

                    <div class="space-y-2">
                        <Label for="refresh_token">Refresh Token</Label>
                        <Input
                            id="refresh_token"
                            type="password"
                            bind:value={formData.refresh_token}
                            placeholder="Refresh token (optional)"
                        />
                    </div>

                    <div class="space-y-2">
                        <Label for="expiry">Expiry Date</Label>
                        <Popover.Root>
                            <Popover.Trigger
                                class={cn(
                                    buttonVariants({
                                        variant: "outline",
                                        class: "w-full justify-start text-left font-normal",
                                    }),
                                    !expiryDate && "text-muted-foreground",
                                )}
                            >
                                <CalendarIcon class="w-4 h-4 mr-2" />
                                {expiryDate
                                    ? df.format(
                                          expiryDate.toDate(getLocalTimeZone()),
                                      )
                                    : "Pick a date"}
                            </Popover.Trigger>
                            <Popover.Content
                                bind:ref={contentRef}
                                class="w-auto p-0"
                            >
                                <Calendar
                                    bind:value={expiryDate}
                                    mode="single"
                                />
                            </Popover.Content>
                        </Popover.Root>
                    </div>

                    <div class="space-y-2">
                        <Label for="scope">Scope</Label>
                        <Input
                            id="scope"
                            bind:value={formData.scope}
                            placeholder="Token scope (optional)"
                        />
                    </div>

                    <Dialog.Footer>
                        <Button type="submit">Create Token</Button>
                    </Dialog.Footer>
                </form>
            </Dialog.Content>
        </Dialog.Root>
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <RefreshCcw class="w-8 h-8 animate-spin" />
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each tokens as token}
                <Card>
                    <CardHeader>
                        <CardTitle class="flex items-center justify-between">
                            <div class="flex items-center gap-2">
                                <Key class="w-4 h-4" />
                                <span>{token.provider}</span>
                            </div>
                            <Badge
                                variant={token.is_active
                                    ? "default"
                                    : "secondary"}
                            >
                                {token.is_active ? "Active" : "Disabled"}
                            </Badge>
                        </CardTitle>
                        <CardDescription>{token.account}</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div class="space-y-2">
                            <div class="text-sm">
                                <span class="font-medium">Created:</span>
                                {formatDate(token.created)}
                            </div>
                            {#if token.expiry}
                                <div class="text-sm">
                                    <span class="font-medium">Expires:</span>
                                    {formatDate(token.expiry)}
                                </div>
                            {/if}
                            {#if token.scope}
                                <div class="text-sm">
                                    <span class="font-medium">Scope:</span>
                                    {token.scope}
                                </div>
                            {/if}
                        </div>
                    </CardContent>
                    <CardFooter class="justify-between">
                        <div class="flex items-center gap-2">
                            <Switch
                                checked={token.is_active}
                                onCheckedChange={() =>
                                    toggleTokenStatus(
                                        token.id,
                                        token.is_active,
                                    )}
                            />
                            <span class="text-sm">
                                {token.is_active ? "Active" : "Disabled"}
                            </span>
                        </div>
                        <Button
                            variant="destructive"
                            size="icon"
                            on:click={() => handleDeleteClick(token.id)}
                        >
                            <Trash2 class="w-4 h-4" />
                        </Button>
                    </CardFooter>
                </Card>
            {/each}
        </div>
    {/if}
</div>
