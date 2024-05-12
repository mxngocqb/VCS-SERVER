package main

import (
	"log"
	"net"

	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report"
	pb "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report/proto"
	"google.golang.org/grpc"
)

var addr string = "127.0.0.1:50051"

type Server struct {
	pb.ReportServiceServer
}

func main() {
	// Schedule daily report
	report.ScheduleDailyReport()

	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	} else{
		log.Printf("Listening on %v", addr)
	}

	s := grpc.NewServer()
	pb.RegisterReportServiceServer(s , &Server{})

	if err := s.Serve(lis); err!= nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	

	// Keep the program running
	select {}
}
