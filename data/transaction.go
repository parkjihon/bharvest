package data

import (
	"encoding/hex"
	"errors"
	"log"
	"math/big"
	"reflect"
)

/* 트랜잭션 요청 변수 */
type Transaction struct {
	Sender string `json:"sender"`
	Seq    uint   `json:"seq"`
	Type   string `json:"type"`
	Amount string `json:"amount"`
	Target string `json:"target,omitempty"`
}

const (
	TYPE_ADD  = "add"
	TYPE_SUB  = "sub"
	TYPE_SEND = "send"
)

/*
A. 트랜잭션 입력값 유효성 검증
  - sender & target addr 형식이 유효한지
  - seq 가 uint 값인지 확인
  - type 이 "add", "sub", "send" 중 하나인지 확인
  - amount 가 양의 실수인지
*/
func (tx *Transaction) TranInputValidator() error {
	if r := checkAddress(tx); !r {
		return errors.New("[입력값 유효성 검증] 올바른 Account 주소 형식이 아닙니다")
	}

	if r := checkSeq(tx); !r {					// 불필요
		return errors.New("[입력값 유효성 검증] 올바른 Seq 형식(uint) 이 아닙니다")
	}

	if r := checkType(tx); !r {
		return errors.New("[입력값 유효성 검증] 올바른 Type 형식이 아닙니다")
	}

	if r := checkAmount(tx) ; !r {
		return errors.New("[입력값 유효성 검증] 올바른 Amount 형식(양의 실수) 이 아닙니다")
	}

	return nil
}

// sender 및 target addr 가 올바른 형식인지 확인
func checkAddress(tx *Transaction) bool {
	result := isValidAddressForm(tx.Sender)
	if !result {
		return false
	}
	if tx.Type == TYPE_SEND {
		if result = isValidAddressForm(tx.Target); !result {
			return false
		}
	}
	return true
}

// seq 가 uint 값인지 확인
func checkSeq(tx *Transaction) bool {
	if reflect.TypeOf(tx.Seq).String() == "uint" {
		return true
	} else {
		return false
	}
}

// type 이 올바른지 확인
func checkType(tx *Transaction) bool {
	if tx.Type == TYPE_ADD || tx.Type == TYPE_SUB || tx.Type == TYPE_SEND {
		return true
	} else {
		return false
	}
}

// decimal은 정수 또는 실수를 string타입으로 받으며, 정수는 20자리, 실수는 총 32자리까지 가능
func checkAmount(tx *Transaction) bool {
	fAmount := new(big.Float).SetPrec(128)					// amount 가 정수 또는 실수인지 확인
	_, _, err := fAmount.Parse(tx.Amount, 0)
	//str := fAmount.Text('f', 32)
	//log.Println("실수 TEXT", str)
	if err != nil {
		log.Println("[에러] amount 값이 정수 또는 실수가 아닙니다")
		return false
	}

	if fAmount.Cmp(new(big.Float)) <= 0 {
		log.Println("[에러] 입력 amount 는 0보다 커야 합니다")
		return false
	}

	if tx.Type == TYPE_ADD || tx.Type == TYPE_SUB || tx.Type == TYPE_SEND {
		return true
	} else {
		return false
	}
}

/*
1. 총 길이가 13자리인지 확인
2. "acc" 로 시작하는지 확인
3. 뒤 10자리가 헥사값을 가지는지 확인
*/
func isValidAddressForm(account string) bool {
	if len(account) == 13 {
		if account[:3] == "acc" {
			_, err := hex.DecodeString(account[3:])
			if err != nil {
				log.Println(err)
				return false
			} else {
				return true
			}
		} else {
			log.Println("[에러] address 유효성 검증 : address 가 acc 로 시작하지 않습니다")
			return false
		}
	} else {
		log.Println("[에러] address 유효성 검증 : address 길이가 13자리가 아닙니다")
		return false
	}
}


