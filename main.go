package main

import (
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
	host, username, password, fileOutput, otheroptions := arguments("project")

	credential := authorizationHeader(*username, *password)

	switch {
	case otheroptions["unlistedApp"] == true:
		lengthProjectPage := projectSearchApiLength(*host, credential, "TRK")

		projectList := applyOwnerInformation(
			qualityGateofProject(
				applyBranchDetail(
					projectFiltering(
						listProject(
							*host,
							credential,
							lengthProjectPage,
						),
						*host,
						credential,
						0,
					),
					*host,
					credential,
				),
				*host,
				credential,
			),
			*host,
			credential,
		)
		if *fileOutput != "" {
			createCSVFile(*fileOutput, projectList)
		} else {
			// clearScreen()
			printStructTable(projectList, "Key", "Name", "Branch", "Loc", "Owner")
		}

	case otheroptions["listedApp"] == true:
		lengthProjectPage := projectSearchApiLength(*host, credential, "TRK")

		projectList := applyOwnerInformation(
			qualityGateofProject(
				applyBranchDetail(
					projectFiltering(
						listProject(
							*host,
							credential,
							lengthProjectPage,
						),
						*host,
						credential,
						1,
					),
					*host,
					credential,
				),
				*host,
				credential,
			),
			*host,
			credential,
		)

		if *fileOutput != "" {
			createCSVFile(*fileOutput, projectList)
		} else {
			// clearScreen()
			printStructTable(projectList, "Key", "Name", "Branch", "Loc", "Owner")
		}
	default:
		lengthProjectPage := projectSearchApiLength(*host, credential, "TRK")

		// fmt.Println(lengthProject)
		projectList := applyOwnerInformation(
			qualityGateofProject(
				applyBranchDetail(
					listProject(
						*host,
						credential,
						lengthProjectPage,
					),
					*host,
					credential,
				),
				*host,
				credential,
			),
			*host,
			credential,
		)

		if *fileOutput != "" {
			createCSVFile(*fileOutput, projectList)
		} else {
			// clearScreen()

			printStructTable(projectList, "Key", "Name", "Branch", "Loc", "Owner")
			// printStructTable(projectList)
		}
	}

}
