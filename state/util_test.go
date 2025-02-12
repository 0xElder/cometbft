package state_test

import (
	"errors"
	"testing"
	"time"

	"github.com/cometbft/cometbft/state"
	"github.com/cometbft/cometbft/state/mocks"
	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/assert"
)

func TestStateError(t *testing.T) {
	blockStore := &mocks.Store{}

	genesisTime := time.Now()
	blockTime := 1 * time.Second

	// empty state
	mockState := state.State{}

	consensusParams := types.ConsensusParams{
		Block: types.BlockParams{
			BlockTime: blockTime,
		},
	}

	height := int64(10)

	header := types.Header{
		Height: height,
		Time:   genesisTime.Add(9 * blockTime),
	}

	currentBlock := &types.Block{
		Header: header,
	}

	blockStore.On("Load").Return(mockState, errors.New("error"))
	blockStore.On("LoadConsensusParams", (height-1)).Return(consensusParams, nil)

	delay := state.CalculateDelay(blockStore, currentBlock)

	// when there is any error we expect the 0 because it is handled itseld later on
	expectedDelay := time.Duration(0)

	assert.Equal(t, expectedDelay, delay)
}

func TestConsensusParamError(t *testing.T) {
	blockStore := &mocks.Store{}

	genesisTime := time.Now()
	blockTime := 1 * time.Second

	mockState := state.State{
		GenesisTime: genesisTime,
	}

	// empty consensus params
	consensusParams := types.ConsensusParams{}

	height := int64(10)

	header := types.Header{
		Height: height,
		Time:   genesisTime.Add(9 * blockTime),
	}

	currentBlock := &types.Block{
		Header: header,
	}

	blockStore.On("Load").Return(mockState, nil)
	blockStore.On("LoadConsensusParams", (height-1)).Return(consensusParams, errors.New("error"))

	delay := state.CalculateDelay(blockStore, currentBlock)

	// when there is any error we expect the 0 because it is handled itseld later on
	expectedDelay := time.Duration(0)

	assert.Equal(t, expectedDelay, delay)
}

func TestCalculateDelayWhenDiffIsZero(t *testing.T) {
	blockStore := &mocks.Store{}

	genesisTime := time.Now()
	blockTime := 1 * time.Second

	mockState := state.State{
		GenesisTime: genesisTime,
	}

	consensusParams := types.ConsensusParams{
		Block: types.BlockParams{
			BlockTime: blockTime,
		},
	}

	height := int64(10)

	header := types.Header{
		Height: height,
		Time:   genesisTime.Add(9 * time.Second),
	}

	currentBlock := &types.Block{
		Header: header,
	}

	blockStore.On("Load").Return(mockState, nil)
	blockStore.On("LoadConsensusParams", (height-1)).Return(consensusParams, nil)

	delay := state.CalculateDelay(blockStore, currentBlock)

	// ideal next block delay (assuming 10% of block time lost in execution and network delays)
	idealExpectedDelay := blockTime - blockTime/state.MsgProcDelayRatio

	assert.Equal(t, idealExpectedDelay, delay)
}

func TestCalculateDelayWhenDiffIsMorePositiveThanMaxDeflection(t *testing.T) {
	blockStore := &mocks.Store{}

	genesisTime := time.Now()
	blockTime := 1 * time.Second

	mockState := state.State{
		GenesisTime: genesisTime,
	}

	consensusParams := types.ConsensusParams{
		Block: types.BlockParams{
			BlockTime: blockTime,
		},
	}

	// ideal next block delay (assuming 10% of block time lost in execution and network delays)
	idealExpectedDelay := blockTime - blockTime/state.MsgProcDelayRatio

	// maxdepliction is 20% of idealExpectedDelay
	maxDeflection := idealExpectedDelay / state.DeflectionRatio

	height := int64(10)
	addedDiff := maxDeflection + maxDeflection/2

	header := types.Header{
		Height: height,
		Time:   genesisTime.Add(9 * time.Second).Add(addedDiff),
	}

	currentBlock := &types.Block{
		Header: header,
	}

	blockStore.On("Load").Return(mockState, nil)
	blockStore.On("LoadConsensusParams", (height-1)).Return(consensusParams, nil)

	delay := state.CalculateDelay(blockStore, currentBlock)

	// as addedDiff is more than the maxDeflection we expect the idealExpectedDelay - maxDeflection
	expectedDelay := idealExpectedDelay - maxDeflection

	assert.Equal(t, expectedDelay, delay)
}

