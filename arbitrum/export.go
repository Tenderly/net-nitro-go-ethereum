package arbitrum

import (
	"context"

	"github.com/tenderly/net-nitro-go-ethereum/common/hexutil"
	"github.com/tenderly/net-nitro-go-ethereum/core"
	"github.com/tenderly/net-nitro-go-ethereum/notinternal/ethapi"
	"github.com/tenderly/net-nitro-go-ethereum/rpc"
)

type TransactionArgs = ethapi.TransactionArgs

func EstimateGas(ctx context.Context, b ethapi.Backend, args TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *ethapi.StateOverride, gasCap uint64) (hexutil.Uint64, error) {
	return ethapi.DoEstimateGas(ctx, b, args, blockNrOrHash, overrides, gasCap)
}

func NewRevertReason(result *core.ExecutionResult) error {
	return ethapi.NewRevertError(result)
}
