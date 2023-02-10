package tests

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/idena-network/idena-wasm-binding/lib"
	"github.com/idena-network/idena-wasm-binding/tests/testdata"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
	"math/big"
	"testing"
)

type contractValue struct {
	value   []byte
	removed bool
}

type ContractData struct {
	Code []byte
}

type ContractContext struct {
	caller       lib.Address
	originCaller lib.Address
	contractAddr lib.Address
	payAmount    *big.Int
}

func (ctx *ContractContext) ContractAddr() lib.Address {
	return ctx.contractAddr
}

type MockDb struct {
	db *db.MemDB
}

func (db *MockDb) GetContractValue(contract lib.Address, key []byte) []byte {
	formattedKey := append(append([]byte{0x5}, contract[:]...), key...)
	v, _ := db.db.Get(formattedKey)
	return v
}

type MockHostEnv struct {
	parent *MockHostEnv
	ctx    *ContractContext
	db     *MockDb

	contractStoreCache    map[lib.Address]map[string]*contractValue
	balancesCache         map[lib.Address]*big.Int
	deployedContractCache map[lib.Address]ContractData
	contractStakeCache    map[lib.Address]*big.Int
}

func (e *MockHostEnv) GlobalState(meter *lib.GasMeter) []byte {
	//TODO implement me
	panic("implement me")
}

func (e *MockHostEnv) BlockHeader(meter *lib.GasMeter, height uint64) []byte {
	panic("implement me")
}

func (e *MockHostEnv) Keccak256(meter *lib.GasMeter, data []byte) []byte {
	panic("implement me")
}

func (e *MockHostEnv) IsDebug() bool {
	return true
}

func NewMockHostEnv() *MockHostEnv {
	return &MockHostEnv{
		db: &MockDb{
			db: db.NewMemDB(),
		},
		ctx: &ContractContext{
			contractAddr: lib.Address{0x1},
			caller:       lib.Address{0x2},
			originCaller: lib.Address{0x3},
			payAmount:    big.NewInt(10),
		},
		deployedContractCache: map[lib.Address]ContractData{},
		contractStakeCache:    map[lib.Address]*big.Int{},
		balancesCache:         map[lib.Address]*big.Int{},
		contractStoreCache:    map[lib.Address]map[string]*contractValue{},
	}
}

func (e *MockHostEnv) SetStorage(meter *lib.GasMeter, key []byte, value []byte) {
	ctx := e.ctx
	if len(key) > 32 {
		panic("key is too big")
	}
	addr := ctx.ContractAddr()
	var cache map[string]*contractValue
	var ok bool
	if cache, ok = e.contractStoreCache[addr]; !ok {
		cache = make(map[string]*contractValue)
		e.contractStoreCache[addr] = cache
	}
	cache[string(key)] = &contractValue{
		value:   value,
		removed: false,
	}
	meter.ConsumeGas(uint64(10 * (len(key) + len(value))))
}

func (e *MockHostEnv) GetStorage(meter *lib.GasMeter, key []byte) []byte {
	value := e.readContractData(e.ctx.ContractAddr(), key)
	meter.ConsumeGas(uint64(10 * len(value)))
	return value
}

func (e *MockHostEnv) RemoveStorage(meter *lib.GasMeter, key []byte) {
	addr := e.ctx.ContractAddr()
	var cache map[string]*contractValue
	var ok bool
	if cache, ok = e.contractStoreCache[addr]; !ok {
		cache = map[string]*contractValue{}
		e.contractStoreCache[addr] = cache
	}
	cache[string(key)] = &contractValue{removed: true}
	meter.ConsumeGas(10)
}

func (e *MockHostEnv) readContractData(contractAddr lib.Address, key []byte) []byte {
	if cache, ok := e.contractStoreCache[contractAddr]; ok {
		if value, ok := cache[string(key)]; ok {
			if value.removed {
				return nil
			}
			return value.value
		}
	}

	if e.parent != nil {
		return e.parent.readContractData(contractAddr, key)
	}

	value := e.db.GetContractValue(contractAddr, key)

	return value
}

