package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

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

	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executableDir := filepath.Dir(executablePath)
	fontDirStr := filepath.Join(executableDir, "fonts")

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:    "mm",
		Size:       gofpdf.SizeType{Wd: pageWidth, Ht: pageHeight},
		FontDirStr: fontDirStr,
	})

	pdf.SetMargins(0, 0, 0)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddFont(fontHelvetica, "", "helvetica_1251.json")
	pdf.AddFont(fontHelveticaBold, "", "helveticab.json")
	pdf.AddPage()

	translator := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	vmarginTop := 22.0
	vmarginBottom := 50.0
	hmargin := 7.0

	if len(name) != 0 {
		// pdf.SetTextColor(0xC0, 0x00, 0x00)
		pdf.SetTextColor(0, 0, 0)
		lineHeight := setFont(pdf, fontHelvetica, "", 40)
		pdf.SetY(vmarginTop + 7)
		pdf.SetX(hmargin)
		// pdf.MultiCell(pageWidth, lineHeight, translator(name), "", "C", false)
		pdf.MultiCell(pageWidth-hmargin, lineHeight, translator(name), "", "L", false)
	}

	mul := 0.6
	withTeamHei := 276.0
	withNameHei := 324.0
	noTeamNameShortHei := 380.0
	noTeamNameLongHei := 320.0

	if len(number) != 0 {
		pdf.SetTextColor(0, 0, 0)
		if len(team) != 0 {
			setFont(pdf, fontHelveticaBold, "", withTeamHei*mul)
			pdf.SetY(vmarginTop + (pageHeight-vmarginBottom)/2 + 6)
			pdf.SetX(0)
			pdf.MultiCell(pageWidth, 0, translator(number), "", "C", false)
		} else if len(name) != 0 {
			setFont(pdf, fontHelveticaBold, "", withNameHei*mul)
			pdf.SetY(vmarginTop + (pageHeight-vmarginBottom)/2 + 14)
			pdf.SetX(0)
			pdf.MultiCell(pageWidth, 0, translator(number), "", "C", false)
		} else {
			if len(number) <= 2 {
				setFont(pdf, fontHelveticaBold, "", noTeamNameShortHei*mul)
				pdf.SetY(vmarginTop + (pageHeight-vmarginBottom)/2 + 6.5)
				pdf.SetX(0)
				pdf.MultiCell(pageWidth, 0, translator(number), "", "C", false)
			} else {
				setFont(pdf, fontHelveticaBold, "", noTeamNameLongHei*mul)
				pdf.SetY(vmarginTop + (pageHeight-vmarginBottom)/2 + 6.5)
				pdf.SetX(0)
				pdf.MultiCell(pageWidth, 0, translator(number), "", "C", false)
			}
		}
	}

	if len(team) != 0 {
		pdf.SetTextColor(0, 0, 0)
		// pdf.SetTextColor(0xC0, 0x00, 0x00)
		setFont(pdf, fontHelvetica, "", 32)
		pdf.SetY(pageHeight - vmarginBottom + 11)
		pdf.SetX(hmargin)
		// pdf.MultiCell(pageWidth, 0, translator(team), "", "C", false)
		pdf.MultiCell(pageWidth-hmargin*2, 0, translator(team), "", "R", false)
	}

	err = pdf.OutputFileAndClose(fileName)
	if err != nil {
		log.Fatalln(err)
	}
}
