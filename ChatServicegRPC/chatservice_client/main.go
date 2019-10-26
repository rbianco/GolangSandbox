package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/rbianco/GolangSandbox/ChatServicegRPC/chatservice"
	"google.golang.org/grpc"
)

func registerContact(c pb.ChatServiceClient, email *string, name *string) {
	contact := &pb.Contact{
		Email: *email,
		Name:  *name,
	}

	r, err := c.Register(context.Background(), &pb.AddClientRequest{Contact: contact})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Id: %s", r.GetId())
	log.Printf("Code: %s", r.GetCode())
}

//ShowContacts(ctx context.Context, in *ShowContactsRequest, opts ...grpc.CallOption) (ChatService_ShowContactsClient, error)
func printContacts(c pb.ChatServiceClient) {
	log.Printf("Lista de contactos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := c.ShowContacts(ctx, &pb.ShowContactsRequest{Page: 1, PageSize: 20})

	if err != nil {
		log.Fatalf("Error al obtener la lista de contactos, %v", err)
	}

	for {
		userContact, err := stream.Recv()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatalf("Error al obtener la lista de contactos, %s", err)
		}

		log.Printf("Contact %s con email: %s ", userContact.Contact.Name, userContact.Contact.Email)
	}

}

func receiveMessages(c pb.ChatServiceClient, email *string) {
	stream, err := c.ReceiveMessage(context.Background(), &pb.GetMessagesRequest{Email: *email})
	if err != nil {
		log.Fatalf("Failed to receive chat messages : %v", err)
	}

	chatc := make(chan struct{})
	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				// read done.
				log.Print("Deja de escuchar")
				close(chatc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive chat messages : %v", err)
			}
			log.Printf("%s --> %s: %s", message.MessageFrom, message.MessageTo, message.Text)
		}
	}()
	stream.CloseSend()
	<-chatc
}

func sendMessages(c pb.ChatServiceClient, email *string) {
	reader := bufio.NewReader(os.Stdin)
	stream, err := c.SendMessage(context.Background())
	if err != nil {
		log.Fatalf("Failed to send chat messages : %v", err)
	}
	for {
		log.Printf("Enter a contact")
		to, _ := reader.ReadString('\n')
		to = strings.Trim(to, "\r\n")

		log.Printf("Enter a message")
		message, _ := reader.ReadString('\n')
		message = strings.Trim(message, "\r\n")

		sendErr := stream.Send(&pb.ChatMessage{MessageFrom: *email, MessageTo: to, Text: message})

		if sendErr != nil && sendErr != io.EOF {
			log.Fatalf("Failed to send chat messages : %v", sendErr)
		}
		// log.Printf("%s --> %s: %s", message.MessageFrom, message.MessageTo, message.Text)
	}
}

func main() {
	// Contact the server and print out its response.
	var (
		register = flag.Bool("register", true, "Indica si se debe registrar al usuario")
		email    = flag.String("email", "test@test.com", "Email del usuario a registrar")
		name     = flag.String("name", "Juan", "Nombre del usuario a registrar")
	)
	flag.Parse()

	//172.20.24.104
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := pb.NewChatServiceClient(conn)

	if *register {
		registerContact(c, email, name)
	}

	printContacts(c)

	chatc := make(chan struct{})
	go func() {
		sendMessages(c, email)
	}()
	receiveMessages(c, email)
	<-chatc

}
