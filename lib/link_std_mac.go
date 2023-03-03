//go:build darwin

package lib

// #cgo LDFLAGS: -L${SRCDIR} -lidena_wasm_darwin_amd64
import "C"
