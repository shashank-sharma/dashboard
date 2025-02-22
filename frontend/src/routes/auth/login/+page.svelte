<script lang="ts">
	import { currentUser, pb } from "$lib/pocketbase";
	import { Button } from "$lib/components/ui/button";
	import { Input } from "$lib/components/ui/input";
	import {
		Card,
		CardContent,
		CardDescription,
		CardFooter,
		CardHeader,
		CardTitle,
	} from "$lib/components/ui/card";
	import { Label } from "$lib/components/ui/label";
	import { Alert, AlertDescription } from "$lib/components/ui/alert";
	import { Loader2 } from "lucide-svelte";
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";

	onMount(() => {
		// Check if user is already logged in
		if (pb.authStore.isValid) {
			goto("/dashboard"); // or wherever your home page is
		}
	});

	let username: string = "";
	let password: string = "";
	let isLoading: boolean = false;
	let error: string | null = null;

	async function login() {
		if (isLoading) return;
		isLoading = true;
		error = null;
		try {
			await pb.collection("users").authWithPassword(username, password);
		} catch (err) {
			console.error(err);
			error = "Invalid username or password";
		} finally {
			isLoading = false;
		}
	}

	async function signUp() {
		if (isLoading) return;
		isLoading = true;
		error = null;
		try {
			const data = {
				username,
				password,
				passwordConfirm: password,
				name: username,
			};
			await pb.collection("users").create(data);
			await login();
		} catch (err) {
			console.error(err);
			error = "Error creating account. Username might be taken.";
		} finally {
			isLoading = false;
		}
	}

	function signOut() {
		pb.authStore.clear();
	}
</script>

<div
	class="flex flex-col items-center justify-center min-h-screen bg-background"
>
	{#if $currentUser}
		<Card class="w-[350px]">
			<CardHeader>
				<CardTitle>Welcome, {$currentUser.username}!</CardTitle>
				<CardDescription>You are currently signed in.</CardDescription>
			</CardHeader>
			<CardFooter>
				<Button on:click={signOut} variant="outline" class="w-full"
					>Sign Out</Button
				>
			</CardFooter>
		</Card>
	{:else}
		<Card class="w-[350px]">
			<CardHeader>
				<CardTitle>Welcome</CardTitle>
				<CardDescription
					>Sign in to your account or create a new one.</CardDescription
				>
			</CardHeader>
			<CardContent>
				<form on:submit|preventDefault>
					<div class="grid w-full items-center gap-4">
						<div class="flex flex-col space-y-1.5">
							<Label for="username">Username</Label>
							<Input
								id="username"
								bind:value={username}
								disabled={isLoading}
							/>
						</div>
						<div class="flex flex-col space-y-1.5">
							<Label for="password">Password</Label>
							<Input
								id="password"
								type="password"
								bind:value={password}
								disabled={isLoading}
							/>
						</div>
					</div>
				</form>
				{#if error}
					<Alert variant="destructive" class="mt-4">
						<AlertDescription>{error}</AlertDescription>
					</Alert>
				{/if}
			</CardContent>
			<CardFooter class="flex justify-between">
				<Button
					on:click={signUp}
					variant="outline"
					disabled={isLoading}
				>
					{#if isLoading}
						<Loader2 class="mr-2 h-4 w-4 animate-spin" />
					{/if}
					Sign Up
				</Button>
				<Button on:click={login} disabled={isLoading}>
					{#if isLoading}
						<Loader2 class="mr-2 h-4 w-4 animate-spin" />
					{/if}
					Login
				</Button>
			</CardFooter>
		</Card>
	{/if}
</div>
