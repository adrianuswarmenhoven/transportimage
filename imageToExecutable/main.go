/*
imageToExecutable
*/
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	inFile  string
	workDir string = "output"
)

const (
	templateDir   = "template"
	imageCodeFile = "drawImage.go"
)

func main() {
	args := os.Args
	if len(args) == 3 {
		inFile = args[1]
		workDir = args[2]
	} else {
		slog.Info("Usage: imageToExecutable <inputfile> <outputdir>")
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

	os.Mkdir(workDir, os.FileMode(int(0775)))

	dir, err := os.ReadDir(workDir)
	if err != nil {
		slog.Error("Error reading working directory", "error", err)
		os.Exit(1)
	}
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{workDir, d.Name()}...))
	}

	dir, err = os.ReadDir(templateDir)
	if err != nil {
		slog.Error("Error reading template directory", "error", err)
		os.Exit(1)
	}
	for _, d := range dir {
		if strings.Contains(d.Name(), "thisFileIsNotUsed") {
			continue
		}
		var srcfd *os.File
		var dstfd *os.File
		var srcinfo os.FileInfo

		src := filepath.Join(templateDir, d.Name())
		dst := filepath.Join(workDir, d.Name())

		if srcfd, err = os.Open(src); err != nil {
			slog.Error("Error opening template file", "error", err)
			os.Exit(1)
		}

		if dstfd, err = os.Create(dst); err != nil {
			slog.Error("Error creating output file", "error", err)
			srcfd.Close()
			os.Exit(1)
		}

		if _, err = io.Copy(dstfd, srcfd); err != nil {
			slog.Error("Error copying template file", "file", d.Name(), "error", err)
			srcfd.Close()
			dstfd.Close()
			os.Exit(1)
		}
		if srcinfo, err = os.Stat(src); err != nil {
			slog.Error("Error getting file info", "file", d.Name(), "error", err)
			srcfd.Close()
			dstfd.Close()
			os.Exit(1)
		}
		err = os.Chmod(dst, srcinfo.Mode())
		if err != nil {
			slog.Error("Error setting file permissions", "file", d.Name(), "error", err)
			srcfd.Close()
			dstfd.Close()
			os.Exit(1)
		}
	}

	paramsData := fmt.Sprintf(templateParameter, width, height)
	err = os.WriteFile(path.Join(workDir, "params.go"), []byte(paramsData), os.FileMode(int(0664)))
	if err != nil {
		slog.Error("Error writing params file", "error", err)
		os.Exit(1)
	}

	pixelstring := strings.Builder{}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, a1 := img.At(x, y).RGBA()
			r := r1 / 256
			g := g1 / 256
			b := b1 / 256
			a := a1 / 256
			pixelstring.Write([]byte(fmt.Sprintf("pixBuffer.SetRGBA(%d, %d, color.RGBA{%d, %d, %d, %d})\n", x, y, r, g, b, a)))
		}
	}
	err = os.WriteFile(path.Join(workDir, imageCodeFile), []byte(fmt.Sprintf(templateDrawImage, pixelstring.String())), os.FileMode(int(0664)))
	if err != nil {
		slog.Error("Error writing drawImage file", "error", err)
		os.Exit(1)
	}

}

const (
	templateParameter = `
	package main

var (
	winTitle = "Awesome Image"

	winWidth, winHeight = %d, %d
)
`

	templateDrawImage = `
package main

import (
	"image/color"

)

func drawScene() {
 %s
}
	`
)
