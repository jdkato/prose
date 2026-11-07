<script lang="ts">
	import { resolve } from '$app/paths';
	import { docs } from '#lib/docs';

	let { data } = $props();
	const Content = $derived(data.doc.component);
</script>

<svelte:head>
	<title>{data.doc.title} · prose</title>
	<meta name="description" content={data.doc.summary} />
</svelte:head>

<div class="grid gap-10 lg:grid-cols-[14rem_minmax(0,1fr)]">
	<aside class="lg:sticky lg:top-6 lg:self-start">
		<nav class="flex flex-row flex-wrap gap-1 text-sm lg:flex-col">
			{#each docs as d (d.slug)}
				<a
					href={resolve('/docs/[slug]', { slug: d.slug })}
					aria-current={d.slug === data.doc.slug ? 'page' : undefined}
					class="rounded px-2 py-1 hover:bg-accent-soft {d.slug === data.doc.slug
						? 'bg-accent-soft font-medium'
						: 'text-muted'}"
				>
					{d.title}
				</a>
			{/each}
		</nav>
	</aside>

	<article class="doc min-w-0 max-w-3xl">
		<Content />

		<nav class="mt-12 flex justify-between gap-4 border-t border-line pt-6 text-sm">
			{#if data.prev}
				<a href={resolve('/docs/[slug]', { slug: data.prev.slug })} class="text-muted hover:text-fg"
					>← {data.prev.title}</a
				>
			{:else}<span></span>{/if}
			{#if data.next}
				<a href={resolve('/docs/[slug]', { slug: data.next.slug })} class="text-muted hover:text-fg"
					>{data.next.title} →</a
				>
			{/if}
		</nav>
	</article>
</div>
