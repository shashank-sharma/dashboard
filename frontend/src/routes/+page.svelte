<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import SettingsModal from "$lib/components/SettingsModal.svelte";
	import {
		showInstallPrompt,
		installPrompt,
		triggerInstallBanner,
	} from "$lib/features/pwa/services";
	import PwaDebugTools from "$lib/features/pwa/components/PwaDebugTools.svelte";
	import PwaInstallBanner from "$lib/features/pwa/components/PwaInstallBanner.svelte";
	import { browser } from "$app/environment";

	let pwaInstallBanner: any;
	let canInstallPwa = false;

	// Function to force show the install banner regardless of dismissal status
	function forceShowInstallBanner() {
		if (browser) {
			console.log(
				"[PWA] Forcing install banner display from layout button",
			);

			sessionStorage.removeItem("pwa-banner-dismissed");

			if (
				pwaInstallBanner &&
				typeof pwaInstallBanner.forceShow === "function"
			) {
				pwaInstallBanner.forceShow();
				return;
			}

			const unsubscribe = installPrompt.subscribe((value) => {
				if (value) {
					console.log("[PWA] Using real prompt");
					showInstallPrompt();
				} else {
					console.log("[PWA] No real prompt, using mock");
					triggerInstallBanner();
				}
			});
			unsubscribe();
		}
		const unsubInstalled = isPwaInstalled.subscribe((value) => {
			canInstallPwa = !value;
		});

		// Listen for force show events from other components
		window.addEventListener(
			"force-show-pwa-banner",
			forceShowInstallBanner,
		);

		return () => {
			unsubInstalled();
			window.removeEventListener(
				"force-show-pwa-banner",
				forceShowInstallBanner,
			);
		};
	}
</script>

<div class="bottom-0 left-0 fixed m-10">
	<SettingsModal />
</div>

<!-- PWA Debug Tools -->
<PwaDebugTools />

{#if browser}
	<PwaInstallBanner bind:this={pwaInstallBanner} />
	<!-- Always available install button (fixed position) -->
	{#if canInstallPwa && !isDebugMode}
		<div class="fixed bottom-4 right-4 z-50">
			<Button
				variant="default"
				size="sm"
				class="shadow-lg gap-1"
				on:click={forceShowInstallBanner}
			>
				<Download class="w-4 h-4" /> Install App
			</Button>
		</div>
	{/if}
{/if}
