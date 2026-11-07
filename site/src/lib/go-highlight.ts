// Just enough Go highlighting for the playground's code panel: comments,
// strings, keywords and a few builtin types. Returns HTML.

const keywords = new Set([
	'package', 'import', 'func', 'var', 'const', 'for', 'range', 'if', 'else', 'return', 'type',
	'struct', 'map', 'chan', 'go', 'defer', 'switch', 'case', 'default', 'break', 'continue', 'nil',
	'true', 'false', 'string', 'int', 'float64', 'bool', 'error', 'byte', 'rune'
]);

const escape = (s: string) =>
	s.replaceAll('&', '&amp;').replaceAll('<', '&lt;').replaceAll('>', '&gt;');

const pattern = /(\/\/[^\n]*)|(`[^`]*`|"(?:[^"\\\n]|\\.)*")|([A-Za-z_][A-Za-z0-9_]*)/g;

export function highlightGo(code: string): string {
	let out = '';
	let last = 0;
	for (const m of code.matchAll(pattern)) {
		out += escape(code.slice(last, m.index));
		const [whole, comment, str, ident] = m;
		if (comment) out += `<span class="c">${escape(comment)}</span>`;
		else if (str) out += `<span class="s">${escape(str)}</span>`;
		else if (ident && keywords.has(ident)) out += `<span class="k">${ident}</span>`;
		else out += escape(whole);
		last = m.index + whole.length;
	}
	return out + escape(code.slice(last));
}
