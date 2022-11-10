package data

import "math/big"

/* 트랜잭션 요청 변수 */
type ResStat struct {
	Add  uint64 `json:"add"`
	Sub  uint64 `json:"sub"`
	Send uint64 `json:"send"`
}

/* POST 결과 데이터 */
type ResAccount struct {
	Id       uint		`json:"id"`
	Address  string     `json:"address"`
	Seq      uint       `json:"seq"`
	Balance  *big.Float `json:"balance"`
	//IsLocked bool
}

/* [GET] balance 결과 데이터 */
type ResBalance struct {
	Balance *big.Float `json:"balance"`
}

/* [GET] seq 결과 데이터 */
type ResSeq struct {
	Seq uint `json:"seq"`
}

/* http response */
type ResHttp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

const (
	CODE_OK         = 200
	CODE_ERROR_400  = 400
	CODE_ERROR_412  = 412
	CODE_ERROR_TRAN = 1000
	MSG_OK          = "ok"
	MSG_ERROR       = "error"
	ERR_MSG_REQ_DECODE = "POST 요청 데이터 형식이 맞지 않습니다."
	ERR_MSG_SENDER_NOT_EXIST = "[에러] sender 계정이 존재하지 않으므로 revoke 되었습니다"
	ERR_MSG_TARGET_NOT_EXIST = "[에러] target 계정이 존재하지 않으므로 revoke 되었습니다"
	ERR_MSG_REQ_WRONG_ID = "유효한 형식의 id 가 아닙니다"
	ERR_MSG_ID_NOT_EXIST = "id 가 존재하지 않습니다"
	ERR_MSG_ACCOUNT_NOTHING = "아직 생성된 계정이 하나도 없습니다"
	ERR_MSG_SEND_WRONG_AMNT = "[에러] sender 밸런스보다 amount 가 더 크므로 sub 할 수 없습니다"
)




