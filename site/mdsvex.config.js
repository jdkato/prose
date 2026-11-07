import { escapeSvelte } from 'mdsvex';
import { createHighlighter } from 'shiki';
import rehypeSlug from 'rehype-slug';
import rehypeAutolinkHeadings from 'rehype-autolink-headings';

const highlighter = await createHighlighter({
	themes: ['github-light', 'github-dark'],
	langs: ['go', 'console', 'bash', 'ini', 'json', 'text']
});

/** @type {import('mdsvex').MdsvexOptions} */
export default {
	extensions: ['.md'],
	highlight: {
		highlighter: (code, lang) => {
			const html = highlighter.codeToHtml(code, {
				lang: lang && highlighter.getLoadedLanguages().includes(lang) ? lang : 'text',
				themes: { light: 'github-light', dark: 'github-dark' },
				defaultColor: false
			});
			return `{@html ${JSON.stringify(html)}}`;
		}
	},
	rehypePlugins: [rehypeSlug, [rehypeAutolinkHeadings, { behavior: 'wrap' }]]
};
