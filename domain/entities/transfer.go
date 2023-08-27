package entities

import "math/big"

type Transfer struct {
	fromAddress string
	toAddress   string
	amount      *big.Int
	log         Log
}

func NewTransfer(fromAddress string, toAddress string, amount *big.Int, log Log) Transfer {
	return Transfer{
		fromAddress,
		toAddress,
		amount,
		log,
	}
}

func (e Transfer) FromAddress() string {
	return e.fromAddress
}

func (e Transfer) ToAddress() string {
	return e.toAddress
}

func (e Transfer) Amount() *big.Int {
	return e.amount
}

func (e Transfer) Log() Log {
	return e.log
}
