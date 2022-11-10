package main

import (
	d "bharvest/jihon/data"
	"encoding/hex"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

// seq 유효성 검증
func isValidSeq(tx *d.Transaction, tranType string) error {
	curAccount := accounts[tx.Sender]

	if curAccount == nil && strings.ToLower(tranType) == "add" {
		if tx.Seq != 1 {
			return errors.New("[에러] 올바른 Seq 값이 아닙니다. 계정 생성 시 seq 는 1 이어야 합니다")
		}
	} else {
		curSeq := curAccount.GetSeq()
		if (curSeq + 1) == tx.Seq {
			return nil
		} else {
			return errors.New(`[에러] 올바른 Seq 값이 아닙니다. 현재 Seq : ` + strconv.Itoa((int(curSeq))) + ", 요청 Seq : " + strconv.Itoa((int(tx.Seq))))
		}
	}
	return nil
}

// account 존재 여부
func isAccountExist(account string) bool {
	curAccount := (accounts[account])
	return curAccount != nil
}

// account 형식이 올바른지 확인
/*
1. 총 길이가 13자리인가
2. "acc" 로 시작하는가
3. 뒤 10자리가 hex 인가
*/
func isValidAddressForm(account string) bool {
	if len(account) == 13 {
		if account[:3] == "acc" {
			_, err := hex.DecodeString(account[3:])
			if err != nil {
				log.Println(err)
				log.Println("유효한 hex 값이 아닙니다")
				return false
			} else {
				log.Println("[확인] address 유효성 검증 : 올바른 형식의 address 입니다.")
				return true
			}
		} else {
			log.Println("[에러] address 유효성 검증 : address 가 acc 로 시작하지 않읍니다")
			return false
		}
	} else {
		log.Println("[에러] address 유효성 검증 : address 길이가 13자리가 아닙니다")
		return false
	}
}

/* 트랜잭션 involved 계정 정보를 조회 */
func getAccount(tx *d.Transaction) map[string]*d.ResAccount {
	var result = map[string]*d.ResAccount{}
	var sender = accounts[tx.Sender]
	var senderAddr = sender.GetAddress()
	var senderData = d.ResAccount{
		Id    : sender.GetId(),
		Address: senderAddr,
		Seq    : sender.GetSeq(),
		Balance: sender.GetBalance(),
	}
	result["sender"] = &senderData

	if tx.Type == d.TYPE_SEND {
		var target = accounts[tx.Target]
		var targetAddr = target.GetAddress()
		var targetData = d.ResAccount{
			Id    : target.GetId(),
			Address: targetAddr,
			Seq    : target.GetSeq(),
			Balance: target.GetBalance(),
		}
		result["target"] = &targetData
	}

	return result
}

func startTranProcessing() {
	// 트랜잭션의 processing time 이 2초라고 가정
	time.Sleep(2 * time.Second)
}
