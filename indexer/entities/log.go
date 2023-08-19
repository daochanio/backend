package entities

import "math/big"

type Log struct {
	blockNumber   *big.Int
	transactionId string
	index         uint32
}

func NewLog(blockNumber *big.Int, transactionId string, index uint32) Log {
	return Log{
		blockNumber,
		transactionId,
		index,
	}
}

func (e Log) BlockNumber() *big.Int {
	return e.blockNumber
}

func (e Log) TransactionId() string {
	return e.transactionId
}

func (e Log) Index() uint32 {
	return e.index
}
