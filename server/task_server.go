package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"task-tracker/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedTaskServiceServer
	tasks []proto.TaskResponse
}

func (s *server) AddTask(ctx context.Context, req *proto.TaskRequest) (*proto.TaskResponse, error) {
	id := fmt.Sprintf("task-%d", len(s.tasks)+1)
	newTask := proto.TaskResponse{Id: id, Title: req.Title, Description: req.Description}
	s.tasks = append(s.tasks, newTask)
	return &newTask, nil
}

func (s *server) GetTasks(ctx context.Context, req *proto.Empty) (*proto.TaskListResponse, error) {
	var tasksWithPointers []*proto.TaskResponse
	for i := range s.tasks {
		task := s.tasks[i]
		tasksWithPointers = append(tasksWithPointers, &task)
	}

	return &proto.TaskListResponse{Tasks: tasksWithPointers}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterTaskServiceServer(s, &server{})
	reflection.Register(s)

	fmt.Println("gRPC Server listening on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
