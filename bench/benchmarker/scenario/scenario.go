package scenario

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/guregu/null/v5"
	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/agent"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/isucon/isucon14/bench/benchmarker/scenario/worldclient"
	"github.com/isucon/isucon14/bench/benchmarker/webapp"
	"github.com/isucon/isucon14/bench/benchmarker/webapp/api"
	"github.com/isucon/isucon14/bench/benchmarker/world"
	"github.com/isucon/isucon14/bench/benchrun"
	"github.com/isucon/isucon14/bench/benchrun/gen/isuxportal/resources"
	"github.com/isucon/isucon14/bench/payment"
)

// Scenario はシナリオを表す構造体
// 以下の関数を実装して b.AddScenario() で追加したあと b.Run() で実行される
// - Prepare(context.Context, *BenchmarkStep) error
//   - シナリオの初期化処理を行う
//   - Initialize した後に Validation を呼ぶことが多いっぽい
//   - Initialize -> Validation -> Initialize してもいいかも？（13とかはそうしてそう）
//
// - Load(context.Context, *BenchmarkStep) error
//   - シナリオのメイン処理を行う
//
// - Validation(context.Context, *BenchmarkStep) error
//   - シナリオの結果検証処理を行う
//   - 料金の整合性をみたいかも
type Scenario struct {
	language                  string
	target                    string
	addr                      string
	paymentURL                string
	contestantLogger          *slog.Logger
	world                     *world.World
	worldCtx                  *world.Context
	paymentServer             *payment.Server
	paymentErrChan            chan error
	step                      *isucandar.BenchmarkStep
	reporter                  benchrun.Reporter
	meter                     metric.Meter
	prepareOnly               bool
	skipStaticFileSanityCheck bool
	finalScore                null.Int64
	evaluationMap             [4]int
	evalScoreMap              [5]int
	completedRequests         int
	evaluationMapLock         sync.RWMutex
	failed                    bool
	sendResultWait            sync.WaitGroup
}

