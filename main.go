package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myapp <subcommand> [options]")
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
	// projectList := listProject(*host, credential, lengthProject)

	// projectList = branchDetailOfProjects(*host, credential, projectList)

	// projectList = ownerProject(*host, credential, projectList)
	// fmt.Printf("%+v\n", projectList)

	// // projectList := listProject(*host, credential, lengthProject)

	// // projectList = branchDetailOfProjects(*host, credential, projectList)

	// // projectList = ownerProject(*host, credential, projectList)

	// // var dataProjectSearch []ProjectSearchList
	// // for pagesIndex := 1; pagesIndex <= lengthProject; pagesIndex++ {
	// // 	dataprojectSearchPageRaw := httpRequest(*host+projectScrapeApi+strconv.Itoa(pagesIndex),
	// // 		credential)

	// // 	var dataprojectSearchPageParsed ProjectSearchPage
	// // 	_ = dataParse(dataprojectSearchPageRaw, &dataprojectSearchPageParsed)
	// // 	for projectsIndex := range dataprojectSearchPageParsed.Components {

	// // 		// lastAnalysisDate, err := time.Parse(timeParseFormat,
	// // 		// 	dataprojectSearchPageParsed.Components[projectsIndex].LastAnalysisDate)

	// // 		// if err != nil {
	// // 		// 	fmt.Println("Error:", err)
	// // 		// 	return
	// // 		// }
	// // 		dataProjectSearch = append(dataProjectSearch, ProjectSearchList{
	// // 			Key:  dataprojectSearchPageParsed.Components[projectsIndex].Key,
	// // 			Name: dataprojectSearchPageParsed.Components[projectsIndex].Name,
	// // 			// LastAnalysisDate: lastAnalysisDate,
	// // 		})
	// // 	}
	// // }

	// // for i := range projectList {
	// // 	dataBranchesListRaw := httpRequest(*host+projectBranchesApi+dataProjectSearch[i].Key, credential)
	// // 	var dataBranchesListParsed ProjectBranchesList

	// // 	_ = dataParse(dataBranchesListRaw, &dataBranchesListParsed)

	// // 	var compareLOC []int
	// // 	for y := range dataBranchesListParsed.Branches {
	// // 		dataMeasuresRaw := httpRequest(
	// // 			*host+ProjectBranchesLocApi+dataProjectSearch[i].Key+"&branch="+dataBranchesListParsed.Branches[y].Name,
	// // 			credential)

	// // 		var dataMeasuresParsed ProjectMeasures
	// // 		_ = dataParse(dataMeasuresRaw, &dataMeasuresParsed)

	// // 		var stringifyLoc int

	// // 		if len(dataMeasuresParsed.Component.Measures) == 0 {
	// // 			stringifyLoc = 0
	// // 		} else {
	// // 			stringifyLoc, _ = strconv.Atoi(dataMeasuresParsed.Component.Measures[0].Value)
	// // 		}

	// // 		compareLOC = append(compareLOC, stringifyLoc)

	// // 	}

	// // 	branchHighestLoc := findIndexOfHighestValue(compareLOC)

	// // 	dataProjectSearch[i].HighestBranch = dataBranchesListParsed.Branches[branchHighestLoc].Name
	// // 	dataProjectSearch[i].Loc = strconv.Itoa(compareLOC[branchHighestLoc])

	// // 	dataPermissionsListRaw := httpRequest(*host+ProjectUserPermissionsApi+dataProjectSearch[i].Key, credential)
	// // 	var dataPermissionsListParsed ProjectPermissions
	// // 	_ = dataParse(dataPermissionsListRaw, &dataPermissionsListParsed)

	// // 	dataProjectSearch[i].Owner = dataPermissionsListParsed.Users[0].Name
	// // 	dataProjectSearch[i].Email = dataPermissionsListParsed.Users[0].Email
	// // }
	// fmt.Println("Finish Project Listing")
	// // This is Application Scrape

	// // host, username, password, fileOutput := arguments()

	// // credential := authorizationHeader(*username, *password)
	// dataApl := httpRequest(*host+aplIndexApi, credential)

	// var aplSearchPage ProjectSearchPage
	// _ = dataParse(dataApl, &aplSearchPage)

	// aplIndexPageNumber := indexPageNumberCounter(aplSearchPage)
	// var dataAplSearch []string
	// for pagesIndex := 1; pagesIndex <= aplIndexPageNumber; pagesIndex++ {
	// 	dataAplSearchPageRaw := httpRequest(*host+apltScrapeApi+strconv.Itoa(pagesIndex),
	// 		credential)

	// 	var dataAplSearchPageParsed ProjectSearchPage
	// 	_ = dataParse(dataAplSearchPageRaw, &dataAplSearchPageParsed)
	// 	for aplIndex := range dataAplSearchPageParsed.Components {

	// 		// lastAnalysisDate, err := time.Parse(timeParseFormat,
	// 		// 	dataprojectSearchPageParsed.Components[projectsIndex].LastAnalysisDate)

	// 		// if err != nil {
	// 		// 	fmt.Println("Error:", err)
	// 		// 	return
	// 		// }
	// 		dataAplSearch = append(dataAplSearch,
	// 			dataAplSearchPageParsed.Components[aplIndex].Key)
	// 		// fmt.Println(dataAplSearchPageParsed.Components[aplIndex].Key)
	// 	}
	// }

	// fmt.Println("Finish Applications Listing")
	// fmt.Println(dataAplSearch)

	// var projectKeyListedApl []string
	// for aplKey := range dataAplSearch {
	// 	dataProjApl := httpRequest(*host+projectIndexAplApi+"&component="+dataAplSearch[aplKey], credential)
	// 	var projAplSearch ProjectSearchOfApplication
	// 	_ = dataParse(dataProjApl, &projAplSearch)
	// 	fmt.Println("On" + dataAplSearch[aplKey] + "application")
	// 	projAplIndexPageNumber := indexPageNumberCounter(projAplSearch)
	// 	fmt.Println(projAplIndexPageNumber, dataAplSearch[aplKey])
	// 	for pagesIndex := 1; pagesIndex <= projAplIndexPageNumber; pagesIndex++ {
	// 		fmt.Println("Paging" + dataAplSearch[aplKey])
	// 		dataProjAplSearchPageRaw := httpRequest(*host+projectScrapeAplApi+strconv.Itoa(pagesIndex)+"&component="+dataAplSearch[aplKey],
	// 			credential)

	// 		var dataProjAplSearchPageParsed ProjectSearchOfApplication
	// 		_ = dataParse(dataProjAplSearchPageRaw, &dataProjAplSearchPageParsed)
	// 		for aplIndex := range dataProjAplSearchPageParsed.Components {

	// 			// lastAnalysisDate, err := time.Parse(timeParseFormat,
	// 			// 	dataprojectSearchPageParsed.Components[projectsIndex].LastAnalysisDate)

	// 			// if err != nil {
	// 			// 	fmt.Println("Error:", err)
	// 			// 	return
	// 			// }
	// 			projectKeyListedApl = append(projectKeyListedApl, dataProjAplSearchPageParsed.Components[aplIndex].Key)
	// 		}
	// 	}
	// }
	// fmt.Println("Finish Project Key List")
	// fmt.Println(projectKeyListedApl)
	// dataProjectSearchUpdate := removeByKeys(dataProjectSearch, projectKeyListedApl)
	// fmt.Println(len(dataProjectSearchUpdate))

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
