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

// BuilderStakingMetaData contains all meta data concerning the BuilderStaking contract.
var BuilderStakingMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"builder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minimalStake\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minimalSubsriptionPeriod\",\"type\":\"uint256\"}],\"name\":\"BuilderUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"builder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subscriptionEnd\",\"type\":\"uint256\"}],\"name\":\"StakeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"builder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"builders\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"minimalStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minimalSubscriptionPeriod\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_builder\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_commitment\",\"type\":\"bytes32\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_builder\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_stakeAmount\",\"type\":\"uint256\"}],\"name\":\"getSubscriptionPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_builder\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_commitment\",\"type\":\"bytes32\"}],\"name\":\"hasMinimalStake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"stakes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subscriptionEnd\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timeLocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"initialAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockDuration\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"timeLocksCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minimalStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minimalSubscriptionPeriod\",\"type\":\"uint256\"}],\"name\":\"updateBuilder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawableAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BuilderStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use BuilderStakingMetaData.ABI instead.
var BuilderStakingABI = BuilderStakingMetaData.ABI

// BuilderStaking is an auto generated Go binding around an Ethereum contract.
type BuilderStaking struct {
	BuilderStakingCaller     // Read-only binding to the contract
	BuilderStakingTransactor // Write-only binding to the contract
	BuilderStakingFilterer   // Log filterer for contract events
}

// BuilderStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type BuilderStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BuilderStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BuilderStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BuilderStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BuilderStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BuilderStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BuilderStakingSession struct {
	Contract     *BuilderStaking   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BuilderStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BuilderStakingCallerSession struct {
	Contract *BuilderStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// BuilderStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BuilderStakingTransactorSession struct {
	Contract     *BuilderStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// BuilderStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type BuilderStakingRaw struct {
	Contract *BuilderStaking // Generic contract binding to access the raw methods on
}

// BuilderStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BuilderStakingCallerRaw struct {
	Contract *BuilderStakingCaller // Generic read-only contract binding to access the raw methods on
}

// BuilderStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BuilderStakingTransactorRaw struct {
	Contract *BuilderStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBuilderStaking creates a new instance of BuilderStaking, bound to a specific deployed contract.
func NewBuilderStaking(address common.Address, backend bind.ContractBackend) (*BuilderStaking, error) {
	contract, err := bindBuilderStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BuilderStaking{BuilderStakingCaller: BuilderStakingCaller{contract: contract}, BuilderStakingTransactor: BuilderStakingTransactor{contract: contract}, BuilderStakingFilterer: BuilderStakingFilterer{contract: contract}}, nil
}

// NewBuilderStakingCaller creates a new read-only instance of BuilderStaking, bound to a specific deployed contract.
func NewBuilderStakingCaller(address common.Address, caller bind.ContractCaller) (*BuilderStakingCaller, error) {
	contract, err := bindBuilderStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BuilderStakingCaller{contract: contract}, nil
}

// NewBuilderStakingTransactor creates a new write-only instance of BuilderStaking, bound to a specific deployed contract.
func NewBuilderStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*BuilderStakingTransactor, error) {
	contract, err := bindBuilderStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BuilderStakingTransactor{contract: contract}, nil
}

// NewBuilderStakingFilterer creates a new log filterer instance of BuilderStaking, bound to a specific deployed contract.
func NewBuilderStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*BuilderStakingFilterer, error) {
	contract, err := bindBuilderStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BuilderStakingFilterer{contract: contract}, nil
}

// bindBuilderStaking binds a generic wrapper to an already deployed contract.
func bindBuilderStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BuilderStakingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BuilderStaking *BuilderStakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BuilderStaking.Contract.BuilderStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BuilderStaking *BuilderStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BuilderStaking.Contract.BuilderStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BuilderStaking *BuilderStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BuilderStaking.Contract.BuilderStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BuilderStaking *BuilderStakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BuilderStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BuilderStaking *BuilderStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BuilderStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BuilderStaking *BuilderStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BuilderStaking.Contract.contract.Transact(opts, method, params...)
}

