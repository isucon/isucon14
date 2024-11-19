package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var erroredUpstream = errors.New("errored upstream")

type paymentGatewayPostPaymentRequest struct {
	Amount int `json:"amount"`
}

type paymentGatewayGetPaymentsResponseOne struct {
	Amount int    `json:"amount"`
	Status string `json:"status"`
}

func requestPaymentGatewayPostPayment(paymentGatewayURL string, token string, param *paymentGatewayPostPaymentRequest, retrieveRidesOrderByCreatedAtAsc func() ([]Ride, error)) error {
	b, err := json.Marshal(param)
	if err != nil {
		return err
	}

	// 失敗したらとりあえずリトライ
	// FIXME: 社内決済マイクロサービスのインフラに異常が発生していて、同時にたくさんリクエストすると変なことになる可能性あり
	retry := 0
	for {
		err := func() error {
			req, err := http.NewRequest(http.MethodPost, paymentGatewayURL+"/payments", bytes.NewBuffer(b))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusNoContent {
				// エラーが返ってきても成功している場合があるので、社内決済マイクロサービスに問い合わせ
				getReq, err := http.NewRequest(http.MethodGet, paymentGatewayURL+"/payments", bytes.NewBuffer([]byte{}))
				if err != nil {
					return err
				}
				getReq.Header.Set("Authorization", "Bearer "+token)

				getRes, err := http.DefaultClient.Do(getReq)
				if err != nil {
					return err
				}
				defer res.Body.Close()

				// GET /payments は障害と関係なく200が返るので、200以外は回復不能なエラーとする
				if getRes.StatusCode != http.StatusOK {
					return fmt.Errorf("[GET /payments] unexpected status code (%d)", getRes.StatusCode)
				}
				var payments []paymentGatewayGetPaymentsResponseOne
				if err := json.NewDecoder(getRes.Body).Decode(&payments); err != nil {
					return err
				}

				rides, err := retrieveRidesOrderByCreatedAtAsc()
				if err != nil {
					return err
				}

				if len(rides) != len(payments) {
					return fmt.Errorf("unexpected number of payments: %d != %d. %w", len(rides), len(payments), erroredUpstream)
				}

				return nil
			}
			return nil
		}()
		if err != nil {
			if retry < 5 {
				retry++
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				return err
			}
		}
		break
	}

	return nil
}
