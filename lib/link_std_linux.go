//go:build linux

package lib

// #cgo LDFLAGS: -L${SRCDIR} -lidena_wasm_linux -lgcc_s -lutil -lrt -lpthread -lm -ldl -lc
import "C"
