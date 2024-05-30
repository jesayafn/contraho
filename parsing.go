package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func arguments(subcommand int) (host *string,
	credential string, authMode int,
	fileOutput *string, pagingOutput *bool,
	additionalOptions map[string]interface{}) {
	flagSet := flag.NewFlagSet("project", flag.ExitOnError)
	host = flagSet.String("host", "localhost", "Host of Sonarqube server. It is can be FQDN, or IP address")
	username := flagSet.String("username", "", "Username will be used for authentication to Sonarqube server")
	password := flagSet.String("password", "", "Password will be used for authentication to Sonarqube server")
	token := flagSet.String("token", "", "Token will be used for authentication to Sonarqube server")
	fileOutput = flagSet.String("filename", "", "CSV filename will be used for CSV output file")
	pagingOutput = flagSet.Bool("paging", false, "Pagination output using pager (only available on Linux and macOS)")
	additionalOptions = make(map[string]interface{})

	switch subcommand {
	case 0:

		unlistedApp := flagSet.Bool("unlisted-on-app", false, "List only not listed projects on any application")
		listedApp := flagSet.Bool("listed-on-app", false, "List only listed projects on any application")
		app := flagSet.String("app", "", "List only listed on the specified application")
		flagSet.Parse(os.Args[2:])
		if *pagingOutput && runtime.GOOS == "windows" {
			fmt.Println("Error: --paging is not supported on Windows")
			os.Exit(1)
		}
		if *username != "" && *password != "" && *token == "" {
			authMode = 1
		} else {
			authMode = 0
		}
		switch authMode {
		case 0:
			credential = *token
		case 1:
			credential = authorizationHeader(*username, *password)
		}
		// fmt.Println(*username, *password, credential, authMode)

		raw := navigationGlobalApi(*host, credential, authMode)
		// raw := []byte(`{"canAdmin":true,"globalPages":[],"settings":{"sonar.lf.enableGravatar":"false","sonar.developerAggregatedInfo.disabled":"false","sonar.lf.gravatarServerUrl":"https://secure.gravatar.com/avatar/{EMAIL_MD5}.jpg?s\u003d{SIZE}\u0026d\u003didenticon","sonar.technicalDebt.ratingGrid":"0.05,0.1,0.2,0.5","sonar.updatecenter.activate":"false"},"qualifiers":["TRK"],"version":"9.5 (build 56709)","productionDatabase":true,"branchesEnabled":false,"instanceUsesDefaultAdminCredentials":false,"multipleAlmEnabled":false,"projectImportFeatureEnabled":false,"regulatoryReportFeatureEnabled":false,"edition":"community","needIssueSync":false,"standalone":true}`)

		var sonarqubeInfo NavigationGlobal
		err := dataParse(raw, &sonarqubeInfo)

		handleErr(err)

		if (*unlistedApp || *listedApp || *app != "") && sonarqubeInfo.Edition == "community" {
			fmt.Println("Error: --unlisted-on-app or --listed-on-app cannot be used on Sonarqube Community Edition.")
			os.Exit(1)
		}

		if *app != "" {
			appArray := strings.Split(*app, ",")
			var notFound []string
			for index := range appArray {
				_, checkStatusCode := applicationsShowApi(*host, appArray[index], "", credential, authMode)
				if checkStatusCode == 404 {
					notFound = append(notFound, appArray[index])
				}
			}

			if len(notFound) >= 1 {
				fmt.Printf("Application not found: %v\nPlease check the requested application key(s). \n", strings.Join(notFound, ", "))
				os.Exit(1)
			}
		}

		if *unlistedApp && *listedApp {
			fmt.Println("Error: --unlisted-on-app and --listed-on-app cannot be used simultaneously.")
			os.Exit(1)
		}

		additionalOptions["unlistedApp"] = *unlistedApp
		additionalOptions["listedApp"] = *listedApp
		additionalOptions["app"] = *app
		// fmt.Println(*unlistedApp)
	case 1:
		flagSet.Parse(os.Args[2:])
		if *pagingOutput && runtime.GOOS == "windows" {
			fmt.Println("Error: --paging is not supported on Windows")
			os.Exit(1)
		}
		if *username != "" && *password != "" && *token == "" {
			authMode = 1
		} else {
			authMode = 0
		}
		switch authMode {
		case 0:
			credential = *token
		case 1:
			credential = authorizationHeader(*username, *password)
		}
		// fmt.Println(*username, *password, credential, authMode)

		raw := navigationGlobalApi(*host, credential, authMode)
		// raw := []byte(`{"canAdmin":true,"globalPages":[],"settings":{"sonar.lf.enableGravatar":"false","sonar.developerAggregatedInfo.disabled":"false","sonar.lf.gravatarServerUrl":"https://secure.gravatar.com/avatar/{EMAIL_MD5}.jpg?s\u003d{SIZE}\u0026d\u003didenticon","sonar.technicalDebt.ratingGrid":"0.05,0.1,0.2,0.5","sonar.updatecenter.activate":"false"},"qualifiers":["TRK"],"version":"9.5 (build 56709)","productionDatabase":true,"branchesEnabled":false,"instanceUsesDefaultAdminCredentials":false,"multipleAlmEnabled":false,"projectImportFeatureEnabled":false,"regulatoryReportFeatureEnabled":false,"edition":"community","needIssueSync":false,"standalone":true}`)

		var sonarqubeInfo NavigationGlobal
		err := dataParse(raw, &sonarqubeInfo)

		handleErr(err)
		if sonarqubeInfo.Edition == "community" {
			fmt.Println("Error: Unavailable on Sonarqube Community. ")
			os.Exit(1)
		}
	default:
		fmt.Println("tesuto")
	}

	return host, credential, authMode, fileOutput, pagingOutput, additionalOptions

}

