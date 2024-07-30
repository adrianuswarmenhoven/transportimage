package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strings"
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

	hexToUnicode = map[rune]rune{
		'0': '\u0020', // SPACE
		'1': '\u00A0', // NO-BREAK SPACE
		'2': '\u2000', // EN QUAD
		'3': '\u2001', // EM QUAD
		'4': '\u2002', // EN SPACE
		'5': '\u2003', // EM SPACE
		'6': '\u2004', // THREE-PER-EM SPACE
		'7': '\u2005', // FOUR-PER-EM SPACE
		'8': '\u2006', // SIX-PER-EM SPACE
		'9': '\u2007', // FIGURE SPACE
		'a': '\u2008', // PUNCTUATION SPACE
		'b': '\u2009', // THIN SPACE
		'c': '\u200A', // HAIR SPACE
		'd': '\u202F', // NARROW NO-BREAK SPACE
		'e': '\u205F', // MEDIUM MATHEMATICAL SPACE
		'f': '\u3000', // IDEOGRAPHIC SPACE
	}
)

var (
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
		slog.Info("Usage: imageToWhitespace <inputfile> <medium> <outputfile>")
		os.Exit(1)
	}

	data, err := os.ReadFile(mediumFile)
	if err != nil {
		slog.Error("Error reading medium file", "error", err)
		os.Exit(1)
	}

	img, err := os.ReadFile(inFile)
	if err != nil {
		slog.Error("Error reading input file", "error", err)
		os.Exit(1)
	}

	// Encode the data of the file to string
	// Yes, in this case we cheat because we do not get all the individual pixels
	// but we get the whole original image file as a hexadecimal string
	encs := hex.EncodeToString(img)

	cursor := 0
	startMarkerWritten := false
	endMarkerWritten := false

	outData := strings.Builder{}
	for _, r := range string(data) {
		or := r
		if r == ' ' && cursor < len(encs) {
			if !startMarkerWritten {
				or = startMarker
				startMarkerWritten = true
			} else {
				or = hexToUnicode[rune(encs[cursor])]
				cursor++
			}
		} else if !endMarkerWritten && r == ' ' && cursor >= len(encs) {
			or = endMarker
			endMarkerWritten = true
		}
		outData.Write([]byte(fmt.Sprintf("%c", or)))
	}

	// Make a trailing tail of whitespace, not ideal, but works
	// If the medium does not have enough whitespace, we will have to add it to the end.
	// This might make the output more suspicious
	if cursor < len(encs) {
		for _, r := range encs[cursor:] {
			outData.Write([]byte(fmt.Sprintf("%c", hexToUnicode[r])))
		}
	}
	err = os.WriteFile(outFile, []byte(outData.String()), 0644)
	if err != nil {
		slog.Error("Error writing output file", "error", err)
		os.Exit(1)
	}
}
