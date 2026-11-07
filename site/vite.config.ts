import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';
import { mdsvex } from 'mdsvex';
import { defineConfig } from 'vite';
import mdsvexOptions from './mdsvex.config.js';

// GitHub Pages serves the site under the repository name.
const base = process.env.BASE_PATH as `/${string}` | undefined;

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit({
			preprocess: [mdsvex(mdsvexOptions), vitePreprocess()],
			extensions: ['.svelte', '.md'],
			adapter: adapter({ fallback: '404.html' }),
			paths: { base }
		})
	],
	server: {
		// The models page reads tag/lang/MODELS.md from the repository.
		fs: { allow: ['..'] }
	}
});
