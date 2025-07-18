package excel

import (
	"fmt"
	"service/internal/pkg/constant"
	errorpkg "service/internal/pkg/error"
	"service/internal/pkg/model"
	"time"

	xtremecore "github.com/globalxtreme/go-core/v2"
	"github.com/xuri/excelize/v2"
)

type TransactionReport struct {
	model.Transaction
	TotalCakes int `gorm:"column:totalCakes;not null"`
}

type TransactionExcel struct {
	Transactions []TransactionReport
}

func (ex TransactionExcel) Generate() (string, error) {
	sheets, properties := ex.setSheetsAndProperties()

	excel := ex.newFile(sheets, properties)
	excel = ex.modifySheet(excel)

	filename := fmt.Sprintf("transactions_%s.xlsx", time.Now().Format("20060102_150405"))
	err := excel.Save(constant.PathExcelTransaction(), filename)
	if err != nil {
		errorpkg.ErrXtremeTransactionExcelGenerate(err.Error())
	}

	return constant.PathExcelTransaction() + filename, nil
}

func (ex TransactionExcel) newFile(sheets []string, properties [][][]interface{}) xtremecore.Excel {
	excel := xtremecore.Excel{
		Sheets:     sheets,
		Properties: properties,
		IsPublic:   true,
	}

	excel.NewFile()

	return excel
}

func (ex TransactionExcel) modifySheet(excel xtremecore.Excel) xtremecore.Excel {
	// Header style for first sheet (Transactions)
	excel.SetStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Family: "Arial",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	}, "A1:F1")

	// Data borders for first sheet
	dataRange := fmt.Sprintf("A2:F%d", len(ex.Transactions)+1)
	excel.SetStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	}, dataRange)

	// Set column widths for first sheet
	excel.SetWidthCols([]xtremecore.ColWidth{
		{Cells: "A", Width: 8},  // ID
		{Cells: "B", Width: 20}, // Transaction Date
		{Cells: "C", Width: 15}, // Total Amount
		{Cells: "D", Width: 20}, // Total Item
		{Cells: "E", Width: 20}, // Created At
		{Cells: "F", Width: 20}, // Updated At
	})

	// Set row heights for first sheet
	excel.SetHeightRows([]xtremecore.RowHeight{
		{Row: 1, Height: 25}, // Header row
	})

	return excel
}

func (ex TransactionExcel) setSheetsAndProperties() ([]string, [][][]interface{}) {
	var sheets []string
	var properties [][][]interface{}

	sheets = append(sheets, "Transactions")

	// Transaction Summary Sheet
	transactionHeaders := []interface{}{
		"ID",
		"Transaction Date",
		"Total Amount",
		"Total Item",
		"Created At",
		"Updated At",
	}

	transactionProperties := [][]interface{}{transactionHeaders}

	// Transaction data rows
	for _, transaction := range ex.Transactions {
		row := []interface{}{
			transaction.ID,
			transaction.Date.Format("02/01/2006"),
			transaction.TotalPrice,
			transaction.TotalCakes,
			transaction.CreatedAt.Format("02/01/2006 15:04"),
			transaction.UpdatedAt.Format("02/01/2006 15:04"),
		}
		transactionProperties = append(transactionProperties, row)
	}

	properties = append(properties, transactionProperties)

	return sheets, properties
}
