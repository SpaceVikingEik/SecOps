package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	SecOps "github.com/SpaceVikingEik/SecOps/Security/grpc"
	"google.golang.org/grpc"
)

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:            ownPort,
		amountOfPings: make(map[int32]int32),
		clients:       make(map[int32]SecOps.SecOpsClient),
		ctx:           ctx,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	grpcServer := grpc.NewServer()
	SecOps.RegisterSecOpsServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	for i := 0; i < 4; i++ {
		port := int32(5000) + int32(i)

		if port == ownPort {
			continue
		}

		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := SecOps.NewSecOpsClient(conn)
		p.clients[port] = c
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		p.sendPingToAll()
	}
}

type peer struct {
	SecOps.UnimplementedSecOpsServer
	id            int32
	amountOfPings map[int32]int32
	clients       map[int32]SecOps.SecOpsClient
	ctx           context.Context
}

// mustEmbedUnimplementedSecOpsServer implements SecOps.SecOpsServer.
func (*peer) mustEmbedUnimplementedSecOpsServer() {
	panic("unimplemented")
}

func (p *peer) Ping(ctx context.Context, req *SecOps.Request) (*SecOps.Reply, error) {
	id := req.Id
	p.amountOfPings[id] += 1

	rep := &SecOps.Reply{Amount: p.amountOfPings[id]}
	return rep, nil
}

func (p *peer) sendPingToAll() {
	request := &SecOps.Request{Id: p.id}
	for id, client := range p.clients {
		reply, err := client.Ping(p.ctx, request)
		if err != nil {
			fmt.Println("something went wrong")
		}
		fmt.Printf("Got reply from id %v: %v\n", id, reply.Amount)
	}
}
