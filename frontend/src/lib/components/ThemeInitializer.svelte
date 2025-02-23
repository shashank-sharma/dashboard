<script lang="ts">
    import { onMount } from "svelte";
    import { browser } from "$app/environment";

    onMount(() => {
        if (!browser) return;

        if (
            localStorage.theme === "dark" ||
            (!("theme" in localStorage) &&
                window.matchMedia("(prefers-color-scheme: dark)").matches)
        ) {
            document.documentElement.classList.add("dark");
        } else {
            document.documentElement.classList.remove("dark");
        }
    });
</script>
