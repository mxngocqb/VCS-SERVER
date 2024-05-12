package main

import (
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report/proto"
)

var addr string = "127.0.0.1:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}	
	
	defer conn.Close()

	c := pb.NewReportServiceClient(conn)

	doSendReport(c, []string{"mxn111333@gmail.com"}, "2021-01-01T00:00:00Z", "2021-01-02T00:00:00Z")
}