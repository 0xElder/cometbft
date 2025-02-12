package state

import (
	"time"

	"github.com/cometbft/cometbft/types"
)

// This file contains utility functions for the state package.
const (
	// msgProcDelayRatio is the ratio of block time to message processing delay
	// assuming block time lost in execution and network delays
	MsgProcDelayRatio = 10

	// deflectionRation is the ratio of the ideal block time to the threshold for the difference
	DeflectionRatio = 5
)

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

	msgProcDelay := blockTime / MsgProcDelayRatio

	idealNextBlockDelay := blockTime - msgProcDelay

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

	// calculate the threshold for the difference using deflection ratio
	threshold := idealNextBlockDelay / DeflectionRatio

	// check if the absolute value of the difference is more than the threshold
	if diff.Abs() > threshold {
		if diff > 0 {
			return idealNextBlockDelay - threshold
		} else {
			return idealNextBlockDelay + threshold
		}
	}

	// return the ideal block delay adjusted by the difference
	return idealNextBlockDelay - diff
}
