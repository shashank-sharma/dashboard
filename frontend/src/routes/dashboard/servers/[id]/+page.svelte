<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { page } from "$app/stores";
    import { pb } from "$lib/config/pocketbase";
    import { toast } from "svelte-sonner";
    import {
        Loader2,
        Server,
        Key,
        Monitor,
        AlertCircle,
        CheckCircle2,
    } from "lucide-svelte";
    import SSHTerminal from "$lib/components/terminal/SSHTerminal.svelte";
    import { Button } from "$lib/components/ui/button";
    import { Badge } from "$lib/components/ui/badge";
    import {
        Card,
        CardContent,
        CardDescription,
        CardFooter,
        CardHeader,
        CardTitle,
    } from "$lib/components/ui/card";
    import {
        Tabs,
        TabsContent,
        TabsList,
        TabsTrigger,
    } from "$lib/components/ui/tabs";

    let server: any = null;
    let securityKey: any = null;
    let isLoading = true;
    let error = "";
    let connectionId = "";
    let isConnecting = false;

    const serverId = $page.params.id;

    onMount(async () => {
        try {
            const serverRecord = await pb
                .collection("servers")
                .getOne(serverId, {
                    expand: "security_key",
                });

            server = serverRecord;

            if (serverRecord.security_key) {
                securityKey = serverRecord.expand?.security_key;
            }

            isLoading = false;
        } catch (err: any) {
            error = "Failed to load server details";
            toast.error(error);
            isLoading = false;
        }
    });

    onDestroy(() => {
        if (connectionId) {
            disconnectFromServer();
        }
    });

    async function connectToServer() {
        if (!server) return;

        isConnecting = true;
        error = "";

        try {
            const response = await pb.send("/api/servers/connect", {
                method: "POST",
                body: JSON.stringify({
                    server_id: serverId,
                }),
            });

            if (response && response.connection_id) {
                connectionId = response.connection_id;
                toast.success("Connected to server");
            } else {
                throw new Error("Failed to connect to server");
            }
        } catch (err: any) {
            console.error("Connection error:", err);
            error = err.message || "Failed to connect to server";
            toast.error(error);
        } finally {
            isConnecting = false;
        }
    }

    function disconnectFromServer() {
        if (!connectionId) return;

        pb.send("/api/ssh/disconnect", {
            method: "POST",
            body: JSON.stringify({
                connection_id: connectionId,
            }),
        }).catch((err) => {
            console.log("Failed to disconnect:", err);
        });

        connectionId = "";
    }

    function handleTerminalDisconnect() {
        connectionId = "";
    }
</script>

<svelte:head>
    <title>{server ? server.name : "Server"} | Dashboard</title>
</svelte:head>

