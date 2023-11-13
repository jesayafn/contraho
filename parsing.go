package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func arguments(subcommand string) (host *string, username *string, password *string, fileOutput *string, additionalOptions map[string]interface{}) {
	flagSet := flag.NewFlagSet("project", flag.ExitOnError)
	host = flagSet.String("host", "localhost", "Host of Sonarqube server. It is can be FQDN, or IP address")
	username = flagSet.String("username", "admin", "Username will be used for authentication to Sonarqube server")
	password = flagSet.String("password", "admin", "Password will be used for authentication to Sonarqube server")
	fileOutput = flagSet.String("filename", "contraho.csv", "CSV filename will be used for CSV output file")
	// flag.Parse()
	// flagSet.Parse(os.Args[2:])

	additionalOptions = make(map[string]interface{})

	switch subcommand {
	case "project":

		unlistedApp := flagSet.Bool("unlisted-on-app", false, "This is UOA option")
		listedApp := flagSet.Bool("listed-on-app", false, "This is LOA option")
		flagSet.Parse(os.Args[2:])

		if *unlistedApp && *listedApp {
			fmt.Println("Error: --unlisted-on-app and --listed-on-app cannot be used simultaneously.")
			os.Exit(1)
		}
		additionalOptions["unlistedApp"] = *unlistedApp
		additionalOptions["listedApp"] = *listedApp
		// fmt.Println(*unlistedApp)

	default:
		fmt.Println("tesuto")
	}

	return host, username, password, fileOutput, additionalOptions

}

func createCSVFile(fileOutput string, data interface{}) {
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
