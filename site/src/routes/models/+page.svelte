<script lang="ts">
	let { data } = $props();
</script>

<svelte:head>
	<title>Models · prose</title>
	<meta name="description" content="Training data, license, accuracy and size of each prose language model." />
</svelte:head>

<article class="doc max-w-4xl">
	<h1>Language models</h1>
	<p>
		English is built into <code>prose</code>; the English tagger uses the Penn Treebank tag set and
		the entity model is English only. Each language below is a package under
		<code>tag/lang</code> and <code>segment/lang</code>, linked only when imported.
	</p>
	<p>
		Each tagging model is an averaged perceptron trained on the Universal Dependencies treebanks
		named here, emitting the universal tag set. Accuracy is measured on the treebank's own
		development and test sections through the production tagger, so it is what a caller of
		<code>tag.ForLanguage</code> gets on text of that kind.
	</p>

	<div class="overflow-x-auto">
		<table>
			<thead>
				<tr>
					<th>Language</th>
					<th>Code</th>
					<th>Training data</th>
					<th>License</th>
					<th>Dev</th>
					<th>Test</th>
					<th>Size</th>
				</tr>
			</thead>
			<tbody>
				{#each data.models as m (m.code)}
					<tr>
						<td>{m.language}</td>
						<td><code>{m.code}</code></td>
						<td>
							{#each m.data as d, i (d)}
								{#if i > 0},{/if}
								<a href="https://github.com/UniversalDependencies/{d}" rel="noopener">{d}</a>
							{/each}
						</td>
						<td>{m.license}</td>
						<td>{m.dev}</td>
						<td>{m.test}</td>
						<td>{m.size}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<h2>Training your own</h2>
	<p>
		<code>cmd/train-tagger</code> trains a model from any CoNLL-U treebank, and
		<code>tag/lang/build.sh</code> retrains every model above and rewrites this table.
	</p>
</article>
