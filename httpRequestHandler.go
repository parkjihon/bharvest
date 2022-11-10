package main

import (
	d "bharvest/jihon/data"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// http 응답을 json 형태로 변환
func ConvertToJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func ResponseOk(rw http.ResponseWriter, msg interface{}) {				// 성공 응답
	json.NewEncoder(rw).Encode(d.ResHttp{
		Code: d.CODE_OK,
		Message: d.MSG_OK,
		Result: msg,
	})
}

func ResponseError(rw http.ResponseWriter, code int, msg string) {		// 실패 응답
	log.Println(msg)
	json.NewEncoder(rw).Encode(d.ResHttp{
		Code: code,
		Message: d.MSG_ERROR,
		Error: msg,
	})
}

// [POST] /tx
func funcTxHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		log.Println("[POST] /tx requested...")								// POST
		var tx d.Transaction

		if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {			// 디코드
			ResponseError(rw, d.CODE_ERROR_400, d.ERR_MSG_REQ_DECODE)
			return
		}

		if err := tx.TranInputValidator(); err != nil {						// 입력값 유효성 검증
			ResponseError(rw, d.CODE_ERROR_400, err.Error())
			return
		}

		if r := isAccountExist(tx.Sender); !r {								// Sender 존재 검증
			if strings.ToLower(tx.Type) != d.TYPE_ADD {
				ResponseError(rw, d.CODE_ERROR_412, d.ERR_MSG_SENDER_NOT_EXIST)
				return
			}
		}
	
		if err := isValidSeq(&tx, strings.ToLower(tx.Type)); err != nil {	// seq 유효성 검증
			ResponseError(rw, d.CODE_ERROR_400, err.Error())
			return
		}

		// 계정이 Lock 상태인지 확인 -> 대기열에 저장 및 실행제어 -> 순서가 오면 재실행
		checkLockAndWait(tx.Sender)

		/* 트랜잭션 수행 */
		var tranError error
		switch strings.ToLower(tx.Type) {
		case d.TYPE_ADD:
			tranError = tranAdd(&tx)
			stats.CountAdd()
		case d.TYPE_SUB:
			tranError = tranSub(&tx)
			stats.CountSub()
		case d.TYPE_SEND:
			if r := isAccountExist(tx.Target); !r {							// Target 존재 검증
				ResponseError(rw, d.CODE_ERROR_412, d.ERR_MSG_TARGET_NOT_EXIST)
				return
			}
			// 받는 계정(target) 이 Lock 상태인지 확인
			checkLockAndWait(tx.Target)
			tranError = tranSend(&tx)
			stats.CountSend()
		}

		unLockAccount(&tx)
		checkWaitings(&tx)		// 대기큐 체크 -> 존재하는 경우 Pop -> 재실행 정보 전달 									

		if tranError == nil {
			ResponseOk(rw, getAccount(&tx))
			log.Println("transaction committed")
		} else {
			ResponseError(rw, d.CODE_ERROR_TRAN, tranError.Error())
		}
	}
}

// [GET] /{id}/balance
func funcBalanceHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("[GET] /{id}/balance requested...")
		vars := mux.Vars(r)
		account := vars["id"]
		if !isValidAddressForm(account) {
			ResponseError(rw, d.CODE_ERROR_400, d.ERR_MSG_REQ_WRONG_ID)
			return
		}

		if accounts[account] != nil {
			ResponseOk(rw, d.ResBalance{Balance: accounts[account].GetBalance()})	
		} else {
			ResponseError(rw, d.CODE_ERROR_412, d.ERR_MSG_ID_NOT_EXIST)
		}
	}
}

// [GET] /{id}/seq
func funcSeqHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("[GET] /{id}/seq requested...")
		vars := mux.Vars(r)
		account := vars["id"]
		if !isValidAddressForm(account) {
			ResponseError(rw, d.CODE_ERROR_400, d.ERR_MSG_REQ_WRONG_ID)
			return
		}

		if accounts[account] != nil {
			ResponseOk(rw, d.ResSeq{Seq: accounts[account].GetSeq() + 1})	
		} else {
			ResponseError(rw, d.CODE_ERROR_412, d.ERR_MSG_ID_NOT_EXIST)
		}
	}
}

// [GET] /stat/total
func funcStatTotalHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("[GET] /stat/total requested...")
		if len(accounts) != 0 {
			var balances = map[string]*d.ResBalance{}
			for _, account := range accounts {
				var bal d.ResBalance
				bal.Balance = account.GetBalance()
				balances[account.GetAddress()] = &bal
			}
			ResponseOk(rw, balances)	
		} else {
			ResponseError(rw, d.CODE_ERROR_412, d.ERR_MSG_ACCOUNT_NOTHING)
		}
	}
}

// [GET] /stat/tx_stat
func funcStatTxHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("/stat/tx_stat requested...")
		txStat := stats.QueryStat()
		ResponseOk(rw, txStat)
	}
}
