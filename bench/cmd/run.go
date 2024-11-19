package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucon14/bench/benchmarker/metrics"
	"github.com/isucon/isucon14/bench/benchmarker/scenario"
	"github.com/isucon/isucon14/bench/benchrun"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

var (
	// ベンチマークターゲット(URL)
	targetURL string
	// ベンチマークターゲット(ip:port)
	targetAddr string
	// ペイメントサーバのURL
	paymentURL string
	// 負荷走行秒数 (0のときは負荷走行を実行せずprepareのみ実行する)
	loadTimeoutSeconds int64
	// エラーが発生した際に非0のexit codeを返すかどうか
	failOnError bool
	// メトリクスを出力するかどうか
	exportMetrics bool
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a benchmark",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// supervisorで起動された場合は、targetを上書きする
		if benchrun.GetTargetAddress() != "" {
			targetURL = "https://trial.isucon14.net"
			targetAddr = fmt.Sprintf("%s:%d", benchrun.GetTargetAddress(), 443)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		contestantLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))

		var reporter benchrun.Reporter
		if fd, err := benchrun.GetReportFD(); err != nil {
			reporter = &benchrun.NullReporter{}
		} else {
			if reporter, err = benchrun.NewFDReporter(fd); err != nil {
				return fmt.Errorf("failed to create reporter: %w", err)
			}
		}

		exporter, err := metrics.Setup(!exportMetrics)
		if err != nil {
			return fmt.Errorf("failed to create meter: %w", err)
		}
		defer exporter.Shutdown(context.Background())

		slog.Debug("target", slog.String("targetURL", targetURL), slog.String("targetAddr", targetAddr), slog.String("benchrun.GetTargetAddress()", benchrun.GetTargetAddress()))

		s := scenario.NewScenario(targetURL, targetAddr, paymentURL, contestantLogger, reporter, otel.Meter("isucon14_benchmarker"), loadTimeoutSeconds == 0)

		b, err := isucandar.NewBenchmark(
			isucandar.WithoutPanicRecover(),
			isucandar.WithLoadTimeout(time.Duration(loadTimeoutSeconds)*time.Second),
		)
		if err != nil {
			return fmt.Errorf("failed to create benchmark: %w", err)
		}
		b.AddScenario(s)

		var errors []error
		if loadTimeoutSeconds == 0 {
			contestantLogger.Info("prepareのみを実行します")
			result := b.Start(context.Background())
			errors = result.Errors.All()
			contestantLogger.Info("prepareが終了しました",
				slog.Any("errors", errors),
			)
		} else {
			contestantLogger.Info("負荷走行を開始します")
			result := b.Start(context.Background())
			errors = result.Errors.All()
			contestantLogger.Info("負荷走行が終了しました",
				slog.Int64("score", s.Score()),
				slog.Any("errors", errors),
			)
		}

		if failOnError && len(errors) > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	runCmd.Flags().StringVar(&targetURL, "target", "http://localhost:8080", "benchmark target url")
	runCmd.Flags().StringVar(&targetAddr, "addr", "", "benchmark target ip:port")
	runCmd.Flags().StringVar(&paymentURL, "payment-url", "http://localhost:12345", "payment server URL")
	runCmd.Flags().Int64VarP(&loadTimeoutSeconds, "load-timeout", "t", 60, "load timeout in seconds (When this value is 0, load does not run and only prepare is run)")
	runCmd.Flags().BoolVar(&failOnError, "fail-on-error", false, "fail on error")
	runCmd.Flags().BoolVar(&exportMetrics, "metrics", false, "whether to output metrics")
	rootCmd.AddCommand(runCmd)
}