// Builders is a free data retrieval call binding the contract method 0x144ddefc.
//
// Solidity: function builders(address ) view returns(uint256 minimalStake, uint256 minimalSubscriptionPeriod)
func (_BuilderStaking *BuilderStakingCaller) Builders(opts *bind.CallOpts, arg0 common.Address) (struct {
	MinimalStake              *big.Int
	MinimalSubscriptionPeriod *big.Int
}, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "builders", arg0)

	outstruct := new(struct {
		MinimalStake              *big.Int
		MinimalSubscriptionPeriod *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MinimalStake = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.MinimalSubscriptionPeriod = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Builders is a free data retrieval call binding the contract method 0x144ddefc.
//
// Solidity: function builders(address ) view returns(uint256 minimalStake, uint256 minimalSubscriptionPeriod)
func (_BuilderStaking *BuilderStakingSession) Builders(arg0 common.Address) (struct {
	MinimalStake              *big.Int
	MinimalSubscriptionPeriod *big.Int
}, error) {
	return _BuilderStaking.Contract.Builders(&_BuilderStaking.CallOpts, arg0)
}

// Builders is a free data retrieval call binding the contract method 0x144ddefc.
//
// Solidity: function builders(address ) view returns(uint256 minimalStake, uint256 minimalSubscriptionPeriod)
func (_BuilderStaking *BuilderStakingCallerSession) Builders(arg0 common.Address) (struct {
	MinimalStake              *big.Int
	MinimalSubscriptionPeriod *big.Int
}, error) {
	return _BuilderStaking.Contract.Builders(&_BuilderStaking.CallOpts, arg0)
}

// GetSubscriptionPeriod is a free data retrieval call binding the contract method 0xb0c346b8.
//
// Solidity: function getSubscriptionPeriod(address _builder, uint256 _stakeAmount) view returns(uint256)
func (_BuilderStaking *BuilderStakingCaller) GetSubscriptionPeriod(opts *bind.CallOpts, _builder common.Address, _stakeAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "getSubscriptionPeriod", _builder, _stakeAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSubscriptionPeriod is a free data retrieval call binding the contract method 0xb0c346b8.
//
// Solidity: function getSubscriptionPeriod(address _builder, uint256 _stakeAmount) view returns(uint256)
func (_BuilderStaking *BuilderStakingSession) GetSubscriptionPeriod(_builder common.Address, _stakeAmount *big.Int) (*big.Int, error) {
	return _BuilderStaking.Contract.GetSubscriptionPeriod(&_BuilderStaking.CallOpts, _builder, _stakeAmount)
}

// GetSubscriptionPeriod is a free data retrieval call binding the contract method 0xb0c346b8.
//
// Solidity: function getSubscriptionPeriod(address _builder, uint256 _stakeAmount) view returns(uint256)
func (_BuilderStaking *BuilderStakingCallerSession) GetSubscriptionPeriod(_builder common.Address, _stakeAmount *big.Int) (*big.Int, error) {
	return _BuilderStaking.Contract.GetSubscriptionPeriod(&_BuilderStaking.CallOpts, _builder, _stakeAmount)
}

// HasMinimalStake is a free data retrieval call binding the contract method 0x0ff55db7.
//
// Solidity: function hasMinimalStake(address _builder, bytes32 _commitment) view returns(bool)
func (_BuilderStaking *BuilderStakingCaller) HasMinimalStake(opts *bind.CallOpts, _builder common.Address, _commitment [32]byte) (bool, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "hasMinimalStake", _builder, _commitment)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasMinimalStake is a free data retrieval call binding the contract method 0x0ff55db7.
//
// Solidity: function hasMinimalStake(address _builder, bytes32 _commitment) view returns(bool)
func (_BuilderStaking *BuilderStakingSession) HasMinimalStake(_builder common.Address, _commitment [32]byte) (bool, error) {
	return _BuilderStaking.Contract.HasMinimalStake(&_BuilderStaking.CallOpts, _builder, _commitment)
}

// HasMinimalStake is a free data retrieval call binding the contract method 0x0ff55db7.
//
// Solidity: function hasMinimalStake(address _builder, bytes32 _commitment) view returns(bool)
func (_BuilderStaking *BuilderStakingCallerSession) HasMinimalStake(_builder common.Address, _commitment [32]byte) (bool, error) {
	return _BuilderStaking.Contract.HasMinimalStake(&_BuilderStaking.CallOpts, _builder, _commitment)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BuilderStaking *BuilderStakingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BuilderStaking *BuilderStakingSession) Owner() (common.Address, error) {
	return _BuilderStaking.Contract.Owner(&_BuilderStaking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BuilderStaking *BuilderStakingCallerSession) Owner() (common.Address, error) {
	return _BuilderStaking.Contract.Owner(&_BuilderStaking.CallOpts)
}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(uint256 subscriptionEnd, uint256 stake)
func (_BuilderStaking *BuilderStakingCaller) Stakes(opts *bind.CallOpts, arg0 [32]byte) (struct {
	SubscriptionEnd *big.Int
	Stake           *big.Int
}, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "stakes", arg0)

	outstruct := new(struct {
		SubscriptionEnd *big.Int
		Stake           *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubscriptionEnd = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Stake = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(uint256 subscriptionEnd, uint256 stake)
func (_BuilderStaking *BuilderStakingSession) Stakes(arg0 [32]byte) (struct {
	SubscriptionEnd *big.Int
	Stake           *big.Int
}, error) {
	return _BuilderStaking.Contract.Stakes(&_BuilderStaking.CallOpts, arg0)
}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(uint256 subscriptionEnd, uint256 stake)
func (_BuilderStaking *BuilderStakingCallerSession) Stakes(arg0 [32]byte) (struct {
	SubscriptionEnd *big.Int
	Stake           *big.Int
}, error) {
	return _BuilderStaking.Contract.Stakes(&_BuilderStaking.CallOpts, arg0)
}

// TimeLocks is a free data retrieval call binding the contract method 0x2c479711.
//
// Solidity: function timeLocks(address , uint256 ) view returns(uint256 initialAmount, uint256 remainingAmount, uint256 startTime, uint256 lockDuration)
func (_BuilderStaking *BuilderStakingCaller) TimeLocks(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	InitialAmount   *big.Int
	RemainingAmount *big.Int
	StartTime       *big.Int
	LockDuration    *big.Int
}, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "timeLocks", arg0, arg1)

	outstruct := new(struct {
		InitialAmount   *big.Int
		RemainingAmount *big.Int
		StartTime       *big.Int
		LockDuration    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.InitialAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.RemainingAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartTime = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LockDuration = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// TimeLocks is a free data retrieval call binding the contract method 0x2c479711.
//
// Solidity: function timeLocks(address , uint256 ) view returns(uint256 initialAmount, uint256 remainingAmount, uint256 startTime, uint256 lockDuration)
func (_BuilderStaking *BuilderStakingSession) TimeLocks(arg0 common.Address, arg1 *big.Int) (struct {
	InitialAmount   *big.Int
	RemainingAmount *big.Int
	StartTime       *big.Int
	LockDuration    *big.Int
}, error) {
	return _BuilderStaking.Contract.TimeLocks(&_BuilderStaking.CallOpts, arg0, arg1)
}

// TimeLocks is a free data retrieval call binding the contract method 0x2c479711.
//
// Solidity: function timeLocks(address , uint256 ) view returns(uint256 initialAmount, uint256 remainingAmount, uint256 startTime, uint256 lockDuration)
func (_BuilderStaking *BuilderStakingCallerSession) TimeLocks(arg0 common.Address, arg1 *big.Int) (struct {
	InitialAmount   *big.Int
	RemainingAmount *big.Int
	StartTime       *big.Int
	LockDuration    *big.Int
}, error) {
	return _BuilderStaking.Contract.TimeLocks(&_BuilderStaking.CallOpts, arg0, arg1)
}

// TimeLocksCount is a free data retrieval call binding the contract method 0x2f7ce4ad.
//
// Solidity: function timeLocksCount(address _address) view returns(uint256)
func (_BuilderStaking *BuilderStakingCaller) TimeLocksCount(opts *bind.CallOpts, _address common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "timeLocksCount", _address)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TimeLocksCount is a free data retrieval call binding the contract method 0x2f7ce4ad.
//
// Solidity: function timeLocksCount(address _address) view returns(uint256)
func (_BuilderStaking *BuilderStakingSession) TimeLocksCount(_address common.Address) (*big.Int, error) {
	return _BuilderStaking.Contract.TimeLocksCount(&_BuilderStaking.CallOpts, _address)
}

// TimeLocksCount is a free data retrieval call binding the contract method 0x2f7ce4ad.
//
// Solidity: function timeLocksCount(address _address) view returns(uint256)
func (_BuilderStaking *BuilderStakingCallerSession) TimeLocksCount(_address common.Address) (*big.Int, error) {
	return _BuilderStaking.Contract.TimeLocksCount(&_BuilderStaking.CallOpts, _address)
}

// WithdrawableAmount is a free data retrieval call binding the contract method 0x951303f5.
//
// Solidity: function withdrawableAmount() view returns(uint256)
func (_BuilderStaking *BuilderStakingCaller) WithdrawableAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "withdrawableAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawableAmount is a free data retrieval call binding the contract method 0x951303f5.
//
// Solidity: function withdrawableAmount() view returns(uint256)
func (_BuilderStaking *BuilderStakingSession) WithdrawableAmount() (*big.Int, error) {
	return _BuilderStaking.Contract.WithdrawableAmount(&_BuilderStaking.CallOpts)
}

// WithdrawableAmount is a free data retrieval call binding the contract method 0x951303f5.
//
// Solidity: function withdrawableAmount() view returns(uint256)
func (_BuilderStaking *BuilderStakingCallerSession) WithdrawableAmount() (*big.Int, error) {
	return _BuilderStaking.Contract.WithdrawableAmount(&_BuilderStaking.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xb9e1aa03.
//
// Solidity: function deposit(address _builder, bytes32 _commitment) payable returns()
func (_BuilderStaking *BuilderStakingTransactor) Deposit(opts *bind.TransactOpts, _builder common.Address, _commitment [32]byte) (*types.Transaction, error) {
	return _BuilderStaking.contract.Transact(opts, "deposit", _builder, _commitment)
}

// Deposit is a paid mutator transaction binding the contract method 0xb9e1aa03.
//
// Solidity: function deposit(address _builder, bytes32 _commitment) payable returns()
func (_BuilderStaking *BuilderStakingSession) Deposit(_builder common.Address, _commitment [32]byte) (*types.Transaction, error) {
	return _BuilderStaking.Contract.Deposit(&_BuilderStaking.TransactOpts, _builder, _commitment)
}

// Deposit is a paid mutator transaction binding the contract method 0xb9e1aa03.
//
// Solidity: function deposit(address _builder, bytes32 _commitment) payable returns()
func (_BuilderStaking *BuilderStakingTransactorSession) Deposit(_builder common.Address, _commitment [32]byte) (*types.Transaction, error) {
	return _BuilderStaking.Contract.Deposit(&_BuilderStaking.TransactOpts, _builder, _commitment)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BuilderStaking *BuilderStakingTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BuilderStaking.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BuilderStaking *BuilderStakingSession) RenounceOwnership() (*types.Transaction, error) {
	return _BuilderStaking.Contract.RenounceOwnership(&_BuilderStaking.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BuilderStaking *BuilderStakingTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BuilderStaking.Contract.RenounceOwnership(&_BuilderStaking.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BuilderStaking *BuilderStakingTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BuilderStaking.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BuilderStaking *BuilderStakingSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BuilderStaking.Contract.TransferOwnership(&_BuilderStaking.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BuilderStaking *BuilderStakingTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BuilderStaking.Contract.TransferOwnership(&_BuilderStaking.TransactOpts, newOwner)
}

// UpdateBuilder is a paid mutator transaction binding the contract method 0x963a40e2.
//
// Solidity: function updateBuilder(uint256 _minimalStake, uint256 _minimalSubscriptionPeriod) returns()
func (_BuilderStaking *BuilderStakingTransactor) UpdateBuilder(opts *bind.TransactOpts, _minimalStake *big.Int, _minimalSubscriptionPeriod *big.Int) (*types.Transaction, error) {
	return _BuilderStaking.contract.Transact(opts, "updateBuilder", _minimalStake, _minimalSubscriptionPeriod)
}

// UpdateBuilder is a paid mutator transaction binding the contract method 0x963a40e2.
//
// Solidity: function updateBuilder(uint256 _minimalStake, uint256 _minimalSubscriptionPeriod) returns()
func (_BuilderStaking *BuilderStakingSession) UpdateBuilder(_minimalStake *big.Int, _minimalSubscriptionPeriod *big.Int) (*types.Transaction, error) {
	return _BuilderStaking.Contract.UpdateBuilder(&_BuilderStaking.TransactOpts, _minimalStake, _minimalSubscriptionPeriod)
}

// UpdateBuilder is a paid mutator transaction binding the contract method 0x963a40e2.
//
// Solidity: function updateBuilder(uint256 _minimalStake, uint256 _minimalSubscriptionPeriod) returns()
func (_BuilderStaking *BuilderStakingTransactorSession) UpdateBuilder(_minimalStake *big.Int, _minimalSubscriptionPeriod *big.Int) (*types.Transaction, error) {
	return _BuilderStaking.Contract.UpdateBuilder(&_BuilderStaking.TransactOpts, _minimalStake, _minimalSubscriptionPeriod)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_BuilderStaking *BuilderStakingTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BuilderStaking.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_BuilderStaking *BuilderStakingSession) Withdraw() (*types.Transaction, error) {
	return _BuilderStaking.Contract.Withdraw(&_BuilderStaking.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_BuilderStaking *BuilderStakingTransactorSession) Withdraw() (*types.Transaction, error) {
	return _BuilderStaking.Contract.Withdraw(&_BuilderStaking.TransactOpts)
}

// BuilderStakingBuilderUpdatedIterator is returned from FilterBuilderUpdated and is used to iterate over the raw logs and unpacked data for BuilderUpdated events raised by the BuilderStaking contract.
type BuilderStakingBuilderUpdatedIterator struct {
	Event *BuilderStakingBuilderUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BuilderStakingBuilderUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuilderStakingBuilderUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BuilderStakingBuilderUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BuilderStakingBuilderUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuilderStakingBuilderUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuilderStakingBuilderUpdated represents a BuilderUpdated event raised by the BuilderStaking contract.
type BuilderStakingBuilderUpdated struct {
	Builder                  common.Address
	MinimalStake             *big.Int
	MinimalSubsriptionPeriod *big.Int
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterBuilderUpdated is a free log retrieval operation binding the contract event 0x676767c63431a9f3ee32e335f01479a12fadf779b27ff442e88722fb80a77c7e.
//
// Solidity: event BuilderUpdated(address builder, uint256 minimalStake, uint256 minimalSubsriptionPeriod)
func (_BuilderStaking *BuilderStakingFilterer) FilterBuilderUpdated(opts *bind.FilterOpts) (*BuilderStakingBuilderUpdatedIterator, error) {

	logs, sub, err := _BuilderStaking.contract.FilterLogs(opts, "BuilderUpdated")
	if err != nil {
		return nil, err
	}
	return &BuilderStakingBuilderUpdatedIterator{contract: _BuilderStaking.contract, event: "BuilderUpdated", logs: logs, sub: sub}, nil
}

// WatchBuilderUpdated is a free log subscription operation binding the contract event 0x676767c63431a9f3ee32e335f01479a12fadf779b27ff442e88722fb80a77c7e.
//
// Solidity: event BuilderUpdated(address builder, uint256 minimalStake, uint256 minimalSubsriptionPeriod)
func (_BuilderStaking *BuilderStakingFilterer) WatchBuilderUpdated(opts *bind.WatchOpts, sink chan<- *BuilderStakingBuilderUpdated) (event.Subscription, error) {

	logs, sub, err := _BuilderStaking.contract.WatchLogs(opts, "BuilderUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuilderStakingBuilderUpdated)
				if err := _BuilderStaking.contract.UnpackLog(event, "BuilderUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBuilderUpdated is a log parse operation binding the contract event 0x676767c63431a9f3ee32e335f01479a12fadf779b27ff442e88722fb80a77c7e.
//
// Solidity: event BuilderUpdated(address builder, uint256 minimalStake, uint256 minimalSubsriptionPeriod)
func (_BuilderStaking *BuilderStakingFilterer) ParseBuilderUpdated(log types.Log) (*BuilderStakingBuilderUpdated, error) {
	event := new(BuilderStakingBuilderUpdated)
	if err := _BuilderStaking.contract.UnpackLog(event, "BuilderUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BuilderStakingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BuilderStaking contract.
type BuilderStakingOwnershipTransferredIterator struct {
	Event *BuilderStakingOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BuilderStakingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuilderStakingOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BuilderStakingOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BuilderStakingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuilderStakingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuilderStakingOwnershipTransferred represents a OwnershipTransferred event raised by the BuilderStaking contract.
type BuilderStakingOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BuilderStaking *BuilderStakingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BuilderStakingOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BuilderStaking.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BuilderStakingOwnershipTransferredIterator{contract: _BuilderStaking.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BuilderStaking *BuilderStakingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BuilderStakingOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BuilderStaking.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuilderStakingOwnershipTransferred)
				if err := _BuilderStaking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BuilderStaking *BuilderStakingFilterer) ParseOwnershipTransferred(log types.Log) (*BuilderStakingOwnershipTransferred, error) {
	event := new(BuilderStakingOwnershipTransferred)
	if err := _BuilderStaking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BuilderStakingStakeUpdatedIterator is returned from FilterStakeUpdated and is used to iterate over the raw logs and unpacked data for StakeUpdated events raised by the BuilderStaking contract.
type BuilderStakingStakeUpdatedIterator struct {
	Event *BuilderStakingStakeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BuilderStakingStakeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuilderStakingStakeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BuilderStakingStakeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BuilderStakingStakeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuilderStakingStakeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuilderStakingStakeUpdated represents a StakeUpdated event raised by the BuilderStaking contract.
type BuilderStakingStakeUpdated struct {
	Builder         common.Address
	Commitment      [32]byte
	Stake           *big.Int
	SubscriptionEnd *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterStakeUpdated is a free log retrieval operation binding the contract event 0x336b0d4818c64063eba1267244b015a6454a97222039ba75b092b2a6c522d334.
//
// Solidity: event StakeUpdated(address builder, bytes32 commitment, uint256 stake, uint256 subscriptionEnd)
func (_BuilderStaking *BuilderStakingFilterer) FilterStakeUpdated(opts *bind.FilterOpts) (*BuilderStakingStakeUpdatedIterator, error) {

	logs, sub, err := _BuilderStaking.contract.FilterLogs(opts, "StakeUpdated")
	if err != nil {
		return nil, err
	}
	return &BuilderStakingStakeUpdatedIterator{contract: _BuilderStaking.contract, event: "StakeUpdated", logs: logs, sub: sub}, nil
}

// WatchStakeUpdated is a free log subscription operation binding the contract event 0x336b0d4818c64063eba1267244b015a6454a97222039ba75b092b2a6c522d334.
//
// Solidity: event StakeUpdated(address builder, bytes32 commitment, uint256 stake, uint256 subscriptionEnd)
func (_BuilderStaking *BuilderStakingFilterer) WatchStakeUpdated(opts *bind.WatchOpts, sink chan<- *BuilderStakingStakeUpdated) (event.Subscription, error) {

	logs, sub, err := _BuilderStaking.contract.WatchLogs(opts, "StakeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuilderStakingStakeUpdated)
				if err := _BuilderStaking.contract.UnpackLog(event, "StakeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStakeUpdated is a log parse operation binding the contract event 0x336b0d4818c64063eba1267244b015a6454a97222039ba75b092b2a6c522d334.
//
// Solidity: event StakeUpdated(address builder, bytes32 commitment, uint256 stake, uint256 subscriptionEnd)
func (_BuilderStaking *BuilderStakingFilterer) ParseStakeUpdated(log types.Log) (*BuilderStakingStakeUpdated, error) {
	event := new(BuilderStakingStakeUpdated)
	if err := _BuilderStaking.contract.UnpackLog(event, "StakeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BuilderStakingWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the BuilderStaking contract.
type BuilderStakingWithdrawalIterator struct {
	Event *BuilderStakingWithdrawal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BuilderStakingWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuilderStakingWithdrawal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BuilderStakingWithdrawal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BuilderStakingWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuilderStakingWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuilderStakingWithdrawal represents a Withdrawal event raised by the BuilderStaking contract.
type BuilderStakingWithdrawal struct {
	Builder common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address builder, uint256 amount)
func (_BuilderStaking *BuilderStakingFilterer) FilterWithdrawal(opts *bind.FilterOpts) (*BuilderStakingWithdrawalIterator, error) {

	logs, sub, err := _BuilderStaking.contract.FilterLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return &BuilderStakingWithdrawalIterator{contract: _BuilderStaking.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address builder, uint256 amount)
func (_BuilderStaking *BuilderStakingFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *BuilderStakingWithdrawal) (event.Subscription, error) {

	logs, sub, err := _BuilderStaking.contract.WatchLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuilderStakingWithdrawal)
				if err := _BuilderStaking.contract.UnpackLog(event, "Withdrawal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawal is a log parse operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address builder, uint256 amount)
func (_BuilderStaking *BuilderStakingFilterer) ParseWithdrawal(log types.Log) (*BuilderStakingWithdrawal, error) {
	event := new(BuilderStakingWithdrawal)
	if err := _BuilderStaking.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
