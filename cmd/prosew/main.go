//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"
)

func errResult(msg string) any {
	return map[string]any{"error": msg}
}

func jsonResult(v any) any {
	data, err := json.Marshal(v)
	if err != nil {
		return errResult(err.Error())
	}
	return map[string]any{"result": string(data)}
}

// ProseAnalyze(text) returns {result: <Analysis as JSON>} or {error}.
func ProseAnalyze(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return errResult("ProseAnalyze(text) requires 1 argument")
	}
	out, err := analyze(args[0].String())
	if err != nil {
		return errResult(err.Error())
	}
	return jsonResult(out)
}

// ProseCase(text) returns {result: <Cases as JSON>}.
func ProseCase(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return errResult("ProseCase(text) requires 1 argument")
	}
	return jsonResult(convertCase(args[0].String()))
}

func main() {
	js.Global().Set("ProseAnalyze", js.FuncOf(ProseAnalyze))
	js.Global().Set("ProseCase", js.FuncOf(ProseCase))
	select {}
}
