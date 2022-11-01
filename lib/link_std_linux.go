//go:build linux

package lib

// #cgo LDFLAGS: -L${SRCDIR} -lidena_wasm_linux
import "C"