func NewScenario(target, addr, paymentURL string, logger *slog.Logger, reporter benchrun.Reporter, meter metric.Meter, prepareOnly bool, skipStaticFileSanityCheck bool) *Scenario {
	completedRequestChan := make(chan *world.Request, 1000)
	worldClient := worldclient.NewWorldClient(context.Background(), webapp.ClientConfig{
		TargetBaseURL:         target,
		TargetAddr:            addr,
		ClientIdleConnTimeout: 10 * time.Second,
	}, skipStaticFileSanityCheck)
	w := world.NewWorld(30*time.Millisecond, completedRequestChan, worldClient, logger)

	worldCtx := world.NewContext(w)

	paymentErrChan := make(chan error, 1000)
	paymentServer := payment.NewServer(w.PaymentDB, 30*time.Millisecond, paymentErrChan)
	go func() {
		http.ListenAndServe(":12345", paymentServer)
	}()

	usersAttributeSets := map[world.UserState]attribute.Set{
		world.UserStateInactive:                  attribute.NewSet(attribute.Int("state", int(world.UserStateInactive))),
		world.UserStateActive:                    attribute.NewSet(attribute.Int("state", int(world.UserStateActive))),
		world.UserStatePaymentMethodsNotRegister: attribute.NewSet(attribute.Int("state", int(world.UserStatePaymentMethodsNotRegister))),
	}
	requestsAttributeSets := map[world.RequestStatus]attribute.Set{
		world.RequestStatusMatching:    attribute.NewSet(attribute.Int("status", int(world.RequestStatusMatching))),
		world.RequestStatusDispatching: attribute.NewSet(attribute.Int("status", int(world.RequestStatusDispatching))),
		world.RequestStatusDispatched:  attribute.NewSet(attribute.Int("status", int(world.RequestStatusDispatched))),
		world.RequestStatusCarrying:    attribute.NewSet(attribute.Int("status", int(world.RequestStatusCarrying))),
		world.RequestStatusArrived:     attribute.NewSet(attribute.Int("status", int(world.RequestStatusArrived))),
		world.RequestStatusCompleted:   attribute.NewSet(attribute.Int("status", int(world.RequestStatusCompleted))),
	}

	lo.Must1(meter.Int64ObservableCounter("world.time", metric.WithDescription("Time"), metric.WithUnit("1"), metric.WithInt64Callback(func(ctx context.Context, o metric.Int64Observer) error {
		o.Observe(w.Time)
		return nil
	})))
	lo.Must1(meter.Int64ObservableCounter("world.timeout", metric.WithDescription("Timeout"), metric.WithUnit("1"), metric.WithInt64Callback(func(ctx context.Context, o metric.Int64Observer) error {
		o.Observe(int64(w.TimeoutTickCount))
		return nil
	})))
	lo.Must1(meter.Int64ObservableGauge("world.users.num", metric.WithDescription("Number of users"), metric.WithUnit("1"), metric.WithInt64Callback(func(ctx context.Context, o metric.Int64Observer) error {
		for _, r := range w.Regions {
			counts := lo.CountValuesBy(r.UsersDB.ToSlice(), func(u *world.User) world.UserState { return u.State })
			for state, set := range usersAttributeSets {
				o.Observe(int64(counts[state]), metric.WithAttributeSet(set), metric.WithAttributes(attribute.String("region", r.Name)))
			}
		}
		return nil
	})))
	lo.Must1(meter.Int64ObservableGauge("world.chairs.num", metric.WithDescription("Number of chairs"), metric.WithUnit("1"), metric.WithInt64Callback(func(ctx context.Context, o metric.Int64Observer) error {
		for _, p := range w.OwnerDB.Iter() {
			chairs := p.ChairDB.ToSlice()
			insideRegion := lo.CountBy(chairs, func(c *world.Chair) bool { return c.Location.Current().Within(p.Region) })
			o.Observe(int64(insideRegion), metric.WithAttributes(attribute.Int("owner", int(p.ID)), attribute.String("region", p.Region.Name), attribute.Bool("inside_region", true)))
			o.Observe(int64(len(chairs)-insideRegion), metric.WithAttributes(attribute.Int("owner", int(p.ID)), attribute.String("region", p.Region.Name), attribute.Bool("inside_region", false)))
		}
		return nil
	})))
	lo.Must1(meter.Int64ObservableCounter("world.owners.num", metric.WithDescription("Number of owners"), metric.WithUnit("1"), metric.WithInt64Callback(func(ctx context.Context, o metric.Int64Observer) error {
		o.Observe(int64(w.OwnerDB.Size()))
		return nil
	})))
	lo.Must1(meter.Int64ObservableCounter("world.owners.sales", metric.WithDescription("Sales of owner"), metric.WithUnit("1"), metric.WithInt64Callback(func(ctx context.Context, o metric.Int64Observer) error {
		for _, p := range w.OwnerDB.Iter() {
			o.Observe(p.TotalSales.Load(), metric.WithAttributes(attribute.Int("owner", int(p.ID)), attribute.String("region", p.Region.Name)))
		}
		return nil
	})))
	lo.Must1(meter.Int64ObservableGauge("world.requests.num", metric.WithDescription("Number of requests"), metric.WithUnit("1"), metric.WithInt64Callback(func(ctx context.Context, o metric.Int64Observer) error {
		counts := lo.CountValuesBy(w.RequestDB.ToSlice(), func(r *world.Request) world.RequestStatus { return r.Statuses.Desired })
		for status, set := range requestsAttributeSets {
			o.Observe(int64(counts[status]), metric.WithAttributeSet(set))
		}
		return nil
	})))
	requestsRecorder := lo.Must1(meter.Int64Counter("world.requests.evaluations", metric.WithDescription("Counter of request's evaluations"), metric.WithUnit("1")))
	matchingLatency := lo.Must1(meter.Int64Histogram("world.request.matching_latency", metric.WithDescription("Histogram of matching latency"), metric.WithUnit("1")))
	dispatchingLatency := lo.Must1(meter.Int64Histogram("world.request.dispatching_latency", metric.WithDescription("Histogram of dispatching latency"), metric.WithUnit("1")))
	carryingLatency := lo.Must1(meter.Int64Histogram("world.request.carrying_latency", metric.WithDescription("Histogram of carrying latency"), metric.WithUnit("1")))

	s := &Scenario{
		target:                    target,
		addr:                      addr,
		paymentURL:                paymentURL,
		contestantLogger:          logger,
		world:                     w,
		worldCtx:                  worldCtx,
		paymentServer:             paymentServer,
		paymentErrChan:            paymentErrChan,
		reporter:                  reporter,
		meter:                     meter,
		prepareOnly:               prepareOnly,
		skipStaticFileSanityCheck: skipStaticFileSanityCheck,
	}
	go func() {
		for req := range completedRequestChan {
			eval := req.CalculateEvaluation()
			intervals := req.Intervals()
			requestsRecorder.Add(context.Background(), 1, metric.WithAttributes(attribute.Int("score", eval.Score()), attribute.Bool("matching", eval.Matching), attribute.Bool("dispatch", eval.Dispatch), attribute.Bool("pickup", eval.Pickup), attribute.Bool("drive", eval.Drive)))
			matchingLatency.Record(context.Background(), intervals[0])
			dispatchingLatency.Record(context.Background(), intervals[1])
			carryingLatency.Record(context.Background(), intervals[2])
			s.evaluationMapLock.Lock()
			s.completedRequests++
			s.evalScoreMap[eval.Score()-1]++
			for i, ok := range eval.Map() {
				if ok {
					s.evaluationMap[i]++
				}
			}
			s.evaluationMapLock.Unlock()
		}
	}()

	return s
}