func createCSVFile(fileOutput string, startTime time.Time, data interface{}) {
	// Open the CSV file
	file, err := os.Create(fileOutput)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Use reflection to get field names and values
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		fmt.Println("Input data is not a slice")
		return
	}

	// Write the CSV header (field names)
	if val.Len() > 0 {
		header := getStructFieldNames(val.Index(0).Interface())
		writer.Write(header)
	}

	// Write the data rows
	for i := 0; i < val.Len(); i++ {
		row := getStructFieldValues(val.Index(i).Interface())
		writer.Write(row)
	}

	fmt.Println("CSV file generated successfully!")
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Seconds()

	fmt.Printf("Execution Time: %.3f seconds\n", elapsedTime)

}

func authorizationHeader(username string, password string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	headerAuthValue := "Basic " + encoded

	return headerAuthValue
}

func dataParse(data []byte, v interface{}) error {

	return json.Unmarshal(data, v)

}

func findIndexOfHighestValue(numbers []int) int {
	if len(numbers) == 0 {
		return -1
	}

	maxIndex := 0
	maxValue := numbers[0]

	for i, num := range numbers {
		if num > maxValue {
			maxValue = num
			maxIndex = i
		}
	}

	return maxIndex
}

func getStructFieldNames(v interface{}) []string {
	var fields []string
	value := reflect.ValueOf(v)

	// Make sure the input is a struct
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			fields = append(fields, value.Type().Field(i).Name)
		}
	}
	return fields
}

func getStructFieldValues(v interface{}) []string {
	var values []string
	value := reflect.ValueOf(v)
	// fmt.Println(reflect.ValueOf(v), reflect.Struct)

	// Make sure the input is a struct
	// if value.Kind() == reflect.Struct {
	for i := 0; i < value.NumField(); i++ {
		// for i := range value.Nu {
		fieldValue := value.Field(i)
		// Handle numeric fields as strings to preserve leading zeros
		if fieldValue.Kind() == reflect.Int {
			values = append(values, strconv.Itoa(int(fieldValue.Int())))
		} else {
			values = append(values, fmt.Sprintf("%v", fieldValue.Interface()))
		}
	}
	return values
}

func removeRedundantValues(arr []string) []string {
	// Create a map to store unique values
	uniqueValues := make(map[string]bool)

	// Create a new array to store non-redundant values
	uniqueArray := make([]string, 0)

	// Iterate through the original array
	for _, val := range arr {
		// If the value is not in the map, add it to the new array and mark it as seen
		if !uniqueValues[val] {
			uniqueArray = append(uniqueArray, val)
			uniqueValues[val] = true
		}
	}

	return uniqueArray
}

func deleteProjectsByKeys(projects []ProjectSearchList, keysToDelete []string) []ProjectSearchList {
	var updatedProjects []ProjectSearchList

	// Create a map for faster lookup of keys to delete
	keysToDeleteMap := make(map[string]bool)
	for _, key := range keysToDelete {
		keysToDeleteMap[key] = true
	}

	// Iterate through the original projects and keep only those not in the keysToDeleteMap
	for _, project := range projects {
		if !keysToDeleteMap[project.Key] {
			updatedProjects = append(updatedProjects, project)
		}
	}

	return updatedProjects
}

func keepProjectsByKeys(projects []ProjectSearchList, keysToKeep []string) []ProjectSearchList {
	var updatedProjects []ProjectSearchList

	// Create a map for faster lookup of keys to keep
	keysToKeepMap := make(map[string]bool)
	for _, key := range keysToKeep {
		keysToKeepMap[key] = true
	}

	// Iterate through the original projects and keep only those in the keysToKeepMap
	for _, project := range projects {
		if keysToKeepMap[project.Key] {
			updatedProjects = append(updatedProjects, project)
		}
	}

	return updatedProjects
}
