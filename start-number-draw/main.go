package main

import (
	"flag"
	"log"

	"github.com/jung-kurt/gofpdf"
)

const (
	pageWidth         = 200
	pageHeight        = 140.7
	fontHelvetica     = "Helvetica"
	fontHelveticaBold = "Helvetica-Bold"
)

var (
	number   string
	name     string
	team     string
	fileName string
)

func setFont(pdf *gofpdf.Fpdf, family, style string, fontSize float64) (lineHeight float64) {
	pdf.SetFont(family, style, fontSize)
	return pdf.PointConvert(fontSize)
}

func main() {

	flag.StringVar(&number, "number", "", "")
	flag.StringVar(&name, "name", "", "")
	flag.StringVar(&team, "team", "", "")
	flag.StringVar(&fileName, "o", "out.pdf", "Output filename")
	flag.Parse()

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:    "mm",
		Size:       gofpdf.SizeType{Wd: pageWidth, Ht: pageHeight},
		FontDirStr: "/Users/ivan/go/src/github.com/jung-kurt/gofpdf/font",
	})

	pdf.SetMargins(0, 0, 0)
	pdf.AddFont(fontHelvetica, "", "helvetica_1251.json")
	pdf.AddFont(fontHelveticaBold, "", "helveticab.json")
	pdf.AddPage()

	translator := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	if len(number) != 0 {
		pdf.SetTextColor(0, 0, 0)
		lineHeight := setFont(pdf, fontHelveticaBold, "", 220)
		pdf.SetY((pageHeight-lineHeight)/2 + 5)
		pdf.MultiCell(pageWidth, lineHeight, translator(number), "", "C", false)
	}

	// hmargin := 5.0

	if len(name) != 0 {
		pdf.SetTextColor(0xDD, 0x54, 0x24)
		lineHeight := setFont(pdf, fontHelvetica, "", 40)
		pdf.SetY(26)
		//pdf.SetX(hmargin)
		pdf.MultiCell(pageWidth, lineHeight, translator(name), "", "C", false)
		//pdf.MultiCell(pageWidth-hmargin, lineHeight, translator(*namePtr), "", "L", false)
	}

	// hmargin = 3

	if len(team) != 0 {
		lineHeight := setFont(pdf, fontHelvetica, "", 28)
		pdf.SetY(pageHeight - 38)
		// pdf.SetX(hmargin)
		pdf.MultiCell(pageWidth, lineHeight, translator(team), "", "C", false)
	}

	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		log.Fatalln(err)
	}
}
