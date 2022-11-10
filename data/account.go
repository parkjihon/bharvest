package data

import (
	"math/big"
)

/*
계정(account) 자료구조
*/
type Account struct { // 계정 구조체
	id       uint
	address  string
	seq      uint
	balance  *big.Float
	isLocked bool
}

func (a *Account) Init(addr string) {
	a.address	= addr
	a.seq		= 0
	a.balance	= new(big.Float)
	a.isLocked	= false
}

func (a *Account) SetAddress(addr string) {
	a.address = addr
}

func (a *Account) SetId(id uint) {
	a.id = id
}

func (a *Account) AddSeq() {
	a.seq++
}

func (a *Account) AddBalance(amount *big.Float) {
	a.balance.Add(amount, a.balance)
}

func (a *Account) SubBalance(amount *big.Float) {
	a.balance.Sub(a.balance, amount)
}

func (a *Account) Lock() {
	a.isLocked = true
}

func (a *Account) UnLock() {
	a.isLocked = false
}

func (a *Account) GetAddress() string {
	return a.address
}

func (a *Account) GetId() uint {
	return a.id
}

func (a *Account) GetSeq() uint {
	return a.seq
}

func (a *Account) GetBalance() *big.Float {
	return a.balance
}

func (a *Account) IsLock() bool {
	return a.isLocked
}












