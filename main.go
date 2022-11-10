package main

import (
	"bharvest/jihon/data"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// 데이터 초기화
var accounts = map[string]*data.Account{} 	// 계정 리스트
var waitings = map[string]*data.Queue{}		// 대기큐 리스트
var stats data.Stat							// 통계 객체

func main() {
	log.Println("Web server starting..")

	stats.Init()
	
	txHandler := ConvertToJson(http.HandlerFunc(funcTxHandler))
	balanceHandler := ConvertToJson(http.HandlerFunc(funcBalanceHandler))
	seqHandler := ConvertToJson(http.HandlerFunc(funcSeqHandler))
	statTotalHandler := ConvertToJson(http.HandlerFunc(funcStatTotalHandler))
	statTxHandler := ConvertToJson(http.HandlerFunc(funcStatTxHandler))

	//mux := http.NewServeMux()
	r := mux.NewRouter()
	r.Handle("/tx", txHandler)
	r.Handle("/{id}/balance", balanceHandler)
	r.Handle("/{id}/seq", seqHandler)
	r.Handle("/stat/total", statTotalHandler)
	r.Handle("/stat/tx_stat", statTxHandler)
	http.Handle("/", r)

	//log.Fatal(http.ListenAndServe(":8080", mux))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
