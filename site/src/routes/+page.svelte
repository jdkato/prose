<script lang="ts">
	import { resolve } from '$app/paths';

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
		content="Tokenization, sentence segmentation, part-of-speech tagging, named entities and readability for Go. Pure Go, no cgo, no services."
	/>
</svelte:head>

<section class="py-10">
	<h1 class="max-w-2xl text-4xl font-bold tracking-tight sm:text-5xl">
		Natural language processing for Go.
	</h1>
	<p class="mt-5 max-w-2xl text-lg text-muted">
		Pure Go, no cgo, no external services. Every token, sentence and entity knows its byte offset
		in the text you gave it.
	</p>
	<div class="mt-8 flex flex-wrap gap-3">
		<a
			href={resolve('/docs/[slug]', { slug: 'getting-started' })}
			class="rounded-md bg-accent px-4 py-2 font-medium text-black hover:opacity-90"
			>Get started</a
		>
		<a
			href={resolve('/playground')}
			class="rounded-md border border-line px-4 py-2 font-medium hover:bg-accent-soft"
			>Try the playground</a
		>
	</div>
	<pre
		class="mt-10 max-w-2xl overflow-x-auto rounded-md border border-line bg-card p-4 font-mono text-sm"><code
			>go get github.com/jdkato/prose/v3</code
		></pre>
</section>

<section class="grid gap-4 py-10 sm:grid-cols-2 lg:grid-cols-3">
	{#each features as [name, blurb] (name)}
		<div class="rounded-lg border border-line bg-card p-5">
			<h2 class="font-semibold">{name}</h2>
			<p class="mt-1 text-sm text-muted">{blurb}</p>
		</div>
	{/each}
</section>

<section class="grid gap-10 py-10 lg:grid-cols-2">
	<div>
		<h2 class="text-xl font-semibold">Pay only for what you use</h2>
		<p class="mt-3 text-muted">
			Each model is its own package. Import the tagger alone and the entity model never enters your
			binary. Other languages are packages too, so a program links only the ones it imports.
		</p>
		<a href={resolve('/docs/[slug]', { slug: 'binary-size' })} class="mt-3 inline-block text-sm text-accent underline"
			>How it works</a
		>
	</div>
	<div>
		<h2 class="text-xl font-semibold">Twelve languages</h2>
		<p class="mt-3 text-muted">
			Tagging and segmentation models for Danish, Dutch, English, Finnish, French, German, Italian,
			Norwegian, Portuguese, Spanish, Swedish and Turkish, trained on Universal Dependencies.
		</p>
		<a href={resolve('/models')} class="mt-3 inline-block text-sm text-accent underline"
			>Training data and accuracy</a
		>
	</div>
</section>
