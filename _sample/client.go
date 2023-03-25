package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	hellopb "main/pkg/grpc"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Data struct {
	Name string
	Message string
	CreatedAt string
}

var (
	scanner bufio.Scanner
	client hellopb.GreetingServiceClient
	list []Data
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
		fmt.Println("6: Chat")
		fmt.Println("7: Chat from Browser")
		fmt.Println("8: exit")
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
			fmt.Println("チャットからexitするには\\qを入力")
			go GetMessage(cancel)

			CreateMessage()
		case "6":
			fmt.Println("チャットからexitするには\\qを入力")
			fmt.Printf("please enter your name ->")
			scanner.Scan()
			name := scanner.Text()
			Chat(name)
		case "7":
			Serve()
		case "8":
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

		if msg == "\\q" {
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

func Chat(name string) {
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Println(err)
	}

	inputch := make(chan *hellopb.MessageRequest)
	go input(inputch,os.Stdin,name)

	rcvCh := make(chan *hellopb.MessageResponse)
	go receive(rcvCh,stream)

	ctx,cancel := context.WithCancel(context.TODO())
	defer cancel()

	for {
		select {
		case <- ctx.Done():
			fmt.Println("Done")
			stream.CloseSend()
			return
		case v := <-rcvCh:
			if (*v).GetMessage() == "" {
				continue
			}
			fmt.Printf(">[%s] : %s \n",(*v).Name,(*v).Message)
		case v:= <-inputch:
			if (*v).GetMessage() == "\\q" {
				return
			}
			if err := stream.Send(&hellopb.MessageRequest{
				Name: (*v).GetName(),
				Message: (*v).GetMessage(),
				CreatedAt: (*v).GetCreatedAt(),
			}); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func receive(ch chan<- *hellopb.MessageResponse,stream hellopb.GreetingService_ChatClient) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Println("ok")
			close(ch)
			return
		}

		if err != nil {
			log.Fatal(err)
		}
		ch<- in
	}
}

func input(ch chan<- *hellopb.MessageRequest,r io.Reader,name string) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		input := hellopb.MessageRequest{
			Name: name,
			Message: s.Text(),
			CreatedAt: timestamppb.Now(),
		}
		ch<- &input
	}
}

func inputFromBrowser(ch chan<- *hellopb.MessageRequest,r *http.Request,name string) {
	s := r.FormValue("message")

	input := hellopb.MessageRequest{
		Name: name,
		Message: s,
		CreatedAt: timestamppb.Now(),
	}
	ch<- &input

}

func Serve() {
	port := "3000"

	mux := http.NewServeMux()
	mux.HandleFunc("/",Index)
	mux.HandleFunc("/chat",ChatHandler)
	
	go func() {
		log.Printf("Started Server: port : %v\n",port)
		http.ListenAndServe(fmt.Sprintf(":%s",port),mux)	
	}()

	quit := make(chan os.Signal,1)
	signal.Notify(quit,os.Interrupt)
	<- quit
	log.Printf("Stopping Web Server")
}

func ChatHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method != "POST" {
		return
	}
	
	//ブラウザから入力されたデータ
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	res := &Data{}
	err = json.Unmarshal(body, res)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	name := res.Name
	data := &Data{
		Name: name,
		Message: res.Message,
		CreatedAt: res.CreatedAt,
	}

	list = append(list,*data)
	//gRPC通信部分ここから
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Println(err)
	}

	inputCh := make(chan *hellopb.MessageRequest)
	go inputFromBrowser(inputCh,r,name)

	rcvCh := make(chan *hellopb.MessageResponse)
	go receive(rcvCh,stream)

	ctx,cancel := context.WithCancel(context.TODO())
	defer cancel()

	for {
		select {
		case <- ctx.Done():
			fmt.Println("Done")
			stream.CloseSend()
			return
		case v := <-rcvCh:
			if (*v).GetMessage() == "" {
				continue
			}

			data = &Data{
				Name: (*v).GetName(),
				Message: (*v).GetMessage(),
				CreatedAt: (*v).CreatedAt.String(),
			}
			list = append(list, *data)
		case v := <- inputCh:
			if (*v).GetMessage() == "\\q" {
				return
			}
			if err := stream.Send(&hellopb.MessageRequest{
				Name: name,
				Message: (*v).GetMessage(),
				CreatedAt: timestamppb.Now(),
			}); err != nil {
				log.Fatal(err)
			}
		}
	}
	//gRPC通信部分ここまで
}

func Index(w http.ResponseWriter,r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../template/chat.html"))
	tmpl.Execute(w,list)
}