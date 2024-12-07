package payment

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"
)

type PostPaymentRequest struct {
	Amount int `json:"amount"`
}

func (r *PostPaymentRequest) IsSamePayload(token string, p *Payment) bool {
	return token == p.Token && r.Amount == p.Amount
}

func getTokenFromAuthorizationHeader(r *http.Request) (string, error) {
	auth := r.Header.Get(AuthorizationHeader)
	prefix := AuthorizationHeaderPrefix
	if !strings.HasPrefix(auth, prefix) {
		return "", fmt.Errorf("不正な値がAuthorization headerにセットされています。expected: Bearer ${token}. got: %s", auth)
	}
	return auth[len(prefix):], nil
}

func (s *Server) PostPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		p          *Payment
		newPayment bool
	)

	idk := r.Header.Get(IdempotencyKeyHeader)
	if len(idk) > 0 {
		p, newPayment = s.knownKeys.GetOrSetDefault(idk, func() *Payment { return NewPayment(idk) })
		if !newPayment && p.locked.Load() {
			writeJSON(w, http.StatusConflict, map[string]string{"message": "既に処理中です"})
			return
		}
	} else {
		p = NewPayment("")
		newPayment = true
	}

	token, err := getTokenFromAuthorizationHeader(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	var req PostPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "不正なリクエスト形式です"})
		return
	}
	if !newPayment {
		if !req.IsSamePayload(token, p) {
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"message": "リクエストペイロードがサーバーに記録されているものと異なります"})
			return
		}
		writeResponse(w, p.Status)
		return
	} else {
		p.Token = token
		p.Amount = req.Amount
	}

	failureCount, _ := s.failureCounts.GetOrSetDefault(token, func() int { return 0 })
	if p.locked.CompareAndSwap(false, true) {
		if failureCount >= 1 {
			alreadyProcessed := false
			if !newPayment {
				alreadyProcessedIdks := map[string]struct{}{}
				for _, processed := range s.processedPayments.ToSlice() {
					_, exist := alreadyProcessedIdks[processed.payment.IdempotencyKey]
					if exist {
						slog.Error("idempotency key が重複しています", "idk", processed.payment.IdempotencyKey)
					}
					alreadyProcessedIdks[processed.payment.IdempotencyKey] = struct{}{}
					if processed.payment.IdempotencyKey == p.IdempotencyKey {
						alreadyProcessed = true
						break
					}
				}
			}
			if !alreadyProcessed {
				p.Status = s.verifier.Verify(p)
				if p.Status.Type == StatusSuccess {
					s.processedPayments.Append(&processedPayment{payment: p, processedAt: time.Now()})
				}
				if p.Status.Err != nil {
					s.errChan <- p.Status.Err
				}
				s.failureCounts.Delete(token)
				p.locked.Store(false) // idenpotency key が同じリクエストが来たときにエラーを返さないように
			}
			writeResponse(w, p.Status)
			return
		} else {
			s.failureCounts.Set(token, failureCount+1)
		}
	}
	writeRandomError(w)
	return

	time.Sleep(s.processTime)
	// (直近3秒で処理された payment の数) / 100 の確率で処理を失敗させる(最大50%)
	var recentProcessedCount int
	for _, processed := range s.processedPayments.BackwardIter() {
		if time.Since(processed.processedAt) > 3*time.Second {
			break
		}
		recentProcessedCount++
	}
	failurePercentage := recentProcessedCount
	if failurePercentage > 50 {
		failurePercentage = 50
	}
	// failureCount, _ := s.failureCounts.GetOrSetDefault(token, func() int { return 0 })
	if rand.IntN(100) > failurePercentage || failureCount >= 4 {
		// lock はここでしか触らない。lock が true の場合は idempotency key が同じリクエストが処理中の場合のみ
		if p.locked.CompareAndSwap(false, true) {
			alreadyProcessed := false
			if !newPayment {
				alreadyProcessedIdks := map[string]struct{}{}
				for _, processed := range s.processedPayments.ToSlice() {
					_, exist := alreadyProcessedIdks[processed.payment.IdempotencyKey]
					if exist {
						slog.Error("idempotency key が重複しています", "idk", processed.payment.IdempotencyKey)
					}
					alreadyProcessedIdks[processed.payment.IdempotencyKey] = struct{}{}
					if processed.payment.IdempotencyKey == p.IdempotencyKey {
						alreadyProcessed = true
						break
					}
				}
			}
			if !alreadyProcessed {
				p.Status = s.verifier.Verify(p)
				if p.Status.Type == StatusSuccess {
					s.processedPayments.Append(&processedPayment{payment: p, processedAt: time.Now()})
				}
				if p.Status.Err != nil {
					s.errChan <- p.Status.Err
				}
				s.failureCounts.Delete(token)
				p.locked.Store(false) // idenpotency key が同じリクエストが来たときにエラーを返さないように
			}
			if rand.IntN(100) > failurePercentage || failureCount >= 4 {
				writeResponse(w, p.Status)
			} else {
				writeRandomError(w)
			}
			return
		}
	} else {
		s.failureCounts.Set(token, failureCount+1)
	}

	// 不安定なエラーを再現
	writeRandomError(w)
}

func writeRandomError(w http.ResponseWriter) {
	switch rand.IntN(3) {
	case 0:
		w.WriteHeader(http.StatusInternalServerError)
	case 1:
		w.WriteHeader(http.StatusBadGateway)
	case 2:
		w.WriteHeader(http.StatusGatewayTimeout)
	}
}

type ResponsePayment struct {
	Amount int    `json:"amount"`
	Status string `json:"status"`
}

func NewResponsePayment(p *Payment) ResponsePayment {
	return ResponsePayment{
		Amount: p.Amount,
		Status: p.Status.Type.String(),
	}
}

func (s *Server) GetPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(300 * time.Millisecond)
	token, err := getTokenFromAuthorizationHeader(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	payments := s.processedPayments.ToSlice()

	res := []ResponsePayment{}
	for _, p := range payments {
		if p.payment.Token != token {
			continue
		}
		res = append(res, NewResponsePayment(p.payment))
	}
	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error(err.Error())
	}
}

func writeResponse(w http.ResponseWriter, paymentStatus Status) {
	switch paymentStatus.Type {
	case StatusInitial:
		w.WriteHeader(http.StatusNoContent)
	case StatusSuccess:
		w.WriteHeader(http.StatusNoContent)
	case StatusInvalidAmount:
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "決済額が不正です"})
	case StatusInvalidToken:
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "決済トークンが無効です"})
	}
}
