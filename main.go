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
	case "project", "proj", "p":
		projectSearch()
	case "application", "app", "a":
		appSearch()
	default:
		fmt.Println("Invalid subcommand:", subcommand)
		os.Exit(1)
	}
}

func projectSearch() {
	startTime := time.Now()

	host, credential, authMode, fileOutput, otheroptions := arguments(0)

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
	return metricProject(
		ownerProject(
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
					).([]ProjectSearchList),
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
		).([]ProjectSearchList),
		host,
		credential,
		authMode,
	).([]ProjectSearchList)
}

func appSearch() {
	startTime := time.Now()

	host, credential, authMode, fileOutput, _ := arguments(1)

	lengthAppPage := projectSearchApiLength(*host, credential, "APP", authMode)

	appListInterface := listApp(*host, credential, lengthAppPage, authMode, 1)

	// Type assertion to convert the interface to []AppList
	appList, ok := appListInterface.([]AppList)
	if !ok {
		// Handle the case where the returned value is not []AppList
		fmt.Println("Error: Unexpected return type from listApp")
		return
	}
	appList = languageofApp(appList, *host, credential, authMode)
	appList = ownerProject(appList, *host, credential, authMode).([]AppList)
	appList = branchDetailOfProjects(appList, *host, credential, authMode).([]AppList)
	appList = metricProject(appList, *host, credential, authMode).([]AppList)

	if *fileOutput != "" {
		createCSVFile(*fileOutput, appList)
	} else {
		printStructTable(appList, "Key", "Name", "Loc")
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Seconds()

	fmt.Printf("Execution Time: %.3f seconds\n", elapsedTime)

}
