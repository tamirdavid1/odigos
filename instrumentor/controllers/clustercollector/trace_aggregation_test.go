package clustercollector

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/odigos-io/odigos/api/k8sconsts"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/common/api/sampling"
	"github.com/odigos-io/odigos/common/config"
	odigosconsts "github.com/odigos-io/odigos/common/consts"
	commonlogger "github.com/odigos-io/odigos/common/logger"
	pipelinegen "github.com/odigos-io/odigos/common/pipelinegen"
)

func ptr[T any](v T) *T { return &v }

func TestTraceAggregationOptions(t *testing.T) {
	gateway := func(ta *common.TraceAggregationConfiguration) *common.OdigosConfiguration {
		return &common.OdigosConfiguration{CollectorGateway: &common.CollectorGatewayConfiguration{TraceAggregation: ta}}
	}
	withLegacy := func(cfg *common.OdigosConfiguration, wait string) *common.OdigosConfiguration {
		cfg.Sampling = &common.SamplingConfiguration{TailSampling: &sampling.TailSamplingConfiguration{TraceAggregationWaitDuration: &wait}}
		return cfg
	}

	tests := []struct {
		name string
		cfg  *common.OdigosConfiguration
		want common.TraceAggregationConfiguration
	}{
		{
			name: "nothing set keeps every processor default",
			cfg:  &common.OdigosConfiguration{},
		},
		{
			name: "new fields pass through",
			cfg: gateway(&common.TraceAggregationConfiguration{
				SyncWait: ptr("20s"), AsyncWait: ptr("45s"), RootSettle: ptr("5s"), MaxTraces: ptr(1000),
			}),
			want: common.TraceAggregationConfiguration{
				SyncWait: ptr("20s"), AsyncWait: ptr("45s"), RootSettle: ptr("5s"), MaxTraces: ptr(1000),
			},
		},
		{
			name: "deprecated wait fills both waits",
			cfg:  withLegacy(&common.OdigosConfiguration{}, "2s"),
			want: common.TraceAggregationConfiguration{SyncWait: ptr("2s"), AsyncWait: ptr("2s")},
		},
		{
			name: "new fields win over the deprecated wait",
			cfg:  withLegacy(gateway(&common.TraceAggregationConfiguration{SyncWait: ptr("10s")}), "2s"),
			want: common.TraceAggregationConfiguration{SyncWait: ptr("10s"), AsyncWait: ptr("2s")},
		},
		{
			name: "invalid values fall back to the defaults",
			cfg: gateway(&common.TraceAggregationConfiguration{
				SyncWait: ptr("soon"), AsyncWait: ptr("0s"), RootSettle: ptr("-1s"), MaxTraces: ptr(0),
			}),
		},
		{
			name: "zero root settle disables it",
			cfg:  gateway(&common.TraceAggregationConfiguration{RootSettle: ptr("0s")}),
			want: common.TraceAggregationConfiguration{RootSettle: ptr("0s")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := traceAggregationOptions(&commonlogger.ContextLogger{Logger: logr.Discard()}, tt.cfg)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestGatewayConfig_GroupByWritesOnlySetFields(t *testing.T) {
	on := true
	tier := common.OnPremOdigosTier
	insightsCfg := effectiveInsightsConfig(&common.InsightsConfiguration{Enabled: &on}, tier)
	ext := k8sconsts.OdigosConfigK8sExtensionType
	tailSampling := false
	gatewayOptions := pipelinegen.GatewayConfigOptions{
		OdigosNamespace:           "odigos-system",
		OdigosConfigExtensionName: &ext,
		TailSamplingEnabled:       &tailSampling,
		Insights:                  insightsCfg,
		InsightsOtlpEndpoint:      k8sconsts.InsightsOtlpGrpcDNSEndpoint("odigos-system"),
		TraceAggregation:          &common.TraceAggregationConfiguration{RootSettle: ptr("3s"), MaxTraces: ptr(5000)},
	}

	cfg, err, _, _ := pipelinegen.CalculateGatewayConfig(nil, nil,
		func(c *config.Config, destinationPipelineNames []string, signalsRootPipelines []string) error {
			if err := addSelfTelemetryPipeline(c, 8888, destinationPipelineNames, signalsRootPipelines); err != nil {
				return err
			}
			return addInsightsGatewayExporter(c, "odigos-system", insightsCfg)
		},
		nil, &gatewayOptions)
	require.NoError(t, err)

	assert.Equal(t, config.GenericMap{"root_settle": "3s", "max_traces": 5000}, cfg.Processors[odigosconsts.GroupByTraceProcessor])
}
