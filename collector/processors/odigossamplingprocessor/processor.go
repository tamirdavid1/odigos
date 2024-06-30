package odigossamplingprocessor

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/odigos/processor/odigossamplingprocessor/internal/sampling"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type samplingProcessor struct {
	logger *zap.Logger
	config *Config
}

func (sp *samplingProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	for _, rule := range sp.config.Rules {
		return ptrace.Traces{}, fmt.Errorf("processing rule %s, %s", rule.Name, rule.Type)

		switch rule.RuleDetails.(type) {
		case *sampling.TraceLatencyRule:
			if rule.RuleDetails.(*sampling.TraceLatencyRule).TraceDropDecision(td) {
				sp.removeAllSpans(&td)
				return td, nil
			}
		case *sampling.ErrorRule:
			if rule.RuleDetails.(*sampling.ErrorRule).TraceDropDecision(td) {
				sp.removeAllSpans(&td)
				return td, nil
			}
		default:
			sp.logger.Error("Unknown rule details type", zap.String("rule", rule.Name))
		}
	}
	return td, nil
}
func (sp *samplingProcessor) removeAllSpans(td *ptrace.Traces) {
	td.ResourceSpans().RemoveIf(func(rs ptrace.ResourceSpans) bool { return true })
}
