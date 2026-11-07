import { asset } from '$app/paths';

export interface Token {
	text: string;
	tag: string;
	label?: string;
	start: number;
	end: number;
}
export interface Span {
	text: string;
	start: number;
	end: number;
	// Flesch-Kincaid grade of the sentence on its own.
	grade: number;
}
export interface Entity extends Span {
	label: string;
}
export interface Readability {
	words: number;
	sentences: number;
	syllables: number;
	characters: number;
	fleschKincaid: number;
	fleschReadingEase: number;
	gunningFog: number;
	smog: number;
	colemanLiau: number;
	automatedReadability: number;
	daleChall: number;
	lix: number;
}
export interface Analysis {
	tokens: Token[];
	sentences: Span[];
	entities: Entity[];
	readability: Readability;
}
export interface Cases {
	ap: string;
	chicago: string;
	sentence: string;
}

let loading: Promise<void> | undefined;

function script(src: string): Promise<void> {
	return new Promise((resolve, reject) => {
		const el = document.createElement('script');
		el.src = src;
		el.onload = () => resolve();
		el.onerror = () => reject(new Error(`failed to load ${src}`));
		document.head.appendChild(el);
	});
}

// Loads prose.wasm once, installing ProseAnalyze and ProseCase on window.
export function load(): Promise<void> {
	loading ??= (async () => {
		await script(asset('wasm_exec.js'));
		const go = new Go();
		const { instance } = await WebAssembly.instantiateStreaming(
			fetch(asset('prose.wasm')),
			go.importObject
		);
		void go.run(instance);
	})();
	return loading;
}

function unwrap<T>(out: { result?: string; error?: string }): T {
	if (out.error) throw new Error(out.error);
	return JSON.parse(out.result!) as T;
}

export const analyze = (text: string) => unwrap<Analysis>(ProseAnalyze(text));
export const convertCase = (text: string) => unwrap<Cases>(ProseCase(text));
