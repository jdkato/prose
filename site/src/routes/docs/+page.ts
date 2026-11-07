import { redirect } from '@sveltejs/kit';
import { resolve } from '$app/paths';

export const load = () => {
	redirect(307, resolve('/docs/[slug]', { slug: 'getting-started' }));
};