func (e MockHostEnv) BlockNumber(meter *lib.GasMeter) uint64 {
	panic("implement me")
}

func (e MockHostEnv) BlockTimestamp(meter *lib.GasMeter) int64 {
	panic("implement me")
}

func (e MockHostEnv) Send(meter *lib.GasMeter, address lib.Address, b *big.Int) error {
	panic("implement me")
}

func (e MockHostEnv) MinFeePerGas(meter *lib.GasMeter) *big.Int {
	panic("implement me")
}

func (e MockHostEnv) Balance(meter *lib.GasMeter) *big.Int {
	panic("implement me")
}

func (e MockHostEnv) BlockSeed(meter *lib.GasMeter) []byte {
	panic("implement me")
}

func (e MockHostEnv) NetworkSize(meter *lib.GasMeter) uint64 {
	panic("implement me")
}

func (e MockHostEnv) IdentityState(meter *lib.GasMeter, address lib.Address) byte {
	panic("implement me")
}

func (e MockHostEnv) Identity(meter *lib.GasMeter, address lib.Address) []byte {
	panic("implement me")
}

func (e MockHostEnv) CreateSubEnv(contract lib.Address, method string, payAmount *big.Int, isDeploy bool) (lib.HostEnv, error) {
	panic("implement me")
}

func (e MockHostEnv) GetCode(addr lib.Address) []byte {
	panic("implement me")
}

func (e MockHostEnv) Commit() {

}

func (e MockHostEnv) Clear() {
	panic("implement me")
}

func (e MockHostEnv) Caller(meter *lib.GasMeter) lib.Address {
	panic("implement me")
}

func (e MockHostEnv) OriginalCaller(meter *lib.GasMeter) lib.Address {
	panic("implement me")
}

func (e MockHostEnv) SubBalance(meter *lib.GasMeter, amount *big.Int) error {
	panic("implement me")
}

func (e MockHostEnv) AddBalance(meter *lib.GasMeter, address lib.Address, bytes *big.Int) {
	panic("implement me")
}

func (e MockHostEnv) ContractAddress(meter *lib.GasMeter) lib.Address {
	return e.ctx.ContractAddr()
}

func (e MockHostEnv) ContractCode(meter *lib.GasMeter, addr lib.Address) []byte {
	panic("implement me")
}

func (e MockHostEnv) ContractAddr(meter *lib.GasMeter, code []byte, args []byte, nonce []byte) lib.Address {
	panic("implement me")
}

func (e MockHostEnv) Deploy(code []byte) {
	panic("implement me")
}

func (e MockHostEnv) ContractAddrByHash(meter *lib.GasMeter, hash []byte, args []byte, nonce []byte) lib.Address {
	panic("implement me")
}

func (e MockHostEnv) OwnCode(meter *lib.GasMeter) []byte {
	panic("implement me")
}

func (e MockHostEnv) CodeHash(meter *lib.GasMeter) []byte {
	panic("implement me")
}

func (e MockHostEnv) Event(meter *lib.GasMeter, name string, args ...[]byte) {
	panic("implement me")
}

func (e MockHostEnv) ReadContractData(meter *lib.GasMeter, address lib.Address, key []byte) []byte {
	panic("implement me")
}

func (e MockHostEnv) Epoch(meter *lib.GasMeter) uint16 {
	panic("implement me")
}

func (e MockHostEnv) ContractCodeHash(addr lib.Address) *[]byte {
	panic("implement me")
}

func (e MockHostEnv) PayAmount(meter *lib.GasMeter) *big.Int {
	panic("implement me")
}

func ToBytes(value interface{}) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func TestSum(t *testing.T) {
	code, _ := testdata.Sum()

	api := lib.NewGoAPI(NewMockHostEnv(), &lib.GasMeter{})

	_, _, err := lib.Deploy(api, code, [][]byte{ToBytes(uint64(1))}, lib.Address{}, 10000000, true)
	require.NoError(t, err)
	_, _, err = lib.Execute(api, code, "compute", [][]byte{ToBytes(uint64(10))}, lib.Address{}, 1000000, true)
	require.NoError(t, err)
}
