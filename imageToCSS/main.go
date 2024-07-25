package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"strings"
)

const (

	definitionPlaceholder = "<$CSSDEFINITION$>"
	instancePlaceholder   = "<$CSSINSTANCES$>"
)

var (

	cssDefs         = make(map[color.Color]string)
	cssSpriteDefIDX = make(map[color.Color]int)
	cssInst         = make([]string, 0)

	inFile string
	outFile string
)

func main() {
	args := os.Args
	if len(args) == 3 {
		inFile = args[1]
		outFile = args[2]
	} else {
		slog.Info("Usage: imageToCSS <inputfile> <outputfile>")
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

	// Run through the image, get the colors and create the css definitions
	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixCol := img.At(x, y)
			r1, g1, b1, _ := pixCol.RGBA()
			r := r1 / 256
			g := g1 / 256
			b := b1 / 256
			if _, done := cssDefs[pixCol]; !done {
				cssDefs[pixCol] = fmt.Sprintf(cssSpriteDefinitionTemplate, idx, r, g, b)
				cssSpriteDefIDX[pixCol] = idx
				idx++
			}
			useIdx := cssSpriteDefIDX[pixCol]

			newline := ""
			if x == width-1 {
				newline = "<br/>\n"
			}
			cssInst = append(cssInst, fmt.Sprintf(cssSpriteInstanceTemplate, useIdx, newline))

		}
	}

	CSSDEFINITIONS := ""
	for _, v := range cssDefs {
		CSSDEFINITIONS += v
	}

	CSSINSTANCES := ""
	for _, v := range cssInst {
		CSSINSTANCES += v
	}

	html := htmlTemplate
	html = strings.Replace(html, definitionPlaceholder, CSSDEFINITIONS, -1)
	html = strings.Replace(html, instancePlaceholder, CSSINSTANCES, -1)

	fmt.Fprint(htmlFile, html)
}

const (
	// Basically a simple HTML file
	htmlTemplate = `<html>
		<head>
		<title>An awesome image</title>		
		<style type="text/css">
		body {background-color: black;}
		span {display: inline-block; width: 1px; height: 1px;}

		<$CSSDEFINITION$>

		</style>
		</head>
		<body backgroundcolour>
		<div style="height:80px;">&nbsp;</div>
		<div align="center">
		
		<$CSSINSTANCES$>

		<br/>
		</div>
		</body>
	</html>`

	cssSpriteDefinitionTemplate = `
	#c%d {  width: 1px;  height: 1px; background-color: rgb(%d, %d, %d); }`
	cssSpriteInstanceTemplate = `<span id="c%d">&nbsp;</span>%s`
)
