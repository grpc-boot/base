package monitor

import (
	"github.com/grpc-boot/base/v2/utils"

	"google.golang.org/grpc/codes"
)

type Options struct {
	Name         string       `json:"name" yaml:"name"`
	Gauges       []string     `json:"gauges" yaml:"gauges"`
	CodeGauges   []string     `json:"CodeGauges" yaml:"CodeGauges"`
	ResetSeconds int          `json:"resetSeconds" yaml:"resetSeconds"`
	ResetTimes   int          `json:"resetTimes" yaml:"resetTimes"`
	CodeList     []codes.Code `json:"codeList" yaml:"codeList"`
}

var (
	DefaultOptions = func() Options {
		return Options{
			Name:         "bm",
			Gauges:       []string{GaugePanicCount},
			CodeGauges:   []string{GaugeRequestCount, GaugeResponseCount, GaugeRequestLen, GaugeResponseLen},
			ResetSeconds: defaultResetSeconds,
			ResetTimes:   defaultResetTimes,
			CodeList: []codes.Code{
				utils.OK,
				utils.CodeCanceled,
				utils.CodeUnknown,
				utils.CodeInvalidArgument,
				utils.CodeDeadlineExceeded,
				utils.CodeNotFound,
				utils.CodeAlreadyExists,
				utils.CodePermissionDenied,
				utils.CodeResourceExhausted,
				utils.CodeFailedPrecondition,
				utils.CodeAborted,
				utils.CodeOutRange,
				utils.CodeUnimplemented,
				utils.CodeInternal,
				utils.CodeUnavailable,
				utils.CodeDataLoss,
				utils.CodeUnauthenticated,
			},
		}
	}
)

func formatOpt(opt Options) Options {
	defOpt := DefaultOptions()

	if opt.Name == "" {
		opt.Name = defOpt.Name
	}

	if opt.ResetSeconds < 1 {
		opt.ResetSeconds = defOpt.ResetSeconds
	}

	if opt.ResetTimes < 1 {
		opt.ResetTimes = defOpt.ResetTimes
	}

	if len(opt.Gauges) < 1 {
		opt.Gauges = defOpt.Gauges
	}

	if len(opt.CodeGauges) < 1 {
		opt.CodeGauges = defOpt.CodeGauges
	}

	if len(opt.CodeList) < 1 {
		opt.CodeList = defOpt.CodeList
	}

	return opt
}
