package main

import (
	"log"

	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report"
	pb "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report/proto"
)

var addr string = "0.0.0.0:50052"

type Server struct {
	pb.ReportServiceServer
	reportService *report.ElasticService
}

func main() {
	
	// load config
	cfgPath := "./conf.yaml"
	cfg, err := config.Load(cfgPath)
	
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	err = report.Start(cfg)
	if err != nil {
		log.Fatalf("Error starting report service: %v", err)
	}	
}
