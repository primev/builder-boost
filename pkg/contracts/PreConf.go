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
	_ = abi.ConvertType
)

// PreConfCommitmentStorePreConfCommitment is an auto generated low-level Go binding around an user-defined struct.
type PreConfCommitmentStorePreConfCommitment struct {
	TxnHash             string
	Bid                 uint64
	BlockNumber         uint64
	BidHash             string
	BidSignature        string
	CommitmentSignature string
}

// PreConfCommitmentStoreMetaData contains all meta data concerning the PreConfCommitmentStore contract.
var PreConfCommitmentStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"commitments\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"txnHash\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"bid\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"bidHash\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"bidSignature\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"commitmentSignature\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getCommitment\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"txnHash\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"bid\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"bidHash\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"bidSignature\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"commitmentSignature\",\"type\":\"string\"}],\"internalType\":\"structPreConfCommitmentStore.PreConfCommitment\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"txnHash\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"bid\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"bidHash\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"bidSignature\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"commitmentSignature\",\"type\":\"string\"}],\"name\":\"storeCommitment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// PreConfCommitmentStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use PreConfCommitmentStoreMetaData.ABI instead.
var PreConfCommitmentStoreABI = PreConfCommitmentStoreMetaData.ABI

// PreConfCommitmentStore is an auto generated Go binding around an Ethereum contract.
type PreConfCommitmentStore struct {
	PreConfCommitmentStoreCaller     // Read-only binding to the contract
	PreConfCommitmentStoreTransactor // Write-only binding to the contract
	PreConfCommitmentStoreFilterer   // Log filterer for contract events
}

// PreConfCommitmentStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type PreConfCommitmentStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreConfCommitmentStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PreConfCommitmentStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreConfCommitmentStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PreConfCommitmentStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreConfCommitmentStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PreConfCommitmentStoreSession struct {
	Contract     *PreConfCommitmentStore // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PreConfCommitmentStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PreConfCommitmentStoreCallerSession struct {
	Contract *PreConfCommitmentStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// PreConfCommitmentStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PreConfCommitmentStoreTransactorSession struct {
	Contract     *PreConfCommitmentStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// PreConfCommitmentStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type PreConfCommitmentStoreRaw struct {
	Contract *PreConfCommitmentStore // Generic contract binding to access the raw methods on
}

// PreConfCommitmentStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PreConfCommitmentStoreCallerRaw struct {
	Contract *PreConfCommitmentStoreCaller // Generic read-only contract binding to access the raw methods on
}

// PreConfCommitmentStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PreConfCommitmentStoreTransactorRaw struct {
	Contract *PreConfCommitmentStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPreConfCommitmentStore creates a new instance of PreConfCommitmentStore, bound to a specific deployed contract.
func NewPreConfCommitmentStore(address common.Address, backend bind.ContractBackend) (*PreConfCommitmentStore, error) {
	contract, err := bindPreConfCommitmentStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PreConfCommitmentStore{PreConfCommitmentStoreCaller: PreConfCommitmentStoreCaller{contract: contract}, PreConfCommitmentStoreTransactor: PreConfCommitmentStoreTransactor{contract: contract}, PreConfCommitmentStoreFilterer: PreConfCommitmentStoreFilterer{contract: contract}}, nil
}

// NewPreConfCommitmentStoreCaller creates a new read-only instance of PreConfCommitmentStore, bound to a specific deployed contract.
func NewPreConfCommitmentStoreCaller(address common.Address, caller bind.ContractCaller) (*PreConfCommitmentStoreCaller, error) {
	contract, err := bindPreConfCommitmentStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PreConfCommitmentStoreCaller{contract: contract}, nil
}

// NewPreConfCommitmentStoreTransactor creates a new write-only instance of PreConfCommitmentStore, bound to a specific deployed contract.
func NewPreConfCommitmentStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*PreConfCommitmentStoreTransactor, error) {
	contract, err := bindPreConfCommitmentStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PreConfCommitmentStoreTransactor{contract: contract}, nil
}

// NewPreConfCommitmentStoreFilterer creates a new log filterer instance of PreConfCommitmentStore, bound to a specific deployed contract.
func NewPreConfCommitmentStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*PreConfCommitmentStoreFilterer, error) {
	contract, err := bindPreConfCommitmentStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PreConfCommitmentStoreFilterer{contract: contract}, nil
}

