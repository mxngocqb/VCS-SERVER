package main

import (
	"context"
	"log"
	"time"

	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report"
	pb "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report/proto"
)

func (s *Server) Report(ctx context.Context, in *pb.SendReportRequest) (*pb.SendReportResponse, error) {
	log.Printf("Received: %v", in)
	
	// Parse start and end time
	startTime, err := time.Parse(time.RFC3339, in.Start)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(time.RFC3339, in.End)
	if err != nil {
		return nil, err
	}
	
	// Send report
	res := report.SendReport(in.Mail, startTime, endTime)

	if res != nil {
		return &pb.SendReportResponse{Message: "Error sending report: " + res.Error()}, nil
	} else {
		return &pb.SendReportResponse{Message: "Report sent successfully"}, nil
	}
}
