package lib

// #include <stdlib.h>
// #include "bindings.h"
import "C"
import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	models "github.com/idena-network/idena-wasm-binding/lib/protobuf"
	"unsafe"
)

// Value types
type cint = C.int
type cbool = C.bool
type cusize = C.size_t
type cu8 = C.uint8_t
type cu16 = C.uint16_t
type cu32 = C.uint32_t
type cu64 = C.uint64_t
type ci8 = C.int8_t
type ci32 = C.int32_t
type ci64 = C.int64_t

// Pointers
type cu8_ptr = *C.uint8_t

const ArgsPlainFormat = 0x0
const ArgsProtobufFormat = 0x1

func copyAndDestroyUnmanagedVector(v C.UnmanagedVector) []byte {
	var out []byte
	if v.is_none {
		out = nil
	} else if v.cap == cusize(0) {
		// There is no allocation we can copy
		out = []byte{}
	} else {
		// C.GoBytes create a copy (https://stackoverflow.com/a/40950744/2013738)
		out = C.GoBytes(unsafe.Pointer(v.ptr), cint(v.len))
	}
	//C.destroy_unmanaged_vector(v)
	return out
}

func errorWithMessage(err int, b C.UnmanagedVector) error {
	// this checks for out of gas as a special case
	if err == 2 {
		return errors.New("out of gas")
	}
	msg := copyAndDestroyUnmanagedVector(b)
	if msg == nil {
		return errors.New("error without description")
	}
	return fmt.Errorf("%s", string(msg))
}

func Execute(api *GoAPI, code []byte, method string, args [][]byte, contractAddr Address, gasLimit uint64, is_debug bool) (uint64, []byte, error) {

	gas, actionResult, err := execute(api, code, []byte(method), PackArguments(args), []byte{}, contractAddr, gasLimit, is_debug)
	return gas, actionResult, err
}

func Deploy(api *GoAPI, code []byte, args [][]byte, contractAddr Address, gasLimit uint64, is_debug bool) (uint64, []byte, error) {
	gas, actionResult, err := deploy(api, code, PackArguments(args), contractAddr, gasLimit, is_debug)
	return gas, actionResult, err
}

func execute(api *GoAPI, code []byte, method []byte, args []byte, invocationContext []byte, contractAddr Address, gasLimit uint64, is_debug bool) (uint64, []byte, error) {

	actionResult := newUnmanagedVector(nil)

	var gasUsed cu64
	C.execute(buildAPI(api), makeView(code), makeView(method), makeView(args), makeView(invocationContext), makeView(contractAddr[:]), cu64(gasLimit), &gasUsed, &actionResult, cbool(is_debug))
	actionResultBytes := copyAndDestroyUnmanagedVector(actionResult)
	protoModel := models.ActionResult{}
	if err := proto.Unmarshal(actionResultBytes, &protoModel); err != nil {
		return uint64(gasUsed), actionResultBytes, err
	}
	truncateActionResultData(&protoModel, 1)
	if protoModel.Success {
		return protoModel.GasUsed, actionResultBytes, nil
	}
	return protoModel.GasUsed, actionResultBytes, errors.New(protoModel.Error)
}

const MaxArgsLength = 100

func truncateActionResultData(actionResult *models.ActionResult, depth int) {
	if depth > 5 {
		return
	}
	if actionResult.InputAction != nil {
		if len(actionResult.InputAction.Args) > MaxArgsLength {
			actionResult.InputAction.Args = actionResult.InputAction.Args[:MaxArgsLength]
		}
	}
	for _, subAction := range actionResult.SubActionResults {
		truncateActionResultData(subAction, depth+1)
	}
}

func deploy(api *GoAPI, code []byte, args []byte, contractAddr Address, gasLimit uint64, is_debug bool) (uint64, []byte, error) {
	actionResult := newUnmanagedVector(nil)

	var gasUsed cu64
	C.deploy(buildAPI(api), makeView(code), makeView(args), makeView(contractAddr[:]), cu64(gasLimit), &gasUsed, &actionResult, cbool(is_debug))
	actionResultBytes := copyAndDestroyUnmanagedVector(actionResult)
	protoModel := models.ActionResult{}
	if err := proto.Unmarshal(actionResultBytes, &protoModel); err != nil {
		return uint64(gasUsed), actionResultBytes, err
	}
	truncateActionResultData(&protoModel, 1)
	if protoModel.Success {
		return protoModel.GasUsed, actionResultBytes, nil
	}
	return protoModel.GasUsed, actionResultBytes, errors.New(protoModel.Error)
}

func PackArguments(args [][]byte) []byte {
	argsMsg := models.ProtoArgs{}
	for idx := range args {
		arg := args[idx]
		if arg == nil {
			argsMsg.Args = append(argsMsg.Args, &models.ProtoArgs_Argument{
				IsNil: true,
			})
		} else {
			argsMsg.Args = append(argsMsg.Args, &models.ProtoArgs_Argument{
				Value: arg,
			})
		}
	}

	argsBytes, _ := proto.Marshal(&argsMsg)
	return append([]byte{ArgsProtobufFormat}, argsBytes...)
}

func UnpackArguments(args []byte) [][]byte {
	if len(args) == 0 {
		return [][]byte{}
	}
	format := args[0]

	switch format {
	case ArgsProtobufFormat:
		argsProto := models.ProtoArgs{}
		if err := proto.Unmarshal(args[1:], &argsProto); err != nil {
			return [][]byte{}
		}
		result := make([][]byte, 0, len(argsProto.GetArgs()))
		for _, arg := range argsProto.GetArgs() {
			if arg.IsNil {
				result = append(result, nil)
			} else {
				result = append(result, arg.Value)
			}
		}
		return result
	case ArgsPlainFormat:
		return [][]byte{args[1:]}
	default:
		return [][]byte{}
	}
}