<div class="container mx-auto py-6 space-y-6">
    {#if isLoading}
        <div class="flex justify-center items-center h-64">
            <Loader2 class="h-8 w-8 animate-spin text-primary" />
        </div>
    {:else if error}
        <div
            class="bg-destructive/10 border border-destructive p-4 rounded-md flex items-start gap-3"
        >
            <AlertCircle
                class="h-5 w-5 text-destructive flex-shrink-0 mt-0.5"
            />
            <div>
                <h3 class="font-medium text-destructive">
                    Error loading server
                </h3>
                <p class="text-destructive/90 text-sm">{error}</p>
            </div>
        </div>
    {:else if server}
        <div class="flex justify-between items-start">
            <div>
                <h1 class="text-2xl font-bold flex items-center gap-2">
                    <Server class="h-6 w-6" />
                    {server.name}
                </h1>
                <p class="text-muted-foreground">{server.ip}:{server.port}</p>
            </div>
            <div class="flex gap-2">
                <Badge
                    variant={server.is_active ? "outline" : "secondary"}
                    class={server.is_active
                        ? "bg-green-100 text-green-800 hover:bg-green-100"
                        : ""}
                >
                    {server.is_active ? "Active" : "Inactive"}
                </Badge>
                <Badge
                    variant={server.is_reachable ? "outline" : "secondary"}
                    class={server.is_reachable
                        ? "bg-blue-100 text-blue-800 hover:bg-blue-100"
                        : ""}
                >
                    {server.is_reachable ? "Reachable" : "Unreachable"}
                </Badge>
            </div>
        </div>

        <Tabs value="terminal" class="w-full">
            <TabsList>
                <TabsTrigger value="terminal">Terminal</TabsTrigger>
                <TabsTrigger value="details">Server Details</TabsTrigger>
            </TabsList>

            <TabsContent value="terminal" class="mt-4">
                <Card>
                    <CardHeader>
                        <CardTitle>SSH Terminal</CardTitle>
                        <CardDescription>
                            Connect to your server and run commands
                        </CardDescription>
                    </CardHeader>

                    <CardContent>
                        {#if !connectionId && !isConnecting}
                            <div
                                class="flex flex-col items-center justify-center h-64 space-y-4"
                            >
                                <div class="text-center">
                                    {#if !server.ssh_enabled}
                                        <AlertCircle
                                            class="h-10 w-10 text-destructive mx-auto mb-2"
                                        />
                                        <h3 class="text-lg font-medium">
                                            SSH is disabled
                                        </h3>
                                        <p class="text-muted-foreground">
                                            SSH access is disabled for this
                                            server.
                                        </p>
                                    {:else if !server.security_key}
                                        <Key
                                            class="h-10 w-10 text-warning mx-auto mb-2"
                                        />
                                        <h3 class="text-lg font-medium">
                                            No security key
                                        </h3>
                                        <p class="text-muted-foreground">
                                            This server doesn't have a security
                                            key configured.
                                        </p>
                                    {:else if securityKey && !securityKey.is_active}
                                        <AlertCircle
                                            class="h-10 w-10 text-warning mx-auto mb-2"
                                        />
                                        <h3 class="text-lg font-medium">
                                            Inactive security key
                                        </h3>
                                        <p class="text-muted-foreground">
                                            The security key for this server is
                                            inactive.
                                        </p>
                                    {:else}
                                        <Monitor
                                            class="h-10 w-10 text-primary mx-auto mb-2"
                                        />
                                        <h3 class="text-lg font-medium">
                                            Ready to connect
                                        </h3>
                                        <p class="text-muted-foreground">
                                            Click the button below to start an
                                            SSH session.
                                        </p>
                                    {/if}
                                </div>

                                <Button
                                    on:click={connectToServer}
                                    disabled={!server.ssh_enabled ||
                                        !server.security_key ||
                                        (securityKey && !securityKey.is_active)}
                                >
                                    Connect to Server
                                </Button>
                            </div>
                        {:else if isConnecting}
                            <div class="flex justify-center items-center h-64">
                                <div class="text-center">
                                    <Loader2
                                        class="h-8 w-8 animate-spin text-primary mx-auto mb-2"
                                    />
                                    <p>Connecting to {server.name}...</p>
                                </div>
                            </div>
                        {:else}
                            <div class="h-[500px]">
                                <SSHTerminal
                                    {connectionId}
                                    serverName={server.name}
                                    serverAddress={server.id}
                                    onDisconnect={handleTerminalDisconnect}
                                />
                            </div>
                        {/if}
                    </CardContent>
                </Card>
            </TabsContent>

            <TabsContent value="details" class="mt-4">
                <Card>
                    <CardHeader>
                        <CardTitle>Server Information</CardTitle>
                    </CardHeader>

                    <CardContent>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div class="space-y-2">
                                <h3
                                    class="font-medium text-sm text-muted-foreground"
                                >
                                    Basic Information
                                </h3>
                                <div class="grid grid-cols-3 gap-2">
                                    <div class="text-sm font-medium">Name</div>
                                    <div class="text-sm col-span-2">
                                        {server.name}
                                    </div>

                                    <div class="text-sm font-medium">
                                        Provider
                                    </div>
                                    <div class="text-sm col-span-2">
                                        {server.provider || "Not specified"}
                                    </div>

                                    <div class="text-sm font-medium">
                                        IP Address
                                    </div>
                                    <div class="text-sm col-span-2">
                                        {server.ip}
                                    </div>

                                    <div class="text-sm font-medium">
                                        SSH Port
                                    </div>
                                    <div class="text-sm col-span-2">
                                        {server.port}
                                    </div>

                                    <div class="text-sm font-medium">
                                        Username
                                    </div>
                                    <div class="text-sm col-span-2">
                                        {server.username}
                                    </div>
                                </div>
                            </div>

                            <div class="space-y-2">
                                <h3
                                    class="font-medium text-sm text-muted-foreground"
                                >
                                    Status Information
                                </h3>
                                <div class="grid grid-cols-3 gap-2">
                                    <div class="text-sm font-medium">
                                        SSH Enabled
                                    </div>
                                    <div
                                        class="text-sm col-span-2 flex items-center gap-1"
                                    >
                                        {#if server.ssh_enabled}
                                            <CheckCircle2
                                                class="h-4 w-4 text-green-600"
                                            />
                                            <span>Enabled</span>
                                        {:else}
                                            <AlertCircle
                                                class="h-4 w-4 text-destructive"
                                            />
                                            <span>Disabled</span>
                                        {/if}
                                    </div>

                                    <div class="text-sm font-medium">
                                        Active Status
                                    </div>
                                    <div
                                        class="text-sm col-span-2 flex items-center gap-1"
                                    >
                                        {#if server.is_active}
                                            <CheckCircle2
                                                class="h-4 w-4 text-green-600"
                                            />
                                            <span>Active</span>
                                        {:else}
                                            <AlertCircle
                                                class="h-4 w-4 text-destructive"
                                            />
                                            <span>Inactive</span>
                                        {/if}
                                    </div>

                                    <div class="text-sm font-medium">
                                        Reachable
                                    </div>
                                    <div
                                        class="text-sm col-span-2 flex items-center gap-1"
                                    >
                                        {#if server.is_reachable}
                                            <CheckCircle2
                                                class="h-4 w-4 text-green-600"
                                            />
                                            <span>Reachable</span>
                                        {:else}
                                            <AlertCircle
                                                class="h-4 w-4 text-destructive"
                                            />
                                            <span>Unreachable</span>
                                        {/if}
                                    </div>

                                    <div class="text-sm font-medium">
                                        Security Key
                                    </div>
                                    <div class="text-sm col-span-2">
                                        {#if server.security_key && securityKey}
                                            <span
                                                class="flex items-center gap-1"
                                            >
                                                <Key class="h-4 w-4" />
                                                {securityKey.name}
                                                {#if securityKey.is_active}
                                                    <Badge
                                                        variant="outline"
                                                        class="bg-green-100 text-green-800 hover:bg-green-100 text-xs"
                                                        >Active</Badge
                                                    >
                                                {:else}
                                                    <Badge
                                                        variant="outline"
                                                        class="bg-red-100 text-red-800 hover:bg-red-100 text-xs"
                                                        >Inactive</Badge
                                                    >
                                                {/if}
                                            </span>
                                        {:else}
                                            <span class="text-muted-foreground"
                                                >No security key configured</span
                                            >
                                        {/if}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </TabsContent>
        </Tabs>
    {/if}
</div>

<style>
    /* Add custom styles if needed */
</style>
