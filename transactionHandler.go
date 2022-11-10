package main

import (
	d "bharvest/jihon/data"
	"errors"
	"log"
	"math/big"
)

/*
	[add]
	1. 계정이 존재하는 경우
	- 계산 (sender에 amount+, seq++)
	2. 계정이 존재하지 않는 경우
	- 계정 생성 및 초기 세팅(balance == amount, seq = 1)
*/
func tranAdd(tx *d.Transaction) error {
	curAccount := (accounts[tx.Sender])
	fReqAmount := new(big.Float).SetPrec(128) 						// amount 값을 big.Float 로 변환
	fReqAmount.Parse(tx.Amount, 0)

	if curAccount != nil {											// 계정이 존재하는 경우
		curAccount.Lock()

		/* 트랜잭션 처리 시작 */
		startTranProcessing()
		curAccount.AddBalance(fReqAmount)							// 밸런스 계산
		log.Println("[확인] sender 계정에 amount 가 추가되었습니다")
	} else {
		var account d.Account
		// 계정 생성
		account.Init(tx.Sender)
		accounts[account.GetAddress()] = &account
		account.Lock()

		/* 트랜잭션 처리 시작 */
		startTranProcessing()
		account.SetId(uint(len(accounts)))
		account.AddBalance(fReqAmount)
		log.Println("[확인] 신규 계정이 생성되었습니다. : " + account.GetAddress())
	}
	return nil
}

/*
[sub]
  - nonce 체크
  - lock 체크 (잠겨 있으면 대기큐, 없으면 lock)
  - 계산 (sender계정에 벨런스-, seq++)
*/
func tranSub(tx *d.Transaction) error {
	curAccount := (accounts[tx.Sender])
	curAccount.Lock()

	/* 트랜잭션 처리 시작 */
	startTranProcessing()
	fReqAmount := new(big.Float).SetPrec(128)
	fReqAmount.Parse(tx.Amount, 0)
	if fReqAmount.Cmp(curAccount.GetBalance()) > 0 {
		return errors.New(d.ERR_MSG_SEND_WRONG_AMNT)
	}

	curAccount.SubBalance(fReqAmount)
	log.Println("[확인] sender 계정에 amount 가 sub 되었습니다")
	return nil
}

/*
[send]
  - sender 의 보유량이 충분한지 체크
  - 받는사람(target) 이 존재하는지 체크
  - 받는사람(target) lock 체크
  - sender 계정에 밸런스- , seq++
  - target 계정에 밸런스+
*/
func tranSend(tx *d.Transaction) error {
	curAccount := accounts[tx.Sender]
	tgtAccount := accounts[tx.Target]
	curAccount.Lock()
	tgtAccount.Lock()
	
	/* 트랜잭션 처리 시작 */
	startTranProcessing()
	fReqAmount := new(big.Float).SetPrec(128)
	fReqAmount.Parse(tx.Amount, 0)
	if fReqAmount.Cmp(curAccount.GetBalance()) > 0 { 			// sender 의 보유량이 충분한지 체크
		return errors.New(d.ERR_MSG_SEND_WRONG_AMNT)
	}

	curAccount.SubBalance(fReqAmount)
	tgtAccount.AddBalance(fReqAmount)
	log.Println("[확인] 요청한 amount 가 target 계정에 전송되었습니다")
	return nil
}
