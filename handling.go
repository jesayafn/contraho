package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

const (
	pageSize = 500
)

var (
	loadingCh = make(chan bool)
)

func indexPageNumberCounter(getProjectSearchPage interface{}) int {
	switch page := getProjectSearchPage.(type) {
	case ProjectSearchPage:
		pages := page.Paging.Total / pageSize
		if page.Paging.Total%pageSize > 0 {
			pages++
		}
		if pages < 1 {
			pages = 1
		}
		return pages
	case ProjectSearchOfApplicationPage:
		pages := page.Paging.Total / pageSize
		if page.Paging.Total%pageSize > 0 {
			pages++
		}
		if pages < 1 {
			pages = 1
		}
		return pages
	default:
		return 0
	}
	// pages := getProjectSearchPage.Paging.Total / pageSize
	// if getProjectSearchPage.Paging.Total%pageSize > 0 {
	// 	pages++
	// }
	// if pages < 1 {
	// 	pages = 1
	// }
	// return pages
}

func removeByKeys(list []ProjectSearchList, keysToRemove []string) []ProjectSearchList {
	// Create a map to store the keys that need to be removed for efficient lookup.
	keySet := make(map[string]bool)
	for _, key := range keysToRemove {
		keySet[key] = true
	}

	// Create a new list to store the updated items.
	var updatedList []ProjectSearchList

	// Iterate through the original list.
	for _, item := range list {
		// Check if the item's key is not in the set of keys to remove.
		if !keySet[item.Key] {
			// If not, add the item to the updated list.
			updatedList = append(updatedList, item)
		}
	}

	return updatedList
}

func projectSearchApiLength(host string, credential string, projectType string, authMode int) int {
	// data := httpRequest(host+projectIndexApi, credential)
	// var data []byte

	data := projectSearchApi(host, projectType, 1, 1, credential, authMode)

	// fmt.Println(string(data))
	var projectSearchPage ProjectSearchPage
	err := dataParse(data, &projectSearchPage)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(projectSearchPage)
	indexPageNumber := indexPageNumberCounter(projectSearchPage)

	return indexPageNumber
}

func listProject(host string, credential string, lengthProject int, authMode int) []ProjectSearchList {

	displayJob("list project", "start")
	go displayLoading(loadingCh)
	var projectList []ProjectSearchList
	for pageIndex := 1; pageIndex <= lengthProject; pageIndex++ {
		// raw := httpRequest(host+projectScrapeApi+strconv.Itoa(pageIndex),
		// 	credential)
		raw := projectSearchApi(host, "TRK", 500, pageIndex, credential, authMode)

		var structured ProjectSearchPage

		err := dataParse(raw, &structured)
		// fmt.Println(structured)
		if err != nil {
			fmt.Println(err)
		}

		for projectIndex := range structured.Components {
			projectList = append(projectList, ProjectSearchList{
				Key:  structured.Components[projectIndex].Key,
				Name: structured.Components[projectIndex].Name,
			})
		}
	}
	loadingCh <- true
	displayJob("list project", "end")
	return projectList

}

