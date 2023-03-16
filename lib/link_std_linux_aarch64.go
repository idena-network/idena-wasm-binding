//go:build linux && arm64

package lib

// #cgo LDFLAGS: -L${SRCDIR} -lidena_wasm_linux_aarch64 -lgcc_s -lutil -lrt -lpthread -lm -ldl -lc
import "C"
