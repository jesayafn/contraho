package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: contraho <subcommand> [options]")
		os.Exit(1)
	}

	subcommand := os.Args[1]

	switch subcommand {
	case "project":
		projectSearch()
	default:
		fmt.Println("Invalid subcommand:", subcommand)
		os.Exit(1)
	}
}

func projectSearch() {
	// This is project scrape
	host, username, password, fileOutput := arguments()

	credential := authorizationHeader(*username, *password)

	lengthProject := projectLength(*host, credential)

	// fmt.Println(lengthProject)
	projectList := applyOwnerInformation(
		applyBranchDetail(
			listProject(*host, credential, lengthProject),
			*host, credential,
		),
		*host, credential,
	)

	file, err := os.Create(*fileOutput)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the CSV header (field names)
	header := getStructFieldNames(projectList[0])
	writer.Write(header)

	// Write the data rows

	for _, i := range projectList {
		row := getStructFieldValues(i)
		writer.Write(row)
	}

	fmt.Println("CSV file generated successfully!")
}
