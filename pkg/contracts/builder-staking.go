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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"builder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minimalStake\",\"type\":\"uint256\"}],\"name\":\"MinimalStakeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"builder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"}],\"name\":\"StakeUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_builder\",\"type\":\"address\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_searcher\",\"type\":\"address\"}],\"name\":\"getStakeAsBuilder\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_builder\",\"type\":\"address\"}],\"name\":\"getStakeAsSearcher\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_builder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_searcher\",\"type\":\"address\"}],\"name\":\"hasMinimalStake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"minimalStakes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minimalStake\",\"type\":\"uint256\"}],\"name\":\"setMinimalStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_builder\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// GetStakeAsBuilder is a free data retrieval call binding the contract method 0x5785bbd5.
//
// Solidity: function getStakeAsBuilder(address _searcher) view returns(uint256)
func (_BuilderStaking *BuilderStakingCaller) GetStakeAsBuilder(opts *bind.CallOpts, _searcher common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "getStakeAsBuilder", _searcher)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakeAsBuilder is a free data retrieval call binding the contract method 0x5785bbd5.
//
// Solidity: function getStakeAsBuilder(address _searcher) view returns(uint256)
func (_BuilderStaking *BuilderStakingSession) GetStakeAsBuilder(_searcher common.Address) (*big.Int, error) {
	return _BuilderStaking.Contract.GetStakeAsBuilder(&_BuilderStaking.CallOpts, _searcher)
}

// GetStakeAsBuilder is a free data retrieval call binding the contract method 0x5785bbd5.
//
// Solidity: function getStakeAsBuilder(address _searcher) view returns(uint256)
func (_BuilderStaking *BuilderStakingCallerSession) GetStakeAsBuilder(_searcher common.Address) (*big.Int, error) {
	return _BuilderStaking.Contract.GetStakeAsBuilder(&_BuilderStaking.CallOpts, _searcher)
}

// GetStakeAsSearcher is a free data retrieval call binding the contract method 0xe75fa583.
//
// Solidity: function getStakeAsSearcher(address _builder) view returns(uint256)
func (_BuilderStaking *BuilderStakingCaller) GetStakeAsSearcher(opts *bind.CallOpts, _builder common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "getStakeAsSearcher", _builder)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakeAsSearcher is a free data retrieval call binding the contract method 0xe75fa583.
//
// Solidity: function getStakeAsSearcher(address _builder) view returns(uint256)
func (_BuilderStaking *BuilderStakingSession) GetStakeAsSearcher(_builder common.Address) (*big.Int, error) {
	return _BuilderStaking.Contract.GetStakeAsSearcher(&_BuilderStaking.CallOpts, _builder)
}

// GetStakeAsSearcher is a free data retrieval call binding the contract method 0xe75fa583.
//
// Solidity: function getStakeAsSearcher(address _builder) view returns(uint256)
func (_BuilderStaking *BuilderStakingCallerSession) GetStakeAsSearcher(_builder common.Address) (*big.Int, error) {
	return _BuilderStaking.Contract.GetStakeAsSearcher(&_BuilderStaking.CallOpts, _builder)
}

// HasMinimalStake is a free data retrieval call binding the contract method 0xa6419e28.
//
// Solidity: function hasMinimalStake(address _builder, address _searcher) view returns(bool)
func (_BuilderStaking *BuilderStakingCaller) HasMinimalStake(opts *bind.CallOpts, _builder common.Address, _searcher common.Address) (bool, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "hasMinimalStake", _builder, _searcher)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasMinimalStake is a free data retrieval call binding the contract method 0xa6419e28.
//
// Solidity: function hasMinimalStake(address _builder, address _searcher) view returns(bool)
func (_BuilderStaking *BuilderStakingSession) HasMinimalStake(_builder common.Address, _searcher common.Address) (bool, error) {
	return _BuilderStaking.Contract.HasMinimalStake(&_BuilderStaking.CallOpts, _builder, _searcher)
}

// HasMinimalStake is a free data retrieval call binding the contract method 0xa6419e28.
//
// Solidity: function hasMinimalStake(address _builder, address _searcher) view returns(bool)
func (_BuilderStaking *BuilderStakingCallerSession) HasMinimalStake(_builder common.Address, _searcher common.Address) (bool, error) {
	return _BuilderStaking.Contract.HasMinimalStake(&_BuilderStaking.CallOpts, _builder, _searcher)
}

