<script lang="ts">
	import {
		Toaster as Sonner,
		type ToasterProps as SonnerProps,
	} from "svelte-sonner";
	import { mode } from "mode-watcher";
	import { browser } from "$app/environment";
	import { onMount } from "svelte";

	type $$Props = SonnerProps;

	let isMobile = false;
	let position:
		| "top-right"
		| "top-left"
		| "top-center"
		| "bottom-right"
		| "bottom-left"
		| "bottom-center" = "bottom-right";

	// Function to check if we're on mobile based on window width
	function checkMobileView() {
		isMobile = window.innerWidth < 768;
		position = isMobile ? "top-center" : "bottom-right";
	}

	onMount(() => {
		if (browser) {
			checkMobileView();
			window.addEventListener("resize", checkMobileView);

			return () => {
				window.removeEventListener("resize", checkMobileView);
			};
		}
	});
</script>

<Sonner
	theme={$mode}
	class="toaster group"
	{position}
	toastOptions={{
		classes: {
			toast: "group toast group-[.toaster]:bg-background group-[.toaster]:text-foreground group-[.toaster]:border-border group-[.toaster]:shadow-lg",
			description: "group-[.toast]:text-muted-foreground",
			actionButton:
				"group-[.toast]:bg-primary group-[.toast]:text-primary-foreground",
			cancelButton:
				"group-[.toast]:bg-muted group-[.toast]:text-muted-foreground",
		},
	}}
	{...$$restProps}
/>
