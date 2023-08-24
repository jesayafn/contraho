package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	host, username, password := arguments()

	credential := authorizationHeader(*username, *password)
	data := httpRequest(*host+"/api/projects/search?qualifiers=TRK&ps=1&p=1", credential)

	// projectSearchPage := projectSearchParse(data)\
	var projectSearchPage ProjectSearchPage
	_ = dataParse(data, &projectSearchPage)

	indexPageNumber := indexPageNumberCounter(projectSearchPage)

	var dataProjectSearch []ProjectSearchList
	for i := 1; i <= indexPageNumber; i++ {
		dataprojectSearchPageRaw := httpRequest(*host+"/api/projects/search?qualifiers=TRK&ps=500&p="+strconv.Itoa(i), credential)
		// fmt.Println(i)
		// fmt.Println(indexPageNumber)
		// dataprojectSearchPageParsed := projectSearchParse(dataprojectSearchPageRaw)
		var dataprojectSearchPageParsed ProjectSearchPage
		_ = dataParse(dataprojectSearchPageRaw, &dataprojectSearchPageParsed)
		for i := range dataprojectSearchPageParsed.Components {

			lastAnalysisDate, err := time.Parse("2006-01-02T15:04:05-0700", dataprojectSearchPageParsed.Components[i].LastAnalysisDate)
			// lastAnalysisDate := formatTime(lastAnalysisDateRaw)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			dataProjectSearch = append(dataProjectSearch, ProjectSearchList{
				Key:              dataprojectSearchPageParsed.Components[i].Key,
				Name:             dataprojectSearchPageParsed.Components[i].Name,
				LastAnalysisDate: lastAnalysisDate,
			})
		}
		// fmt.Println(dataprojectSearchPageParsed.Components[i].LastAnalysisDate)
		// fmt.Println(len(dataProjectSearch))

		// fmt.Println(len(dataSearch))
		// fmt.Println(dataSearch)

		// respBody := bytes.NewBuffer(data).String()
		// fmt.Println(data)
	}

	for i := range dataProjectSearch {
		dataBranchesListRaw := httpRequest(*host+"/api/project_branches/list?project="+dataProjectSearch[i].Key, credential)
		var dataBranchesListParsed ProjectBranchesList

		_ = dataParse(dataBranchesListRaw, &dataBranchesListParsed)

		// fmt.Println(dataProjectSearch[i].Name, dataBranchesListParsed)
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
		// fmt.Println(dataPermissionsListParsed.Users[0].Name)
	}

	file, err := os.Create("test.csv")
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

	// fmt.Println(reflect.ValueOf(dataProjectSearch[0]).Kind())

	for _, i := range dataProjectSearch {
		row := getStructFieldValues(i)
		writer.Write(row)
	}

	fmt.Println("CSV file generated successfully!")
}
