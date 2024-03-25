package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
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
func branchDetailOfProjects(projectList interface{}, host string, credential string, authMode int) interface{} {
	// Type assertion to determine the type of projectList
	displayJob("obtain branch data", "start")
	go displayLoading(loadingCh)

	switch projectList := projectList.(type) {
	case []ProjectSearchList:
		// Handle ProjectSearchList type

		for index := range projectList {
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

			branchCalculatedNloc := findIndexOfHighestValue(compareNloc)
			lastAnalysisDate, err := findIndexOfLatestDate(compareLastDate)
			handleErr(err)

			projectList[index].Branch = structured.Branches[branchCalculatedNloc].Name
			projectList[index].Loc = strconv.Itoa(compareNloc[branchCalculatedNloc])
			projectList[index].LastAnalysisDate = compareLastDate[lastAnalysisDate]
			projectList[index].LastAnalysisBranch = structured.Branches[lastAnalysisDate].Name
		}
		loadingCh <- true
		displayJob("obtain branch data", "end")
		return projectList
	case []AppList:
		// Handle AppList type

		for index := range projectList {
			raw := projectBranchesListApi(host, projectList[index].Key, credential, authMode)
			var structured ProjectBranchesList
			err := dataParse(raw, &structured)
			handleErr(err)

			var (
				mainBranch string
				loc        int
			)

			for _, branch := range structured.Branches {
				if branch.IsMain {
					mainBranch = branch.Name
					break
				}
			}
			nlocRaw := measuresComponentApi(host, projectList[index].Key,
				mainBranch, "ncloc", credential, authMode)

			var nlocStructured ProjectMeasures

			err = dataParse(nlocRaw, &nlocStructured)
			handleErr(err)

			if len(nlocStructured.Component.Measures) == 0 {
				loc = 0
			} else {
				loc, err = strconv.Atoi(nlocStructured.Component.Measures[0].Value)
				handleErr(err)
			}
			projectList[index].MainBranch = mainBranch
			projectList[index].Loc = strconv.Itoa(loc)
		}
		loadingCh <- true
		displayJob("obtain branch data", "end")
		return projectList
	default:
		// Handle unsupported types
		panic("unsupported type for projectList")
	}
}

