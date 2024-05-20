package report

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

func ScheduleDailyReport(elasticService *ElasticService) {
	loc, _ := time.LoadLocation("Asia/Bangkok") // Ensure timezone consistency with server logs
	c := cron.New(cron.WithLocation(loc))
	// Send daily report at 8:00 AM
	_, err := c.AddFunc("44 9 * * *", func() {
		// _, err := c.AddFunc("@every 10m", func() {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 0, 1)
		err1 := SendReport([]string{"mxn111333@gmail.com"}, start, end, elasticService)
		if err1 != nil {
			log.Printf("Error sending daily report: %v", err1)
		} else {
			log.Printf("Daily report sent successfully")
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling daily report: %v", err)
	}
	c.Start()
}

// SendReport sends a daily server report to the specified email addresses.
// It fetches server information for the given time range and constructs an email
// containing the average uptime, number of online and offline servers, and the total number of servers.
// The email is sent using the SMTP protocol with the provided Gmail account credentials.
// If any error occurs during the process, it is returned. Otherwise, nil is returned.
func SendReport(email []string, start, end time.Time, elasticService *ElasticService) error {
	avgUptime, online, offline, totalServers, err := elasticService.elastic.FetchServersInfo(start, end)
	if err != nil {
		log.Printf("Error fetching server info: %v", err)
		return err
	}

	fmt.Println("Email:", email)

	m := gomail.NewMessage()
	m.SetHeader("From", "mxngocqb@gmail.com")
	m.SetHeader("To", email...)
	m.SetHeader("Subject", "Daily Server Report")
	m.SetBody("text/html", fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				color: #333;
			}
			.container {
				margin: 20px;
				padding: 20px;
				border: 1px solid #ddd;
				border-radius: 5px;
				background-color: #f9f9f9;
			}
			h2 {
				color: #4CAF50;
			}
			.data {
				margin-top: 20px;
			}
			.data strong {
				color: #555;
			}
			.footer {
				margin-top: 30px;
				font-size: 0.9em;
				color: #777;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Daily Server Report</h2>
			<p>Hello,</p>
			<p>Here is the daily server report for your servers:</p>
			<p><strong>From:</strong> %s<br>
			<strong>To:</strong> %s</p>
			<div class="data">
				<p><strong>Average Uptime:</strong> %.2f hours</p>
				<p><strong>Online Servers:</strong> %d</p>
				<p><strong>Offline Servers:</strong> %d</p>
				<p><strong>Total Servers:</strong> %d</p>
			</div>
			<div class="footer">
				<p>Best regards,<br>
				Server Management Systems - VCS Team</p>
			</div>
		</div>
	</body>
	</html>`,
		start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), avgUptime, online, offline, totalServers))

	d := gomail.NewDialer("smtp.gmail.com", 587, "mxngocqb@gmail.com", "xftw lchz hruo ojkq")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	} else {
		log.Println("Email sent successfully")
		return nil
	}
}
