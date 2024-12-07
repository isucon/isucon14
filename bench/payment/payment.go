package payment

import (
	"fmt"
)

type Status struct {
	Type StatusType
	Err  error
}
type StatusType int

const (
	StatusInitial StatusType = iota
	StatusSuccess
	StatusInvalidAmount
	StatusInvalidToken
)

func (s StatusType) String() string {
	switch s {
	case StatusInitial:
		return "決済処理中"
	case StatusSuccess:
		return "成功"
	case StatusInvalidAmount:
		return "決済額が不正"
	case StatusInvalidToken:
		return "決済トークンが無効"
	default:
		panic(fmt.Sprintf("unknown payment status: %d", s))
	}
}

type Payment struct {
	IdempotencyKey string
	Token          string
	Amount         int
	Status         Status
}

func NewPayment(idk string) *Payment {
	p := &Payment{
		IdempotencyKey: idk,
		Status:         Status{Type: StatusInitial, Err: nil},
	}
	return p
}
