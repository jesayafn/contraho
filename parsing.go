package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"reflect"
	"strconv"
)

func arguments() (host *string, username *string, password *string) {
	host = flag.String("host", "localhost", "Host of Sonarqube server. It is can be FQDN, or IP address")
	username = flag.String("username", "admin", "Username will be used for authentication to Sonarqube server")
	password = flag.String("password", "admin", "Password will be used for authentication to Sonarqube server")

	flag.Parse()
	return host, username, password

}

func authorizationHeader(username string, password string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	headerAuthValue := "Basic " + encoded
	// fmt.Println(headerAuthValue)``

	return headerAuthValue
}

// func projectSearchParse(projectSearchData []byte) (projectSearchPage ProjectSearchPage) {
// 	// func responseParse(rawData []byte) (any) {
// 	json.Unmarshal(projectSearchData, &projectSearchPage)
// 	return projectSearchPage
// 	// json.Unmarshal(rawData, &any)
// }

func dataParse(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func findIndexOfHighestValue(numbers []int) int {
	if len(numbers) == 0 {
		// Return an appropriate default value or error handling
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

// func getStructFieldValues(v interface{}) []string {
// 	var values []string
// 	value := reflect.ValueOf(v)

// 	// Make sure the input is a struct
// 	if value.Kind() == reflect.Struct {
// 		for i := 0; i 	< value.NumField(); i++ {
// 			values = append(values, fmt.Sprintf("%v", value.Field(i).Interface()))
// 		}
// 	}
// 	return values
// }

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
	// }
	return values
}

// func stringToTime(input string, layout string) (time.Time, error) {
// 	parsedTime, err := time.Parse(layout, input)
// 	if err != nil {
// 		return time.Time{}, err
// 	}
// 	return parsedTime, nil
// }

// func formatTime(parsedTime time.Time) string {
// 	return parsedTime.Format("2006-01-02 15:04:05")
// }
