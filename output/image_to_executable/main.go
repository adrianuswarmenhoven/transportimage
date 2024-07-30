package main

import (
	"image"
	"log/slog"
	"os"
	"time"

	"golang.org/x/mobile/event/lifecycle"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/size"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

var (
	sizeEvent size.Event

	screenSize   = image.Point{winWidth, winHeight}
	screenBuffer screen.Buffer
	pixBuffer    *image.RGBA
	s            screen.Screen
)

func main() {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  winTitle,
			Width:  winWidth,
			Height: winHeight,
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
			drawScene()
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
