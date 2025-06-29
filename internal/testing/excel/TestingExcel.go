package excel

import (
	xtremecore "github.com/globalxtreme/go-core/v2"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/xuri/excelize/v2"
	"service/internal/pkg/constant"
)

type TestingExcel struct {
}

func (ex TestingExcel) Generate() error {
	sheets, properties := ex.setSheetsAndProperties()

	excel := ex.newFile(sheets, properties)
	excel = ex.modifySheet(excel)

	err := excel.Save(constant.PathExcelTesting(), "testing.xlsx")
	if err != nil {
		return err
	}

	return nil
}

func (ex TestingExcel) newFile(sheets []string, properties [][][]interface{}) xtremecore.Excel {
	excel := xtremecore.Excel{
		Sheets:     sheets,
		Properties: properties,
		IsPublic:   false,
	}

	excel.NewFile()

	return excel
}

func (ex TestingExcel) modifySheet(excel xtremecore.Excel) xtremecore.Excel {
	excel.SetStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "860A35", Style: 1},
			{Type: "right", Color: "860A35", Style: 1},
			{Type: "top", Color: "860A35", Style: 1},
			{Type: "bottom", Color: "860A35", Style: 1},
		},
	}, "A1:C1")

	excel.SetStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "A3B763", Style: 2},
			{Type: "right", Color: "A3B763", Style: 2},
			{Type: "top", Color: "A3B763", Style: 2},
			{Type: "bottom", Color: "A3B763", Style: 2},
		},
	}, "A3:C4")

	excel.MergeCells("A5:C5", "A6:A7")

	excel.SetWidthCols([]xtremecore.ColWidth{
		{Cells: "A", Width: 5},
		{Cells: "B", Width: 15},
		{Cells: "C", Width: 5},
	})

	excel.SetHeightRows([]xtremecore.RowHeight{
		{Row: 1, Height: 15},
	})

	return excel
}

func (TestingExcel) setSheetsAndProperties() ([]string, [][][]interface{}) {
	var sheets []string
	var properties [][][]interface{}

	sheets = append(sheets, "Testing 1", "Testing 2")
	for _ = range sheets {
		properties = append(properties, [][]interface{}{
			{"ID", "Name", "Age"},
		})
	}

	dataProperties := [][]int{
		{29, 23, 12},
		{19, 26, 14},
	}

	for sKey, property := range dataProperties {
		for pKey, val := range property {
			properties[sKey] = append(properties[sKey], []interface{}{pKey + 1, xtremepkg.RandomString(10), val})
		}
	}

	return sheets, properties
}
