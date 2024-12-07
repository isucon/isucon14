package scenario

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/isucon/isucandar"
)

// Validation はシナリオの結果検証処理を行う
func (s *Scenario) Validation(ctx context.Context, step *isucandar.BenchmarkStep) error {
	// 負荷走行終了後、payment server へのリクエストが届くかもしれないので5秒だけ待つ
	time.Sleep(5 * time.Second)
	s.paymentServer.Close()
	s.sendResultWait.Wait()

	actual := s.world.PaymentDB.TotalPayment() + s.TotalDiscount()
	expected := s.TotalSales()
	if actual != expected {
		s.contestantLogger.Error("決済サーバーで決済された額とユーザーが支払うべき額が一致していません", slog.Int64("diff(actual-expected)", actual-expected), slog.Int64("total_discount", s.TotalDiscount()), slog.Int64("total_payment", s.world.PaymentDB.TotalPayment()), slog.Int64("total_sales", s.TotalSales()))
	}

	for _, region := range s.world.Regions {
		s.contestantLogger.Info("最終地域情報",
			slog.String("名前", region.Name),
			slog.Int("ユーザー登録数", region.UsersDB.Len()),
			slog.Int("アクティブユーザー数", region.ActiveUserNum()),
		)
	}
	for _, owner := range s.world.OwnerDB.Iter() {
		s.contestantLogger.Info("最終オーナー情報",
			slog.String("名前", owner.RegisteredData.Name),
			slog.Int64("売上", owner.TotalSales.Load()),
			slog.Int("椅子数", owner.ChairDB.Len()),
		)
	}
	if s.completedRequests > 0 {
		s.contestantLogger.Info(fmt.Sprintf("%.1f%%のライドは椅子がマッチされるまでの時間、%.1f%%のライドはマッチされた椅子が乗車地点までに掛かる時間、%.1f%%のライドは椅子の実移動時間に不満がありました",
			(1-float64(s.evaluationMap[0])/float64(s.completedRequests))*100,
			(1-float64(s.evaluationMap[1])/float64(s.completedRequests))*100,
			(1-float64(s.evaluationMap[2]+s.evaluationMap[3])/float64(s.completedRequests*2))*100,
		))
	}
	s.contestantLogger.Info("結果", slog.Bool("pass", !s.failed), slog.Int64("スコア", s.Score(true)), slog.Any("種別エラー数", s.world.ErrorCounter.Count()))
	return sendResult(s, true, !s.failed)
}
