package main

import (
	"fmt"
	"os"
	"time"
)

const (
	emptyString = ""
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

	host, credential, authMode, csvOutput, pdfOutput, pagingOutput, otherOptions := arguments(0)

	lengthProjectPage := projectSearchApiLength(*host, credential, "TRK", authMode)

	var projectList []ProjectSearchList
	switch {
	case otherOptions["unlistedApp"] == true:
		projectList = project(*host, credential, authMode, lengthProjectPage, 0, emptyString)
	case otherOptions["listedApp"] == true:
		projectList = project(*host, credential, authMode, lengthProjectPage, 1, emptyString)
	case otherOptions["app"] != "":
		projectList = project(*host, credential, authMode, lengthProjectPage, 1, otherOptions["app"].(string))
	default:
		projectList = project(*host, credential, authMode, lengthProjectPage, -1, emptyString)
	}
	if *csvOutput != "" {
		createCSVFile(*csvOutput, startTime, projectList)
	} else if *pdfOutput != "" {
		generatePDF(*pdfOutput, projectList, "Key", "Name", "Branch", "Loc", "Owner")
	} else {
		printStructTable(projectList, startTime, *pagingOutput, "Key", "Name", "Branch", "Loc", "Owner")

		// printStructAsTable(projectList, []string{"Key", "Name", "Branch", "Loc", "Owner"})
	}

}

func project(host string, credential string, authMode int, lengthProjectPage int, filterMode int, appName string) []ProjectSearchList {
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
							appName,
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

	host, credential, authMode, csvOutput, pdfOutput, pagingOutput, _ := arguments(1)

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

	if *csvOutput != "" {
		createCSVFile(*csvOutput, startTime, appList)
	} else if *pdfOutput != "" {
		err := generatePDF(*pdfOutput, appList,
			"Key", "Name", "Loc", "Email", "Owner")
		handleErr(err)
	} else {
		printStructTable(appList, startTime, *pagingOutput, "Key", "Name", "Loc")
	}
}
