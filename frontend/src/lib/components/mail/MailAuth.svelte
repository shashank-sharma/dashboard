<!-- src/lib/components/mail/MailAuth.svelte -->
<script lang="ts">
    import { onMount } from "svelte";
    import { Button } from "$lib/components/ui/button";
    import {
        Card,
        CardHeader,
        CardTitle,
        CardDescription,
        CardContent,
    } from "$lib/components/ui/card";
    import { Mail, Check, Loader2 } from "lucide-svelte";
    import { mailStore } from "$lib/stores/mailStore";
    import { toast } from "svelte-sonner";

    let authWindow: Window | null = null;
    let pollInterval: number;

    async function handleAuthClick() {
        const url = await mailStore.startAuth();
        if (!url) return;

        // Open auth window
        authWindow = window.open(
            url,
            "Mail Authentication",
            "width=600,height=600",
        );

        // Start polling for authentication status
        pollInterval = window.setInterval(async () => {
            const isAuthed = await mailStore.checkStatus();

            if (isAuthed) {
                clearInterval(pollInterval);
                if (authWindow) {
                    authWindow.close();
                }
                toast.success("Mail authentication successful!");
            }
        }, 2000);
    }

    // Cleanup on component unmount
    onMount(() => {
        mailStore.checkStatus();
        return () => {
            if (pollInterval) {
                clearInterval(pollInterval);
            }
            if (authWindow) {
                authWindow.close();
            }
        };
    });
</script>

<Card class="w-full max-w-md mx-auto">
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <Mail class="h-5 w-5" />
            Mail Integration
        </CardTitle>
        <CardDescription>
            Connect your email account to enable mail sync
        </CardDescription>
    </CardHeader>

    <CardContent>
        {#if $mailStore.isLoading}
            <div class="flex justify-center items-center py-8">
                <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
        {:else if $mailStore.isAuthenticated}
            <div
                class="flex items-center gap-4 text-green-600 dark:text-green-500"
            >
                <Check class="h-6 w-6" />
                <span>Mail sync is active</span>
            </div>
        {:else}
            <div class="space-y-4">
                <p class="text-sm text-muted-foreground">
                    Click below to start the authentication process. You'll be
                    redirected to sign in with your email provider.
                </p>
                <Button
                    class="w-full"
                    on:click={handleAuthClick}
                    disabled={$mailStore.isAuthenticating}
                >
                    {#if $mailStore.isAuthenticating}
                        <Loader2 class="mr-2 h-4 w-4 animate-spin" />
                        Authenticating...
                    {:else}
                        Connect Mail Account
                    {/if}
                </Button>
            </div>
        {/if}
    </CardContent>
</Card>
