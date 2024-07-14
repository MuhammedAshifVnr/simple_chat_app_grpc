package service

import (
	pb "chatappp/proto"
	"sync"
)

type Service struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[pb.ChatService_ChatServer]struct{}
}

func NewChatService() *Service {
	return &Service{
		clients: make(map[pb.ChatService_ChatServer]struct{}),
	}
}

func (s *Service) Chat(stream pb.ChatService_ChatServer) error {
	s.mu.Lock()
	s.clients[stream] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, stream)
		s.mu.Unlock()
	}()

	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		s.broadcast(msg)
	}
}

func (s *Service) broadcast(msg *pb.ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for client := range s.clients {
		if err := client.Send(msg); err != nil {
			delete(s.clients, client)
		}
	}
}
