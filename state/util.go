package state

import (
	"time"

	"github.com/cometbft/cometbft/types"
)

// This file contains utility functions for the state package.

func CalculateDelay(store Store, currentBlock *types.Block) time.Duration {
	// load the state from the store
	state, err := store.Load()
	if err != nil {
		return 0
	}

	// load the consensus params from the store for the one block less than current block
	consensusParams, err := store.LoadConsensusParams(currentBlock.Height - 1)
	if err != nil {
		return 0
	}

	// get the block time
	blockTime := consensusParams.Block.BlockTime

	// ideal next block delay (assuming 10% of block time lost in execution and network delays)
	idealNextBlockDelay := blockTime - blockTime/10

	// get the genesis time, it is a time for the  block=1
	genesisTime := state.GenesisTime
	if genesisTime.IsZero() {
		return 0
	}

	// get the duration from the genesis time to the current block time
	gotDuration := currentBlock.Time.Sub(genesisTime)

	// calculate the expected duration
	expected := blockTime.Nanoseconds() * (int64(currentBlock.Height) - int64(1))
	expectedDuration := time.Duration(expected)

	// calculate the difference between the expected duration and the got duration
	diff := gotDuration - expectedDuration

	// calculate the threshold for the difference
	threshold := idealNextBlockDelay / 5

	// check if the absolute value of the difference is more than 20% of the ideal block delay
	if diff.Abs() > idealNextBlockDelay/5 {
		if diff > 0 {
			return idealNextBlockDelay - threshold
		} else {
			return idealNextBlockDelay + threshold
		}
	}

	// return the ideal block delay adjusted by the difference
	return idealNextBlockDelay - diff
}
