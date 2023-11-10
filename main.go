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

	if unlistedApp, ok := otheroptions["unlistedApp"].(bool); ok {

		if unlistedApp {

			lengthProjectPage := projectSearchApiLength(*host, credential, "TRK")

			// fmt.Println(lengthProject)
			projectList := listProject(*host, credential, lengthProjectPage)

			lengthAppPage := projectSearchApiLength(*host, credential, "APP")

			applicationList := listApp(*host, credential, lengthAppPage)
			var lisedProjectOnApp []string
			for i := range applicationList {
				lengthProjectOfAppPage := applicationProjectSearchApiLength(*host, applicationList[i], credential)
				// fmt.Println(lengthProjectOfAppPage)
				lisedProjectOnApp = listProjectofApplication(*host, lisedProjectOnApp, applicationList[i],
					lengthProjectOfAppPage, credential)
				fmt.Println(lisedProjectOnApp)
			}

			lisedProjectOnApp = removeRedundantValues(lisedProjectOnApp)

			// fmt.Println(lisedProjectOnApp)

			projectList = deleteProjectsByKeys(projectList, lisedProjectOnApp)

			// fmt.Println(projectList)

			createCSVFile(*fileOutput, projectList)
		} else {
			lengthProjectPage := projectSearchApiLength(*host, credential, "TRK")

			// fmt.Println(lengthProject)
			projectList := applyOwnerInformation(
				applyBranchDetail(
					listProject(*host, credential, lengthProjectPage),
					*host, credential,
				),
				*host, credential,
			)

			createCSVFile(*fileOutput, projectList)
		}
	}
}
