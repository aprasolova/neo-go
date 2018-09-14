package mempool

import "time"

type Config struct {

	// This is the maximum amount
	// of transactions that we will allow in the mempool
	MaxNumOfTX uint64

	// FreeTX defines the maximum amount of free txs that can be in the mempool at one time
	//  Default is 20
	FreeTX uint32

	// MinTXFee is a number in Fixed8 format. If set at 1GAS, minTXFee would equal 1e8
	// The mineTXFee is used to set the floor, it defaults to zero meaning we will allow all transactions
	// with a fee of 0 or more
	MinTXFee uint64

	// MaxTXSize is the maximum number of bytes a tx can have to be entered into the pool
	MaxTXSize uint64

	// // DynamicFee will allow the MiXTXFee to be dynamic
	// DynamicFee bool

	// // FeeEstimator is a stub function which takes in the current fee,
	// // the number of txs already in the mempool and returns a new TXFee
	// FeeEstimator func(currentTXFee int64, numOfTX int64) int64

	// TXTTL is the duration to which we should keep an item in the mempool before removing it
	// HMM: Should this be amount of blocks instead? For when blocks take time a long time
	// to process?
	TXTTL time.Duration

	// SigLimit is the maximum amount of signatures
	// that we will allow a tx to have, default will be 20
	SigLimit uint8

	// // RateLimit set to true will limit the amount of TXs that are <= the MinTXFee
	// RateLimit bool

	// // TXRate is the rate at which we should accept one TX below the MinTXFee
	// // E.g. If we want 10TXS every 10 seconds. TXRate = 1 Second
	// TXRate time.Time
}
