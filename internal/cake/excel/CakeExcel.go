package excel

import (
	"fmt"
	"service/internal/pkg/model/cake"
	cakeparser "service/internal/pkg/parser/cake"

	"github.com/xuri/excelize/v2"
)

type CakeExcel struct {
	Data []cake.Cake
}

func (excel CakeExcel) Export() *excelize.File {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new worksheet
	index, err := f.NewSheet("Cakes")
	if err != nil {
		fmt.Println(err)
		return f
	}

	// Set headers
	headers := []string{"ID", "Name", "Description", "Margin", "Sell Price", "Unit", "Stock", "Created At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue("Cakes", cell, header)
	}

	// Parse data and populate rows
	parser := cakeparser.CakeParser{Array: excel.Data}
	parsedData := parser.Get()

	for i, item := range parsedData {
		rowNum := i + 2
		data := item.(map[string]interface{})

		f.SetCellValue("Cakes", fmt.Sprintf("A%d", rowNum), data["id"])
		f.SetCellValue("Cakes", fmt.Sprintf("B%d", rowNum), data["name"])
		f.SetCellValue("Cakes", fmt.Sprintf("C%d", rowNum), data["description"])
		f.SetCellValue("Cakes", fmt.Sprintf("D%d", rowNum), data["margin"])
		f.SetCellValue("Cakes", fmt.Sprintf("E%d", rowNum), data["sellPrice"])
		f.SetCellValue("Cakes", fmt.Sprintf("F%d", rowNum), data["unit"])
		f.SetCellValue("Cakes", fmt.Sprintf("G%d", rowNum), data["stock"])
		f.SetCellValue("Cakes", fmt.Sprintf("H%d", rowNum), data["createdAt"])
	}

	// Set active sheet
	f.SetActiveSheet(index)

	return f
}

type IngredientExcel struct {
	Data []cake.Ingredient
}

func (excel IngredientExcel) Export() *excelize.File {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new worksheet
	index, err := f.NewSheet("Ingredients")
	if err != nil {
		fmt.Println(err)
		return f
	}

	// Set headers
	headers := []string{"ID", "Name", "Description", "Unit Price", "Unit", "Created At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue("Ingredients", cell, header)
	}

	// Parse data and populate rows
	parser := cakeparser.IngredientParser{Array: excel.Data}
	parsedData := parser.Get()

	for i, item := range parsedData {
		rowNum := i + 2
		data := item.(map[string]interface{})

		f.SetCellValue("Ingredients", fmt.Sprintf("A%d", rowNum), data["id"])
		f.SetCellValue("Ingredients", fmt.Sprintf("B%d", rowNum), data["name"])
		f.SetCellValue("Ingredients", fmt.Sprintf("C%d", rowNum), data["description"])
		f.SetCellValue("Ingredients", fmt.Sprintf("D%d", rowNum), data["unitPrice"])
		f.SetCellValue("Ingredients", fmt.Sprintf("E%d", rowNum), data["unit"])
		f.SetCellValue("Ingredients", fmt.Sprintf("F%d", rowNum), data["createdAt"])
	}

	// Set active sheet
	f.SetActiveSheet(index)

	return f
}
