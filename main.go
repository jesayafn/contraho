package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	host, username, password, fileOutput := arguments()

	credential := authorizationHeader(*username, *password)
	data := httpRequest(*host+"/api/projects/search?qualifiers=TRK&ps=1&p=1", credential)

	// projectSearchPage := projectSearchParse(data)\
	var projectSearchPage ProjectSearchPage
	_ = dataParse(data, &projectSearchPage)

	indexPageNumber := indexPageNumberCounter(projectSearchPage)

	var dataProjectSearch []ProjectSearchList
	for pagesIndex := 1; pagesIndex <= indexPageNumber; pagesIndex++ {
		dataprojectSearchPageRaw := httpRequest(*host+"/api/projects/search?qualifiers=TRK&ps=500&p="+strconv.Itoa(pagesIndex),
			credential)

		var dataprojectSearchPageParsed ProjectSearchPage
		_ = dataParse(dataprojectSearchPageRaw, &dataprojectSearchPageParsed)
		for projectsIndex := range dataprojectSearchPageParsed.Components {

			lastAnalysisDate, err := time.Parse("2006-01-02T15:04:05-0700",
				dataprojectSearchPageParsed.Components[projectsIndex].LastAnalysisDate)
			// lastAnalysisDate := formatTime(lastAnalysisDateRaw)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			dataProjectSearch = append(dataProjectSearch, ProjectSearchList{
				Key:              dataprojectSearchPageParsed.Components[projectsIndex].Key,
				Name:             dataprojectSearchPageParsed.Components[projectsIndex].Name,
				LastAnalysisDate: lastAnalysisDate,
			})
		}
	}

	for i := range dataProjectSearch {
		dataBranchesListRaw := httpRequest(*host+"/api/project_branches/list?project="+dataProjectSearch[i].Key, credential)
		var dataBranchesListParsed ProjectBranchesList

		_ = dataParse(dataBranchesListRaw, &dataBranchesListParsed)

		var compareLOC []int
		for y := range dataBranchesListParsed.Branches {
			dataMeasuresRaw := httpRequest(
				*host+"/api/measures/component?metricKeys=ncloc&component="+dataProjectSearch[i].Key+"&branch="+dataBranchesListParsed.Branches[y].Name,
				credential)

			var dataMeasuresParsed ProjectMeasures
			_ = dataParse(dataMeasuresRaw, &dataMeasuresParsed)

			var stringifyLoc int

			if len(dataMeasuresParsed.Component.Measures) == 0 {
				stringifyLoc = 0
			} else {
				stringifyLoc, _ = strconv.Atoi(dataMeasuresParsed.Component.Measures[0].Value)
			}

			compareLOC = append(compareLOC, stringifyLoc)

		}

		branchHighestLoc := findIndexOfHighestValue(compareLOC)

		dataProjectSearch[i].HighestBranch = dataBranchesListParsed.Branches[branchHighestLoc].Name
		dataProjectSearch[i].Loc = strconv.Itoa(compareLOC[branchHighestLoc])

		dataPermissionsListRaw := httpRequest(*host+"/api/permissions/users?projectKey="+dataProjectSearch[i].Key, credential)
		var dataPermissionsListParsed ProjectPermissions
		_ = dataParse(dataPermissionsListRaw, &dataPermissionsListParsed)

		dataProjectSearch[i].Owner = dataPermissionsListParsed.Users[0].Name
		dataProjectSearch[i].Email = dataPermissionsListParsed.Users[0].Email
	}
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
	header := getStructFieldNames(dataProjectSearch[0])
	writer.Write(header)

	// Write the data rows

	for _, i := range dataProjectSearch {
		row := getStructFieldValues(i)
		writer.Write(row)
	}

	fmt.Println("CSV file generated successfully!")
}
