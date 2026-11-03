package clustercollector

import (
	"errors"
	"sync"
	"time"

	odigoscommon "github.com/odigos-io/odigos/common"
	commonlogger "github.com/odigos-io/odigos/common/logger"
)

var (
	errInvalidTraceAggregationValue = errors.New("invalid collectorGateway.traceAggregation value, using the default")

	// The gateway config is rebuilt on every reconcile; warn about the deprecated
	// field once per process instead of on each one.
	legacyWaitDurationWarning sync.Once
)

// traceAggregationOptions resolves collectorGateway.traceAggregation for the
// gateway's groupby processor. The deprecated
// sampling.tailSampling.traceAggregationWaitDuration fills syncWait and asyncWait
// when they are unset. Invalid values are dropped so the processor default applies
// instead of the gateway failing to load its config.
func traceAggregationOptions(logger *commonlogger.ContextLogger, odigosCfg *odigoscommon.OdigosConfiguration) *odigoscommon.TraceAggregationConfiguration {
	var in odigoscommon.TraceAggregationConfiguration
	if odigosCfg != nil && odigosCfg.CollectorGateway != nil && odigosCfg.CollectorGateway.TraceAggregation != nil {
		in = *odigosCfg.CollectorGateway.TraceAggregation
	}
	if legacy := legacyTraceAggregationWait(odigosCfg); legacy != nil {
		legacyWaitDurationWarning.Do(func() {
			logger.Info("sampling.tailSampling.traceAggregationWaitDuration is deprecated, use collectorGateway.traceAggregation.syncWait and asyncWait")
		})
		if in.SyncWait == nil {
			in.SyncWait = legacy
		}
		if in.AsyncWait == nil {
			in.AsyncWait = legacy
		}
	}

	out := &odigoscommon.TraceAggregationConfiguration{
		SyncWait:   validDuration(logger, "syncWait", in.SyncWait, false),
		AsyncWait:  validDuration(logger, "asyncWait", in.AsyncWait, false),
		RootSettle: validDuration(logger, "rootSettle", in.RootSettle, true),
	}
	if in.MaxTraces != nil {
		if *in.MaxTraces > 0 {
			out.MaxTraces = in.MaxTraces
		} else {
			logger.Error(errInvalidTraceAggregationValue, "ignoring traceAggregation setting", "field", "maxTraces", "value", *in.MaxTraces)
		}
	}
	return out
}

func legacyTraceAggregationWait(odigosCfg *odigoscommon.OdigosConfiguration) *string {
	if odigosCfg == nil || odigosCfg.Sampling == nil || odigosCfg.Sampling.TailSampling == nil {
		return nil
	}
	//nolint:staticcheck // still honored while deprecated, see TraceAggregationConfiguration
	return odigosCfg.Sampling.TailSampling.TraceAggregationWaitDuration
}

func validDuration(logger *commonlogger.ContextLogger, field string, v *string, allowZero bool) *string {
	if v == nil || *v == "" {
		return nil
	}
	d, err := time.ParseDuration(*v)
	if err != nil || d < 0 || (d == 0 && !allowZero) {
		logger.Error(errInvalidTraceAggregationValue, "ignoring traceAggregation setting", "field", field, "value", *v)
		return nil
	}
	return v
}
