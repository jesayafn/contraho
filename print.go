package main

import (
	"errors"
	"fmt"
	"reflect"
)

func printStructTable(data interface{}, selectedColumns ...string) {
	slice := reflect.ValueOf(data)

	if err := validateInput(slice); err != nil {
		fmt.Println(err)
		return
	}

	columnWidths := calculateColumnWidths(slice, selectedColumns)
	// fmt.Println(columnWidths)

	printHeader(slice, columnWidths, selectedColumns)
	printValues(slice, columnWidths, selectedColumns)
}

func validateInput(slice reflect.Value) error {
	if slice.Kind() != reflect.Slice {
		return errors.New("Input is not a slice")
	}

	if slice.Len() == 0 {
		return errors.New("Slice is empty")
	}

	return nil
}

func calculateColumnWidths(slice reflect.Value, selectedColumns []string) []int {
	columnWidths := make([]int, slice.Index(0).Type().NumField())

	for fieldIndex := 0; fieldIndex < slice.Index(0).Type().NumField(); fieldIndex++ {
		fieldName := slice.Index(0).Type().Field(fieldIndex).Name
		if len(selectedColumns) == 0 || contains(selectedColumns, fieldName) {
			updateColumnWidths(slice, fieldIndex, columnWidths)
		}
	}

	return columnWidths
}

func updateColumnWidths(slice reflect.Value, fieldIndex int, columnWidths []int) {
	for rowIndex := 0; rowIndex < slice.Len(); rowIndex++ {
		cellValue := fmt.Sprintf("%v", slice.Index(rowIndex).Field(fieldIndex).Interface())
		cellWidth := len(cellValue)
		if cellWidth > columnWidths[fieldIndex] {
			columnWidths[fieldIndex] = cellWidth
		}
	}
}

func printHeader(slice reflect.Value, columnWidths []int, selectedColumns []string) {
	for fieldIndex := 0; fieldIndex < slice.Index(0).Type().NumField(); fieldIndex++ {
		fieldName := slice.Index(0).Type().Field(fieldIndex).Name
		if len(selectedColumns) == 0 || contains(selectedColumns, fieldName) {
			fmt.Printf("%-*s", columnWidths[fieldIndex]+2, fieldName)
		}
	}
	fmt.Println()
}

func printValues(slice reflect.Value, columnWidths []int, selectedColumns []string) {
	for rowIndex := 0; rowIndex < slice.Len(); rowIndex++ {
		for fieldIndex := 0; fieldIndex < slice.Index(0).Type().NumField(); fieldIndex++ {
			fieldName := slice.Index(0).Type().Field(fieldIndex).Name
			if len(selectedColumns) == 0 || contains(selectedColumns, fieldName) {
				cellValue := fmt.Sprintf("%v", slice.Index(rowIndex).Field(fieldIndex).Interface())
				fmt.Printf("%-*s", columnWidths[fieldIndex]+2, cellValue)
			}
		}
		fmt.Println()
	}
}

func contains(slice []string, s string) bool {
	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}
