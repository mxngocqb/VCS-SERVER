package util

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/xuri/excelize/v2"
)

type ParseErrors struct {
	Errors []string
}

func (e *ParseErrors) Error() string {
	return fmt.Sprintf("parse errors: %v", e.Errors)
}

// ParseExcel parses the uploaded Excel file into a slice of Server models, similar to ParseCSV.
func ParseExcel(file io.Reader) ([]model.Server, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	// Check available sheets
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in the Excel file")
	} else{
		log.Printf("Found %d sheets in the Excel file", len(sheets))
	}

	// Selecting the first sheet as default or a specific sheet by name
	sheetName := sheets[0] // default to the first sheet

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet '%s': %v", sheetName, err)
	} else{
		log.Printf("Successfully read %d rows from sheet '%s'", len(rows), sheetName)
	}

	var servers []model.Server
	var errors []string

	for i, row := range rows {
		if i == 0 { // Skip header row
			continue
		}

		if len(row) != 3 {
			errors = append(errors, fmt.Sprintf("line %d: incorrect number of fields (%d)", i+1, len(row)))
			continue // Skip rows with incorrect number of fields
		}

		status, err := strconv.ParseBool(row[1]) // Safely parse boolean status
		if err != nil {
			errors = append(errors, fmt.Sprintf("line %d: error parsing status for '%s'", i+1, row[1]))
			continue // Add error and continue with next line
		}

		if !IsValidIPv4(row[2]) {
			errors = append(errors, fmt.Sprintf("line %d: invalid IPv4 address '%s'", i+1, row[2]))
			continue // Skip rows with invalid IP addresses
		}

		servers = append(servers, model.Server{
			Name:   row[0],
			Status: status,
			IP:     row[2],
		})
		log.Printf("Parsed server: %s, %t, %s", row[0], status, row[2])
	}

	if len(errors) > 0 {
		return servers, &ParseErrors{Errors: errors}
	}
	return servers, nil
}


// createExcelWithServers creates an Excel file with server data
func CreateExcelFile(servers []model.Server) (*excelize.File, error) {
	f := excelize.NewFile()
	index, err := f.NewSheet("Servers")
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)

	// Create header
	headers := []string{"ID", "Name", "Status", "IP", "Created At", "Updated At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Servers", cell, header)
	}

	// Fill data
	for i, server := range servers {
		row := i + 2 // Starting from the second row
		f.SetCellValue("Servers", "A"+strconv.Itoa(row), server.ID)
		f.SetCellValue("Servers", "B"+strconv.Itoa(row), server.Name)
		f.SetCellValue("Servers", "C"+strconv.Itoa(row), server.Status)
		f.SetCellValue("Servers", "D"+strconv.Itoa(row), server.IP)
		f.SetCellValue("Servers", "E"+strconv.Itoa(row), server.CreatedAt.Format(time.RFC3339))
		f.SetCellValue("Servers", "F"+strconv.Itoa(row), server.UpdatedAt.Format(time.RFC3339))
	}

	return f, nil
}