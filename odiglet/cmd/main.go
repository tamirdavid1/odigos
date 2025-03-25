package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/odigos-io/odigos/distros"
	"github.com/odigos-io/odigos/odiglet"
	"github.com/odigos-io/odigos/odiglet/pkg/ebpf"
	"github.com/odigos-io/odigos/odiglet/pkg/ebpf/sdks"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/odigos-io/odigos/common"
	commonInstrumentation "github.com/odigos-io/odigos/instrumentation"
	pb "github.com/odigos-io/odigos/k8sutils/pkg/instrumented_process/proto"
	"github.com/odigos-io/odigos/odiglet/pkg/env"
	"github.com/odigos-io/odigos/odiglet/pkg/instrumentation"
	"github.com/odigos-io/odigos/odiglet/pkg/instrumentation/instrumentlang"
	"github.com/odigos-io/odigos/odiglet/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	_ "net/http/pprof"
)

func main() {
	// Init Kubernetes clientset
	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	// Increase the QPS and Burst to avoid client throttling
	// Observed mainly in large clusters once updating big amount of instrumentationInstances
	cfg.QPS = 200
	cfg.Burst = 200

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	// If started in init mode
	if len(os.Args) == 2 && os.Args[1] == "init" {
		odiglet.OdigletInitPhase(clientset)
	}

	if err := log.Init(); err != nil {
		panic(err)
	}

	log.Logger.V(0).Info("Starting odiglet")

	// Load env
	if err := env.Load(); err != nil {
		log.Logger.Error(err, "Failed to load env")
		os.Exit(1)
	}

	dg, err := distros.NewCommunityGetter()
	if err != nil {
		log.Logger.Error(err, "Failed to create distro getter")
		os.Exit(1)
	}

	instrumentationManagerOptions := ebpf.InstrumentationManagerOptions{
		Factories:          ebpfInstrumentationFactories(),
		DistributionGetter: dg,
	}

	address := "ui.odigos-system.svc.cluster.local:50051"
	log.Logger.V(0).Info("Connecting to gRPC server", "address", address)

	conn, err := grpc.Dial(address, grpc.WithInsecure()) // ⚠️ For production: use grpc.WithTransportCredentials
	if err != nil {
		log.Logger.V(0).Error(err, "Failed to connect to gRPC server")
		return
	}
	defer conn.Close()

	client := pb.NewProcessServiceClient(conn)

	stream, err := client.SendProcessStream(context.Background())
	if err != nil {
		log.Logger.V(0).Error(err, "Could not open gRPC stream")
		return
	}

	// Simulate sending a batch of processes
	batch := &pb.InstrumentedProcessBatch{}
	for i := 0; i < 3; i++ {
		process := &pb.InstrumentedProcess{
			OdigletName:          "odiglet-demo",
			K8SNodeName:          "node-1",
			WorkloadName:         "demo-app",
			WorkloadKind:         "Deployment",
			K8SPodName:           fmt.Sprintf("demo-pod-%d", i),
			ProcessPid:           int32(1000 + i),
			K8SNamespaceName:     "default",
			K8SContainerName:     "app-container",
			CreatedAt:            timestamppb.New(time.Now()),
			TelemetrySdkLanguage: "go",
			ServiceInstanceId:    "svc-demo",
			Healthy:              true,
			HealthyReason:        "Running fine",
		}
		batch.Processes = append(batch.Processes, process)
	}

	log.Logger.V(0).Info("Sending batch of InstrumentedProcesses", "count", len(batch.Processes))
	if err := stream.Send(batch); err != nil {
		log.Logger.V(0).Error(err, "Failed to send process batch")
		return
	}

	resp, err := stream.Recv()
	if err != nil {
		log.Logger.V(0).Error(err, "Failed to receive response from server")
		return
	}

	log.Logger.V(0).Info("Received response from gRPC server",
		"status", resp.Status,
		"message", resp.Message,
		"processed_count", resp.ProcessedCount)

	o, err := odiglet.New(clientset, deviceInjectionCallbacks(), instrumentationManagerOptions)
	if err != nil {
		log.Logger.Error(err, "Failed to initialize odiglet")
		os.Exit(1)
	}

	ctx := signals.SetupSignalHandler()
	o.Run(ctx)

	log.Logger.V(0).Info("odiglet exiting")
}

func deviceInjectionCallbacks() instrumentation.OtelSdksLsf {
	return map[common.ProgrammingLanguage]map[common.OtelSdk]instrumentation.LangSpecificFunc{
		common.GoProgrammingLanguage: {
			common.OtelSdkEbpfCommunity: instrumentlang.Go,
		},
		common.JavaProgrammingLanguage: {
			common.OtelSdkNativeCommunity: instrumentlang.Java,
		},
		common.PythonProgrammingLanguage: {
			common.OtelSdkNativeCommunity: instrumentlang.Python,
		},
		common.JavascriptProgrammingLanguage: {
			common.OtelSdkNativeCommunity: instrumentlang.NodeJS,
		},
		common.DotNetProgrammingLanguage: {
			common.OtelSdkNativeCommunity: instrumentlang.DotNet,
		},
		common.NginxProgrammingLanguage: {
			common.OtelSdkNativeCommunity: instrumentlang.Nginx,
		},
	}
}

func ebpfInstrumentationFactories() map[commonInstrumentation.OtelDistribution]commonInstrumentation.Factory {
	return map[commonInstrumentation.OtelDistribution]commonInstrumentation.Factory{
		commonInstrumentation.OtelDistribution{
			Language: common.GoProgrammingLanguage,
			OtelSdk:  common.OtelSdkEbpfCommunity,
		}: sdks.NewGoInstrumentationFactory(),
	}
}
