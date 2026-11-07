declare global {
	namespace App {}

	// Installed by wasm_exec.js and cmd/prosew.
	class Go {
		importObject: WebAssembly.Imports;
		run(instance: WebAssembly.Instance): Promise<void>;
	}
	function ProseAnalyze(text: string): { result?: string; error?: string };
	function ProseCase(text: string): { result?: string; error?: string };
}

export {};
