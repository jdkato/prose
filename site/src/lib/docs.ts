import type { Component } from 'svelte';

export interface Doc {
	slug: string;
	title: string;
	summary: string;
	order: number;
	component: Component;
}

interface Module {
	default: Component;
	metadata: { title: string; summary?: string; order?: number };
}

const modules = import.meta.glob<Module>('/src/content/docs/*.md', { eager: true });

// Every page under src/content/docs, in sidebar order.
export const docs: Doc[] = Object.entries(modules)
	.map(([path, mod]) => ({
		slug: path.split('/').pop()!.replace(/\.md$/, ''),
		title: mod.metadata.title,
		summary: mod.metadata.summary ?? '',
		order: mod.metadata.order ?? 99,
		component: mod.default
	}))
	.sort((a, b) => a.order - b.order);
