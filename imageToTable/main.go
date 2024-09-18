/*
imageToTable
*/
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log/slog"
	"os"
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
		slog.Info("Usage: imageToTable <inputfile> <outputfile>")
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

	htmlFile, err := os.Create(outFile)
	if err != nil {
		slog.Error("Error creating output file", "error", err)
		os.Exit(1)
	}
	defer htmlFile.Close()

	fmt.Fprint(htmlFile, htmlPreamble)

	for y := 0; y < height; y++ {
		fmt.Fprint(htmlFile, htmlTableRowStart)
		for x := 0; x < width; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r := r1 / 256
			g := g1 / 256
			b := b1 / 256
			fmt.Fprintf(htmlFile, htmlTableCellData, r, g, b)
		}
		fmt.Fprint(htmlFile, htmlTableRowEnd)
	}
	fmt.Fprint(htmlFile, htmlPostamble)

}

const (
	// Basically a simple HTML file
	htmlPreamble = `<html>
		<head>
		<title>An awesome image</title>		
		<style type="text/css">
		body {background-color: black;}

		.table{
			width:960px;
			border: 0px;
		}
		td{
			width: 1px;
			height: 1px;
			border: 0px;
		}
		div {
			height: 500px;
			-webkit-align-content: center;
			align-content: center;
		}
		</style>
		</head>
		<body backgroundcolour>
		<div style="height:80px;">&nbsp;</div>
		<div align="center">
		<table class="imgtable" border="0" cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; font-size:0px;">
			   `
	htmlTableRowStart = `<tr>`
	htmlTableRowEnd   = `</tr>`
	// Each cell is a pixel
	htmlTableCellData = `<td style="background-color:rgb(%d,%d,%d);"></td>`

	htmlPostamble = `</table><br/>
		</div>
		</body>
	</html>`
)
