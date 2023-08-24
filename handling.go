package main

func indexPageNumberCounter(getProjectSearchPage ProjectSearchPage) int {
	if getProjectSearchPage.Paging.Total%500 == 0 && getProjectSearchPage.Paging.Total > 500 {
		return getProjectSearchPage.Paging.Total / 500
		// fmt.Println("a")
	} else if getProjectSearchPage.Paging.Total%500 >= 0 && getProjectSearchPage.Paging.Total > 500 {
		return (getProjectSearchPage.Paging.Total / 500) + 1
		// fmt.Println("b")
	} else {
		return 1
		// fmt.Println("c")
	}
}
