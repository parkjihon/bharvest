package data

import (
	"sync"
	"sync/atomic"
)

/*
통계(stat) 자료구조
*/
type Stat struct {
	add  uint64
	sub  uint64
	send uint64
}

func (s *Stat) Init() {
	s.add  = 0
	s.sub  = 0
	s.send = 0
}

var mutex = &sync.Mutex{}

// func (s *Stat) CountAdd() {
// 	wg := new(sync.WaitGroup)
// 	wg.Add(1)
// 	go func() {
// 		atomic.AddUint64(&s.add, 1)
// 		wg.Done()
// 	}()
// 	wg.Wait()
// }
func (s *Stat) CountAdd() {
	mutex.Lock()
	atomic.AddUint64(&s.add, 1)
	mutex.Unlock()
}

func (s *Stat) CountSub() {
	mutex.Lock()
	atomic.AddUint64(&s.sub, 1)
	mutex.Unlock()
}

func (s *Stat) CountSend() {
	mutex.Lock()
	atomic.AddUint64(&s.send, 1)
	mutex.Unlock()
}

// 통계 조회
func (s *Stat) QueryStat() ResStat {
	var statOutput ResStat
	// wg := new(sync.WaitGroup)
	// wg.Add(1)
	// go func() {
	// 	statOutput.Add  = atomic.LoadUint64(&s.add)
	// 	statOutput.Sub  = atomic.LoadUint64(&s.sub)
	// 	statOutput.Send = atomic.LoadUint64(&s.send)
	// 	wg.Done()
	// }()
	// wg.Wait()
	mutex.Lock()
	statOutput.Add  = atomic.LoadUint64(&s.add)
	statOutput.Sub  = atomic.LoadUint64(&s.sub)
	statOutput.Send = atomic.LoadUint64(&s.send)
	mutex.Unlock()
	return statOutput
}



