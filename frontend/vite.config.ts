import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	css: {
		postcss: true
	},
	optimizeDeps: {
		exclude: ['fsevents']
	},
	server: {
        host: true
    }
});
