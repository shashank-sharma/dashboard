<script lang="ts">
    import { currentUser } from '$lib/pocketbase';
</script>
<nav>
	<a href="/">Home</a>
	<a href="/about">About</a>
	<a href="/settings">Settings</a>
    <h3>{JSON.stringify($currentUser)}</h3>
{#if $currentUser}
	<a href="/logout">{$currentUser.username} Logout</a>
{:else}
	<a href="/login">Login</a>
{/if}
</nav>

<slot></slot>
