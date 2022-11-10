package main

import (
	d "bharvest/jihon/data"
	"log"
	"strings"
	"sync"
)

func checkLockAndWait(account string) bool {			// 재귀 함수
	curAccount := accounts[account]
	if curAccount != nil && curAccount.IsLock() {
		//log.Println("처리 중인 트랜잭션이 존재하여 해당 요청은 pending 상태로 처리됩니다")
		var wg sync.WaitGroup
		wg.Add(1)
		var queueData d.QueueData
		queueData.Wg = &wg

		if waitings[account] != nil {					//대기 큐에 등록
			waitings[account].Push(&queueData)
		} else {
			queue := d.NewQueue()						//새로운 queue 생성
			queue.Push(&queueData)
			waitings[account] = queue
		}
		log.Println("대기 queue 에 추가, 대기큐 갯수: ", waitings[account].Len())

		wg.Wait()										// 내 순서가 올 때까지 대기
		return checkLockAndWait(account)				// 그럴 일은 없겠지만 또 pending 상태일 수 있으므로 재확인
	} else {
		return true
	}
}

func unLockAccount(tx *d.Transaction) {
	accounts[tx.Sender].AddSeq()
	accounts[tx.Sender].UnLock()
	if strings.ToLower(tx.Type) == d.TYPE_SEND {
		accounts[tx.Target].UnLock()
	}
}

// 대기큐 체크 -> 존재하는 경우 Pop -> 재실행 정보 전달
func checkWaitings(tx *d.Transaction) {
	if waitings[tx.Sender] != nil {
		if waitings[tx.Sender].Len() > 0 {
			queueData := waitings[tx.Sender].Pop() // 큐에서 첫번째 요청을 꺼낸다
			invokePendingTran(queueData.(*d.QueueData))
		}
	}
	if waitings[tx.Target] != nil && strings.ToLower(tx.Type) == d.TYPE_SEND {
		if waitings[tx.Target].Len() > 0 {
			queueData := waitings[tx.Target].Pop()
			invokePendingTran(queueData.(*d.QueueData))
		}
	}
}

func invokePendingTran(q *d.QueueData) {
	//log.Println("[정상] 최상위 대기큐에게 실행할 준비가 되었다고 전달")
	q.Wg.Done()
}