func TestCalculateDelayWhenDiffIsLessPositiveThanMaxDeflection(t *testing.T) {
	blockStore := &mocks.Store{}

	genesisTime := time.Now()
	blockTime := 1 * time.Second

	mockState := state.State{
		GenesisTime: genesisTime,
	}

	consensusParams := types.ConsensusParams{
		Block: types.BlockParams{
			BlockTime: blockTime,
		},
	}

	// ideal next block delay (assuming 10% of block time lost in execution and network delays)
	idealExpectedDelay := blockTime - blockTime/state.MsgProcDelayRatio

	// maxdepliction is 20% of idealExpectedDelay
	maxDeflection := idealExpectedDelay / state.DeflectionRatio

	height := int64(10)
	addedDiff := maxDeflection - maxDeflection/2

	header := types.Header{
		Height: height,
		Time:   genesisTime.Add(9 * time.Second).Add(addedDiff),
	}

	currentBlock := &types.Block{
		Header: header,
	}

	blockStore.On("Load").Return(mockState, nil)
	blockStore.On("LoadConsensusParams", (height-1)).Return(consensusParams, nil)

	delay := state.CalculateDelay(blockStore, currentBlock)

	// as addedDiff is less than the maxDeflection we expect the idealExpectedDelay - addedDiff
	expectedDelay := idealExpectedDelay - addedDiff

	assert.Equal(t, expectedDelay, delay)
}

func TestCalculateDelayWhenDiffIsMoreNegativeThanMaxDeflection(t *testing.T) {
	blockStore := &mocks.Store{}

	genesisTime := time.Now()
	blockTime := 1 * time.Second

	mockState := state.State{
		GenesisTime: genesisTime,
	}

	consensusParams := types.ConsensusParams{
		Block: types.BlockParams{
			BlockTime: blockTime,
		},
	}

	// ideal next block delay (assuming 10% of block time lost in execution and network delays)
	idealExpectedDelay := blockTime - blockTime/state.MsgProcDelayRatio

	// maxdepliction is 20% of idealExpectedDelay
	maxDeflection := idealExpectedDelay / state.DeflectionRatio

	height := int64(10)
	addedDiff := -1 * (maxDeflection + maxDeflection/2)

	header := types.Header{
		Height: height,
		Time:   genesisTime.Add(9 * time.Second).Add(addedDiff),
	}

	currentBlock := &types.Block{
		Header: header,
	}

	blockStore.On("Load").Return(mockState, nil)
	blockStore.On("LoadConsensusParams", (height-1)).Return(consensusParams, nil)

	delay := state.CalculateDelay(blockStore, currentBlock)

	// as addedDiff is more than the maxDeflection we expect the idealExpectedDelay - maxDeflection
	expectedDelay := idealExpectedDelay + maxDeflection

	assert.Equal(t, expectedDelay, delay)
}

func TestCalculateDelayWhenDiffIsLessNegativeThanMaxDeflection(t *testing.T) {
	blockStore := &mocks.Store{}

	genesisTime := time.Now()
	blockTime := 1 * time.Second

	mockState := state.State{
		GenesisTime: genesisTime,
	}

	consensusParams := types.ConsensusParams{
		Block: types.BlockParams{
			BlockTime: blockTime,
		},
	}

	// ideal next block delay (assuming 10% of block time lost in execution and network delays)
	idealExpectedDelay := blockTime - blockTime/state.MsgProcDelayRatio

	// maxdepliction is 20% of idealExpectedDelay
	maxDeflection := idealExpectedDelay / state.DeflectionRatio

	height := int64(10)
	addedDiff := -1 * (maxDeflection - maxDeflection/2)

	header := types.Header{
		Height: height,
		Time:   genesisTime.Add(9 * time.Second).Add(addedDiff),
	}

	currentBlock := &types.Block{
		Header: header,
	}

	blockStore.On("Load").Return(mockState, nil)
	blockStore.On("LoadConsensusParams", (height-1)).Return(consensusParams, nil)

	delay := state.CalculateDelay(blockStore, currentBlock)

	// as addedDiff is less than the maxDeflection we expect the idealExpectedDelay - addedDiff
	expectedDelay := idealExpectedDelay - addedDiff

	assert.Equal(t, expectedDelay, delay)
}
