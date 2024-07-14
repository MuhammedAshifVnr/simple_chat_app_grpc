package main

import (
	"bufio"
	"context"
	"log"
	"os"

	chatpb "chatappp/proto"

	"github.com/Pallinder/go-randomdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := chatpb.NewChatServiceClient(conn)
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("failed to create stream: %v", err)
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Fatalf("error receiving message: %v", err)
			}
			log.Printf("%s: %s", msg.GetUser(), msg.GetMessage())
		}
	}()
	
	randName := randomdata.FirstName(randomdata.RandomGender)
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		stream.Send(&chatpb.ChatMessage{User: randName, Message: text})
	}
}
