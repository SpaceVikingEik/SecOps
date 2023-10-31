package main

//This is largely a template of these sources combined, which establishes a p2p TLS communication scheme:
// https://think.unblog.ch/en/how-to-install-openssl-on-windows-10-11/
// https://raymii.org/s/tutorials/OpenSSL_generate_self_signed_cert_with_Subject_Alternative_name_oneliner.html
// https://pkg.go.dev/google.golang.org/grpc/credentials
// https://github.com/NaddiNadja/peer-to-peer/blob/main/main.go
// The hospital is initialized on localhost 5000 and i do some hardcoding around that fact below
import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"

	SecOps "github.com/SpaceVikingEik/SecOps/Security/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 64)
	ownPort := int64(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:      ownPort,
		shares:  []int64{0},
		clients: make(map[int64]SecOps.SecOpsClient),
		ctx:     ctx,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	serverCertification, err := credentials.NewServerTLSFromFile("certificate/server.crt", "certificate/priv.key")
	if err != nil {
		log.Fatalln("failed to create cert", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(serverCertification))
	SecOps.RegisterSecOpsServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	for i := 0; i < 4; i++ {
		port := int64(5000) + int64(i)

		if port == ownPort {
			continue
		}

		clientCertification, err := credentials.NewClientTLSFromFile("certificate/server.crt", "")
		if err != nil {
			log.Fatalln("failed to create cert", err)
		}

		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", port), grpc.WithTransportCredentials(clientCertification), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()
		c := SecOps.NewSecOpsClient(conn)
		p.clients[port] = c
		fmt.Printf("%v", p.clients)
	}
	fmt.Printf("\n")
	fmt.Printf("Connections established\n")
	fmt.Printf("If you are a client type in your secret, when ready to send shares of other clients to hospital type send\n")
	fmt.Printf("if you are the hospital and have received the secrets type compute\n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		if p.id == 5000 && scanner.Text() == "compute" {
			sum := int64(0)
			for _, num := range p.shares {
				sum += num
			}
			fmt.Printf("Result: %v \n", sum)
		} else if p.id != 5000 && scanner.Text() == "send" {
			p.sendSharesToHospital()
		} else if p.id != 5000 {
			secret, _ := strconv.ParseInt(scanner.Text(), 10, 64)
			p.sendSharesToAllClients(secret)
		}
	}
}

type peer struct {
	SecOps.UnimplementedSecOpsServer
	id      int64
	shares  []int64
	clients map[int64]SecOps.SecOpsClient
	ctx     context.Context
}

// mustEmbedUnimplementedSecOpsServer implements SecOps.SecOpsServer.
func (*peer) mustEmbedUnimplementedSecOpsServer() {
	panic("unimplemented")
}

func (p *peer) Ping(ctx context.Context, req *SecOps.Share) (*SecOps.Reply, error) {
	share := req.Share
	p.shares = append(p.shares, share)
	fmt.Printf("Current share list \n\n", p.shares)
	rep := &SecOps.Reply{Success: true}
	return rep, nil
}

func (p *peer) sendSharesToAllClients(secret int64) {

	temp := secret
	i := int64(5001)

	reset := len(p.clients) + 5001
	//fmt.Printf("%v", reset)
	for temp > 1 {

		if i == p.id {
			i = p.id + 1
			if i >= int64(reset) {
				i = 5001
			}
		}
		temp2 := rand.Int63n(temp)
		temp = temp - temp2
		if temp == 1 {
			temp2 = temp2 + 1
		}
		//fmt.Printf("%v", temp2)
		request := &SecOps.Share{Share: temp2}

		currentClient := p.clients[int64(i)]
		reply, err := currentClient.Ping(p.ctx, request)
		if err != nil {
			fmt.Println("something went wrong")
		}
		fmt.Printf("Got reply from id %v: %v\n", (i), reply.Success)

		i = i + 1
		if i >= int64(reset) {
			i = 5001
		}
		/*
			for id, client := range p.clients {
				if id != 5000 {
					reply, err := client.Ping(p.ctx, request)
					if err != nil {
						fmt.Println("something went wrong")
					}
					fmt.Printf("Got reply from id %v: %v\n", id, reply.Success)
				}
			} */

	}

}

func (p *peer) sendSharesToHospital() {
	sum := int64(0)
	for _, num := range p.shares {
		sum += num
	}
	request := &SecOps.Share{Share: sum}
	for id, client := range p.clients {
		if id == 5000 {
			reply, err := client.Ping(p.ctx, request)
			if err != nil {
				fmt.Println("something went wrong")
			}
			fmt.Printf("Got reply from id %v: %v\n", id, reply.Success)
		}
	}
}
