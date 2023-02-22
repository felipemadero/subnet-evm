// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package helloworld

import (
	"errors"
	"fmt"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"

	_ "embed"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// Gas costs for each function. These are set to 0 by default.
	// You should set a gas cost for each function in your contract.
	// Generally, you should not set gas costs very low as this may cause your network to be vulnerable to DoS attacks.
	// There are some predefined gas costs in contract/utils.go that you can use.
	// This contract also uses AllowList precompile.
	// You should also increase gas costs of functions that read from AllowList storage.}
	SayHelloGasCost    uint64 = contract.ReadGasCostPerSlot
	SetGreetingGasCost uint64 = contract.WriteGasCostPerSlot + allowlist.ReadAllowListGasCost
)

// Singleton StatefulPrecompiledContract and signatures.
var (
	ErrCannotSetGreeting = errors.New("non-enabled cannot call setGreeting")
	ErrInputExceedsLimit = errors.New("input string is longer than 32 bytes")

	// HelloWorldRawABI contains the raw ABI of HelloWorld contract.
	//go:embed contract.abi
	HelloWorldRawABI string

	HelloWorldABI = contract.ParseABI(HelloWorldRawABI)

	HelloWorldPrecompile = createHelloWorldPrecompile()

	storageKeyHash = common.BytesToHash([]byte("storageKey"))
)

// GetHelloWorldAllowListStatus returns the role of [address] for the HelloWorld list.
func GetHelloWorldAllowListStatus(stateDB contract.StateDB, address common.Address) allowlist.Role {
	return allowlist.GetAllowListStatus(stateDB, ContractAddress, address)
}

// SetHelloWorldAllowListStatus sets the permissions of [address] to [role] for the
// HelloWorld list. Assumes [role] has already been verified as valid.
// This stores the [role] in the contract storage with address [ContractAddress]
// and [address] hash. It means that any reusage of the [address] key for different value
// conflicts with the same slot [role] is stored.
// Precompile implementations must use a different key than [address] for their storage.
func SetHelloWorldAllowListStatus(stateDB contract.StateDB, address common.Address, role allowlist.Role) {
	allowlist.SetAllowListRole(stateDB, ContractAddress, address, role)
}

// PackSayHello packs the include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackSayHello() ([]byte, error) {
	return HelloWorldABI.Pack("sayHello")
}

// PackSayHelloOutput attempts to pack given result of type string
// to conform the ABI outputs.
func PackSayHelloOutput(result string) ([]byte, error) {
	return HelloWorldABI.PackOutput("sayHello", result)
}

// GetGreeting returns the value of the storage key "storageKey" in the contract storage,
// with leading zeroes trimmed.
// This function is mostly used for tests.
func GetGreeting(stateDB contract.StateDB) string {
	// Get the value set at recipient
	value := stateDB.GetState(ContractAddress, storageKeyHash)
	return string(common.TrimLeftZeroes(value.Bytes()))
}

// sayHello is the reader fucntion that returns the value of greeting stored in the contract storage.
func sayHello(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, SayHelloGasCost); err != nil {
		return nil, 0, err
	}
	// no input provided for this function

	// CUSTOM CODE STARTS HERE

	// Get the current state
	currentState := accessibleState.GetStateDB()
	// Get the value set at recipient
	value := GetGreeting(currentState)
	packedOutput, err := PackSayHelloOutput(value)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackSetGreetingInput attempts to unpack [input] into the string type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackSetGreetingInput(input []byte) (string, error) {
	res, err := HelloWorldABI.UnpackInput("setGreeting", input)
	if err != nil {
		return "", err
	}
	unpacked := *abi.ConvertType(res[0], new(string)).(*string)
	return unpacked, nil
}

// PackSetGreeting packs [response] of type string into the appropriate arguments for setGreeting.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackSetGreeting(response string) ([]byte, error) {
	return HelloWorldABI.Pack("setGreeting", response)
}

// StoreGreeting sets the value of the storage key "storageKey" in the contract storage.
func StoreGreeting(stateDB contract.StateDB, input string) {
	inputPadded := common.LeftPadBytes([]byte(input), common.HashLength)
	inputHash := common.BytesToHash(inputPadded)

	stateDB.SetState(ContractAddress, storageKeyHash, inputHash)
}

// setGreeting is a state-changer function that sets the value of the greeting in the contract storage.
func setGreeting(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, SetGreetingGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the SetGreetingInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackSetGreetingInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// Allow list is enabled and SetGreeting is a state-changer function.
	// This part of the code restricts the function to be called only by enabled/admin addresses in the allow list.
	// You can modify/delete this code if you don't want this function to be restricted by the allow list.
	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotSetGreeting, caller)
	}
	// allow list code ends here.

	// CUSTOM CODE STARTS HERE
	// Check if the input string is longer than HashLength
	if len(inputStruct) > common.HashLength {
		return nil, 0, ErrInputExceedsLimit
	}

	// setGreeting is the execution function
	// "SetGreeting(name string)" and sets the storageKey
	// in the string returned by hello world
	StoreGreeting(stateDB, inputStruct)

	// This function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// createHelloWorldPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.
// Access to the getters/setters is controlled by an allow list for ContractAddress.
func createHelloWorldPrecompile() contract.StatefulPrecompiledContract {
	var functions []*contract.StatefulPrecompileFunction
	functions = append(functions, allowlist.CreateAllowListFunctions(ContractAddress)...)

	abiFunctionMap := map[string]contract.RunStatefulPrecompileFunc{
		"sayHello":    sayHello,
		"setGreeting": setGreeting,
	}

	for name, function := range abiFunctionMap {
		method, ok := HelloWorldABI.Methods[name]
		if !ok {
			panic(fmt.Errorf("given method (%s) does not exist in the ABI", name))
		}
		functions = append(functions, contract.NewStatefulPrecompileFunction(method.ID, function))
	}
	// Construct the contract with no fallback function.
	statefulContract, err := contract.NewStatefulPrecompileContract(nil, functions)
	if err != nil {
		panic(err)
	}
	return statefulContract
}