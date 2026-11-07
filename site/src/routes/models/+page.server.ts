import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import type { PageServerLoad } from './$types';

export interface Model {
	language: string;
	code: string;
	data: string[];
	license: string;
	dev: string;
	test: string;
	size: string;
}

// Read at build time, from the repository the site lives in.
const source = resolve(process.cwd(), '../tag/lang/MODELS.md');

// The table in tag/lang/MODELS.md, written by tag/lang/build.sh.
function parse(markdown: string): Model[] {
	return markdown
		.split('\n')
		.filter((line) => line.startsWith('| ') && !line.startsWith('| Language') && !line.startsWith('| ---'))
		.map((line) => {
			const cells = line
				.slice(1, -1)
				.split('|')
				.map((c) => c.trim());
			const [language, code, data, license, dev, test, size] = cells;
			return {
				language,
				code: code.replaceAll('`', ''),
				data: data.split(';').map((d) => d.trim()),
				license,
				dev,
				test,
				size
			};
		});
}

export const load: PageServerLoad = () => {
	return { models: parse(readFileSync(source, 'utf8')) };
};