func findIndexOfLatestDate(dateStrings []string) (int, error) {
	// Common layout for parsing date strings.
	layout := "2006-01-02T15:04:05-0700"

	// Initialize variables to keep track of the latest date and its index.
	var latestDateIndex int
	var latestDate time.Time

	for i, dateStr := range dateStrings {
		parsedDate, err := time.Parse(layout, dateStr)
		if err != nil {
			return -1, err // Return -1 for index and the error.
		}

		// If it's the first valid date or later than the current latestDate, update latestDate and its index.
		if i == 0 || parsedDate.After(latestDate) {
			latestDate = parsedDate
			latestDateIndex = i
		}
	}

	return latestDateIndex, nil
}
func branchDetailOfProjects(projectList []ProjectSearchList, host string, credential string, authMode int) []ProjectSearchList {
	// fmt.Println("Gather Branch Detail")

	displayJob("obtain branch data", "start")
	go displayLoading(loadingCh)

	for index := range projectList {
		// raw := httpRequest(host+projectBranchesApi+projectList[index].Key, credential)
		raw := projectBranchesListApi(host, projectList[index].Key, credential, authMode)
		var structured ProjectBranchesList
		err := dataParse(raw, &structured)

		handleErr(err)
		var (
			compareNloc     []int
			loc             int
			lastDate        string
			compareLastDate []string
		)
		for branchIndex := range structured.Branches {
			// nlocRaw := httpRequest(
			// 	host+ProjectBranchesLocApi+projectList[index].Key+"&branch="+structured.Branches[branchIndex].Name,
			// 	credential)
			nlocRaw := measuresComponentApi(host, projectList[index].Key,
				structured.Branches[branchIndex].Name, "ncloc", credential, authMode)

			var nlocStructured ProjectMeasures

			err := dataParse(nlocRaw, &nlocStructured)
			handleErr(err)

			if len(nlocStructured.Component.Measures) == 0 {
				loc = 0
			} else {
				loc, err = strconv.Atoi(nlocStructured.Component.Measures[0].Value)
				handleErr(err)
			}
			compareNloc = append(compareNloc, loc)

			// lastDateRaw := httpRequest(host+ProjectDateAnalysisApi+projectList[index].Key+"&branch="+structured.Branches[branchIndex].Name,
			// 	credential)
			lastDateRaw := projectAnalysesSearchApi(host, 1, 1, projectList[index].Key,
				structured.Branches[branchIndex].Name, credential, authMode)

			var lastDateStructured ProjectAnalyses
			err = dataParse(lastDateRaw, &lastDateStructured)
			handleErr(err)

			if len(lastDateStructured.Analyses) == 0 {
				lastDate = "0001-01-01T00:00:00+0000"
			} else {
				lastDate = lastDateStructured.Analyses[0].Date
			}

			compareLastDate = append(compareLastDate, lastDate)
		}
		// fmt.Println(projectList[index].Key, compareLastDate)
		branchCalculatedNloc := findIndexOfHighestValue(compareNloc)
		lastAnalysisDate, err := findIndexOfLatestDate(compareLastDate)
		handleErr(err)
		// projectList[index] = ProjectSearchList{
		// 	HighestLinesOfCodeBranch: structured.Branches[branchCalculatedNloc].Name,
		// 	LinesOfCode:              strconv.Itoa(compareNloc[branchCalculatedNloc]),
		// 	LastAnalysisDate:         compareLastDate[lastAnalysisDate],
		// 	LastAnalysisBranch:       structured.Branches[lastAnalysisDate].Name,
		// }
		projectList[index].Branch = structured.Branches[branchCalculatedNloc].Name
		projectList[index].Loc = strconv.Itoa(compareNloc[branchCalculatedNloc])
		projectList[index].LastAnalysisDate = compareLastDate[lastAnalysisDate]
		projectList[index].LastAnalysisBranch = structured.Branches[lastAnalysisDate].Name

	}
	loadingCh <- true
	displayJob("obtain branch data", "end")

	return projectList
}

func ownerProject(projectList []ProjectSearchList, host string, credential string, authMode int) []ProjectSearchList {
	// fmt.Println("Owner func")

	displayJob("obtain project owner", "start")
	go displayLoading(loadingCh)

	for index := range projectList {
		// raw := httpRequest(host+ProjectUserPermissionsApi+projectList[index].Key, credential)
		raw := permissionUsersApi(host, projectList[index].Key, credential, authMode)
		var structured ProjectPermissions

		err := dataParse(raw, &structured)
		// fmt.Println(structured)
		handleErr(err)

		// projectList[index] = ProjectSearchList{
		// 	Owner: structured.Users[0].Name,
		// 	Email: structured.Users[0].Email,
		// }
		projectList[index].Owner = structured.Users[0].Name
		projectList[index].Email = structured.Users[0].Email

	}
	loadingCh <- true
	displayJob("obtain project owner", "end")
	return projectList
}
func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Just hold array of the project key of application
func listApp(host string, credential string, lengthProject int, authMode int) []string {

	displayJob("list app", "start")
	go displayLoading(loadingCh)

	var applicationList []string
	for pageIndex := 1; pageIndex <= lengthProject; pageIndex++ {
		// raw := httpRequest(host+projectScrapeApi+strconv.Itoa(pageIndex),
		// 	credential)
		raw := projectSearchApi(host, "APP", 500, pageIndex, credential, authMode)

		var structured ProjectSearchPage

		err := dataParse(raw, &structured)
		// fmt.Println(structured)
		if err != nil {
			fmt.Println(err)
		}

		for projectIndex := range structured.Components {
			applicationList = append(applicationList, structured.Components[projectIndex].Key)
		}
	}
	loadingCh <- true
	displayJob("list app", "end")
	return applicationList

}

