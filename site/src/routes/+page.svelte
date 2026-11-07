<script lang="ts">
	import { resolve } from '$app/paths';
	import Playground from '#lib/components/Playground.svelte';

	// One run on an Apple M1, Go 1.27, ten copies of the Sherlock fixture
	// (1.08M words): tokens, sentences and tags in 2.9 s.
	const stats = [
		['370k', 'words per second', 'tokens, sentences and tags, one core'],
		['3 + 4 MB', 'models', 'tagger and entity recognizer, linked separately'],
		['12', 'languages', 'tagging and segmentation, one package each'],
		['0', 'cgo', 'pure Go, no services, runs in a browser']
	];

	const features = [
		['Tokenization', 'Words and punctuation, with byte offsets, in any language.'],
		['Segmentation', 'Sentences, via the Punkt algorithm, in 12 languages.'],
		['Tagging', 'Part of speech in 12 languages, from averaged perceptrons.'],
		['Entities', 'People, places and organizations, in English.'],
		['Readability', 'Flesch–Kincaid, Gunning fog, SMOG, Coleman–Liau and more.'],
		['Case conversion', 'Title and sentence case, AP and Chicago styles.']
	];
</script>

<svelte:head>
	<title>prose · natural language processing for Go</title>
	<meta
		name="description"
		content="Tokenization, sentence segmentation, part-of-speech tagging, named entities and readability for Go. Pure Go, no cgo, no services. Try it in your browser."
	/>
</svelte:head>

<section class="pb-8">
	<h1 class="max-w-3xl text-4xl font-bold tracking-tight sm:text-5xl">
		Natural language processing for Go.
	</h1>
	<p class="mt-4 max-w-2xl text-lg text-muted">
		Pure Go, no cgo, no external services. This is the library itself, compiled to WebAssembly
		and running in your browser. Nothing leaves the page.
	</p>
	<dl class="mt-8 grid grid-cols-2 gap-x-8 gap-y-5 sm:grid-cols-4">
		{#each stats as [n, label, detail] (label)}
			<div title={detail}>
				<dt class="text-2xl font-semibold tracking-tight">{n}</dt>
				<dd class="text-sm text-muted">{label}</dd>
			</div>
		{/each}
	</dl>
</section>

<Playground />

<section class="mt-16 flex flex-wrap items-center gap-4 border-t border-line pt-10">
	<pre
		class="overflow-x-auto rounded-md border border-line bg-card px-4 py-2 font-mono text-sm"><code
			>go get github.com/jdkato/prose/v3</code
		></pre>
	<a
		href={resolve('/docs/[slug]', { slug: 'getting-started' })}
		class="rounded-md bg-accent px-4 py-2 font-medium text-black hover:opacity-90"
		>Get started</a
	>
	<a
		href="https://pkg.go.dev/github.com/jdkato/prose/v3"
		class="rounded-md border border-line px-4 py-2 font-medium hover:bg-accent-soft"
		rel="noopener">API reference</a
	>
</section>

<section class="grid gap-4 py-10 sm:grid-cols-2 lg:grid-cols-3">
	{#each features as [name, blurb] (name)}
		<div class="rounded-lg border border-line bg-card p-5">
			<h2 class="font-semibold">{name}</h2>
			<p class="mt-1 text-sm text-muted">{blurb}</p>
		</div>
	{/each}
</section>

<section class="grid gap-10 py-6 lg:grid-cols-3">
	<div>
		<h2 class="text-xl font-semibold">Everything knows where it came from</h2>
		<p class="mt-3 text-muted">
			Every token, sentence and entity carries a byte offset into the text you gave it, and the text
			at that offset is exactly what it says.
		</p>
		<a href={resolve('/docs/[slug]', { slug: 'offsets' })} class="mt-3 inline-block text-sm text-accent underline"
			>Offsets</a
		>
	</div>
	<div>
		<h2 class="text-xl font-semibold">Pay only for what you use</h2>
		<p class="mt-3 text-muted">
			Each model is its own package. Import the tagger alone and the entity model never enters your
			binary. Other languages are packages too, so a program links only the ones it imports.
		</p>
		<a href={resolve('/docs/[slug]', { slug: 'binary-size' })} class="mt-3 inline-block text-sm text-accent underline"
			>Binary size</a
		>
	</div>
	<div>
		<h2 class="text-xl font-semibold">Twelve languages</h2>
		<p class="mt-3 text-muted">
			Tagging and segmentation for Danish, Dutch, English, Finnish, French, German, Italian,
			Norwegian, Portuguese, Spanish, Swedish and Turkish, trained on Universal Dependencies.
		</p>
		<a href={resolve('/models')} class="mt-3 inline-block text-sm text-accent underline"
			>Training data and accuracy</a
		>
	</div>
</section>
