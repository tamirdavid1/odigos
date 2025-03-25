// db/grpc_server.go
package db

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/odigos-io/odigos/k8sutils/pkg/instrumented_process/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedProcessServiceServer
}

func (s *Server) SendProcessStream(stream pb.ProcessService_SendProcessStreamServer) error {
	for {
		batch, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream closed by client")
			return nil
		}
		if err != nil {
			log.Printf("Error receiving: %v", err)
			return err
		}

		count := len(batch.Processes)
		for _, p := range batch.Processes {
			log.Printf("Received process: Pod=%s, PID=%d", p.K8SPodName, p.ProcessPid)
		}

		resp := &pb.ProcessStreamResponse{
			Status:         "received",
			Message:        fmt.Sprintf("Received %d processes", count),
			ProcessedCount: int32(count),
		}

		if err := stream.Send(resp); err != nil {
			log.Printf("Failed to send response: %v", err)
			return err
		}
	}
}

func StartGRPCServer(ctx context.Context) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProcessServiceServer(s, &Server{})

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		log.Println("Shutting down gRPC server...")
		s.GracefulStop()
	}()

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil && ctx.Err() == nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
