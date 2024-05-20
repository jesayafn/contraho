package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func printStructTable(data interface{}, startTime time.Time, selectedColumns ...string) {
	slice := reflect.ValueOf(data)

	if err := validateInput(slice); err != nil {
		fmt.Println(err)
		return
	}

	columnWidths := calculateColumnWidths(slice, selectedColumns)
	// fmt.Println(columnWidths)

	// printHeader(slice, columnWidths, selectedColumns)
	// printValues(slice, columnWidths, selectedColumns)
	var pagerCmd *exec.Cmd
	switch runtime.GOOS {
	// case "windows":
	// 	pagerCmd = exec.Command("more")
	case "darwin", "linux":
		pagerCmd = exec.Command("less")
	default:
		pagerCmd = nil
	}

	if pagerCmd != nil {
		pagerIn, err := pagerCmd.StdinPipe()
		if err == nil {
			pagerCmd.Stdout = os.Stdout
			pagerCmd.Stderr = os.Stderr
			if err := pagerCmd.Start(); err == nil {
				// Print the header and data to the pager
				printHeader(slice, columnWidths, selectedColumns, pagerIn.(io.Writer))
				printValues(slice, columnWidths, selectedColumns, pagerIn.(io.Writer))
				endTime := time.Now()
				elapsedTime := endTime.Sub(startTime).Seconds()

				fmt.Fprintf(pagerIn.(io.Writer), "Execution Time: %.3f seconds\n", elapsedTime)

				pagerIn.Close()
				pagerCmd.Wait()

				return
			}
		}
	}

	// Fallback to standard output if pager is not available
	printHeader(slice, columnWidths, selectedColumns, os.Stdout)
	printValues(slice, columnWidths, selectedColumns, os.Stdout)
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

func printHeader(slice reflect.Value, columnWidths []int, selectedColumns []string, writer io.Writer) {
	for fieldIndex := 0; fieldIndex < slice.Index(0).Type().NumField(); fieldIndex++ {
		fieldName := slice.Index(0).Type().Field(fieldIndex).Name
		if len(selectedColumns) == 0 || contains(selectedColumns, fieldName) {
			// fmt.Printf("%-*s", columnWidths[fieldIndex]+2, strings.ToUpper(fieldName))
			capitalizedFieldName := strings.ToUpper(fieldName)
			fmt.Fprintf(writer, "%-*s", columnWidths[fieldIndex]+2, capitalizedFieldName)
		}
	}
	fmt.Fprintln(writer)
}

func printValues(slice reflect.Value, columnWidths []int, selectedColumns []string, writer io.Writer) {
	for rowIndex := 0; rowIndex < slice.Len(); rowIndex++ {
		for fieldIndex := 0; fieldIndex < slice.Index(0).Type().NumField(); fieldIndex++ {
			fieldName := slice.Index(0).Type().Field(fieldIndex).Name
			if len(selectedColumns) == 0 || contains(selectedColumns, fieldName) {
				cellValue := fmt.Sprintf("%v", slice.Index(rowIndex).Field(fieldIndex).Interface())
				fmt.Fprintf(writer, "%-*s", columnWidths[fieldIndex]+2, cellValue)
			}
		}
		fmt.Fprintln(writer)
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
