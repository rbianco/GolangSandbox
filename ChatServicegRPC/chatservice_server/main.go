package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	pb "github.com/rbianco/GolangSandbox/ChatServicegRPC/chatservice"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedChatServiceServer
	users    []userContact
	messages []userMessage
}

type userContact struct {
	Email string
	Name  string
}

type userMessage struct {
	From        string
	To          string
	Message     string
	hasBeenSent bool
}

func (s *server) Register(ctx context.Context, in *pb.AddClientRequest) (*pb.AddClientResponse, error) {
	log.Printf("Solicitud de registro recibida")

	for _, userContact := range s.users {
		if in.Contact.Email == userContact.Email {
			return &pb.AddClientResponse{Id: in.Contact.Email, Code: "Duplicated"}, nil
		}
	}

	s.users = append(s.users, userContact{Name: in.Contact.Name, Email: in.Contact.Email})
	return &pb.AddClientResponse{Id: in.Contact.Email, Code: "Success"}, nil
}

func (s *server) ShowContacts(request *pb.ShowContactsRequest, stream pb.ChatService_ShowContactsServer) error {
	log.Printf("Solicitud para mostrar contactos recibida")
	for _, userContact := range s.users {
		stream.Send(&pb.ShowContactsResponse{Contact: &pb.Contact{Email: userContact.Email, Name: userContact.Name}})
	}

	return nil
}

//SendMessage(ChatService_SendMessageServer) error
func (s *server) SendMessage(stream pb.ChatService_SendMessageServer) error {
	log.Printf("Solicitud de envío de mensaje recibida")
	for {
		message, err := stream.Recv()
		if err != nil {
			return err
		}

		for _, userContact := range s.users {
			if message.MessageTo == userContact.Email {
				log.Printf("Mensaje %s de %s para %s", message.Text, message.MessageFrom, message.MessageTo)
				s.messages = append(s.messages, userMessage{From: message.MessageFrom, To: message.MessageTo, Message: message.Text})
				break
			}
		}
	}
}

//ReceiveMessage(*GetMessagesRequest, ChatService_ReceiveMessageServer) error
func (s *server) ReceiveMessage(request *pb.GetMessagesRequest, stream pb.ChatService_ReceiveMessageServer) error {
	log.Printf("Solicitud para recibir mensajes recibida")
	for {

		if len(s.messages) == 0 {
			time.Sleep(2 * time.Second)
			continue
		}

		var sentMessagesCounter = 0
		for i, message := range s.messages {
			if message.To == request.Email {
				stream.Send(&pb.ChatMessage{MessageFrom: message.From, MessageTo: message.To, Text: message.Message})
				s.messages[i].hasBeenSent = true
				sentMessagesCounter++
			}
		}

		for i := len(s.messages) - 1; i >= 0; i-- {
			if sentMessagesCounter == 0 {
				break
			}

			if s.messages[i].hasBeenSent {
				s.messages = append(s.messages[:i], s.messages[i+1:]...)
				sentMessagesCounter--
			}
		}
	}
}

func main() {
	// Contact the server and print out its response.
	var (
		port = flag.String("port", ":50051", "Puerto en el que escucha el server")
	)
	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