// Prepare はシナリオの初期化処理を行う
func (s *Scenario) Prepare(ctx context.Context, step *isucandar.BenchmarkStep) error {
	if s.skipStaticFileSanityCheck {
		s.contestantLogger.Info("静的ファイルのチェックをスキップします")
	} else {
		if err := s.validateFrontendFiles(ctx); err != nil {
			s.contestantLogger.Error("静的ファイルのチェックに失敗しました", slog.String("error", err.Error()))
			return err
		}
	}

	client, err := webapp.NewClient(webapp.ClientConfig{
		TargetBaseURL:         s.target,
		TargetAddr:            s.addr,
		ClientIdleConnTimeout: 10 * time.Second,
	}, agent.WithTimeout(30*time.Second))
	if err != nil {
		return err
	}

	if err := s.initializeData(ctx, client); err != nil {
		s.contestantLogger.Error("initializeに失敗しました", slog.String("error", err.Error()))
		return err
	}

	client, err = webapp.NewClient(webapp.ClientConfig{
		TargetBaseURL:         s.target,
		TargetAddr:            s.addr,
		ClientIdleConnTimeout: 10 * time.Second,
	})
	if err != nil {
		return err
	}

	if err := s.prevalidation(ctx, client); err != nil {
		return err
	}

	return nil
}

func (s *Scenario) initializeData(ctx context.Context, client *webapp.Client) error {
	resp, err := client.PostInitialize(ctx, &api.PostInitializeReq{PaymentServer: s.paymentURL})
	if err != nil {
		return err
	}

	// 言語情報を追加
	s.language = resp.Language

	const (
		initialOwnersNum         = 5
		initialChairsNumPerOwner = 4
		initialUsersNum          = 10
	)

	for i := range initialOwnersNum {
		owner, err := s.world.CreateOwner(s.worldCtx, &world.CreateOwnerArgs{
			Region: s.world.Regions[i%len(s.world.Regions)],
		})
		if err != nil {
			return err
		}

		for range initialChairsNumPerOwner {
			_, err := s.world.CreateChair(s.worldCtx, &world.CreateChairArgs{
				Owner:             owner,
				InitialCoordinate: world.RandomCoordinateOnRegionWithRand(owner.Region, owner.Rand),
				Model:             owner.ChairModels[2].Random(),
			})
			if err != nil {
				return err
			}
		}
	}
	for i := range initialUsersNum {
		_, err := s.world.CreateUser(s.worldCtx, &world.CreateUserArgs{Region: s.world.Regions[i%len(s.world.Regions)]})
		if err != nil {
			return err
		}
	}

	return nil
}

// Load はシナリオのメイン処理を行う
func (s *Scenario) Load(ctx context.Context, step *isucandar.BenchmarkStep) error {
	if s.prepareOnly {
		return nil
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	s.sendResultWait.Add(1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.sendResultWait.Done()
				return
			case <-ticker.C:
				if err := sendResult(s, false, false); err != nil {
					slog.Error(err.Error())
				}
			}
		}
	}()

	s.world.RestTicker()
