package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

func createCSVFile(csvOutput string, startTime time.Time, data interface{}) {
	// Open the CSV file
	file, err := os.Create(csvOutput)
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
