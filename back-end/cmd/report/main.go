package main

import (
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report"
)

func main() {
	report.ScheduleDailyReport()

	// Keep the program running
    select {}

}