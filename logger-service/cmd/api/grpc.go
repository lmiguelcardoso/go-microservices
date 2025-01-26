package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"
	"time"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, r *logs.LogRequest) (*logs.LogResponse, error) {
	input := r.GetLogEntry()

	logEntry := data.LogEntry{
		Name:      input.Name,
		Data:      input.Data,
		CreatedAt: time.Now(),
	}

	err := l.Models.LogEntry.Insert(logEntry)

	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "logged",
	}

	return res, err
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))

	if err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", grpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}
}
