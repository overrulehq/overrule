//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/overrulehq/overrule/pkg/appeal"
	"github.com/overrulehq/overrule/pkg/parser"
)

func parseDenialWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return js.ValueOf(map[string]any{"error": "missing text argument"})
		}
		rawText := args[0].String()
		parsed := parser.ParseRawText(rawText)
		c := parsed.ToDenialCase("Patient Appellant", "Attending Physician")
		b, _ := json.Marshal(c)
		return js.ValueOf(string(b))
	})
}

func generateAppealWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return js.ValueOf(map[string]any{"error": "missing text argument"})
		}
		rawText := args[0].String()
		parsed := parser.ParseRawText(rawText)
		c := parsed.ToDenialCase("Patient Appellant", "Attending Physician")
		packet := appeal.GeneratePacket(c)
		return js.ValueOf(packet.ToMarkdown())
	})
}

func main() {
	js.Global().Set("overruleParseDenial", parseDenialWrapper())
	js.Global().Set("overruleGenerateAppeal", generateAppealWrapper())
	select {}
}