// MinimalStakes is a free data retrieval call binding the contract method 0x136fc2b3.
//
// Solidity: function minimalStakes(address ) view returns(uint256)
func (_BuilderStaking *BuilderStakingCaller) MinimalStakes(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BuilderStaking.contract.Call(opts, &out, "minimalStakes", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinimalStakes is a free data retrieval call binding the contract method 0x136fc2b3.
//
// Solidity: function minimalStakes(address ) view returns(uint256)
func (_BuilderStaking *BuilderStakingSession) MinimalStakes(arg0 common.Address) (*big.Int, error) {
	return _BuilderStaking.Contract.MinimalStakes(&_BuilderStaking.CallOpts, arg0)
}

// MinimalStakes is a free data retrieval call binding the contract method 0x136fc2b3.
//
// Solidity: function minimalStakes(address ) view returns(uint256)
func (_BuilderStaking *BuilderStakingCallerSession) MinimalStakes(arg0 common.Address) (*big.Int, error) {
	return _BuilderStaking.Contract.MinimalStakes(&_BuilderStaking.CallOpts, arg0)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _builder) payable returns()
func (_BuilderStaking *BuilderStakingTransactor) Deposit(opts *bind.TransactOpts, _builder common.Address) (*types.Transaction, error) {
	return _BuilderStaking.contract.Transact(opts, "deposit", _builder)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _builder) payable returns()
func (_BuilderStaking *BuilderStakingSession) Deposit(_builder common.Address) (*types.Transaction, error) {
	return _BuilderStaking.Contract.Deposit(&_BuilderStaking.TransactOpts, _builder)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _builder) payable returns()
func (_BuilderStaking *BuilderStakingTransactorSession) Deposit(_builder common.Address) (*types.Transaction, error) {
	return _BuilderStaking.Contract.Deposit(&_BuilderStaking.TransactOpts, _builder)
}

// SetMinimalStake is a paid mutator transaction binding the contract method 0x3d6ec65e.
//
// Solidity: function setMinimalStake(uint256 _minimalStake) returns()
func (_BuilderStaking *BuilderStakingTransactor) SetMinimalStake(opts *bind.TransactOpts, _minimalStake *big.Int) (*types.Transaction, error) {
	return _BuilderStaking.contract.Transact(opts, "setMinimalStake", _minimalStake)
}

// SetMinimalStake is a paid mutator transaction binding the contract method 0x3d6ec65e.
//
// Solidity: function setMinimalStake(uint256 _minimalStake) returns()
func (_BuilderStaking *BuilderStakingSession) SetMinimalStake(_minimalStake *big.Int) (*types.Transaction, error) {
	return _BuilderStaking.Contract.SetMinimalStake(&_BuilderStaking.TransactOpts, _minimalStake)
}

// SetMinimalStake is a paid mutator transaction binding the contract method 0x3d6ec65e.
//
// Solidity: function setMinimalStake(uint256 _minimalStake) returns()
func (_BuilderStaking *BuilderStakingTransactorSession) SetMinimalStake(_minimalStake *big.Int) (*types.Transaction, error) {
	return _BuilderStaking.Contract.SetMinimalStake(&_BuilderStaking.TransactOpts, _minimalStake)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _builder, uint256 _amount) returns()
func (_BuilderStaking *BuilderStakingTransactor) Withdraw(opts *bind.TransactOpts, _builder common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BuilderStaking.contract.Transact(opts, "withdraw", _builder, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _builder, uint256 _amount) returns()
func (_BuilderStaking *BuilderStakingSession) Withdraw(_builder common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BuilderStaking.Contract.Withdraw(&_BuilderStaking.TransactOpts, _builder, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _builder, uint256 _amount) returns()
func (_BuilderStaking *BuilderStakingTransactorSession) Withdraw(_builder common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BuilderStaking.Contract.Withdraw(&_BuilderStaking.TransactOpts, _builder, _amount)
}

// BuilderStakingMinimalStakeUpdatedIterator is returned from FilterMinimalStakeUpdated and is used to iterate over the raw logs and unpacked data for MinimalStakeUpdated events raised by the BuilderStaking contract.
type BuilderStakingMinimalStakeUpdatedIterator struct {
	Event *BuilderStakingMinimalStakeUpdated // Event containing the contract specifics and raw log

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
func (it *BuilderStakingMinimalStakeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuilderStakingMinimalStakeUpdated)
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
		it.Event = new(BuilderStakingMinimalStakeUpdated)
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
func (it *BuilderStakingMinimalStakeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuilderStakingMinimalStakeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuilderStakingMinimalStakeUpdated represents a MinimalStakeUpdated event raised by the BuilderStaking contract.
type BuilderStakingMinimalStakeUpdated struct {
	Builder      common.Address
	MinimalStake *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMinimalStakeUpdated is a free log retrieval operation binding the contract event 0x5980f614bd0d5d26f55eef7ec350fc66b080cd09bd38c87bcac4c4c7853d76f7.
//
// Solidity: event MinimalStakeUpdated(address builder, uint256 minimalStake)
func (_BuilderStaking *BuilderStakingFilterer) FilterMinimalStakeUpdated(opts *bind.FilterOpts) (*BuilderStakingMinimalStakeUpdatedIterator, error) {

	logs, sub, err := _BuilderStaking.contract.FilterLogs(opts, "MinimalStakeUpdated")
	if err != nil {
		return nil, err
	}
	return &BuilderStakingMinimalStakeUpdatedIterator{contract: _BuilderStaking.contract, event: "MinimalStakeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinimalStakeUpdated is a free log subscription operation binding the contract event 0x5980f614bd0d5d26f55eef7ec350fc66b080cd09bd38c87bcac4c4c7853d76f7.
//
// Solidity: event MinimalStakeUpdated(address builder, uint256 minimalStake)
func (_BuilderStaking *BuilderStakingFilterer) WatchMinimalStakeUpdated(opts *bind.WatchOpts, sink chan<- *BuilderStakingMinimalStakeUpdated) (event.Subscription, error) {

	logs, sub, err := _BuilderStaking.contract.WatchLogs(opts, "MinimalStakeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuilderStakingMinimalStakeUpdated)
				if err := _BuilderStaking.contract.UnpackLog(event, "MinimalStakeUpdated", log); err != nil {
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

// ParseMinimalStakeUpdated is a log parse operation binding the contract event 0x5980f614bd0d5d26f55eef7ec350fc66b080cd09bd38c87bcac4c4c7853d76f7.
//
// Solidity: event MinimalStakeUpdated(address builder, uint256 minimalStake)
func (_BuilderStaking *BuilderStakingFilterer) ParseMinimalStakeUpdated(log types.Log) (*BuilderStakingMinimalStakeUpdated, error) {
	event := new(BuilderStakingMinimalStakeUpdated)
	if err := _BuilderStaking.contract.UnpackLog(event, "MinimalStakeUpdated", log); err != nil {
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
	Builder  common.Address
	Searcher common.Address
	Stake    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStakeUpdated is a free log retrieval operation binding the contract event 0xe5c5e2e23207d40d5bc669072459c1fee64395c31fcde90a3d63432ace92c9ec.
//
// Solidity: event StakeUpdated(address builder, address searcher, uint256 stake)
func (_BuilderStaking *BuilderStakingFilterer) FilterStakeUpdated(opts *bind.FilterOpts) (*BuilderStakingStakeUpdatedIterator, error) {

	logs, sub, err := _BuilderStaking.contract.FilterLogs(opts, "StakeUpdated")
	if err != nil {
		return nil, err
	}
	return &BuilderStakingStakeUpdatedIterator{contract: _BuilderStaking.contract, event: "StakeUpdated", logs: logs, sub: sub}, nil
}

// WatchStakeUpdated is a free log subscription operation binding the contract event 0xe5c5e2e23207d40d5bc669072459c1fee64395c31fcde90a3d63432ace92c9ec.
//
// Solidity: event StakeUpdated(address builder, address searcher, uint256 stake)
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

// ParseStakeUpdated is a log parse operation binding the contract event 0xe5c5e2e23207d40d5bc669072459c1fee64395c31fcde90a3d63432ace92c9ec.
//
// Solidity: event StakeUpdated(address builder, address searcher, uint256 stake)
func (_BuilderStaking *BuilderStakingFilterer) ParseStakeUpdated(log types.Log) (*BuilderStakingStakeUpdated, error) {
	event := new(BuilderStakingStakeUpdated)
	if err := _BuilderStaking.contract.UnpackLog(event, "StakeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
