package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func httpRequest(url string, credential string, mode int) (data []byte, statusCode int) {

	// t := http.DefaultTransport.(*http.Transport).Clone()
	// t.MaxIdleConns = 100
	// t.MaxConnsPerHost = 100 // Corrected field name
	// t.MaxIdleConnsPerHost = 100
	// t.IdleConnTimeout = 45 * time.Second
	// t.TLSHandshakeTimeout = 45 * time.Second
	client := &http.Client{
		// Timeout:   60 * time.Second,
		// Transport: t,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	switch mode {
	case 1:
		req.Header.Add("Authorization", credential)
	case 0:
		req.SetBasicAuth(credential, "")
	}
	// fmt.Println(url)
	// req.Header.Add("ContentType", headerContentType)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	// defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	statusCode = resp.StatusCode

	if err != nil {
		fmt.Println("Error reading HTTP response body:", err)
		return
	}
	// fmt.Println(string(data))

	// data = string(body)
	// fmt.Println(resp.StatusCode)
	return data, statusCode

}

func projectSearchApi(host string, qualifiers string, size int, pageNumber int, credential string, mode int) (data []byte) {
	const (
		projectsSearch = "/api/projects/search"
	)
	queryParams := url.Values{}
	queryParams.Add("qualifiers", qualifiers)
	queryParams.Add("ps", fmt.Sprintf("%d", size))
	queryParams.Add("p", fmt.Sprintf("%d", pageNumber))
	encodedQuery := queryParams.Encode()
	fullPath := host + projectsSearch + "?" + encodedQuery

	data, _ = httpRequest(fullPath, credential, mode)
	return data
}

func projectBranchesListApi(host string, projectKey string, credential string, mode int) (data []byte) {
	const (
		projectBranchesList = "/api/project_branches/list"
	)
	queryParams := url.Values{}
	queryParams.Add("project", projectKey)
	encodedQuery := queryParams.Encode()
	fullPath := host + projectBranchesList + "?" + encodedQuery

	data, _ = httpRequest(fullPath, credential, mode)

	return data
}

func measuresComponentApi(host string, projectKey string, branch string, metricKeys string, credential string, mode int) (data []byte) {
	const (
		measuresComponent = "/api/measures/component"
	)
	queryParams := url.Values{}
	queryParams.Add("metricKeys", metricKeys)
	queryParams.Add("component", projectKey)
	queryParams.Add("branch", branch)
	encodedQuery := queryParams.Encode()
	fullPath := host + measuresComponent + "?" + encodedQuery
	// fmt.Println(fullPath)

	data, _ = httpRequest(fullPath, credential, mode)

	return data

}

func projectAnalysesSearchApi(host string, size int, pageNumber int, projectKey string, branch string, credential string, mode int) (data []byte) {
	const (
		projectAnalysesSearch = "/api/project_analyses/search"
	)
	queryParams := url.Values{}
	queryParams.Add("ps", fmt.Sprintf("%d", size))
	queryParams.Add("p", fmt.Sprintf("%d", pageNumber))
	queryParams.Add("branch", branch)
	queryParams.Add("project", projectKey)
	encodedQuery := queryParams.Encode()
	fullPath := host + projectAnalysesSearch + "?" + encodedQuery

	data, _ = httpRequest(fullPath, credential, mode)

	return data

}

func permissionUsersApi(host string, projectKey string, credential string, mode int) (data []byte) {
	const (
		permissionUsers = "/api/permissions/users"
	)

	queryParams := url.Values{}
	queryParams.Add("projectKey", projectKey)

	encodedQuery := queryParams.Encode()
	fullPath := host + permissionUsers + "?" + encodedQuery
	// fmt.Println(fullPath)
	data, _ = httpRequest(fullPath, credential, mode)

	return data

}

func applicationsSearchApi(host string, size int, pageNumber int, applicationKey string, credential string, mode int) (data []byte) {
	const (
		applicationsSearch = "/api/applications/search_projects"
	)
	queryParams := url.Values{}
	queryParams.Add("application", applicationKey)
	queryParams.Add("ps", fmt.Sprintf("%d", size))
	queryParams.Add("p", fmt.Sprintf("%d", pageNumber))
	encodedQuery := queryParams.Encode()
	fullPath := host + applicationsSearch + "?" + encodedQuery
	data, _ = httpRequest(fullPath, credential, mode)
	return data

}

func navigationGlobalApi(host string, credential string, mode int) (data []byte) {
	const (
		navigationGlobal = "/api/navigation/global"
	)
	fullPath := host + navigationGlobal
	// fmt.Println(credential)
	data, _ = httpRequest(fullPath, credential, mode)
	return data
}

func qualityGatesGetByProjectApi(host string, projectKey string, credential string, mode int) (data []byte) {
	const (
		qualityGatesGetByProject = "/api/qualitygates/get_by_project"
	)

	queryParams := url.Values{}
	queryParams.Add("project", projectKey)
	encodedQuery := queryParams.Encode()
	fullPath := host + qualityGatesGetByProject + "?" + encodedQuery
	data, _ = httpRequest(fullPath, credential, mode)
	return data
}

func navigationComponentApi(host string, projectKey string, branch string, credential string, mode int) (data []byte) {
	const (
		navigationComponent = "/api/navigation/component"
	)

	queryParams := url.Values{}
	queryParams.Set("component", projectKey)
	queryParams.Set("branch", branch)
	encodedQuery := queryParams.Encode()
	fullPath := host + navigationComponent + "?" + encodedQuery
	data, _ = httpRequest(fullPath, credential, mode)
	return data
}

func qualityProfilesShowApi(host string, qualityProfileKey string, credential string, mode int) (data []byte) {
	const (
		qualityProfilesShow = "/api/qualityprofiles/show"
	)

	queryParams := url.Values{}
	queryParams.Set("key", qualityProfileKey)
	encodedQuery := queryParams.Encode()
	fullPath := host + qualityProfilesShow + "?" + encodedQuery
	data, _ = httpRequest(fullPath, credential, mode)
	return data
}

func measuresSearchHistoryApi(host string, projectKey string, branch string, metrics string, dateFrom string, dateTo string,
	credential string, mode int) (data []byte) {
	const (
		measuresSearchHistory = "/api/measures/search_history"
	)
	queryParams := url.Values{}
	queryParams.Set("component", projectKey)
	queryParams.Set("branch", branch)
	queryParams.Set("metrics", metrics)
	queryParams.Set("from", dateFrom)
	queryParams.Set("to", dateTo)
	encodedQuery := queryParams.Encode()
	fullPath := host + measuresSearchHistory + "?" + encodedQuery
	data, _ = httpRequest(fullPath, credential, mode)
	return data
}

func applicationsShowApi(host string, applicationKey string, branch string, credential string, mode int) (data []byte, statusCode int) {
	const (
		measuresSearchHistory = "/api/applications/show"
	)
	queryParams := url.Values{}
	queryParams.Set("application", applicationKey)
	queryParams.Set("branch", branch)
	encodedQuery := queryParams.Encode()
	fullPath := host + measuresSearchHistory + "?" + encodedQuery
	// fmt.Println(fullPath)
	data, statusCode = httpRequest(fullPath, credential, mode)
	return data, statusCode
}
