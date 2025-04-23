package vm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
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

	interpreter.evm.StateDB.AddBalance(beneficiary, balance, tracing.BalanceIncreaseSelfdestruct) // Use the beneficiary variable directly
	interpreter.evm.StateDB.SelfDestruct(scope.Contract.Address())

	if beneficiary == scope.Contract.Address() {
		// Arbitrum: calling selfdestruct(this) burns the balance
		interpreter.evm.StateDB.ExpectBalanceBurn(balance.ToBig())
	}

	if tracer := interpreter.evm.Config.Tracer; tracer != nil {
		if tracer.OnEnter != nil {
			tracer.OnEnter(interpreter.evm.depth, byte(SELFDESTRUCT), scope.Contract.Address(), beneficiary, []byte{}, 0, balance.ToBig())
		}
		if tracer.OnExit != nil {
			tracer.OnExit(interpreter.evm.depth, []byte{}, 0, nil, false)
		}
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
	interpreter.evm.StateDB.SubBalance(scope.Contract.Address(), balance, tracing.BalanceDecreaseSelfdestruct)
	interpreter.evm.StateDB.AddBalance(beneficiary, balance, tracing.BalanceIncreaseSelfdestruct)
	interpreter.evm.StateDB.SelfDestruct6780(scope.Contract.Address())
	if tracer := interpreter.evm.Config.Tracer; tracer != nil {
		if tracer.OnEnter != nil {
			tracer.OnEnter(interpreter.evm.depth, byte(SELFDESTRUCT), scope.Contract.Address(), beneficiary, []byte{}, 0, balance.ToBig())
		}
		if tracer.OnExit != nil {
			tracer.OnExit(interpreter.evm.depth, []byte{}, 0, nil, false)
		}
	}

	return nil, errStopToken
}
