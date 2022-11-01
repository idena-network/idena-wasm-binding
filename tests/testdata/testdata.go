package testdata

import "embed"

//go:embed sum.wasm
var content embed.FS

func Sum() ([]byte, error) {
	return content.ReadFile("sum.wasm")
}
