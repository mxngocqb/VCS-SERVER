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

	layout := "2006-01-02"
	location, err := time.LoadLocation("Asia/Bangkok") // Load the GMT+7 timezone

	startTime, err := time.ParseInLocation(layout, in.Start, location)
	if err != nil {
		return &pb.SendReportResponse{Message: "Error parsing time: " + err.Error()}, nil
	} else {
		log.Println("Start time: ", startTime)
	}

	endTime, err := time.ParseInLocation(layout, in.End, location)
	if err != nil {
		return &pb.SendReportResponse{Message: "Error parsing time: " + err.Error()}, nil
	} else {
		log.Println("End time: ", endTime)
	}

	// Send report
	res := report.SendReport(in.Mail, startTime, endTime)

	if res != nil {
		return &pb.SendReportResponse{Message: "Error sending report: " + res.Error()}, nil
	} else {
		return &pb.SendReportResponse{Message: "Report sent successfully"}, nil
	}
}
