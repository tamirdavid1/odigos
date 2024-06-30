package sampling

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/pdata/ptrace"
)

type ErrorRule struct {
	statusCodes []ptrace.StatusCode `mapstructure:"status_codes"`
}

var SupportedStatuses = map[ptrace.StatusCode]struct{}{
	ptrace.StatusCodeOk:    {},
	ptrace.StatusCodeError: {},
	ptrace.StatusCodeUnset: {},
}

func (tlr *ErrorRule) Validate() error {
	for _, statusCode := range tlr.statusCodes {
		if _, found := SupportedStatuses[statusCode]; !found {
			return errors.New("unsupported status code")
		}
	}
	return nil
}

func (tlr *ErrorRule) TraceDropDecision(td ptrace.Traces) bool {
	for i := 0; i < td.ResourceSpans().Len(); i++ {
		rs := td.ResourceSpans().At(i)

		ilss := rs.ScopeSpans()

		for i := 0; i < ilss.Len(); i++ {
			ils := ilss.At(i)

			for j := 0; j < ils.Spans().Len(); j++ {
				span := ils.Spans().At(j)

				for _, statusCode := range tlr.statusCodes {
					if span.Status().Code() == statusCode {
						return false
					}
				}
			}
		}
	}
	fmt.Println("Dropping trace")
	return true
}
