<script lang="ts">
	import { onMount } from 'svelte';
	import {
		load,
		analyze,
		convertCase,
		type Analysis,
		type Cases,
		type Entity,
		type Span,
		type Token
	} from '#lib/prose';
	import { highlightGo } from '#lib/go-highlight';

	const sample = `Marie Curie studied in Paris. She won the Nobel Prize twice, and the Curie Institute still carries her name. Institutional recognition of her contributions nevertheless arrived comparatively late.`;

	let text = $state(sample);
	let title = $state('the quick brown fox jumps over the lazy dog');
	let status = $state<'loading' | 'ready' | 'failed'>('loading');
	let failure = $state('');
	let result = $state<Analysis | null>(null);
	let cases = $state<Cases | null>(null);
	let elapsed = $state(0);
	let copied = $state(false);

	// What is drawn over the text.
	let layers = $state({ tags: true, entities: true, sentences: false, readability: false });
	let hidden = $state(new Set<string>());

	const layerNames = [
		['tags', 'Tags'],
		['entities', 'Entities'],
		['sentences', 'Sentences'],
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

	// --- The rendered view ---------------------------------------------------
	//
	// The text is drawn as its sentences, each a list of groups, each group
	// a run of tokens and the whitespace between them, wrapped in an entity
	// or not. Offsets are bytes, so slicing goes through the encoded text.

	interface Chunk {
		text?: string;
		token?: Token;
	}
	interface Group {
		entity?: Entity;
		chunks: Chunk[];
	}
	interface SentenceView {
		span?: Span;
		groups: Group[];
	}

	const views = $derived.by((): SentenceView[] => {
		if (!result) return [];
		const bytes = new TextEncoder().encode(text);
		const dec = new TextDecoder();
		const slice = (a: number, b: number) => dec.decode(bytes.slice(a, b));
		const tokens = result.tokens;
		const entities = result.entities;

		let ti = 0;
		function view(a: number, b: number, span?: Span): SentenceView | null {
			const groups: Group[] = [];
			let group: Group = { chunks: [] };
			let pos = a;
			for (; ti < tokens.length && tokens[ti].end <= b; ti++) {
				const tok = tokens[ti];
				if (tok.start < a) continue;
				// Whitespace stays inside an entity only between its own
				// tokens; at a boundary it goes to the next group.
				const ent = entities.find((e) => tok.start >= e.start && tok.end <= e.end);
				if (ent !== group.entity) {
					if (group.chunks.length) groups.push(group);
					group = { entity: ent, chunks: [] };
				}
				if (tok.start > pos) group.chunks.push({ text: slice(pos, tok.start) });
				group.chunks.push({ token: tok });
				pos = tok.end;
			}
			if (pos < b) group.chunks.push({ text: slice(pos, b) });
			if (group.chunks.length) groups.push(group);
			return groups.length ? { span, groups } : null;
		}

		const out: SentenceView[] = [];
		let cursor = 0;
		for (const s of result.sentences) {
			if (s.start > cursor) {
				const v = view(cursor, s.start);
				if (v) out.push(v);
			}
			const v = view(s.start, s.end, s);
			if (v) out.push(v);
			cursor = s.end;
		}
		if (cursor < bytes.length) {
			const v = view(cursor, bytes.length);
			if (v) out.push(v);
		}
		return out;
	});

	const legend = $derived.by(() => {
		const counts = new Map<string, number>();
		for (const e of result?.entities ?? []) counts.set(e.label, (counts.get(e.label) ?? 0) + 1);
		return [...counts].sort((a, b) => b[1] - a[1]);
	});

	function toggleLabel(label: string) {
		const next = new Set(hidden);
		if (next.has(label)) next.delete(label);
		else next.add(label);
		hidden = next;
	}

	const difficulty = (grade: number) => (grade > 12 ? 'harder' : grade > 8 ? 'hard' : '');

	const scores = $derived(
		result
			? ([
					['Flesch reading ease', result.readability.fleschReadingEase],
					['Gunning fog', result.readability.gunningFog],
					['SMOG', result.readability.smog],
					['Coleman–Liau', result.readability.colemanLiau],
					['Automated readability', result.readability.automatedReadability],
					['Dale–Chall', result.readability.daleChall],
					['LIX', result.readability.lix]
				] as [string, number][])
			: []
	);

	// --- The code panel ------------------------------------------------------

	function goString(s: string): string {
		if (s.includes('\n') && !s.includes('`')) return '`' + s + '`';
		return JSON.stringify(s);
	}

	function comments(lines: string[], indent: string): string {
		const shown = lines.slice(0, 3).map((l) => `${indent}// ${l.length > 60 ? l.slice(0, 57) + '…' : l}`);
		if (lines.length > 3) shown.push(`${indent}// …`);
		return shown.join('\n');
	}

	const code = $derived.by(() => {
		const r = result;
		const imports = ['"fmt"', '', '"github.com/jdkato/prose/v3"'];
		if (layers.readability) imports.push('"github.com/jdkato/prose/v3/summarize"');
		const body: string[] = [];
		body.push(`\ttext := ${goString(text)}`);
		body.push('', '\tdoc, err := prose.NewDocument(text)', '\tif err != nil {', '\t\tpanic(err)', '\t}');

		if (layers.sentences) {
			body.push('', '\tfor _, sent := range doc.Sentences() {', '\t\tfmt.Println(sent.Text)');
			if (r) body.push(comments(r.sentences.map((s) => s.text), '\t\t'));
			body.push('\t}');
		}
		if (layers.tags || (!layers.entities && !layers.sentences && !layers.readability)) {
			body.push('', '\tfor _, tok := range doc.Tokens() {');
			body.push(layers.tags ? '\t\tfmt.Println(tok.Text, tok.Tag)' : '\t\tfmt.Println(tok.Text)');
			if (r) body.push(comments(r.tokens.map((t) => (layers.tags ? `${t.text} ${t.tag}` : t.text)), '\t\t'));
			body.push('\t}');
		}
		if (layers.entities) {
			body.push('', '\tfor _, ent := range doc.Entities() {', '\t\tfmt.Println(ent.Text, ent.Label)');
			if (r) body.push(comments(r.entities.map((e) => `${e.text} ${e.label}`), '\t\t'));
			body.push('\t}');
		}
		if (layers.readability) {
			body.push('', '\td := summarize.NewDocument(text)', '\tfmt.Println(d.FleschKincaid())');
			if (r) body.push(`\t// ${r.readability.fleschKincaid.toFixed(1)}`);
		}

		return [
			'package main',
			'',
			'import (',
			...imports.map((i) => (i ? `\t${i}` : '')),
			')',
			'',
			'func main() {',
			...body,
			'}',
			''
		].join('\n');
	});

	async function copy() {
		try {
			await navigator.clipboard.writeText(code);
			copied = true;
			setTimeout(() => (copied = false), 1500);
		} catch {
			// Clipboard access can be refused; the code is still selectable.
		}
	}
</script>

<div class="grid gap-8 lg:grid-cols-[minmax(0,3fr)_minmax(0,2fr)]">
	<div class="min-w-0">
		<textarea
			bind:value={text}
			oninput={schedule}
			rows="4"
			aria-label="Text"
			class="w-full rounded-md border border-line bg-card p-3 font-mono text-sm"
			spellcheck="false"
		></textarea>

		<div class="mt-3 flex flex-wrap items-center gap-x-4 gap-y-2 text-sm">
			{#each layerNames as [key, label] (key)}
				<label class="flex items-center gap-1.5">
					<input type="checkbox" bind:checked={layers[key]} class="accent-accent" />
					{label}
				</label>
			{/each}
			{#if layers.entities && legend.length}
				<span class="legend ml-auto flex flex-wrap gap-2">
					{#each legend as [label, n] (label)}
						<button
							onclick={() => toggleLabel(label)}
							class="ent ent-{label} text-xs {hidden.has(label) ? 'opacity-40' : ''}"
							aria-pressed={!hidden.has(label)}
						>
							{n}<span class="chip">{label}</span>
						</button>
					{/each}
				</span>
			{/if}
		</div>

		<p class="mt-3 h-5 text-xs text-muted">
			{#if status === 'loading'}
				Loading prose.wasm (about 15 MB) …
			{:else if status === 'failed'}
				<span class="text-red-600">Could not load the playground: {failure}</span>
			{:else if failure}
				<span class="text-red-600">{failure}</span>
			{:else if result}
				{result.tokens.length} tokens · {result.sentences.length} sentences ·
				{result.entities.length} entities
				{#if layers.readability}
					· grade {result.readability.fleschKincaid.toFixed(1)}
				{/if}
				· {elapsed.toFixed(1)} ms
			{/if}
		</p>

		{#if result}
			<div
				class="rendered mt-2 rounded-md border border-line bg-card px-4 py-3 text-[0.95rem]"
				class:tags={layers.tags}
			>
				{#each views as v, i (i)}
					{#if layers.sentences && v.span}<span class="boundary"></span>{/if}
					<span
						class="sentence"
						class:hard={layers.readability && v.span && difficulty(v.span.grade) === 'hard'}
						class:harder={layers.readability && v.span && difficulty(v.span.grade) === 'harder'}
						title={layers.readability && v.span ? `grade ${v.span.grade.toFixed(1)}` : undefined}
					>
						{#each v.groups as g, j (j)}
							{#snippet chunks()}
								{#each g.chunks as c, k (k)}
									{#if c.token}
										<span class="tok">
											<span>{c.token.text}</span>
											{#if layers.tags}<span class="tag">{c.token.tag}</span>{/if}
										</span>
									{:else}{c.text}{/if}
								{/each}
							{/snippet}
							{#if g.entity && layers.entities && !hidden.has(g.entity.label)}
								<mark class="ent ent-{g.entity.label}">{@render chunks()}<span class="chip">{g.entity.label}</span></mark>
							{:else}
								{@render chunks()}
							{/if}
						{/each}
					</span>
				{/each}
			</div>

			{#if layers.readability}
				<details class="mt-3 text-sm">
					<summary class="cursor-pointer text-muted">
						Sentences above grade 8 are tinted, above 12 more so. All scores
					</summary>
					<table class="doc mt-2 max-w-sm">
						<tbody>
							<tr>
								<td>Flesch–Kincaid grade</td>
								<td class="text-right font-mono">{result.readability.fleschKincaid.toFixed(1)}</td>
							</tr>
							{#each scores as [name, n] (name)}
								<tr><td>{name}</td><td class="text-right font-mono">{n.toFixed(1)}</td></tr>
							{/each}
							<tr>
								<td>Words / sentences / syllables</td>
								<td class="text-right font-mono"
									>{result.readability.words} / {result.readability.sentences} / {result.readability
										.syllables}</td
								>
							</tr>
						</tbody>
					</table>
				</details>
			{/if}
		{/if}
	</div>

	<aside class="min-w-0 lg:sticky lg:top-6 lg:self-start">
		<div class="flex items-center justify-between">
			<h2 class="text-sm font-medium">The same thing in Go</h2>
			<button onclick={copy} class="text-xs text-muted hover:text-fg">{copied ? 'Copied' : 'Copy'}</button>
		</div>
		<pre class="gocode mt-2 rounded-md border border-line bg-card p-4">{@html highlightGo(code)}</pre>
	</aside>
</div>

<section class="mt-12 max-w-2xl">
	<h2 class="text-lg font-semibold">Case conversion</h2>
	<p class="mt-1 text-sm text-muted">
		<code class="rounded bg-accent-soft px-1">strcase</code> writes a string in a style, using the tagger to tell
		which words stay lowercase.
	</p>
	<input
		bind:value={title}
		oninput={schedule}
		aria-label="Title"
		class="mt-3 w-full rounded-md border border-line bg-card p-2 font-mono text-sm"
		spellcheck="false"
	/>
	{#if cases}
		<table class="doc mt-3 w-full text-sm">
			<tbody>
				<tr><td class="w-32 text-muted">AP</td><td>{cases.ap}</td></tr>
				<tr><td class="text-muted">Chicago</td><td>{cases.chicago}</td></tr>
				<tr><td class="text-muted">Sentence</td><td>{cases.sentence}</td></tr>
			</tbody>
		</table>
	{/if}
</section>
