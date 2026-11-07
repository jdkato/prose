import { error } from '@sveltejs/kit';
import { docs } from '#lib/docs';
import type { EntryGenerator, PageLoad } from './$types';

export const entries: EntryGenerator = () => docs.map(({ slug }) => ({ slug }));

export const load: PageLoad = ({ params }) => {
	const doc = docs.find((d) => d.slug === params.slug);
	if (!doc) error(404, 'Not found');
	const i = docs.indexOf(doc);
	return { doc, prev: docs[i - 1] ?? null, next: docs[i + 1] ?? null };
};
