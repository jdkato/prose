<script lang="ts">
	import { onMount } from 'svelte';
	import { load, analyze, convertCase, type Analysis, type Cases } from '#lib/prose';

	const sample = `Marie Curie studied in Paris. She won the Nobel Prize twice, and the Curie Institute still carries her name.`;

	let text = $state(sample);
	let title = $state('the quick brown fox jumps over the lazy dog');
	let status = $state<'loading' | 'ready' | 'failed'>('loading');
	let failure = $state('');
	let result = $state<Analysis | null>(null);
	let cases = $state<Cases | null>(null);
	let view = $state<'tokens' | 'sentences' | 'entities' | 'readability'>('tokens');
	let elapsed = $state(0);

	const views = [
		['tokens', 'Tokens'],
		['sentences', 'Sentences'],
		['entities', 'Entities'],
		['readability', 'Readability']
	] as const;

	function run() {
		if (status !== 'ready') return;
		const t0 = performance.now();
		try {
			result = analyze(text);
			cases = convertCase(title);
			failure = '';
		} catch (e) {
			failure = (e as Error).message;
		}
		elapsed = performance.now() - t0;
	}

	let timer: ReturnType<typeof setTimeout>;
	function schedule() {
		clearTimeout(timer);
		timer = setTimeout(run, 150);
	}

	onMount(async () => {
		try {
			await load();
			status = 'ready';
			run();
		} catch (e) {
			status = 'failed';
			failure = (e as Error).message;
		}
	});

	// Entities as inline marks over the source text.
	const marked = $derived.by(() => {
		if (!result) return [];
		const bytes = new TextEncoder().encode(text);
		const dec = new TextDecoder();
		const parts: { text: string; label?: string }[] = [];
		let at = 0;
		for (const e of result.entities) {
			if (e.start > at) parts.push({ text: dec.decode(bytes.slice(at, e.start)) });
			parts.push({ text: e.text, label: e.label });
			at = e.end;
		}
		if (at < bytes.length) parts.push({ text: dec.decode(bytes.slice(at)) });
		return parts;
	});

	const scores: [string, number][] = $derived(
		result
			? [
					['Flesch–Kincaid grade', result.readability.fleschKincaid],
					['Flesch reading ease', result.readability.fleschReadingEase],
					['Gunning fog', result.readability.gunningFog],
					['SMOG', result.readability.smog],
					['Coleman–Liau', result.readability.colemanLiau],
					['Automated readability', result.readability.automatedReadability],
					['Dale–Chall', result.readability.daleChall],
					['LIX', result.readability.lix]
				]
			: []
	);

	const counts: [string, number][] = $derived(
		result
			? [
					['Words', result.readability.words],
					['Sentences', result.readability.sentences],
					['Syllables', result.readability.syllables],
					['Characters', result.readability.characters]
				]
			: []
	);
</script>

<svelte:head>
	<title>Playground · prose</title>
	<meta name="description" content="Run prose in your browser: tokens, sentences, entities and readability for any English text." />
</svelte:head>

<div class="max-w-4xl">
	<h1 class="text-3xl font-bold tracking-tight">Playground</h1>
	<p class="mt-2 text-muted">
		The same Go code, compiled to WebAssembly and run in your browser. English only; nothing
		leaves the page.
	</p>

	<label class="mt-6 block">
		<span class="text-sm font-medium">Text</span>
		<textarea
			bind:value={text}
			oninput={schedule}
			rows="5"
			class="mt-1 w-full rounded-md border border-line bg-card p-3 font-mono text-sm"
			spellcheck="false"
		></textarea>
	</label>

	<p class="mt-2 h-5 text-xs text-muted">
		{#if status === 'loading'}
			Loading prose.wasm (about 15 MB) …
		{:else if status === 'failed'}
			<span class="text-red-600">Could not load the playground: {failure}</span>
		{:else if failure}
			<span class="text-red-600">{failure}</span>
		{:else if result}
			{result.tokens.length} tokens, {result.sentences.length} sentences, {result.entities.length}
			entities in {elapsed.toFixed(1)} ms
		{/if}
	</p>

	<div class="mt-4 flex flex-wrap gap-1 border-b border-line text-sm">
		{#each views as [key, label] (key)}
			<button
				onclick={() => (view = key)}
				class="-mb-px border-b-2 px-3 py-2 {view === key
					? 'border-accent font-medium'
					: 'border-transparent text-muted hover:text-fg'}"
			>
				{label}
			</button>
		{/each}
	</div>

	{#if result}
		<div class="doc mt-4">
			{#if view === 'tokens'}
				<div class="max-h-[32rem] overflow-auto rounded-md border border-line">
					<table>
						<thead>
							<tr><th>Text</th><th>Tag</th><th>Label</th><th>Start</th><th>End</th></tr>
						</thead>
						<tbody>
							{#each result.tokens as t, i (i)}
								<tr>
									<td><code>{t.text}</code></td>
									<td>{t.tag}</td>
									<td class="text-muted">{t.label ?? ''}</td>
									<td class="font-mono text-muted">{t.start}</td>
									<td class="font-mono text-muted">{t.end}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{:else if view === 'sentences'}
				<ol>
					{#each result.sentences as s, i (i)}
						<li>
							{s.text}
							<span class="font-mono text-xs text-muted">[{s.start}:{s.end}]</span>
						</li>
					{/each}
				</ol>
			{:else if view === 'entities'}
				<p class="rounded-md border border-line bg-card p-4 leading-8">
					{#each marked as part, i (i)}
						{#if part.label}
							<mark class="rounded bg-accent-soft px-1 text-fg">
								{part.text}
								<span class="ml-1 align-super text-[0.6rem] font-semibold text-accent"
									>{part.label}</span
								>
							</mark>
						{:else}{part.text}{/if}
					{/each}
				</p>
				{#if result.entities.length === 0}
					<p class="text-muted">No entities found.</p>
				{/if}
			{:else}
				<div class="grid gap-6 sm:grid-cols-2">
					<table>
						<tbody>
							{#each counts as [name, n] (name)}
								<tr><td>{name}</td><td class="text-right font-mono">{n}</td></tr>
							{/each}
						</tbody>
					</table>
					<table>
						<tbody>
							{#each scores as [name, n] (name)}
								<tr><td>{name}</td><td class="text-right font-mono">{n.toFixed(2)}</td></tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</div>
	{/if}

	<h2 class="mt-12 text-xl font-semibold">Case conversion</h2>
	<label class="mt-3 block">
		<span class="text-sm font-medium">Title</span>
		<input
			bind:value={title}
			oninput={schedule}
			class="mt-1 w-full rounded-md border border-line bg-card p-2 font-mono text-sm"
			spellcheck="false"
		/>
	</label>
	{#if cases}
		<table class="doc mt-3 w-full text-sm">
			<tbody>
				<tr><td class="w-32 text-muted">AP</td><td>{cases.ap}</td></tr>
				<tr><td class="text-muted">Chicago</td><td>{cases.chicago}</td></tr>
				<tr><td class="text-muted">Sentence</td><td>{cases.sentence}</td></tr>
			</tbody>
		</table>
	{/if}
</div>