LOOP:
	for {
		select {
		case <-ctx.Done():
			s.world.Finish()
			break LOOP

		default:
			err := s.world.Tick(s.worldCtx)
			if err != nil {
				s.contestantLogger.Error("クリティカルエラーが発生しました", slog.String("error", err.Error()))
				s.failed = true
				return err
			}

			// paymentErrChan にエラーがあればエラーを返す
			select {
			case err := <-s.paymentErrChan:
				s.contestantLogger.Error("クリティカルエラーが発生しました", slog.String("error", err.Error()))
				return err
			default:
			}

			if s.world.Time%world.LengthOfHour == 0 {
				slog.Debug("時間経過", slog.Int64("tick", s.world.Time))
				s.evaluationMapLock.RLock()
				slog.Debug("eval", slog.Int("reqs", s.completedRequests), slog.Any("score", s.evalScoreMap), slog.Any("criteria", s.evaluationMap))
				s.evaluationMapLock.RUnlock()
			}
			if s.world.Time%(world.LengthOfHour*3) == 0 {
				s.contestantLogger.Info(fmt.Sprintf("これまでに地域内の評判によって%d人、既存ユーザーの招待経由で%d人が新規登録しました", s.world.NotInvitedUserCount.Load(), s.world.InvitedUserCount.Load()))
				if num := s.world.LeavedUserCount.Load(); num > 0 {
					s.contestantLogger.Info(fmt.Sprintf("これまでに低評価なライドによって%d人が利用をやめました", num))
				}
				s.evaluationMapLock.RLock()
				if s.completedRequests > 0 {
					s.contestantLogger.Info(fmt.Sprintf("%.1f%%のライドは椅子がマッチされるまでの時間、%.1f%%のライドはマッチされた椅子が乗車地点までに掛かる時間、%.1f%%のライドは椅子の実移動時間に不満がありました",
						(1-float64(s.evaluationMap[0])/float64(s.completedRequests))*100,
						(1-float64(s.evaluationMap[1])/float64(s.completedRequests))*100,
						(1-float64(s.evaluationMap[2]+s.evaluationMap[3])/float64(s.completedRequests*2))*100,
					))
				}
				s.evaluationMapLock.RUnlock()
			}
		}
	}

	return nil
}

func (s *Scenario) TotalSales() int64 {
	return lo.SumBy(s.world.OwnerDB.ToSlice(), func(p *world.Owner) int64 { return p.TotalSales.Load() })
}

func (s *Scenario) Score(final bool) int64 {
	if s.finalScore.Valid {
		return s.finalScore.Int64
	}
	score := lo.SumBy(s.world.OwnerDB.ToSlice(), func(p *world.Owner) int64 { return p.SubScore.Load() }) / 100
	if final {
		score += lo.SumBy(s.world.RequestDB.ToSlice(), func(r *world.Request) int64 {
			if r.Evaluated.Load() {
				return 0
			}
			return int64(r.PartialScore())
		}) / 100
		s.finalScore = null.IntFrom(score)
	}
	return score
}

func (s *Scenario) TotalDiscount() int64 {
	return lo.SumBy(s.world.RequestDB.ToSlice(), func(r *world.Request) int64 {
		if r.Evaluated.Load() {
			return int64(r.ActualDiscount())
		} else {
			return 0
		}
	})
}

func sendResult(s *Scenario, finished bool, passed bool) error {
	rawScore := s.Score(finished)
	if err := s.reporter.Report(&resources.BenchmarkResult{
		Finished: finished,
		Passed:   passed,
		Score:    rawScore,
		ScoreBreakdown: &resources.BenchmarkResult_ScoreBreakdown{
			Raw:       rawScore,
			Deduction: 0,
		},
		// Reason以外はsupervisorが設定する
		Execution: &resources.BenchmarkResult_Execution{
			Reason: "実行終了",
		},
		SurveyResponse: &resources.SurveyResponse{
			Language: s.language,
		},
	}); err != nil {
		return err
	}

	return nil
}
