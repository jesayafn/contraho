package main

import (
	"fmt"
	"os"
	"time"
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
	startTime := time.Now()

	host, username, password, token, authMode, fileOutput, otheroptions := arguments("project")
	var credential string
	switch authMode {
	case 0:
		credential = *token
	case 1:
		credential = authorizationHeader(*username, *password)
	}
	lengthProjectPage := projectSearchApiLength(*host, credential, "TRK", authMode)

	var projectList []ProjectSearchList
	switch {
	case otheroptions["unlistedApp"] == true:
		projectList = project(*host, credential, authMode, lengthProjectPage, 0)
	case otheroptions["listedApp"] == true:
		projectList = project(*host, credential, authMode, lengthProjectPage, 1)
	default:
		projectList = project(*host, credential, authMode, lengthProjectPage, -1)
	}
	if *fileOutput != "" {
		createCSVFile(*fileOutput, projectList)
	} else {
		printStructTable(projectList, "Key", "Name", "Branch", "Loc", "Owner")
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Seconds()

	fmt.Printf("Execution Time: %.3f seconds\n", elapsedTime)

}

func project(host string, credential string, authMode int, lengthProjectPage int, filterMode int) []ProjectSearchList {
	return metricProject(ownerProject(
		languageofProject(
			qualityGateofProject(
				branchDetailOfProjects(
					projectFiltering(
						listProject(
							host,
							credential,
							lengthProjectPage,
							authMode,
						),
						host,
						credential,
						filterMode,
						authMode,
					),
					host,
					credential,
					authMode,
				),
				host,
				credential,
				authMode,
			),
			host,
			credential,
			authMode,
		),
		host,
		credential,
		authMode,
	),
		host,
		credential,
		authMode,
	)
}
