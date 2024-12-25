<script lang="ts">
    import { fade } from "svelte/transition";
</script>

<div
    in:fade={{ duration: 150, delay: 200 }}
    out:fade={{ duration: 150 }}
    class="fixed inset-0 z-50 flex items-center justify-center bg-background"
>
    <div class="flex flex-col items-center gap-2">
        <div class="h-8 w-8 animate-spin">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                class="h-8 w-8"
            >
                <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                />
                <path
                    class="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                />
            </svg>
        </div>
        <p class="text-sm text-muted-foreground">Loading...</p>
    </div>
</div>