func applicationProjectSearchApiLength(host string, applicationKey string, credential string, authMode int) int {
	// data := httpRequest(host+projectIndexApi, credential)
	// var data []byte

	data := applicationsSearchApi(host, 1, 1, applicationKey, credential, authMode)

	// fmt.Println(string(data))
	var ProjectSearchOfApplicationPage ProjectSearchOfApplicationPage
	err := dataParse(data, &ProjectSearchOfApplicationPage)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(projectSearchPage)
	indexPageNumber := indexPageNumberCounter(ProjectSearchOfApplicationPage)

	return indexPageNumber
}

func listProjectofApplication(host string, projectListed []string, applicationKey string, lengthPage int, credential string, authMode int) []string {
	// var projectListed []string

	for pageIndex := 1; pageIndex <= lengthPage; pageIndex++ {
		data := applicationsSearchApi(host, 500, pageIndex, applicationKey, credential, authMode)

		var projectSearchOfApplicationPage ProjectSearchOfApplicationPage
		err := dataParse(data, &projectSearchOfApplicationPage)

		if err != nil {
			fmt.Println(err)
		}
		for _, project := range projectSearchOfApplicationPage.Projects {
			projectListed = append(projectListed, project.Key)
		}
	}
	return projectListed

}

func projectFiltering(projectList []ProjectSearchList, host string, credential string, option int, authMode int) []ProjectSearchList {
	lengthAppPage := projectSearchApiLength(host, credential, "APP", authMode)

	applicationList := listApp(host, credential, lengthAppPage, authMode)
	var lisedProjectOnApp []string
	for i := range applicationList {
		lengthProjectOfAppPage := applicationProjectSearchApiLength(host, applicationList[i], credential, authMode)
		// fmt.Println(lengthProjectOfAppPage)
		lisedProjectOnApp = listProjectofApplication(host, lisedProjectOnApp, applicationList[i],
			lengthProjectOfAppPage, credential, authMode)
		// fmt.Println(lisedProjectOnApp)
	}

	lisedProjectOnApp = removeRedundantValues(lisedProjectOnApp)

	// fmt.Println(lisedProjectOnApp)
	switch option {
	case 0:
		projectList = deleteProjectsByKeys(projectList, lisedProjectOnApp)
	case 1:
		projectList = keepProjectsByKeys(projectList, lisedProjectOnApp)
	}

	return projectList
}

func qualityGateofProject(projectList []ProjectSearchList, host string, credential string, authMode int) []ProjectSearchList {
	// fmt.Println("Gather Branch Detail")

	displayJob("obtain quality gate data", "start")
	go displayLoading(loadingCh)

	for index := range projectList {
		// raw := httpRequest(host+projectBranchesApi+projectList[index].Key, credential)
		raw := qualityGatesGetByProjectApi(host, projectList[index].Key, credential, authMode)
		var structured QualityGatesGetByProject
		err := dataParse(raw, &structured)

		handleErr(err)

		projectList[index].QualityGateId = structured.QualityGate.ID
		projectList[index].QualityGateName = structured.QualityGate.Name

	}
	displayJob("obtain quality gate data", "end")

	return projectList
}

// func printStructTable(data interface{}, columnsToShow ...string) {
// 	val := reflect.ValueOf(data)

// 	// Check if the input is an array or slice of structs
// 	if val.Kind() != reflect.Slice {
// 		fmt.Println("Input is not a slice")
// 		return
// 	}

// 	// Check if the slice is empty
// 	if val.Len() == 0 {
// 		fmt.Println("Slice is empty")
// 		return
// 	}

// 	// Determine column widths
// 	columnWidths := make([]int, val.Index(0).NumField())
// 	typ := val.Index(0).Type()

