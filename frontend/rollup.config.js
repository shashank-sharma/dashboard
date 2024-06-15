import svelte from 'rollup-plugin-svelte';
import svelteConfig from '/svelte.config.ts'

export default {
  // ...,
  plugins: [
    svelte({
      ...svelteConfig,
      // ...,
    }),
  ],
};
