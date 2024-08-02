package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"log/slog"
	"os"
	"sync"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"

	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	winTitle             = "Image Viewer"
	magicMarker          = "AW"
	magicMarkerSeekDepth = 20
)

var (
	screenBuffer screen.Buffer
	pixBuffer    *image.RGBA
	inFile       string

	deNib chan byte = make(chan byte, 512)

	waitForFile sync.WaitGroup
	doneFile    bool
)

func main() {
	args := os.Args
	if len(args) == 2 {
		inFile = args[1]
	} else {
		slog.Info("Usage: viewer <inputfile>")
		os.Exit(1)
	}

	// Open up the pcap file for reading
	handle, err := pcap.OpenOffline(inFile)
	if err != nil {
		slog.Error("Error opening pcap file", "error", err)
		os.Exit(1)
	}
	defer handle.Close()

	imgBytes := new(bytes.Buffer)
	doneFile = false

	// Fire up a routine to build the content file from the pcap
	// This will check for the magic marker first, then 8 bytes (64 bit int) for the
	// embedded file length, then the embedded file itself
	go func() {
		waitForFile.Add(1)
		markerFound := false
		lengthCnt := 0
		fileLength := 0
		bytesWritten := 0
		b := make([]byte, 8)
		curMarkerSearch := ""
		for {
			highNib := <-deNib
			lowNib := <-deNib

			Byte := highNib<<4 | lowNib
			if !markerFound {
				curMarkerSearch += string(Byte)
				if len(curMarkerSearch) > magicMarkerSeekDepth {
					slog.Error("Magic marker not found")
					os.Exit(1)
				}
				// Found the magic marker complete
				if curMarkerSearch == magicMarker {
					markerFound = true
				}
			} else {
				// Fetch the file length
				if lengthCnt < 9 {
					b[lengthCnt] = Byte
					lengthCnt++
					if lengthCnt == 8 {
						fileLength = int(binary.LittleEndian.Uint64(b))
						lengthCnt = 10
					}
				} else {
					// Write the embedded file to the output file
					// until we reach the file length
					if bytesWritten < fileLength {
						imgBytes.WriteByte(Byte)
						bytesWritten++
					} else if bytesWritten == fileLength {
						waitForFile.Done()
						doneFile = true
						// Make some space in deNib queue to avoid deadlock
						if len(deNib) > 0 {
							<-deNib
						}
						return
					}
				}
			}
		}
	}()

	// Read the packets from the pcap file
	// and extract the data from the TCP headers in nibbles.
	// Then send those nibbles to the deNib channel
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if doneFile {
			break
		}
		tcpxLayer, hasTCP := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
		if hasTCP {
			tmpByte := 0
			if tcpxLayer.ECE {
				tmpByte |= 1
			}
			if tcpxLayer.CWR {
				tmpByte |= 2
			}
			if tcpxLayer.URG {
				tmpByte |= 4
			}
			if tcpxLayer.PSH {
				tmpByte |= 8
			}
			deNib <- byte(tmpByte)
		}
	}
	doneCheckCnt := 0
	for !doneFile {
		time.Sleep(time.Millisecond * 100)
		doneCheckCnt++
		if doneCheckCnt > 100 {
			slog.Error("Error reading pcap file")
			os.Exit(1)
		}
	}

	waitForFile.Wait()

	if imgBytes.Len() == 0 {
		slog.Error("No embedded data found")
		os.Exit(1)
	}

	imageRGB, _, err := image.Decode(bytes.NewReader(imgBytes.Bytes()))
	if err != nil {
		slog.Error("Error decoding image", "error", err)
		os.Exit(1)
	}

	//Basically we now have the image. We could write it to disk, but we want to display it.
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
		var currTime time.Time
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
