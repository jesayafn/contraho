package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func httpRequest(url string, headerAuthValue string) (data []byte) {
	clien := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Add("Authorization", headerAuthValue)
	// req.Header.Add("ContentType", headerContentType)

	resp, err := clien.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	// defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading HTTP response body:", err)
		return
	}
	// fmt.Println(string(data))

	// data = string(body)
	// fmt.Println(resp.StatusCode)
	return data

}

func projectSearchApi(host string, qualifiers string, size int, pageNumber int, authToken string) (data []byte) {
	const (
		projectsSearch = "/api/projects/search"
	)
	queryParams := url.Values{}
	queryParams.Add("qualifiers", qualifiers)
	queryParams.Add("ps", fmt.Sprintf("%d", size))
	queryParams.Add("p", fmt.Sprintf("%d", pageNumber))
	encodedQuery := queryParams.Encode()
	fullPath := host + projectsSearch + "?" + encodedQuery

	data = httpRequest(fullPath, authToken)
	return data
}

func projectBranchesListApi(host string, projectKey string, authToken string) (data []byte) {
	const (
		projectBranchesList = "/api/project_branches/list"
	)
	queryParams := url.Values{}
	queryParams.Add("project", projectKey)
	encodedQuery := queryParams.Encode()
	fullPath := host + projectBranchesList + "?" + encodedQuery

	data = httpRequest(fullPath, authToken)

	return data
}

func measuresComponentApi(host string, projectKey string, branch string, metricKeys string, authToken string) (data []byte) {
	const (
		measuresComponent = "/api/measures/component"
	)
	queryParams := url.Values{}
	queryParams.Add("metricKeys", metricKeys)
	queryParams.Add("component", projectKey)
	queryParams.Add("branch", branch)
	encodedQuery := queryParams.Encode()
	fullPath := host + measuresComponent + "?" + encodedQuery

	data = httpRequest(fullPath, authToken)

	return data

}

func projectAnalysesSearchApi(host string, size int, pageNumber int, projectKey string, branch string, authToken string) (data []byte) {
	const (
		projectAnalysesSearch = "/api/project_analyses/search"
	)
	queryParams := url.Values{}
	queryParams.Add("ps", fmt.Sprintf("%d", size))
	queryParams.Add("p", fmt.Sprintf("%d", pageNumber))
	queryParams.Add("project", projectKey)
	encodedQuery := queryParams.Encode()
	fullPath := host + projectAnalysesSearch + "?" + encodedQuery

	data = httpRequest(fullPath, authToken)

	return data

}

func permissionUsersApi(host string, projectKey string, authToken string) (data []byte) {
	const (
		permissionUsers = "/api/permissions/users"
	)

	queryParams := url.Values{}
	queryParams.Add("projectKey", projectKey)

	encodedQuery := queryParams.Encode()
	fullPath := host + permissionUsers + "?" + encodedQuery
	// fmt.Println(fullPath)
	data = httpRequest(fullPath, authToken)

	return data

}

func applicationsSearchApi(host string, size int, pageNumber int, applicationKey string, authToken string) (data []byte) {
	const (
		applicationsSearch = "/api/applications/search_projects"
	)
	queryParams := url.Values{}
	queryParams.Add("application", applicationKey)
	queryParams.Add("ps", fmt.Sprintf("%d", size))
	queryParams.Add("p", fmt.Sprintf("%d", pageNumber))
	encodedQuery := queryParams.Encode()
	fullPath := host + applicationsSearch + "?" + encodedQuery
	data = httpRequest(fullPath, authToken)
	return data

}

func navigationGlobalApi(host string, authToken string) (data []byte) {
	const (
		navigationGlobal = "/api/navigation/global"
	)
	fullPath := host + navigationGlobal
	// fmt.Println(authToken)
	data = httpRequest(fullPath, authToken)
	return data
}

func qualityGatesGetByProjectApi(host string, projectKey string, authToken string) (data []byte) {
	const (
		qualityGatesGetByProject = "/api/qualitygates/get_by_project"
	)

	queryParams := url.Values{}
	queryParams.Add("project", projectKey)
	encodedQuery := queryParams.Encode()
	fullPath := host + qualityGatesGetByProject + "?" + encodedQuery
	data = httpRequest(fullPath, authToken)
	return data
}
