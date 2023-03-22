package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	hellopb "main/pkg/grpc"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	scanner bufio.Scanner
	client hellopb.GreetingServiceClient
)

func main() {
	cancel := make(chan struct{})

	fmt.Println("start gRPC Client")

	scanner = *bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	
	if err != nil {
		log.Printf("cannot open connection")
		return
	}

	defer conn.Close()

	client = hellopb.NewGreetingServiceClient(conn)

	for {
		fmt.Println("1: Hello")
		fmt.Println("2: HelloServerStream")
		fmt.Println("3: HelloClientStream")
		fmt.Println("4: HelloBiStream")
		fmt.Println("5: Chat")
		fmt.Println("6: exit")
		fmt.Print("please enter -> ")

		scanner.Scan()
		in := scanner.Text()
		
		switch in {
		case "1":
			Hello()
		case "2":
			HelloServerStream()
		case "3":
			HelloClientStream()
		case "4":
			HelloBiStream()
		case "5":
			go GetMessage(cancel)

			CreateMessage()
		case "6":
			fmt.Println("Bye")
			goto M
		}
	}
M:
close(cancel)
}

func Hello() {
	fmt.Print("enter your name ->")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	res, err := client.Hello(context.Background(),req)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloServerStream() {
	fmt.Print("please enter name ->")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}

	stream, err := client.HelloServerStream(context.Background(),req)
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err,io.EOF) {
			fmt.Println("all the response have already received")
			break
		}

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
}

func HelloClientStream() {
	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}

	sendCount := 5
	fmt.Printf("Please enter %d names\n",sendCount)
	for i := 0;i < sendCount; i++ {
		scanner.Scan()
		name := scanner.Text()

		if err := stream.Send(&hellopb.HelloRequest{
			Name: name,
		});err != nil {
			log.Fatal(err)
			return
		}
	}
	
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloBiStream() {
	stream, err := client.HelloBiStream(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	sendNum := 5
	fmt.Printf("Please enter %d names\n",sendNum)

	var sendEnd,rcvEnd bool
	sendCount := 0
	for !(sendEnd && rcvEnd) {
		if !sendEnd {
			scanner.Scan()
			name :=  scanner.Text()

			sendCount++
			if err := stream.Send(&hellopb.HelloRequest{
				Name: name,
			});err != nil {
				log.Println(err)
				sendEnd = true
			}

			if sendCount == sendNum {
				sendEnd = true
				if err := stream.CloseSend();err != nil {
					fmt.Println(err)
				}
			}
		}

		if !rcvEnd {
			if res, err := stream.Recv();err != nil  {
				if errors.Is(err,io.EOF) {
					fmt.Println(err)
				}
				rcvEnd= true
			} else {
				fmt.Println(res.GetMessage())
			}
		}
	}
}

func CreateMessage() {
	fmt.Printf("Please enter your name -> ")
	scanner.Scan()
	name := scanner.Text()

	timestamp := timestamppb.New(time.Now())

	req_msg := &hellopb.MessageRequest{
		Name: name,
		Message: "",
		CreatedAt: timestamp,
	}

	_, err := client.CreateMessage(context.Background(),req_msg)
	if err != nil {
		return
	} else {
		fmt.Println("you are joined room")
	}

	for {
		scanner.Scan()
		msg := scanner.Text()

		if msg == "/q" {
			return
		}

		timestamp := timestamppb.New(time.Now())

		req_msg := &hellopb.MessageRequest{
			Name: name,
			Message: msg,
			CreatedAt: timestamp,
		}

		_, err := client.CreateMessage(context.Background(),req_msg)
		if err != nil {
			return
		}
	}
}

func GetMessage(ch chan struct{}) {
	stream, err := client.GetMessages(context.Background(),&emptypb.Empty{})
	if err != nil {
		return
	}

	for {
		select {
		case <- ch:
			return
		default:
			res, err := stream.Recv()
			if errors.Is(err,io.EOF) {
				break
			}
	
			if err != nil {
				return
			}
			if res.Message != "" {
				fmt.Printf("[%s] : %s (%s)\n",res.Name,res.Message,res.CreatedAt)
			}
		}
	}
}