<script>
    import { onMount } from 'svelte';
    
    let isAuthenticated = false;
  
    onMount(() => {
      // Check for a token or any flag in local storage
      isAuthenticated = !!localStorage.getItem('token');
      
      // Redirect to the login page if not authenticated
      if (!isAuthenticated) {
        window.location.href = "/login"; // or wherever your login page is
      }
    });
  </script>
  
  {#if isAuthenticated}
    <slot></slot> <!-- This will display content of child components if authenticated -->
  {/if}