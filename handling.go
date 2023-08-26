package main

var (
	pageSize = 500
)

func indexPageNumberCounter(getProjectSearchPage ProjectSearchPage) int {
	pages := getProjectSearchPage.Paging.Total / pageSize
	if getProjectSearchPage.Paging.Total%pageSize > 0 {
		pages++
	}
	if pages < 1 {
		pages = 1
	}
	return pages

}
