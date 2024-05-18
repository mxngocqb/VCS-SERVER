package main

import (
	"log"
	"net"

	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/elastic"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report"
	pb "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report/proto"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
	"google.golang.org/grpc"
)

var addr string = "127.0.0.1:50052"

type Server struct {
	pb.ReportServiceServer
	reportService *report.ElasticService
}

func main() {
	
	// load config
	cfgPath := "./conf.yaml"
	cfg, err := config.Load(cfgPath)
	elasticClient, err := elastic.ConnectElasticSearch(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to ElasticSearch: %v", err)
	}
	elasticService := elastic.NewElasticsearch(elasticClient)
	reportEs := report.NewReportService(elasticService)
	// Schedule daily report
	report.ScheduleDailyReport(reportEs)

	logger := util.GRPCLog()
    opts := []grpc.ServerOption{
        grpc.UnaryInterceptor(unaryLoggingInterceptor(logger)),
    }

	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	} else{
		log.Printf("Listening on %v", addr)
	}

	s := grpc.NewServer(opts...)
	pb.RegisterReportServiceServer(s , &Server{
		reportService: reportEs,
	})

	if err := s.Serve(lis); err!= nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	
	// Keep the program running
	select {}
}
