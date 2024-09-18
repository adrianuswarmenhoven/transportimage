/*
imageToSpreadsheet
*/
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

var (
	inFile  string
	outFile string
)

func main() {
	args := os.Args
	if len(args) == 3 {
		inFile = args[1]
		outFile = args[2]
	} else {
		slog.Info("Usage: imageToSpreadsheet <inputfile> <outputfile>")
		os.Exit(1)
	}
	imgfile, err := os.Open(inFile)
	if err != nil {
		slog.Error("Error opening input file", "error", err)
		os.Exit(1)
	}
	defer imgfile.Close()
	imgCfg, _, err := image.DecodeConfig(imgfile)
	if err != nil {
		slog.Error("Error decoding image config", "error", err)
		os.Exit(1)
	}
	width := imgCfg.Width
	height := imgCfg.Height
	imgfile.Seek(0, 0)
	img, _, err := image.Decode(imgfile)
	if err != nil {
		slog.Error("Error decoding image", "error", err)
		os.Exit(1)
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Data")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.SetActiveSheet(index)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r := r1 / 256
			g := g1 / 256
			b := b1 / 256
			cellVal := uint32(r)<<16 | uint32(g)<<8 | uint32(b) // CHECK MAX VALUE OF A CELL AND STACK VALUES
			cellLoc, _ := excelize.CoordinatesToCellName(x, y)
			f.SetCellValue("Data", cellLoc, cellVal)
		}
	}
	// Sav spreadsheet by the given path.
	if err := f.SaveAs(outFile); err != nil {
		slog.Error("Error saving output", "error", err)
		os.Exit(1)
	}

	// Create a macro to display the image in the spreadsheet
	macro := macroTemplate
	macro = strings.Replace(macro, "<$IMAGEFILE$>", inFile, -1)
	macro = strings.Replace(macro, "<$IMAGEWIDTH$>", strconv.Itoa(width), -1)
	macro = strings.Replace(macro, "<$IMAGEHEIGHT$>", strconv.Itoa(height), -1)

	err = os.WriteFile(strings.TrimSuffix(outFile, filepath.Ext(outFile))+".bas", []byte(macro), 0644)
	if err != nil {
		slog.Error("Error writing macro", "error", err)
		os.Exit(1)
	}
}

const (
	macroTemplate = `
	' To use this macro with the data for the image, you must add this code to the spreadsheet as a macro (might take some time to run).
	' The data for the image must be in a sheet named "Data".
	' The spreadsheet must have a sheet called 'View' to display the image. You can just create a new sheet and name it 'View'.
	'
	' This viewer is for the image file <$IMAGEFILE$>
Sub CreateCanvas()
	ScreenUpdatingOff()
	Dim oDrawPage As Object
	Dim oShape As Object
	Dim x As Integer, y As Integer
	Dim canvasWidth As Integer, canvasHeight As Integer
	Dim cellSize As Integer
	Dim offSetX As Integer
	Dim offSetY As Integer

	Dim oSize As New com.sun.star.awt.Size
	Dim oPosition As New com.sun.star.awt.Point

	canvasWidth = <$IMAGEWIDTH$> ' Number of cells wide
	canvasHeight = <$IMAGEHEIGHT$> ' Number of cells high
	cellSize = 40 ' Size of each cell in pixels
	
	offSetX=1000
	offSetY=1000
	
	oDrawPage = ThisComponent.getSheets().getByName("View").getDrawPage()
	dataSheet = ThisComponent.getSheets().getByName("Data")
	
	' Loop through the canvas size
	For x = 0 To canvasWidth - 1
		For y = 0 To canvasHeight - 1
			' Calculate cell positions in pixels
			Dim xPos As Integer
			Dim yPos As Integer
			xPos = x * cellSize
			yPos = y * cellSize
			
			' Create a rectangle shape (simulating a pixel)
			oShape = ThisComponent.createInstance("com.sun.star.drawing.RectangleShape")

			oDrawPage.add(oShape)

			oSize.Width = cellSize
			oSize.Height = cellSize

			oPosition.X = offSetX+(cellSize * x)
			oPosition.Y = offSetY+(cellSize * y)
			
    		Dim cellValue As Long
   			 Dim red As Long, green As Long, blue As Long
			
			cell = dataSheet.getCellByPosition(x,y)
            cellValue = cell.Value
                      ' Ensure the cell value is a valid 32-bit integer
            If IsNumeric(cellValue) Then
                cellValue = CLng(cellValue)
                
                ' Extract the RGB components from the 32-bit integer
                red = (cellValue And 16711680) \ &H10000
                if red>255 then
                   red = 255
                endif
                
                green = (cellValue And 65280) \ &H100
                if green>255 then
                   green=255
                endif
                
                blue = cellValue And 255
                if blue >255 then
                   blue=255
                endif
               
            End If
			With oShape
				.setSize(oSize)
				.setPosition(oPosition)
				.FillColor = RGB(red, green, blue) 
				.FillStyle = com.sun.star.drawing.FillStyle.SOLID ' or .NONE
				.FillTransparence = 0 ' 0-100 from opaque to transparent
				.LineStyle = com.sun.star.drawing.LineStyle.SOLID ' or .NONE
				.LineWidth = 1 
				.LineColor = RGB(red, green, blue)
			End With
		Next y
	Next x
	ScreenUpdatingOn()
End Sub

Private Sub ScreenUpdatingOff()
  
#If VBA6 = 0 Then
  ' Excel ignores all statements up to #End If
  ' Please advise the LO Calc code equivalent of
  ' MS Excel VBA "Application.ScreenUpdating = False" to go
  ' on the next line. Thank you

  Exit Sub
#End If

  ' MS Excel ScreenUpdating code
  Application.ScreenUpdating = False
End Sub


Private Sub ScreenUpdatingOn()
  
#If VBA6 = 0 Then
  ' Excel ignores all statements up to #End If
  ' Please advise the LO Calc code equivalent of
  ' MS Excel VBA "Application.ScreenUpdating = True" to go
  ' on the next line. Thank you

  Exit Sub
#End If

  ' MS Excel ScreenUpdating code
  Application.ScreenUpdating = True
End Sub

`
)
