package client

import (
	"context"
	"log"
	tictactoepb "main/pkg/grpc"
)

const (
	address= "localhost:8080"
	defaultname = "hoge"
)

func greetServer(ctx context.Context,client tictactoepb.GameServiceClient,name string) error {
	r, err := client.Greetserver(ctx,&tictactoepb.GreetRequest{Msg: name})
	if err != nil {
		return err
	}
	log.Printf("%s",r.GetMsg())
	return nil
}