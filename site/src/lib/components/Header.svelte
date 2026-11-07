<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';

	const links = [
		{ href: resolve('/docs/[slug]', { slug: 'getting-started' }), label: 'Docs', prefix: resolve('/docs') },
		{ href: resolve('/models'), label: 'Models', prefix: resolve('/models') }
	];

	const active = (prefix: string) => page.url.pathname.startsWith(prefix);
</script>

<header class="border-b border-line bg-bg/80 backdrop-blur">
	<div class="mx-auto flex h-14 max-w-6xl items-center gap-6 px-5">
		<a href={resolve('/')} class="flex items-center gap-2 font-semibold tracking-tight">
			<span class="inline-block h-5 w-5 rounded bg-accent"></span>
			prose
		</a>
		<nav class="flex items-center gap-4 text-sm">
			{#each links as link (link.href)}
				<a
					href={link.href}
					class="rounded px-2 py-1 hover:bg-accent-soft {active(link.prefix)
						? 'font-medium text-fg'
						: 'text-muted'}"
				>
					{link.label}
				</a>
			{/each}
		</nav>
		<div class="ml-auto flex items-center gap-4 text-sm">
			<a
				href="https://pkg.go.dev/github.com/jdkato/prose/v3"
				class="text-muted hover:text-fg"
				rel="noopener">API</a
			>
			<a href="https://github.com/jdkato/prose" class="text-muted hover:text-fg" rel="noopener"
				>GitHub</a
			>
		</div>
	</div>
</header>
