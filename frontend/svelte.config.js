import adapter from '@sveltejs/adapter-static';
import preprocess from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
    kit: {
        adapter: adapter({
            pages: 'build',
            assets: 'build',
            fallback: 'index.html',
            precompress: false,
            strict: true
        }),
        alias: {
            $components: 'src/lib/components',
            $stores: 'src/lib/stores',
            $services: 'src/lib/services',
            $types: 'src/lib/types'
        },
        paths: {
            base: process.env.NODE_ENV === 'production' ? '/dashboard/frontend' : ''
        }
    },
    preprocess: preprocess()
};

export default config;