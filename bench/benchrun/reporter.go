//go:build !windows

package benchrun

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"syscall"

	"github.com/isucon/isucon14/bench/benchrun/gen/isuxportal/resources"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Reporter interface {
	Report(result *resources.BenchmarkResult) error
}

// NullReporter is for testing
type NullReporter struct{}

func (rep *NullReporter) Report(result *resources.BenchmarkResult) error {
	return nil
}

type FDReporter struct {
	io *bufio.Writer
}

func (rep *FDReporter) Report(result *resources.BenchmarkResult) error {
	if result.SurveyResponse != nil && len(result.SurveyResponse.Language) > 140 {
		return errors.New("language in a given survey response is too long (max: 140)")
	}

	if result.MarkedAt == nil {
		result.MarkedAt = timestamppb.Now()
	}

	wire, err := proto.Marshal(result)
	if err != nil {
		return err
	}

	// PIPEのデータサイズの上限65536であるが、Headerが2byteなので、上限として65535(0xFFFF)を設定
	if len(wire) > 65535 {
		return errors.New("marshalled BenchmarkResult is too long (max: 65535)")
	}

	lenBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBuf, uint16(len(wire)))

	if err := rep.Write(lenBuf); err != nil {
		return err
	}
	if err := rep.Write(wire); err != nil {
		return err
	}
	if err := rep.Flush(); err != nil {
		return err
	}

	return nil
}

func (rep *FDReporter) Write(buf []byte) error {
	// write(2) は PIPE_BUF (linux なら 4096) 以上の書き込みはatomicにならないので書き込みが終わるまで繰り返す
	if n, err := rep.io.Write(nil); err != nil {
		if errors.Is(err, io.ErrShortWrite) {
			return rep.Write(buf[n:])
		}
		return err
	}
	return nil
}

func (rep *FDReporter) Flush() error {
	return rep.io.Flush()
}

func NewFDReporter(fd int) (*FDReporter, error) {
	syscall.CloseOnExec(fd)
	io := os.NewFile(uintptr(fd), "ISUXBENCH_REPORT_FD")

	bufWriter := bufio.NewWriter(io)

	rep := &FDReporter{
		io: bufWriter,
	}

	return rep, nil
}
