package main

import (
	"encoding/binary" // Import the fmt package to print messages to the console.
	"log"             // Import the log package to log errors to the console.
	"log/slog"
	"os"

	"github.com/google/gopacket"        // Import the gopacket package to decode packets.
	"github.com/google/gopacket/layers" // Import the layers package to access the various network layers.
	"github.com/google/gopacket/pcap"   // Import the pcap package to capture packets.
	"github.com/google/gopacket/pcapgo"
)

const (
	magicMarker = "AW"
)

var (
	embedData  chan byte = make(chan byte, 16384)
	inFile     string
	mediumFile string
	outFile    string
)

func main() {
	args := os.Args
	if len(args) == 4 {
		inFile = args[1]
		mediumFile = args[2]
		outFile = args[3]
	} else {
		slog.Info("Usage: imageToPCAP <inputfile> <medium> <outputfile>")
		os.Exit(1)
	}

	// Open up the pcap file for reading
	handle, err := pcap.OpenOffline(mediumFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	f, err := os.Create(outFile)
	if err != nil {
		slog.Error("Could not create output file", "error", err)
		os.Exit(1)
	}
	defer f.Close()

	// Write magic marker to the data to be embedded
	for _, b := range []byte(magicMarker) {
		lownib := b & 0x0F
		highnib := b >> 4
		embedData <- highnib
		embedData <- lownib
	}

	embBytes, err := os.ReadFile(inFile)
	if err != nil {
		slog.Error("Could not read input file", "error", err)
		os.Exit(1)
	}

	// After the magic marker, write the length of the image to be embedded
	lenBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(lenBuf, uint64(len(embBytes)))
	for _, b := range lenBuf {
		lownib := b & 0x0F
		highnib := b >> 4
		embedData <- highnib
		embedData <- lownib
	}

	// Feed the image to be embedded to the packet reader/writer
	// Every byte gets converted into two nibbles
	// Those nibbles are then embedded into the TCP flags
	go func() {
		for _, b := range embBytes {
			lownib := b & 0x0F
			highnib := b >> 4
			embedData <- highnib
			embedData <- lownib
		}
	}()

	w := pcapgo.NewWriter(f)
	err = w.WriteFileHeader(uint32(0), layers.LinkTypeEthernet)
	if err != nil {
		slog.Error("Could not write file header", "error", err)
		os.Exit(1)
	}
	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	storage := 0
	stored := 0
	for packet := range packetSource.Packets() {
		ci := packet.Metadata().CaptureInfo
		tcpModLayer, hasTCP := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
		if len(embedData) > 0 && hasTCP {
			embByte, ok := <-embedData
			if ok {
				tcpModLayer.ECE = ((embByte & 1) == 1)
				tcpModLayer.CWR = ((embByte & 2) == 2)
				tcpModLayer.URG = ((embByte & 4) == 4)
				tcpModLayer.PSH = ((embByte & 8) == 8)
				stored++
			}
			newLayers := []gopacket.SerializableLayer{}
			orgLayers := packet.Layers()

			// We need to get the IP layer to set the checksum
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer == nil {
				ipLayer = packet.Layer(layers.LayerTypeIPv6)
			}

			// We copy all the layers except for the modified TCP layer
			for _, copyLayer := range orgLayers {
				if copyLayer.LayerType() != layers.LayerTypeTCP {
					newLayers = append(newLayers, copyLayer.(gopacket.SerializableLayer))
				} else {
					err = tcpModLayer.SetNetworkLayerForChecksum(ipLayer.(gopacket.NetworkLayer))
					if err != nil {
						slog.Error("Could not set network layer for checksum", "error", err)
						os.Exit(1)
					}
					newLayers = append(newLayers, tcpModLayer)
				}
			}
			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				FixLengths:       true,
				ComputeChecksums: true,
			}
			err = gopacket.SerializeLayers(buf, opts, newLayers...)
			if err != nil {
				slog.Error("Could not serialize packet", "error", err)
				os.Exit(1)
			}
			ci.Length = len(buf.Bytes())
			ci.CaptureLength = len(buf.Bytes())
			err = w.WritePacket(ci, buf.Bytes())
			if err != nil {
				slog.Error("Could not write packet", "error", err)
				os.Exit(1)
			}
		} else {
			if hasTCP {
				storage++
			}

			// Write the packet. If not TCP it was not modified.
			err = w.WritePacket(ci, packet.Data())
			if err != nil {
				slog.Error("Could not write packet", "error", err)
				os.Exit(1)
			}
		}
	}
	if len(embedData) > 0 {
		slog.Error("Not enough space to embed data, needs a PCAP with more TCP packets", "stored", stored, "remaining", len(embedData))
		os.Exit(1)
	}

}
