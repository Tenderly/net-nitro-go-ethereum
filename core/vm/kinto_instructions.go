package vm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

func opSelfdestructKinto(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
	if interpreter.readOnly {
		return nil, ErrWriteProtection
	}

	// Declare the beneficiary variable outside the if-else block
	var beneficiary common.Address

	// Determine the beneficiary based on the block number
	if interpreter.evm.Context.BlockNumber.Cmp(common.KintoHardfork2) > 0 {
		beneficiary = common.HexToAddress(common.SelfDestructWallet)
	} else {
		beneficiaryAddr := scope.Stack.pop()
		beneficiary = common.BytesToAddress(beneficiaryAddr.Bytes())
	}

	balance := interpreter.evm.StateDB.GetBalance(scope.Contract.Address())

	interpreter.evm.StateDB.AddBalance(beneficiary, balance) // Use the beneficiary variable directly
	interpreter.evm.StateDB.SelfDestruct(scope.Contract.Address())

	if beneficiary == scope.Contract.Address() {
		// Arbitrum: calling selfdestruct(this) burns the balance
		interpreter.evm.StateDB.ExpectBalanceBurn(balance.ToBig())
	}

	if tracer := interpreter.evm.Config.Tracer; tracer != nil {
		tracer.CaptureEnter(SELFDESTRUCT, scope.Contract.Address(), beneficiary, []byte{}, 0, balance.ToBig())
		tracer.CaptureExit([]byte{}, 0, nil)
	}
	return nil, errStopToken
}

func opSelfdestruct6780Kinto(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
	if interpreter.readOnly {
		return nil, ErrWriteProtection
	}

	// Arbitrum: revert if acting account is a Stylus program
	actingAddress := scope.Contract.Address()
	if code := interpreter.evm.StateDB.GetCode(actingAddress); state.IsStylusProgram(code) {
		return nil, ErrExecutionReverted
	}

	var beneficiary common.Address
	// Determine the beneficiary based on the block number (opSelfDestruct was added in hf5)
	if interpreter.evm.Context.BlockNumber.Cmp(common.KintoHardfork5) > 0 {
		beneficiary = common.HexToAddress(common.SelfDestructWallet)
	} else {
		beneficiaryAddr := scope.Stack.pop()
		beneficiary = common.BytesToAddress(beneficiaryAddr.Bytes())
	}
	balance := interpreter.evm.StateDB.GetBalance(scope.Contract.Address())
	interpreter.evm.StateDB.SubBalance(scope.Contract.Address(), balance)
	interpreter.evm.StateDB.AddBalance(beneficiary, balance)
	interpreter.evm.StateDB.Selfdestruct6780(scope.Contract.Address())
	if tracer := interpreter.evm.Config.Tracer; tracer != nil {
		tracer.CaptureEnter(SELFDESTRUCT, scope.Contract.Address(), beneficiary, []byte{}, 0, balance.ToBig())
		tracer.CaptureExit([]byte{}, 0, nil)
	}

	return nil, errStopToken
}
