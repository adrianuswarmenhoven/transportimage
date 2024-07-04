/*
This program reads an image file and converts it into an HTML table.
Usage: <program> <input image file> <output HTML file>
It uses the 'no news is good news' philosophy for error handling. So if there is nothing in your terminal, the program ran successfully.

Pro:
It is more or less the simplest way to convert an image into an HTML table.
Each pixel in the image is represented by a cell in the table. The color of the cell is the color of the pixel in the image.
The recipient does not need any special software to view the image. A web browser is enough.

Con:
The huge size increase. A 40 KB image can easily become a 4 MB HTML file.
This is because each pixel is represented by a cell in the table.
The table is a 2D grid of cells, and each cell is a pixel.
So, a 100x100 pixel image will be represented by a 100x100 table, which is 10,000 cells.

File compression helps, but the size increase is still significant.
*/
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log/slog"
	"os"

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
			cellVal := uint32(r)<<16 | uint32(g)<<8 | uint32(b)
			cellLoc, _ := excelize.CoordinatesToCellName(x, y)
			f.SetCellValue("Data", cellLoc, cellVal)
		}
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs(outFile); err != nil {
		slog.Error("Error saving output", "error", err)
		os.Exit(1)
	}

}
