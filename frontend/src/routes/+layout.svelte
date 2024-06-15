<script lang="ts">
    import { currentUser } from "$lib/pocketbase";
    import Notification from "../components/Notification.svelte";
    import { toast } from "$lib/notification";
    import { Card } from "attractions";
    import Sidebar from "../components/Sidebar.svelte";

    toast("Notification do it");
    setTimeout(() => {
        toast("Hello");
    }, 2000);
</script>

<Notification />

<div class="container">
    <div class="card">
        <Sidebar />
        <div class="card-content">
            <slot />
        </div>
    </div>
</div>

<style>
    @font-face {
        font-family: "Gilroy";
        src: url("/fonts/Gilroy.woff2");
    }
    :global(body) {
        background-color: gray;
        height: 100%;
    }

    :global(html) {
        height: 95%;
        margin: 1em;
        font-family: Gilroy, serif;
    }

    .container {
        display: flex;
        height: 100%;
    }

    .card-content {
        flex: 1;
        background-color: #fff;
        padding: 2em;
        overflow: auto;
    }

    .card {
        display: flex;
        border-radius: 8px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        overflow: hidden;
        width: 100%;
        height: 100%;
        overflow: hidden;
    }
</style>