func ownerProject(projectList interface{}, host string, credential string, authMode int) interface{} {
	// fmt.Println("Owner func")

	displayJob("obtain project owner", "start")
	go displayLoading(loadingCh)
	sliceValue := reflect.ValueOf(projectList)
	if sliceValue.Kind() != reflect.Slice {
		panic("Input is not a slice")
	}
	for index := 0; index < sliceValue.Len(); index++ {
		// raw := httpRequest(host+ProjectUserPermissionsApi+projectList[index].Key, credential)
		element := sliceValue.Index(index)
		keyField := element.FieldByName("Key")
		if !keyField.IsValid() {
			panic("Key field not found")
		}
		key := keyField.Interface().(string)
		if !keyField.IsValid() {
			panic("Key field not found")
		}
		raw := permissionUsersApi(host, key, credential, authMode)
		var structured ProjectPermissions

		err := dataParse(raw, &structured)
		// fmt.Println(structured)
		handleErr(err)

		// projectList[index] = ProjectSearchList{
		// 	Owner: structured.Users[0].Name,
		// 	Email: structured.Users[0].Email,
		// }
		ownerField := element.FieldByName("Owner")
		emailField := element.FieldByName("Email")
		if !ownerField.IsValid() || !emailField.IsValid() {
			panic("Owner or Email field not found")
		}
		ownerField.SetString(structured.Users[0].Name)
		emailField.SetString(structured.Users[0].Email)

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
	displayJob("project filtering", "start")
	go displayLoading(loadingCh)

	lengthAppPage := projectSearchApiLength(host, credential, "APP", authMode)

	applicationListInterface := listApp(host, credential, lengthAppPage, authMode, 0)

	// Type assertion to convert the interface to []string
	applicationList, ok := applicationListInterface.([]string)
	if !ok {
		// Handle the case where the returned value is not []string
		fmt.Println("Error: Unexpected return type from listApp")
		return projectList
	}
	var listedProjectOnApp []string
	for i := range applicationList {
		lengthProjectOfAppPage := applicationProjectSearchApiLength(host, applicationList[i], credential, authMode)
		// fmt.Println(lengthProjectOfAppPage)
		listedProjectOnApp = listProjectofApplication(host, listedProjectOnApp, applicationList[i],
			lengthProjectOfAppPage, credential, authMode)
		// fmt.Println(lisedProjectOnApp)
	}

	listedProjectOnApp = removeRedundantValues(listedProjectOnApp)

	// fmt.Println(lisedProjectOnApp)
	switch option {
	case 0:
		projectList = deleteProjectsByKeys(projectList, listedProjectOnApp)
	case 1:
		projectList = keepProjectsByKeys(projectList, listedProjectOnApp)
	default:
	}

	displayJob("project filtering", "end")

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

func languageofProject(projectList []ProjectSearchList, host string, credential string, authMode int) []ProjectSearchList {
	displayJob("obtain language of project", "start")
	go displayLoading(loadingCh)
	const empty = `-`

	for index := range projectList {
		// raw := httpRequest(host+projectBranchesApi+projectList[index].Key, credential)
		raw := navigationComponentApi(host, projectList[index].Key, credential, authMode)
		var structured NavigationComponent
		err := dataParse(raw, &structured)
		handleErr(err)
		qualityProfilesKeys := make([]string, len(structured.QualityProfiles))
		for indexKey, profile := range structured.QualityProfiles {
			qualityProfilesKeys[indexKey] = profile.Key
		}
		// fmt.Println(qualityProfilesKey)

		if len(qualityProfilesKeys) != 0 {
			qualityProfilesNames := make([]string, len(qualityProfilesKeys)) //
			for indexKey, qualityProfileKey := range qualityProfilesKeys {
				raw := qualityProfilesShowApi(host, qualityProfileKey, credential, authMode)
				var structured QualityProfilesShow
				err := dataParse(raw, &structured)
				handleErr(err)
				qualityProfilesNames[indexKey] = structured.Profile.LanguageName
			}
			projectList[index].Language = strings.Join(qualityProfilesNames, ", ")
		} else {
			projectList[index].Language = empty
		}

	}
	displayJob("obtain language of project", "end")
	return projectList
}

func metricProject(projectList interface{}, host string, credential string, authMode int) interface{} {
	displayJob("obtain metric of projects", "start")
	go displayLoading(loadingCh)
	metrics := []string{
		"bugs", "security_hotspots",
		"line_coverage", "duplicated_lines",
		"code_smells", "sqale_index"}

	sliceValue := reflect.ValueOf(projectList)
	if sliceValue.Kind() != reflect.Slice {
		panic("Input is not a slice")
	}

	for index := 0; index < sliceValue.Len(); index++ {
		element := sliceValue.Index(index)
		keyField := element.FieldByName("Key")
		key := keyField.Interface().(string)

		var branch string
		if element.Type().Name() == "ProjectSearchList" {
			branchField := element.FieldByName("Branch")
			branch = branchField.Interface().(string)
		} else if element.Type().Name() == "AppList" {
			// If it's an AppList, use default branch "main"
			branchField := element.FieldByName("MainBranch")
			branch = branchField.Interface().(string)
		}

		var structured ProjectMeasures
		raw := measuresComponentApi(host, key, branch, strings.Join(metrics, ", "), credential, authMode)
		err := dataParse(raw, &structured)
		handleErr(err)

		for _, measure := range structured.Component.Measures {
			switch measure.Metric {
			case "bugs":
				element.FieldByName("Bugs").SetString(measure.Value)
			case "security_hotspots":
				element.FieldByName("SecurityHotspots").SetString(measure.Value)
			case "line_coverage":
				element.FieldByName("LineCoverage").SetString(measure.Value)
			case "duplicated_lines":
				element.FieldByName("DuplicatedLines").SetString(measure.Value)
			case "code_smells":
				element.FieldByName("CodeSmells").SetString(measure.Value)
			case "sqale_index":
				element.FieldByName("DebtInMinute").SetString(measure.Value)
			}
		}
	}

	displayJob("obtain metric of projects", "end")
	return projectList
}
