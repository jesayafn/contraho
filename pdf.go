package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jung-kurt/gofpdf/v2"
)

const (
	defaultHeightCell = 6
)

func generatePDF(filename string, data interface{}, fields ...string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	// Adjust bottom margin
	// newBottomMargin := 10.0
	newBottomMargin := 25.4

	pdf.SetMargins(25.4, 25.4, 25.4)
	pdf.SetAutoPageBreak(true, newBottomMargin)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 11)

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
	leftMargin, topMargin, rightMargin, _ := pdf.GetMargins()
	// fmt.Println(bottomMargin, topMargin)
	// fmt.Printf("Current Margins - Left: %f, Top: %f, Right: %f, Bottom: %f\n", leftMargin, topMargin, rightMargin, bottomMargin)

	pageWidth, pageHeight := pdf.GetPageSize()

	usablePageWidth := pageWidth - (leftMargin + rightMargin)
	usablePageHeight := pageHeight - (topMargin)
	// fmt.Println(usablePageHeight, usablePageWidth, pageHeight, pageWidth)

	// Calculate cell widths based on number of fields
	cellWidth := usablePageWidth / float64(len(fields)) // Total width is 190mm

	// Table content
	err := printTable(pdf, usablePageHeight, cellWidth,
		val, fields, fieldIndices,
		filename, 25.4)

	return err

}

func printTable(pdf *gofpdf.Fpdf, usablePageHeight float64, cellWidth float64,
	val reflect.Value, fields []string, fieldIndices []int,
	filename string, margin float64) error {
	// Function to print table header
	printHeader := func() {
		for _, field := range fields {
			pdf.CellFormat(cellWidth, defaultHeightCell, strings.ToUpper(field), "TB", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Print table header initially
	printHeader()

	var previousY float64
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

			currentPage := pdf.PageNo()

			pdf.SetPage(currentPage - 1)
			currentPage = pdf.PageNo()

			pageWidth, _ := pdf.GetPageSize()
			pdf.Line(margin, previousY,
				pageWidth-margin, previousY)

			pdf.SetPage(currentPage + 1)

			// Print table header again on new page
			printHeader()
		}

		y := pdf.GetY()

		// fmt.Println(pdf.GetY()+maxHeight, usablePageHeight)
		// Print each cell
		for _, fieldIndex := range fieldIndices {
			fieldValue := elem.Field(fieldIndex)
			text := fmt.Sprintf("%v", fieldValue.Interface())

			// fmt.Println(text)
			// fmt.Println(heights[j])
			x := pdf.GetX()

			// Add padding to the cell
			padding := 1.0

			// Calculate the Y position of the next row
			// nextRowY := y + maxHeight

			// Draw bottom edge of rectangle for the last row or if it will exceed the page height
			isLastRow := (i == val.Len()-1)
			// exceedsPage := (pdf.GetY()+maxHeight > usablePageHeight-5)

			// if isLastRow || exceedsPage {
			if isLastRow {
				pdf.Line(x, y+maxHeight, x+cellWidth, y+maxHeight)
			}
			// Calculate the position for text to ensure padding
			textX := x + padding
			textY := y + padding

			// Print the text with proper alignment
			pdf.SetXY(textX, textY)
			pdf.MultiCell(cellWidth, defaultHeightCell, text, "", "LT", false)
			// pdf.SetXY(x+cellWidth, y)

			// Move to the next cell position
			pdf.SetXY(x+cellWidth, y)
		}
		previousY = y + maxHeight
		pdf.Ln(maxHeight)
	}

	// Save the PDF to a file
	return pdf.OutputFileAndClose(filename)
}
