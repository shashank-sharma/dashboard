<script lang="ts">
    import { pb } from "$lib/pocketbase";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";

    onMount(() => {
        pb.authStore.clear();
        goto("/auth/login");
    });
</script>

<div>
    <h1>Logging out...</h1>
    <p>You will be redirected to the login page.</p>
</div>
