<script lang="ts">
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { notifications } from "$lib/notification";
  import { Button, SnackbarContainer } from "attractions";

  let thisSnack: SnackbarContainer;

  onMount(() => {
    const unsubscribe = notifications.subscribe((value) => {
      console.log("Change", value, thisSnack);
      thisSnack.showSnackbar({ props: { text: value } });
    });
  });
</script>

<div class="notification-container">
  <SnackbarContainer let:showSnackbar bind:this={thisSnack}></SnackbarContainer>
</div>

<style>
  .notification-container {
    position: fixed;
    top: 0;
    right: 0;
    padding: 2em;
  }
</style>
