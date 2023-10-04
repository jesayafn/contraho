package main

import (
	"fmt"
	"strconv"
)

var (
	pageSize = 500
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
	case ProjectSearchOfApplication:
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

func projectLength(host string, credential string) int {
	data := httpRequest(host+projectIndexApi, credential)
	var projectSearchPage ProjectSearchPage
	err := dataParse(data, projectSearchPage)
	if err != nil {
		fmt.Println(err)
	}
	indexPageNumber := indexPageNumberCounter(projectSearchPage)

	return indexPageNumber
}

func listProject(host string, credential string, lengthProject int) []ProjectSearchList {
	var projectList []ProjectSearchList
	for pageIndex := 1; pageIndex <= lengthProject; pageIndex++ {
		raw := httpRequest(host+projectScrapeApi+strconv.Itoa(pageIndex),
			credential)
		var structured ProjectSearchPage

		err := dataParse(raw, &structured)
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

	return projectList
}

func branchDetailOfProjects(host string, credential string, projectList []ProjectSearchList) {
	for index := range projectList {
		raw := httpRequest(host+projectBranchesApi+projectList[index].Key, credential)
		var structured ProjectBranchesList
		err := dataParse(raw, &structured)

		if err != nil {
			fmt.Println(err)
		}

	}
}
