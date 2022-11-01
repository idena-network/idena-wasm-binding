//go:build darwin

package lib

// #cgo LDFLAGS: -L${SRCDIR} -lidena_wasm_mac
import "C"
