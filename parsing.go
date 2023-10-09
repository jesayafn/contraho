package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func arguments() (host *string, username *string, password *string, fileOutput *string) {
	flagSet := flag.NewFlagSet("project", flag.ExitOnError)
	host = flagSet.String("host", "localhost", "Host of Sonarqube server. It is can be FQDN, or IP address")
	username = flagSet.String("username", "admin", "Username will be used for authentication to Sonarqube server")
	password = flagSet.String("password", "admin", "Password will be used for authentication to Sonarqube server")
	fileOutput = flagSet.String("filename", "contraho.csv", "CSV filename will be used for CSV output file")
	// flag.Parse()
	flagSet.Parse(os.Args[2:])
	return host, username, password, fileOutput

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
