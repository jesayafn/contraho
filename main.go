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

	switch {
	case otheroptions["unlistedApp"] == true:
		lengthProjectPage := projectSearchApiLength(*host, credential, "TRK", authMode)
		projectList := ownerProject(
			languageofProject(
				qualityGateofProject(
					branchDetailOfProjects(
						projectFiltering(
							listProject(
								*host,
								credential,
								lengthProjectPage,
								authMode,
							),
							*host,
							credential,
							0,
							authMode,
						),
						*host,
						credential,
						authMode,
					),
					*host,
					credential,
					authMode,
				),
				*host,
				credential,
				authMode,
			),
			*host,
			credential,
			authMode,
		)

		if *fileOutput != "" {
			createCSVFile(*fileOutput, projectList)
		} else {
			// clearScreen()
			printStructTable(projectList, "Key", "Name", "Branch", "Loc", "Owner")
		}

	case otheroptions["listedApp"] == true:
		lengthProjectPage := projectSearchApiLength(*host, credential, "TRK", authMode)
		projectList := ownerProject(
			languageofProject(
				qualityGateofProject(
					branchDetailOfProjects(
						projectFiltering(
							listProject(
								*host,
								credential,
								lengthProjectPage,
								authMode,
							),
							*host,
							credential,
							1,
							authMode,
						),
						*host,
						credential,
						authMode,
					),
					*host,
					credential,
					authMode,
				),
				*host,
				credential,
				authMode,
			),
			*host,
			credential,
			authMode,
		)

		if *fileOutput != "" {
			createCSVFile(*fileOutput, projectList)
		} else {
			// clearScreen()
			printStructTable(projectList, "Key", "Name", "Branch", "Loc", "Owner")
		}
	default:
		lengthProjectPage := projectSearchApiLength(*host, credential, "TRK", authMode)
		projectList := ownerProject(
			languageofProject(
				qualityGateofProject(
					branchDetailOfProjects(
						listProject(
							*host,
							credential,
							lengthProjectPage,
							authMode,
						),
						*host,
						credential,
						authMode,
					),
					*host,
					credential,
					authMode,
				),
				*host,
				credential,
				authMode,
			),
			*host,
			credential,
			authMode,
		)

		// fmt.Println(lengthProject)

		if *fileOutput != "" {
			createCSVFile(*fileOutput, projectList)
		} else {
			// clearScreen()

			printStructTable(projectList, "Key", "Name", "Branch", "Loc", "Owner")
			// printStructTable(projectList)
		}
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Seconds()

	fmt.Printf("Execution Time: %.3f seconds\n", elapsedTime)

}
