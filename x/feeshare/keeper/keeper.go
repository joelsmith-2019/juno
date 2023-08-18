package keeper

import (
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmosContracts/juno/v16/x/feeshare/types"
	revtypes "github.com/CosmosContracts/juno/v16/x/feeshare/types"
)

// Keeper of this module maintains collections of feeshares for contracts
// registered to receive transaction fees.
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	bankKeeper    revtypes.BankKeeper
	wasmKeeper    wasmkeeper.Keeper
	accountKeeper revtypes.AccountKeeper

	feeCollectorName string

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper creates new instances of the fees Keeper
func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	bk revtypes.BankKeeper,
	wk wasmkeeper.Keeper,
	ak revtypes.AccountKeeper,
	feeCollector string,
	authority string,
) Keeper {
	return Keeper{
		storeKey:         storeKey,
		cdc:              cdc,
		bankKeeper:       bk,
		wasmKeeper:       wk,
		accountKeeper:    ak,
		feeCollectorName: feeCollector,
		authority:        authority,
	}
}

// GetAuthority returns the x/feeshare module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", revtypes.ModuleName))
}

// Get if the counter is running
func (k Keeper) IsRunning(ctx sdk.Context) (running types.RunningData) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has([]byte("running")) {
		return running
	}

	runningBytes := store.Get([]byte("running"))

	k.cdc.MustUnmarshal(runningBytes, &running)

	return running
}

// Set the count status to running
func (k Keeper) StartCounter(ctx sdk.Context, _ *types.MsgStartCounter) {
	store := ctx.KVStore(k.storeKey)

	runningData := types.RunningData{
		Running: true,
	}

	runningBytes := k.cdc.MustMarshal(&runningData)

	store.Set([]byte("running"), runningBytes)
}

func (k Keeper) StopCounter(ctx sdk.Context, _ *types.MsgStopCounter) {
	store := ctx.KVStore(k.storeKey)

	runningData := types.RunningData{
		Running: false,
	}

	runningBytes := k.cdc.MustMarshal(&runningData)

	store.Set([]byte("running"), runningBytes)
}

// Set the current count in the KV store
func (k Keeper) GetCount(ctx sdk.Context) (count types.CountData) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has([]byte("count")) {
		return count
	}

	countBytes := store.Get([]byte("count"))

	k.cdc.MustUnmarshal(countBytes, &count)

	return count
}

// Retrieve the current count from the KV store
func (k Keeper) SetCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)

	countData := types.CountData{
		Count: count,
	}

	countBytes := k.cdc.MustMarshal(&countData)

	store.Set([]byte("count"), countBytes)
}
