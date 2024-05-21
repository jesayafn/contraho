package main

import (
	"fmt"
	"reflect"

	"github.com/jung-kurt/gofpdf/v2"
)

func generatePDF(filename string, data interface{}, fieldsToShow ...string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 8)

	val := reflect.ValueOf(data)

	// Check if the data is a slice
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data should be a slice of structs")
	}

	// Check if the slice is empty
	if val.Len() == 0 {
		return fmt.Errorf("data slice is empty")
	}

	// Get the type of the slice elements
	elemType := val.Index(0).Type()

	// Map to keep track of field indices
	fieldIndices := make(map[string]int)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		fieldIndices[field.Name] = i
	}

	// Table header with specified fields
	for _, fieldName := range fieldsToShow {
		if index, ok := fieldIndices[fieldName]; ok {
			field := elemType.Field(index)
			pdf.CellFormat(40, 10, field.Name, "1", 0, "C", false, 0, "")
		}
	}
	pdf.Ln(-1)

	// Table content with specified fields
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		for _, fieldName := range fieldsToShow {
			if index, ok := fieldIndices[fieldName]; ok {
				field := elem.Field(index)
				pdf.CellFormat(40, 10, fmt.Sprintf("%v", field.Interface()), "1", 0, "L", false, 0, "")
			}
		}
		pdf.Ln(-1)
	}

	// Save the PDF to a file
	return pdf.OutputFileAndClose(filename)
}