// bindPreConfCommitmentStore binds a generic wrapper to an already deployed contract.
func bindPreConfCommitmentStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PreConfCommitmentStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PreConfCommitmentStore *PreConfCommitmentStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PreConfCommitmentStore.Contract.PreConfCommitmentStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PreConfCommitmentStore *PreConfCommitmentStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PreConfCommitmentStore.Contract.PreConfCommitmentStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PreConfCommitmentStore *PreConfCommitmentStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PreConfCommitmentStore.Contract.PreConfCommitmentStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PreConfCommitmentStore *PreConfCommitmentStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PreConfCommitmentStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PreConfCommitmentStore *PreConfCommitmentStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PreConfCommitmentStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PreConfCommitmentStore *PreConfCommitmentStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PreConfCommitmentStore.Contract.contract.Transact(opts, method, params...)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(string txnHash, uint64 bid, uint64 blockNumber, string bidHash, string bidSignature, string commitmentSignature)
func (_PreConfCommitmentStore *PreConfCommitmentStoreCaller) Commitments(opts *bind.CallOpts, arg0 *big.Int) (struct {
	TxnHash             string
	Bid                 uint64
	BlockNumber         uint64
	BidHash             string
	BidSignature        string
	CommitmentSignature string
}, error) {
	var out []interface{}
	err := _PreConfCommitmentStore.contract.Call(opts, &out, "commitments", arg0)

	outstruct := new(struct {
		TxnHash             string
		Bid                 uint64
		BlockNumber         uint64
		BidHash             string
		BidSignature        string
		CommitmentSignature string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TxnHash = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Bid = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.BlockNumber = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.BidHash = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.BidSignature = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.CommitmentSignature = *abi.ConvertType(out[5], new(string)).(*string)

	return *outstruct, err

}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(string txnHash, uint64 bid, uint64 blockNumber, string bidHash, string bidSignature, string commitmentSignature)
func (_PreConfCommitmentStore *PreConfCommitmentStoreSession) Commitments(arg0 *big.Int) (struct {
	TxnHash             string
	Bid                 uint64
	BlockNumber         uint64
	BidHash             string
	BidSignature        string
	CommitmentSignature string
}, error) {
	return _PreConfCommitmentStore.Contract.Commitments(&_PreConfCommitmentStore.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(string txnHash, uint64 bid, uint64 blockNumber, string bidHash, string bidSignature, string commitmentSignature)
func (_PreConfCommitmentStore *PreConfCommitmentStoreCallerSession) Commitments(arg0 *big.Int) (struct {
	TxnHash             string
	Bid                 uint64
	BlockNumber         uint64
	BidHash             string
	BidSignature        string
	CommitmentSignature string
}, error) {
	return _PreConfCommitmentStore.Contract.Commitments(&_PreConfCommitmentStore.CallOpts, arg0)
}

// GetCommitment is a free data retrieval call binding the contract method 0x69bcdb7d.
//
// Solidity: function getCommitment(uint256 id) view returns((string,uint64,uint64,string,string,string))
func (_PreConfCommitmentStore *PreConfCommitmentStoreCaller) GetCommitment(opts *bind.CallOpts, id *big.Int) (PreConfCommitmentStorePreConfCommitment, error) {
	var out []interface{}
	err := _PreConfCommitmentStore.contract.Call(opts, &out, "getCommitment", id)

	if err != nil {
		return *new(PreConfCommitmentStorePreConfCommitment), err
	}

	out0 := *abi.ConvertType(out[0], new(PreConfCommitmentStorePreConfCommitment)).(*PreConfCommitmentStorePreConfCommitment)

	return out0, err

}

// GetCommitment is a free data retrieval call binding the contract method 0x69bcdb7d.
//
// Solidity: function getCommitment(uint256 id) view returns((string,uint64,uint64,string,string,string))
func (_PreConfCommitmentStore *PreConfCommitmentStoreSession) GetCommitment(id *big.Int) (PreConfCommitmentStorePreConfCommitment, error) {
	return _PreConfCommitmentStore.Contract.GetCommitment(&_PreConfCommitmentStore.CallOpts, id)
}

// GetCommitment is a free data retrieval call binding the contract method 0x69bcdb7d.
//
// Solidity: function getCommitment(uint256 id) view returns((string,uint64,uint64,string,string,string))
func (_PreConfCommitmentStore *PreConfCommitmentStoreCallerSession) GetCommitment(id *big.Int) (PreConfCommitmentStorePreConfCommitment, error) {
	return _PreConfCommitmentStore.Contract.GetCommitment(&_PreConfCommitmentStore.CallOpts, id)
}

// StoreCommitment is a paid mutator transaction binding the contract method 0x18da8c64.
//
// Solidity: function storeCommitment(string txnHash, uint64 bid, uint64 blockNumber, string bidHash, string bidSignature, string commitmentSignature) returns(uint256)
func (_PreConfCommitmentStore *PreConfCommitmentStoreTransactor) StoreCommitment(opts *bind.TransactOpts, txnHash string, bid uint64, blockNumber uint64, bidHash string, bidSignature string, commitmentSignature string) (*types.Transaction, error) {
	return _PreConfCommitmentStore.contract.Transact(opts, "storeCommitment", txnHash, bid, blockNumber, bidHash, bidSignature, commitmentSignature)
}

// StoreCommitment is a paid mutator transaction binding the contract method 0x18da8c64.
//
// Solidity: function storeCommitment(string txnHash, uint64 bid, uint64 blockNumber, string bidHash, string bidSignature, string commitmentSignature) returns(uint256)
func (_PreConfCommitmentStore *PreConfCommitmentStoreSession) StoreCommitment(txnHash string, bid uint64, blockNumber uint64, bidHash string, bidSignature string, commitmentSignature string) (*types.Transaction, error) {
	return _PreConfCommitmentStore.Contract.StoreCommitment(&_PreConfCommitmentStore.TransactOpts, txnHash, bid, blockNumber, bidHash, bidSignature, commitmentSignature)
}

// StoreCommitment is a paid mutator transaction binding the contract method 0x18da8c64.
//
// Solidity: function storeCommitment(string txnHash, uint64 bid, uint64 blockNumber, string bidHash, string bidSignature, string commitmentSignature) returns(uint256)
func (_PreConfCommitmentStore *PreConfCommitmentStoreTransactorSession) StoreCommitment(txnHash string, bid uint64, blockNumber uint64, bidHash string, bidSignature string, commitmentSignature string) (*types.Transaction, error) {
	return _PreConfCommitmentStore.Contract.StoreCommitment(&_PreConfCommitmentStore.TransactOpts, txnHash, bid, blockNumber, bidHash, bidSignature, commitmentSignature)
}
