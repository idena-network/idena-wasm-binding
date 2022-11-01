//go:build windows

package lib

// #cgo LDFLAGS: -L${SRCDIR} -lidena_wasm_win
// #cgo LDFLAGS: -lws2_32 -lbcrypt -luserenv
import "C"
