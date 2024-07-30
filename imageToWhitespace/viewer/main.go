package main

import (
	"bytes"
	"encoding/hex"
	"image"
	"log/slog"
	"os"
	"time"

	"golang.org/x/mobile/event/lifecycle"

	"golang.org/x/mobile/event/key"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"

	"image/color"
	_ "image/jpeg"
	_ "image/png"
)

const (
	winTitle = "Image to Whitespace"
)

var (
	startMarker = '\u2060' // WORD JOINER
	endMarker   = '\u180E' // MONGOLIAN VOWEL SEPARATOR

	unicodeToHex = map[rune]rune{
		'\u0020': '0', // SPACE
		'\u00A0': '1', // NO-BREAK SPACE
		'\u2000': '2', // EN QUAD
		'\u2001': '3', // EM QUAD
		'\u2002': '4', // EN SPACE
		'\u2003': '5', // EM SPACE
		'\u2004': '6', // THREE-PER-EM SPACE
		'\u2005': '7', // FOUR-PER-EM SPACE
		'\u2006': '8', // SIX-PER-EM SPACE
		'\u2007': '9', // FIGURE SPACE
		'\u2008': 'a', // PUNCTUATION SPACE
		'\u2009': 'b', // THIN SPACE
		'\u200A': 'c', // HAIR SPACE
		'\u202F': 'd', // NARROW NO-BREAK SPACE
		'\u205F': 'e', // MEDIUM MATHEMATICAL SPACE
		'\u3000': 'f', // IDEOGRAPHIC SPACE
	}
)

var (
	screenBuffer screen.Buffer
	pixBuffer    *image.RGBA
	inFile       string
)

func main() {
	args := os.Args
	if len(args) == 2 {
		inFile = args[1]
	} else {
		slog.Info("Usage: viewer <inputfile>")
		os.Exit(1)
	}

	data, err := os.ReadFile(inFile)
	if err != nil {
		slog.Error("Error reading file", "error", err)
		os.Exit(1)
	}
	hexdata := ""
	active := false
	for _, r := range string(data) {
		if r == startMarker {
			active = true
			continue
		}
		if r == endMarker {
			active = false
			break
		}
		if char, ok := unicodeToHex[r]; ok && active {
			hexdata += string(char)
		}
	}

	bindata, err := hex.DecodeString(hexdata)
	if err != nil {
		slog.Error("Error decoding hex data", "error", err)
		os.Exit(1)
	}
	imageRGB, _, err := image.Decode(bytes.NewReader(bindata))
	if err != nil {
		slog.Error("Error decoding image", "error", err)
		os.Exit(1)
	}
	iB := imageRGB.Bounds()

	screenSize := image.Point{iB.Dx(), iB.Dy()}

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  winTitle,
			Width:  iB.Dx(),
			Height: iB.Dy(),
		})
		if err != nil {
			slog.Error("Failed to create Window", "error", err)
			return
		}
		defer w.Release()

		screenBuffer, err = s.NewBuffer(screenSize)
		if err != nil {
			slog.Error("Failed to create screen buffer", "error", err)
			os.Exit(1)
		}
		defer screenBuffer.Release()
		pixBuffer = screenBuffer.RGBA()

		var frames = 0
		var startTime time.Time
		var currTime = time.Now()
		for {
			drawScene(imageRGB)
			w.Upload(image.Point{0, 0}, screenBuffer, screenBuffer.Bounds())
			w.Publish()
			time.Sleep(time.Millisecond * 16)

			// print out the ms/frame value
			frames++
			currTime = time.Now()
			if currTime.Sub(startTime).Seconds() >= 1.0 {
				slog.Debug("Rendering at %.5f ms/frame", 1000.0/float64(frames))
				frames = 0
				startTime = currTime
			}

			// Handle window events:
			switch e := w.NextEvent().(type) {

			case key.Event:
				if e.Code == key.CodeEscape {
					return // quit app
				}

			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					// Do any final cleanup or saving here:
					return // quit the application.
				}
			}
		}
	})
}

func drawScene(i image.Image) {
	bounds := i.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := i.At(x, y).RGBA()
			r /= 256
			g /= 256
			b /= 256
			a /= 256
			pixBuffer.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
}
