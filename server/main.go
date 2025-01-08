// Mock server
package main

import (
	"context"
	pb "github.com/usercanal/sdk-go/proto"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTestServiceServer
}

func (s *server) SendMessage(ctx context.Context, event *pb.Event) (*pb.EventResponse, error) {
	// Handle based on event type
	switch event.Type {
	case pb.Event_TRACK:
		if track := event.GetTrack(); track != nil {
			// Enrich the base message
			if track.Base != nil {
				track.Base.ContextId = uuid.New().String()
				track.Base.MessageId = uuid.New().String()
			}
		}
	case pb.Event_IDENTIFY:
		if identify := event.GetIdentify(); identify != nil {
			if identify.Base != nil {
				identify.Base.ContextId = uuid.New().String()
				identify.Base.MessageId = uuid.New().String()
			}
		}
	case pb.Event_GROUP:
		if group := event.GetGroup(); group != nil {
			if group.Base != nil {
				group.Base.ContextId = uuid.New().String()
				group.Base.MessageId = uuid.New().String()
			}
		}
	}

	// Log the enriched message
	log.Printf("Enriched message: %+v", event)

	// Return response with message ID and server timestamp
	return &pb.EventResponse{
		MessageId:       uuid.New().String(),
		ServerTimestamp: time.Now().Unix(),
	}, nil
}

func (s *server) SendBatch(ctx context.Context, request *pb.BatchRequest) (*pb.BatchResponse, error) {
	responses := make([]*pb.EventResponse, 0, len(request.Events))

	// Process each event in the batch
	for _, event := range request.Events {
		response, err := s.SendMessage(ctx, event)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}

	return &pb.BatchResponse{
		Responses: responses,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTestServiceServer(s, &server{})
	log.Printf("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
