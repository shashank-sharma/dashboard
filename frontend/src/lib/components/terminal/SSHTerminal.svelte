<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import Terminal from "./Terminal.svelte";
    import { pb } from "$lib/config/pocketbase";
    import { toast } from "svelte-sonner";
    import { XCircle, Loader2 } from "lucide-svelte";

    export let connectionId: string = "";
    export let serverName: string = "";
    export let serverAddress: string = "";
    export let onDisconnect: () => void = () => {};

    let terminal: InstanceType<typeof Terminal>;
    let isConnected = false;
    let isConnecting = false;
    let errorMessage = "";
    let pingInterval: ReturnType<typeof setInterval>;
    let lastPongTime = 0;
    let reconnectAttempts = 0;
    let maxReconnectAttempts = 3;
    let currentCommand = "";
    let commandHistory: string[] = [];
    let historyIndex = -1;

    let socket: WebSocket | null = null;

    onMount(() => {
        if (connectionId) {
            connectToServer();
        }

        // Keep me alive
        pingInterval = setInterval(() => {
            if (isConnected && socket?.readyState === WebSocket.OPEN) {
                sendWebSocketPing();
            }
        }, 30000);
    });

    onDestroy(() => {
        if (pingInterval) {
            clearInterval(pingInterval);
        }
    });

    async function connectToServer() {
        isConnecting = true;
        errorMessage = "";

        try {
            // If we already have a connection ID, establish WebSocket connection
            if (connectionId) {
                establishWebSocketConnection();
            } else {
                // If we don't have a connection ID, we need to establish it first
                const response = await pb.send("/api/servers/connect", {
                    method: "POST",
                    body: JSON.stringify({
                        server_id: serverAddress,
                    }),
                });

                if (response && response.connection_id) {
                    connectionId = response.connection_id;
                    establishWebSocketConnection();
                } else {
                    throw new Error("Failed to get connection ID");
                }
            }
        } catch (error: any) {
            console.error("Connection error:", error);
            errorMessage = error.message || "Failed to connect to the server";
            toast.error(errorMessage);
            isConnecting = false;
        }
    }

    function establishWebSocketConnection() {
        try {
            // Determine WebSocket protocol based on current page protocol
            const wsProtocol =
                window.location.protocol === "https:" ? "wss:" : "ws:";

            // In development, we need to use the backend server port (8090) not the frontend port
            // TODO: This is a hack, we need to find a better way to handle this
            let wsHost = window.location.host;
            if (
                window.location.host.includes("5173") ||
                window.location.host.includes("3000")
            ) {
                // We're in development, use the backend port
                wsHost = window.location.hostname + ":8090";
            }

            const authToken = pb.authStore.token;
            const wsUrl = `${wsProtocol}//${wsHost}/api/ssh/stream?connection_id=${connectionId}${authToken ? "&token=" + authToken : ""}`;

            socket = new WebSocket(wsUrl);

            socket.onopen = handleSocketOpen;
            socket.onmessage = handleSocketMessage;
            socket.onclose = handleSocketClose;
            socket.onerror = handleSocketError;
        } catch (error) {
            console.error("WebSocket connection error:", error);
            errorMessage = "Failed to establish WebSocket connection";
            isConnecting = false;
        }
    }

    function handleSocketOpen(event: Event) {
        isConnected = true;
        isConnecting = false;
        reconnectAttempts = 0;

        if (terminal) {
            terminal.writeln(
                "\r\n\x1b[32mConnected to " + serverName + "\x1b[0m\r\n",
            );
            terminal.focus();
        }
    }

    function handleSocketMessage(event: MessageEvent) {
        if (terminal && event.data) {
            try {
                const jsonData = JSON.parse(event.data);
                if (jsonData.type === "pong") {
                    lastPongTime = Date.now();
                    return;
                }
            } catch (e) {
                // Not JSON or not a control message, treating as regular terminal output
            }

            terminal.write(event.data);
        }
    }

    function handleSocketClose(event: CloseEvent) {
        isConnected = false;

        // Try to reconnect if not deliberately disconnected and under max attempts
        if (!event.wasClean && reconnectAttempts < maxReconnectAttempts) {
            reconnectAttempts++;

            if (terminal) {
                terminal.writeln(
                    `\r\n\x1b[33mConnection lost. Attempting to reconnect (${reconnectAttempts}/${maxReconnectAttempts})...\x1b[0m`,
                );
            }

            setTimeout(() => {
                if (socket) {
                    socket.close();
                    socket = null;
                }
                connectToServer();
            }, 2000);
        } else {
            if (terminal) {
                terminal.writeln("\r\n\x1b[31mDisconnected from server\x1b[0m");
            }

            if (reconnectAttempts >= maxReconnectAttempts) {
                errorMessage =
                    "Connection lost. Max reconnection attempts reached.";
                toast.error(errorMessage);
            }
        }
    }

    function handleSocketError(event: Event) {
        console.error("WebSocket error event:", event);

        if (terminal) {
            terminal.writeln("\r\n\x1b[31mConnection error\x1b[0m");
        }

        errorMessage = "Connection error";
        isConnecting = false;
    }

    function handleTerminalData(event: CustomEvent<{ data: string }>) {
        const { data } = event.detail;

        // Handle special keys
        if (data === "\r") {
            if (currentCommand.trim() !== "") {
                commandHistory.push(currentCommand);
                commandHistory = commandHistory.slice(-100);
                historyIndex = -1;
            }

            currentCommand = "";
            sendKeystrokeViaWebSocket(data);
        } else if (data === "\x1b[A") {
            // Up arrow - navigate history
            if (commandHistory.length > 0) {
                if (historyIndex < commandHistory.length - 1) {
                    historyIndex++;
                }
                currentCommand =
                    commandHistory[commandHistory.length - 1 - historyIndex];
            }

            sendKeystrokeViaWebSocket(data);
        } else if (data === "\x1b[B") {
            // Down arrow - navigate history
            if (historyIndex > 0) {
                historyIndex--;
                currentCommand =
                    commandHistory[commandHistory.length - 1 - historyIndex];
            } else if (historyIndex === 0) {
                historyIndex = -1;
                currentCommand = "";
            }

            sendKeystrokeViaWebSocket(data);
        } else if (data === "\b" || data === "\x7f") {
            // Backspace - remove last character
            currentCommand = currentCommand.slice(0, -1);
            sendKeystrokeViaWebSocket(data);
        } else if (data.length === 1 && data.charCodeAt(0) >= 32) {
            // Printable character
            currentCommand += data;

            sendKeystrokeViaWebSocket(data);
        } else {
            sendKeystrokeViaWebSocket(data);
        }
    }

    function sendKeystrokeViaWebSocket(key: string) {
        if (!isConnected || !socket || socket.readyState !== WebSocket.OPEN) {
            return;
        }

        try {
            socket.send(
                JSON.stringify({
                    type: "input",
                    data: key,
                }),
            );
        } catch (error) {
            console.error("Failed to send keystroke via WebSocket:", error);

            // Fall back to HTTP if WebSocket fails
            sendCommand(key);
        }
    }

    async function sendCommand(command: string) {
        if (!isConnected || !connectionId) return;

        try {
            await pb.send("/api/ssh/execute", {
                method: "POST",
                body: JSON.stringify({
                    connection_id: connectionId,
                    command: command,
                }),
            });
        } catch (error) {
            console.error("Failed to send command:", error);
            if (terminal) {
                terminal.writeln("\r\n\x1b[31mFailed to send command\x1b[0m");
            }
        }
    }

    function sendWebSocketPing() {
        if (!isConnected || !connectionId) {
            return;
        }

        // If we haven't received a pong in over 40 seconds and we've sent at least two pings,
        // the connection is likely dead despite the socket appearing open
        const pongTimeout = 40000;
        if (lastPongTime > 0 && Date.now() - lastPongTime > pongTimeout) {
            console.warn(
                "WebSocket connection appears dead (no pong response)",
            );
            sendHttpPing().catch(() => {
                if (socket) {
                    socket.close();
                    socket = null;
                }
                isConnected = false;
                connectToServer();
            });
            return;
        }

        if (socket && socket.readyState === WebSocket.OPEN) {
            try {
                socket.send(
                    JSON.stringify({
                        type: "ping",
                        time: Date.now(),
                    }),
                );
                return;
            } catch (error) {
                console.warn(
                    "Failed to send WebSocket ping, trying HTTP fallback:",
                    error,
                );
            }
        }

        sendHttpPing().catch((error) => {
            console.error("HTTP fallback ping also failed:", error);
        });
    }

    async function sendHttpPing(): Promise<void> {
        if (!isConnected || !connectionId) {
            return Promise.reject("Not connected");
        }

        try {
            const response = await pb.send("/api/ssh/ping", {
                method: "POST",
                body: JSON.stringify({
                    connection_id: connectionId,
                }),
            });

            if (response && response.status === "ok") {
                lastPongTime = Date.now();
            }

            return Promise.resolve();
        } catch (error) {
            console.error("Failed to send HTTP ping:", error);
            return Promise.reject(error);
        }
    }

    function disconnectFromServer() {
        if (socket) {
            socket.close();
            socket = null;
        }

        if (connectionId) {
            console.log("Disconnecting from server id", connectionId);
            pb.send("/api/ssh/disconnect", {
                method: "POST",
                body: JSON.stringify({
                    connection_id: connectionId,
                }),
            }).catch((error) => {
                console.error("Failed to disconnect properly:", error);
            });
        }

        isConnected = false;
        onDisconnect();
    }

    function handleTerminalResize(
        event: CustomEvent<{ cols: number; rows: number }>,
    ) {
        const { cols, rows } = event.detail;

        if (isConnected && connectionId) {
            pb.send("/api/ssh/resize", {
                method: "POST",
                body: JSON.stringify({
                    connection_id: connectionId,
                    cols,
                    rows,
                }),
            }).catch((error) => {
                console.error("Failed to resize terminal:", error);
            });
        }
    }
