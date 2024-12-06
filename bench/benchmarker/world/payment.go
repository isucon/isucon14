package world

import (
	"log/slog"

	"github.com/isucon/isucon14/bench/internal/concurrent"
	"github.com/isucon/isucon14/bench/payment"
	"github.com/samber/lo"
)

type invalidPaymentReason string

const (
	// 既に支払い済み
	invalidPaymentReasonAlreadyPaid invalidPaymentReason = "already_paid"
	// Amount が異なる
	invalidPaymentReasonInvalidAmount invalidPaymentReason = "invalid_amount"
	// 進行中のリクエストがない
	invalidPaymentReasonNoRequest invalidPaymentReason = "no_request"
)

type invalidPayment struct {
	Payment *payment.Payment
	Request *Ride
	Reason  invalidPaymentReason
}

type PaymentDB struct {
	PaymentTokens     *concurrent.SimpleMap[string, *User]
	CommittedPayments *concurrent.SimpleSlice[*payment.Payment]
	// TODO: このフィールドに入れるんじゃなくてdeductionとかにまとめる
	invalidPayments    *concurrent.SimpleSlice[*invalidPayment]
	alreadyPaidRequest *concurrent.SimpleMap[*Ride, struct{}]
}

func NewPaymentDB() *PaymentDB {
	return &PaymentDB{
		PaymentTokens:      concurrent.NewSimpleMap[string, *User](),
		CommittedPayments:  concurrent.NewSimpleSlice[*payment.Payment](),
		invalidPayments:    concurrent.NewSimpleSlice[*invalidPayment](),
		alreadyPaidRequest: concurrent.NewSimpleMap[*Ride, struct{}](),
	}
}

func (db *PaymentDB) Verify(p *payment.Payment) payment.Status {
	user, ok := db.PaymentTokens.Get(p.Token)
	if !ok {
		return payment.StatusInvalidToken
	}
	if p.Amount <= 0 && p.Amount > 1_000_000 {
		return payment.StatusInvalidAmount
	}

	// 支払いがリクエストに対して valid かどうかを確認するが、
	// Payment Server 自体はリクエストに対して valid かどうかに関わらず決済を行う
	// このため、invalid だった場合も status には反映しないが invalidPayments に記録する
	req := user.Ride
	if req == nil {
		db.invalidPayments.Append(&invalidPayment{Payment: p, Reason: invalidPaymentReasonNoRequest})
		// TODO: ロギング
		slog.Debug("invalid payment", "payment", p, "reason", invalidPaymentReasonNoRequest)
	}
	if _, alreadyPaid := db.alreadyPaidRequest.Get(req); alreadyPaid {
		db.invalidPayments.Append(&invalidPayment{Payment: p, Request: req, Reason: invalidPaymentReasonAlreadyPaid})
		// TODO: ロギング
		slog.Debug("invalid payment", "payment", p, "request", req, "reason", invalidPaymentReasonAlreadyPaid)
	}
	if p.Amount != req.Fare() {
		db.invalidPayments.Append(&invalidPayment{Payment: p, Request: req, Reason: invalidPaymentReasonInvalidAmount})
		// TODO: ロギング
		slog.Debug("invalid payment", "payment", p, "request", req, "reason", invalidPaymentReasonInvalidAmount)
	}

	db.CommittedPayments.Append(p)
	return payment.StatusSuccess
}

func (db *PaymentDB) TotalPayment() int64 {
	return lo.SumBy(db.CommittedPayments.ToSlice(), func(p *payment.Payment) int64 { return int64(p.Amount) })
}
