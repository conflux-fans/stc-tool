// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// TemplatesMetaData contains all meta data concerning the Templates contract.
var TemplatesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"keys\",\"type\":\"string[]\"}],\"name\":\"AddTemplate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"keys\",\"type\":\"string[]\"}],\"name\":\"AppendTemplateKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"DeleteTemplate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"GetTemplate\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ListTemplates\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"keys\",\"type\":\"string[]\"}],\"name\":\"RemoveTemplateKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TemplatesABI is the input ABI used to generate the binding from.
// Deprecated: Use TemplatesMetaData.ABI instead.
var TemplatesABI = TemplatesMetaData.ABI

// Templates is an auto generated Go binding around an Ethereum contract.
type Templates struct {
	TemplatesCaller     // Read-only binding to the contract
	TemplatesTransactor // Write-only binding to the contract
	TemplatesFilterer   // Log filterer for contract events
}

// TemplatesCaller is an auto generated read-only Go binding around an Ethereum contract.
type TemplatesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TemplatesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TemplatesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TemplatesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TemplatesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TemplatesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TemplatesSession struct {
	Contract     *Templates        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TemplatesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TemplatesCallerSession struct {
	Contract *TemplatesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TemplatesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TemplatesTransactorSession struct {
	Contract     *TemplatesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TemplatesRaw is an auto generated low-level Go binding around an Ethereum contract.
type TemplatesRaw struct {
	Contract *Templates // Generic contract binding to access the raw methods on
}

// TemplatesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TemplatesCallerRaw struct {
	Contract *TemplatesCaller // Generic read-only contract binding to access the raw methods on
}

// TemplatesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TemplatesTransactorRaw struct {
	Contract *TemplatesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTemplates creates a new instance of Templates, bound to a specific deployed contract.
func NewTemplates(address common.Address, backend bind.ContractBackend) (*Templates, error) {
	contract, err := bindTemplates(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Templates{TemplatesCaller: TemplatesCaller{contract: contract}, TemplatesTransactor: TemplatesTransactor{contract: contract}, TemplatesFilterer: TemplatesFilterer{contract: contract}}, nil
}

// NewTemplatesCaller creates a new read-only instance of Templates, bound to a specific deployed contract.
func NewTemplatesCaller(address common.Address, caller bind.ContractCaller) (*TemplatesCaller, error) {
	contract, err := bindTemplates(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TemplatesCaller{contract: contract}, nil
}

// NewTemplatesTransactor creates a new write-only instance of Templates, bound to a specific deployed contract.
func NewTemplatesTransactor(address common.Address, transactor bind.ContractTransactor) (*TemplatesTransactor, error) {
	contract, err := bindTemplates(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TemplatesTransactor{contract: contract}, nil
}

// NewTemplatesFilterer creates a new log filterer instance of Templates, bound to a specific deployed contract.
func NewTemplatesFilterer(address common.Address, filterer bind.ContractFilterer) (*TemplatesFilterer, error) {
	contract, err := bindTemplates(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TemplatesFilterer{contract: contract}, nil
}

// bindTemplates binds a generic wrapper to an already deployed contract.
func bindTemplates(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TemplatesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Templates *TemplatesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Templates.Contract.TemplatesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Templates *TemplatesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Templates.Contract.TemplatesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Templates *TemplatesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Templates.Contract.TemplatesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Templates *TemplatesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Templates.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Templates *TemplatesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Templates.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Templates *TemplatesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Templates.Contract.contract.Transact(opts, method, params...)
}

// GetTemplate is a free data retrieval call binding the contract method 0xa7e2859a.
//
// Solidity: function GetTemplate(string name) view returns(string[])
func (_Templates *TemplatesCaller) GetTemplate(opts *bind.CallOpts, name string) ([]string, error) {
	var out []interface{}
	err := _Templates.contract.Call(opts, &out, "GetTemplate", name)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetTemplate is a free data retrieval call binding the contract method 0xa7e2859a.
//
// Solidity: function GetTemplate(string name) view returns(string[])
func (_Templates *TemplatesSession) GetTemplate(name string) ([]string, error) {
	return _Templates.Contract.GetTemplate(&_Templates.CallOpts, name)
}

// GetTemplate is a free data retrieval call binding the contract method 0xa7e2859a.
//
// Solidity: function GetTemplate(string name) view returns(string[])
func (_Templates *TemplatesCallerSession) GetTemplate(name string) ([]string, error) {
	return _Templates.Contract.GetTemplate(&_Templates.CallOpts, name)
}

// ListTemplates is a free data retrieval call binding the contract method 0xcb0e245f.
//
// Solidity: function ListTemplates() view returns(string[])
func (_Templates *TemplatesCaller) ListTemplates(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Templates.contract.Call(opts, &out, "ListTemplates")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// ListTemplates is a free data retrieval call binding the contract method 0xcb0e245f.
//
// Solidity: function ListTemplates() view returns(string[])
func (_Templates *TemplatesSession) ListTemplates() ([]string, error) {
	return _Templates.Contract.ListTemplates(&_Templates.CallOpts)
}

// ListTemplates is a free data retrieval call binding the contract method 0xcb0e245f.
//
// Solidity: function ListTemplates() view returns(string[])
func (_Templates *TemplatesCallerSession) ListTemplates() ([]string, error) {
	return _Templates.Contract.ListTemplates(&_Templates.CallOpts)
}

// AddTemplate is a paid mutator transaction binding the contract method 0xf39f913a.
//
// Solidity: function AddTemplate(string name, string[] keys) returns()
func (_Templates *TemplatesTransactor) AddTemplate(opts *bind.TransactOpts, name string, keys []string) (*types.Transaction, error) {
	return _Templates.contract.Transact(opts, "AddTemplate", name, keys)
}

// AddTemplate is a paid mutator transaction binding the contract method 0xf39f913a.
//
// Solidity: function AddTemplate(string name, string[] keys) returns()
func (_Templates *TemplatesSession) AddTemplate(name string, keys []string) (*types.Transaction, error) {
	return _Templates.Contract.AddTemplate(&_Templates.TransactOpts, name, keys)
}

// AddTemplate is a paid mutator transaction binding the contract method 0xf39f913a.
//
// Solidity: function AddTemplate(string name, string[] keys) returns()
func (_Templates *TemplatesTransactorSession) AddTemplate(name string, keys []string) (*types.Transaction, error) {
	return _Templates.Contract.AddTemplate(&_Templates.TransactOpts, name, keys)
}

// AppendTemplateKeys is a paid mutator transaction binding the contract method 0xc19c6227.
//
// Solidity: function AppendTemplateKeys(string name, string[] keys) returns()
func (_Templates *TemplatesTransactor) AppendTemplateKeys(opts *bind.TransactOpts, name string, keys []string) (*types.Transaction, error) {
	return _Templates.contract.Transact(opts, "AppendTemplateKeys", name, keys)
}

// AppendTemplateKeys is a paid mutator transaction binding the contract method 0xc19c6227.
//
// Solidity: function AppendTemplateKeys(string name, string[] keys) returns()
func (_Templates *TemplatesSession) AppendTemplateKeys(name string, keys []string) (*types.Transaction, error) {
	return _Templates.Contract.AppendTemplateKeys(&_Templates.TransactOpts, name, keys)
}

// AppendTemplateKeys is a paid mutator transaction binding the contract method 0xc19c6227.
//
// Solidity: function AppendTemplateKeys(string name, string[] keys) returns()
func (_Templates *TemplatesTransactorSession) AppendTemplateKeys(name string, keys []string) (*types.Transaction, error) {
	return _Templates.Contract.AppendTemplateKeys(&_Templates.TransactOpts, name, keys)
}

// DeleteTemplate is a paid mutator transaction binding the contract method 0x7a560b9c.
//
// Solidity: function DeleteTemplate(string name) returns()
func (_Templates *TemplatesTransactor) DeleteTemplate(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _Templates.contract.Transact(opts, "DeleteTemplate", name)
}

// DeleteTemplate is a paid mutator transaction binding the contract method 0x7a560b9c.
//
// Solidity: function DeleteTemplate(string name) returns()
func (_Templates *TemplatesSession) DeleteTemplate(name string) (*types.Transaction, error) {
	return _Templates.Contract.DeleteTemplate(&_Templates.TransactOpts, name)
}

// DeleteTemplate is a paid mutator transaction binding the contract method 0x7a560b9c.
//
// Solidity: function DeleteTemplate(string name) returns()
func (_Templates *TemplatesTransactorSession) DeleteTemplate(name string) (*types.Transaction, error) {
	return _Templates.Contract.DeleteTemplate(&_Templates.TransactOpts, name)
}

// RemoveTemplateKeys is a paid mutator transaction binding the contract method 0x5a65f233.
//
// Solidity: function RemoveTemplateKeys(string name, string[] keys) returns()
func (_Templates *TemplatesTransactor) RemoveTemplateKeys(opts *bind.TransactOpts, name string, keys []string) (*types.Transaction, error) {
	return _Templates.contract.Transact(opts, "RemoveTemplateKeys", name, keys)
}

// RemoveTemplateKeys is a paid mutator transaction binding the contract method 0x5a65f233.
//
// Solidity: function RemoveTemplateKeys(string name, string[] keys) returns()
func (_Templates *TemplatesSession) RemoveTemplateKeys(name string, keys []string) (*types.Transaction, error) {
	return _Templates.Contract.RemoveTemplateKeys(&_Templates.TransactOpts, name, keys)
}

// RemoveTemplateKeys is a paid mutator transaction binding the contract method 0x5a65f233.
//
// Solidity: function RemoveTemplateKeys(string name, string[] keys) returns()
func (_Templates *TemplatesTransactorSession) RemoveTemplateKeys(name string, keys []string) (*types.Transaction, error) {
	return _Templates.Contract.RemoveTemplateKeys(&_Templates.TransactOpts, name, keys)
}
