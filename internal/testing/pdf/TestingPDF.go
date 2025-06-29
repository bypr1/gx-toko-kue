package pdf

import (
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"log"
	"service/internal/pkg/constant"
	"service/internal/pkg/core"
)

type TestingPDF struct {
	Name string
}

func (pdf TestingPDF) Generate() string {
	xpdf := core.PDF{}
	xpdf.NewGenerator("testing_pdf.html", pdf)

	xpdf.PDFG.Dpi.Set(300)
	xpdf.PDFG.PageSize.Set("A4")
	xpdf.PDFG.Orientation.Set(wkhtml.OrientationPortrait)

	path := constant.PathPDFTesting()
	filename := "hasil.pdf"

	err := xpdf.Save(path, filename)
	if err != nil {
		log.Panicf("Unable to generate PDF: %s", err)
	}

	return path + filename
}
