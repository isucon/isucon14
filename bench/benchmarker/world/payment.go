package world

import (
	"fmt"

	"github.com/isucon/isucon14/bench/internal/concurrent"
	"github.com/isucon/isucon14/bench/payment"
	"github.com/samber/lo"
)

type PaymentDB struct {
	PaymentTokens      *concurrent.SimpleMap[string, *User]
	CommittedPayments  *concurrent.SimpleSlice[*payment.Payment]
	alreadyPaidRequest *concurrent.SimpleMap[*Request, struct{}]
}

func NewPaymentDB() *PaymentDB {
	return &PaymentDB{
		PaymentTokens:      concurrent.NewSimpleMap[string, *User](),
		CommittedPayments:  concurrent.NewSimpleSlice[*payment.Payment](),
		alreadyPaidRequest: concurrent.NewSimpleMap[*Request, struct{}](),
	}
}

func (db *PaymentDB) Verify(p *payment.Payment) payment.Status {
	user, ok := db.PaymentTokens.Get(p.Token)
	if !ok {
		return payment.Status{Type: payment.StatusInvalidToken, Err: nil}
	}
	if p.Amount <= 0 && p.Amount > 1_000_000 {
		return payment.Status{Type: payment.StatusInvalidAmount, Err: nil}
	}

	// 支払いがリクエストに対して valid かどうかを確認するが、
	// Payment Server 自体はリクエストに対して valid かどうかに関わらず決済を行う
	// ただこの status に入れた err はベンチマーカーで critical error として扱われて即負荷走行が終了する
	req := user.Request
	status := payment.Status{Type: payment.StatusSuccess, Err: nil}
	if req == nil {
		status.Err = fmt.Errorf("進行中のリクエストがありません。token: %s, amount: %v", p.Token, p.Amount)
	}
	if _, alreadyPaid := db.alreadyPaidRequest.Get(req); alreadyPaid {
		status.Err = fmt.Errorf("既に支払い済みです。token: %s, amount: %v, request id: %s", p.Token, p.Amount, req.ServerID)
	}
	if p.Amount != req.Fare() {
		status.Err = fmt.Errorf("支払い額が不正です。token: %s, expected amount: %v, actual amount: %v, request id: %s", p.Token, req.Fare(), p.Amount, req.ServerID)
	}

	db.CommittedPayments.Append(p)
	return status
}

func (db *PaymentDB) TotalPayment() int64 {
	return lo.SumBy(db.CommittedPayments.ToSlice(), func(p *payment.Payment) int64 { return int64(p.Amount) })
}