// 	// Find the maximum width for each column
// 	for i := 0; i < typ.NumField(); i++ {
// 		columnName := typ.Field(i).Name
// 		if len(columnsToShow) == 0 || contains(columnsToShow, columnName) {
// 			for j := 0; j < val.Len(); j++ {
// 				fieldValue := fmt.Sprintf("%v", val.Index(j).Field(i).Interface())
// 				fieldWidth := len(fieldValue)
// 				if fieldWidth > columnWidths[i] {
// 					columnWidths[i] = fieldWidth
// 				}
// 			}
// 		}
// 	}

// 	// Print header
// 	for i := 0; i < typ.NumField(); i++ {
// 		columnName := typ.Field(i).Name
// 		if len(columnsToShow) == 0 || contains(columnsToShow, columnName) {
// 			fmt.Printf("%-*s", columnWidths[i]+2, columnName)
// 		}
// 	}
// 	fmt.Println()

// 	// Print values
// 	for j := 0; j < val.Len(); j++ {
// 		for i := 0; i < typ.NumField(); i++ {
// 			columnName := typ.Field(i).Name
// 			if len(columnsToShow) == 0 || contains(columnsToShow, columnName) {
// 				fieldValue := fmt.Sprintf("%v", val.Index(j).Field(i).Interface())
// 				fmt.Printf("%-*s", columnWidths[i]+2, fieldValue)
// 			}
// 		}
// 		fmt.Println()
// 	}
// }

func printStructTable(data interface{}, selectedColumns ...string) {
	slice := reflect.ValueOf(data)

	if err := validateInput(slice); err != nil {
		fmt.Println(err)
		return
	}

	columnWidths := calculateColumnWidths(slice, selectedColumns)
	// fmt.Println(columnWidths)

	printHeader(slice, columnWidths, selectedColumns)
	printValues(slice, columnWidths, selectedColumns)
}

func validateInput(slice reflect.Value) error {
	if slice.Kind() != reflect.Slice {
		return errors.New("Input is not a slice")
	}

	if slice.Len() == 0 {
		return errors.New("Slice is empty")
	}

	return nil
}

func calculateColumnWidths(slice reflect.Value, selectedColumns []string) []int {
	columnWidths := make([]int, slice.Index(0).Type().NumField())

	for fieldIndex := 0; fieldIndex < slice.Index(0).Type().NumField(); fieldIndex++ {
		fieldName := slice.Index(0).Type().Field(fieldIndex).Name
		if len(selectedColumns) == 0 || contains(selectedColumns, fieldName) {
			updateColumnWidths(slice, fieldIndex, columnWidths)
		}
	}

	return columnWidths
}

func updateColumnWidths(slice reflect.Value, fieldIndex int, columnWidths []int) {
	for rowIndex := 0; rowIndex < slice.Len(); rowIndex++ {
		cellValue := fmt.Sprintf("%v", slice.Index(rowIndex).Field(fieldIndex).Interface())
		cellWidth := len(cellValue)
		if cellWidth > columnWidths[fieldIndex] {
			columnWidths[fieldIndex] = cellWidth
		}
	}
}

func printHeader(slice reflect.Value, columnWidths []int, selectedColumns []string) {
	for fieldIndex := 0; fieldIndex < slice.Index(0).Type().NumField(); fieldIndex++ {
		fieldName := slice.Index(0).Type().Field(fieldIndex).Name
		if len(selectedColumns) == 0 || contains(selectedColumns, fieldName) {
			fmt.Printf("%-*s", columnWidths[fieldIndex]+2, fieldName)
		}
	}
	fmt.Println()
}

func printValues(slice reflect.Value, columnWidths []int, selectedColumns []string) {
	for rowIndex := 0; rowIndex < slice.Len(); rowIndex++ {
		for fieldIndex := 0; fieldIndex < slice.Index(0).Type().NumField(); fieldIndex++ {
			fieldName := slice.Index(0).Type().Field(fieldIndex).Name
			if len(selectedColumns) == 0 || contains(selectedColumns, fieldName) {
				cellValue := fmt.Sprintf("%v", slice.Index(rowIndex).Field(fieldIndex).Interface())
				fmt.Printf("%-*s", columnWidths[fieldIndex]+2, cellValue)
			}
		}
		fmt.Println()
	}
}

func contains(slice []string, s string) bool {
	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}
