package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jung-kurt/gofpdf/v2"
)

func generatePDF(filename string, data interface{}, fields ...string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 8)

	val := reflect.ValueOf(data)

	// Check if the data is a slice
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data should be a slice of structs")
	}

	// Get the type of the slice elements
	if val.Len() == 0 {
		return fmt.Errorf("data slice is empty")
	}

	elemType := val.Index(0).Type()

	// Verify that the fields exist in the struct
	fieldIndices := make([]int, len(fields))
	for i, field := range fields {
		structField, found := elemType.FieldByName(field)
		if !found {
			return fmt.Errorf("field %s not found in struct", field)
		}
		fieldIndices[i] = structField.Index[0]
	}

	// Calculate the usable page height

	leftMargin, topMargin, rightMargin, bottomMargin := pdf.GetMargins()
	_, pageHeight := pdf.GetPageSize()
	usablePageWidth := 210 - leftMargin - rightMargin
	usablePageHeight := pageHeight - bottomMargin - topMargin

	// Calculate cell widths based on number of fields
	cellWidth := usablePageWidth / float64(len(fields)) // Total width is 190mm
	const (
		defaultHeightCell = 5
		// lineHeightMultiplier = 0.7
	)

	// Function to print table header
	printHeader := func() {
		for _, field := range fields {
			pdf.CellFormat(cellWidth, defaultHeightCell, strings.ToUpper(field), "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Print table header initially
	printHeader()

	// Table content
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		heights := make([]float64, len(fields))

		// Measure the height of each cell
		for j, fieldIndex := range fieldIndices {
			fieldValue := elem.Field(fieldIndex)
			text := fmt.Sprintf("%v", fieldValue.Interface())
			lines := pdf.SplitLines([]byte(text), cellWidth)
			heights[j] = float64(len(lines)) * float64(defaultHeightCell)
		}

		// Determine the maximum height of the row
		maxHeight := 0.0
		for _, h := range heights {
			if h > maxHeight {
				maxHeight = h
			}
		}

		// Check if the row fits on the current page, otherwise add a new page
		if pdf.GetY()+maxHeight > usablePageHeight {
			pdf.AddPage()
			// Print table header again on new page
			printHeader()
		}

		y := pdf.GetY()

		// Print each cell
		for _, fieldIndex := range fieldIndices {
			fieldValue := elem.Field(fieldIndex)
			text := fmt.Sprintf("%v", fieldValue.Interface())

			// fmt.Println(text)
			// fmt.Println(heights[j])
			x := pdf.GetX()

			// Add padding to the cell
			padding := 1.0
			pdf.Rect(x, y, cellWidth, maxHeight, "D")

			// Calculate the position for text to ensure padding
			textX := x + padding
			textY := y + padding

			// Print the text with proper alignment
			pdf.SetXY(textX, textY)
			pdf.MultiCell(cellWidth, defaultHeightCell, text, "0", "LT", false)
			pdf.SetXY(x+cellWidth, y)

			// Move to the next cell position
			pdf.SetXY(x+cellWidth, y)
		}

		pdf.Ln(maxHeight)
	}

	// Save the PDF to a file
	return pdf.OutputFileAndClose(filename)
}
