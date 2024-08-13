package eth

import (
	"context"

	"github.com/tenderly/net-nitro-go-ethereum/core"
	"github.com/tenderly/net-nitro-go-ethereum/core/state"
	"github.com/tenderly/net-nitro-go-ethereum/core/types"
	"github.com/tenderly/net-nitro-go-ethereum/core/vm"
	"github.com/tenderly/net-nitro-go-ethereum/eth/tracers"
	"github.com/tenderly/net-nitro-go-ethereum/ethdb"
)

func NewArbEthereum(
	blockchain *core.BlockChain,
	chainDb ethdb.Database,
) *Ethereum {
	return &Ethereum{
		blockchain: blockchain,
		chainDb:    chainDb,
	}
}

func (eth *Ethereum) StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64) (*core.Message, vm.BlockContext, *state.StateDB, tracers.StateReleaseFunc, error) {
	return eth.stateAtTransaction(ctx, block, txIndex, reexec)
}
