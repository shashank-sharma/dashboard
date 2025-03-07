<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { Check, X, TerminalSquare } from "lucide-svelte";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { pb } from "$lib/config/pocketbase";
    import { toast } from "svelte-sonner";
    import type { Server } from "../types";

    export let open = false;
    export let server: Server | null = null;
    export let onClose: () => void;

    let terminalOutput: string[] = [];
    let command = "";
    let isConnecting = false;
    let isConnected = false;
    let isExecuting = false;
    let terminalRef: HTMLDivElement;
    let connectionId: string | null = null;

    async function connectToServer() {
        if (!server) return;

        isConnecting = true;
        terminalOutput = [];
        addToTerminal(
            `Connecting to ${server.ip}:${server.port} as ${server.username}...`,
        );

        try {
            const response = await pb.send("/api/servers/connect", {
                method: "POST",
                body: JSON.stringify({
                    server_id: server.id,
                }),
            });

            if (response && response.connection_id) {
                connectionId = response.connection_id;
                isConnected = true;
                addToTerminal("Connected successfully.");
                addToTerminal(`Welcome to ${server.name} (${server.ip})`);
                addToTerminal("Type 'exit' to close the connection.");
            } else {
                throw new Error("Failed to establish connection");
            }
        } catch (error: any) {
            console.error("SSH connection error:", error);
            addToTerminal(
                "Connection failed: " + (error.message || "Unknown error"),
            );
            toast.error("Failed to connect to server");
        } finally {
            isConnecting = false;
        }
    }

    /**
     * Execute a command on the server
     */
    async function executeCommand() {
        if (!command.trim() || !isConnected || !connectionId) return;

        const cmd = command.trim();
        addToTerminal(`$ ${cmd}`);
        command = "";

        if (cmd === "exit") {
            disconnectFromServer();
            return;
        }

        isExecuting = true;

        try {
            // Send the command to the SSH server via API
            const response = await pb.send("/api/servers/execute", {
                method: "POST",
                body: JSON.stringify({
                    connection_id: connectionId,
                    command: cmd,
                }),
            });

            if (response && response.output) {
                if (response.output.trim()) {
                    // Split multi-line output
                    const lines = response.output.split("\n");
                    for (const line of lines) {
                        addToTerminal(line);
                    }
                } else {
                    addToTerminal("Command executed with no output.");
                }
            } else {
                addToTerminal("No response from server.");
            }
        } catch (error: any) {
            console.error("Command execution error:", error);
            addToTerminal(
                "Command execution failed: " +
                    (error.message || "Unknown error"),
            );

            // If the connection was lost, attempt to reconnect
            if (
                error.status === 404 &&
                error.message.includes("connection not found")
            ) {
                isConnected = false;
                connectionId = null;
                addToTerminal("Connection lost. Attempting to reconnect...");
                await connectToServer();
            }
        } finally {
            isExecuting = false;
        }
    }

    /**
     * Disconnect from the server
     */
    async function disconnectFromServer() {
        if (!connectionId) return;

        try {
            // Send disconnect request to API
            await pb.send("/api/servers/disconnect", {
                method: "POST",
                body: JSON.stringify({
                    connection_id: connectionId,
                }),
            });
        } catch (error) {
            console.error("Disconnect error:", error);
        } finally {
            isConnected = false;
            connectionId = null;
            addToTerminal("Disconnected from server.");
        }
    }

    /**
     * Add text to the terminal output
     */
    function addToTerminal(text: string) {
        // Split multi-line output
        const lines = text.split("\n");
        terminalOutput = [...terminalOutput, ...lines];

        // Scroll to bottom
        setTimeout(() => {
            if (terminalRef) {
                terminalRef.scrollTop = terminalRef.scrollHeight;
            }
        }, 0);
    }

    /**
     * Handle key press events
     */
    function handleKeyDown(e: KeyboardEvent) {
        if (e.key === "Enter" && !isExecuting) {
            e.preventDefault();
            executeCommand();
        }
    }

    // Connect to server when dialog opens
    $: if (open && server) {
        connectToServer();
    }

    // Clean up when component is destroyed
    onDestroy(() => {
        if (isConnected && connectionId) {
            disconnectFromServer();
        }
    });
</script>

<Dialog.Root
    bind:open
    onOpenChange={() => {
        if (isConnected) {
            disconnectFromServer();
        }
        onClose();
    }}
>
    <Dialog.Content class="sm:max-w-[80%] max-h-[90vh] flex flex-col">
        <Dialog.Header>
            <Dialog.Title class="flex items-center gap-2">
                <TerminalSquare class="w-5 h-5" />
                {server ? `SSH: ${server.name} (${server.ip})` : "SSH Terminal"}
            </Dialog.Title>
            <Dialog.Description>
                {#if isConnecting}
                    Connecting to server...
                {:else if isConnected}
                    Connected to {server?.username}@{server?.ip}:{server?.port}
                {:else}
                    Terminal session
                {/if}
            </Dialog.Description>
        </Dialog.Header>

        <!-- Terminal Output -->
        <div
            class="flex-1 min-h-[300px] max-h-[60vh] overflow-y-auto bg-black text-green-500 p-4 font-mono text-sm rounded-md mb-4"
            bind:this={terminalRef}
        >
            {#each terminalOutput as line}
                <div class="whitespace-pre-wrap break-all">{line}</div>
            {/each}
            {#if isExecuting}
                <div class="animate-pulse">...</div>
            {/if}
        </div>

        <!-- Command Input -->
        <div class="flex gap-2 items-center">
            <div class="flex-1 relative">
                <Input
                    bind:value={command}
                    on:keydown={handleKeyDown}
                    placeholder={isConnected
                        ? "Enter command..."
                        : "Connecting..."}
                    disabled={!isConnected || isExecuting}
                    class="font-mono pl-6"
                />
                <div
                    class="absolute inset-y-0 left-0 flex items-center pl-3 text-muted-foreground"
                >
                    {#if isConnected}
                        <span class="text-xs">$</span>
                    {/if}
                </div>
            </div>
            <Button
                variant="default"
                disabled={!isConnected || isExecuting || !command.trim()}
                on:click={executeCommand}
            >
                <Check class="w-4 h-4 mr-2" />
                Run
            </Button>
        </div>

        <Dialog.Footer class="mt-4">
            <span class="text-xs text-muted-foreground mr-auto">
                {isConnected ? "Connected" : "Disconnected"}
            </span>
            <Button variant="outline" on:click={onClose}>
                <X class="w-4 h-4 mr-2" />
                Close
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>