</script>

<div class="ssh-terminal-container">
    <div class="terminal-header">
        <div class="server-info">
            <span class="server-name">{serverName}</span>
            <span class="connection-status" class:connected={isConnected}>
                {isConnected
                    ? "Connected"
                    : isConnecting
                      ? "Connecting..."
                      : "Disconnected"}
            </span>
        </div>
        <div class="terminal-actions">
            <button
                class="disconnect-button"
                on:click={disconnectFromServer}
                disabled={!isConnected && !isConnecting}
                title="Disconnect"
            >
                <XCircle size={18} />
            </button>
        </div>
    </div>

    <div class="terminal-body">
        {#if isConnecting}
            <div class="terminal-overlay connecting">
                <Loader2 size={32} class="animate-spin" />
                <span>Connecting to {serverName}...</span>
            </div>
        {:else if errorMessage}
            <div class="terminal-overlay error">
                <XCircle size={32} />
                <span>{errorMessage}</span>
                <button class="retry-button" on:click={connectToServer}>
                    Retry Connection
                </button>
            </div>
        {/if}

        <Terminal
            bind:this={terminal}
            on:data={handleTerminalData}
            on:resize={handleTerminalResize}
            initialText={`Connecting to ${serverName}...\r\n`}
        />
    </div>
</div>

<style>
    .ssh-terminal-container {
        display: flex;
        flex-direction: column;
        height: 100%;
        border-radius: 0.5rem;
        overflow: hidden;
        border: 1px solid var(--border, #333);
        background-color: var(--background, #1a1a1a);
        max-width: 100%;
        position: relative;
    }

    .terminal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.5rem 1rem;
        background-color: var(--header-bg, #252525);
        border-bottom: 1px solid var(--border, #333);
        flex-wrap: wrap;
    }

    .server-info {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        flex-wrap: wrap;
    }

    .server-name {
        font-weight: 600;
        color: var(--header-text, #f0f0f0);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 200px;
    }

    .connection-status {
        font-size: 0.75rem;
        padding: 0.25rem 0.5rem;
        border-radius: 1rem;
        background-color: var(--status-bg, #f44336);
        color: var(--status-text, white);
        white-space: nowrap;
    }

    .connection-status.connected {
        background-color: var(--status-connected-bg, #4caf50);
    }

    .terminal-actions {
        display: flex;
        gap: 0.5rem;
    }

    .disconnect-button {
        background: none;
        border: none;
        color: var(--button-color, #e0e0e0);
        cursor: pointer;
        padding: 0.25rem;
        border-radius: 0.25rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .disconnect-button:hover {
        background-color: var(--button-hover, rgba(255, 255, 255, 0.1));
    }

    .disconnect-button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .terminal-body {
        flex: 1;
        position: relative;
        background: var(--terminal-bg, #1e1e1e);
        min-height: 300px;
        height: 100%;
        max-height: calc(100vh - 140px);
        overflow: hidden;
    }

    .terminal-overlay {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        background-color: rgba(0, 0, 0, 0.7);
        z-index: 10;
        color: white;
        gap: 1rem;
        padding: 1rem;
        text-align: center;
    }

    .terminal-overlay span {
        max-width: 100%;
        word-break: break-word;
    }

    .retry-button {
        background-color: var(--primary, #61afef);
        color: white;
        border: none;
        padding: 0.5rem 1rem;
        border-radius: 0.25rem;
        cursor: pointer;
        font-weight: 500;
        margin-top: 1rem;
    }

    .retry-button:hover {
        background-color: var(--primary-hover, #4d8ebd);
    }

    /* Mobile responsiveness */
    @media (max-width: 768px) {
        .ssh-terminal-container {
            border-radius: 0.25rem;
        }

        .terminal-header {
            padding: 0.25rem 0.5rem;
        }

        .server-name {
            max-width: 140px;
            font-size: 0.9rem;
        }

        .connection-status {
            font-size: 0.7rem;
            padding: 0.15rem 0.35rem;
        }

        .terminal-body {
            min-height: 200px;
        }
    }
</style>
