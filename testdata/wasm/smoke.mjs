// Exercises the wasm build's JavaScript API the way the playground does.
// Run with `make test-wasm`.
//
//   node smoke.mjs <wasm_exec.js> <prose.wasm>
import fs from "node:fs";
import path from "node:path";
import { createRequire } from "node:module";
import assert from "node:assert/strict";

const require = createRequire(import.meta.url);
globalThis.require = require;
globalThis.TextEncoder = TextEncoder;
globalThis.TextDecoder = TextDecoder;
globalThis.performance ??= { now: () => Date.now() };
const crypto = await import("node:crypto");
globalThis.crypto ??= crypto.webcrypto;

require(path.resolve(process.argv[2]));
const go = new Go();
const { instance } = await WebAssembly.instantiate(fs.readFileSync(path.resolve(process.argv[3])), go.importObject);
go.run(instance);
await new Promise((r) => setTimeout(r, 50));

for (const name of ["ProseAnalyze", "ProseCase"]) {
  assert.equal(typeof globalThis[name], "function", `${name} is installed`);
}

const src = "Marie Curie studied in Paris. She won the Nobel Prize twice.";
let out = ProseAnalyze(src);
assert.equal(out.error, undefined, out.error);
const doc = JSON.parse(out.result);
assert.equal(doc.sentences.length, 2);
assert.deepEqual(doc.entities.map((e) => [e.text, e.label]), [["Marie Curie", "PERSON"], ["Paris", "GPE"]]);
for (const t of doc.tokens) {
  assert.equal(src.slice(t.start, t.end), t.text);
  assert.ok(t.tag, `${t.text} has a tag`);
}
assert.equal(doc.readability.sentences, 2);
assert.ok(doc.readability.fleschKincaid > 0);

out = ProseCase("the quick brown fox");
assert.equal(JSON.parse(out.result).ap, "The Quick Brown Fox");

assert.match(ProseAnalyze().error, /requires 1 argument/);

console.log("ok");
process.exit(0);
