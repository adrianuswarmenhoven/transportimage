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

func main() {
	data, err := os.ReadFile("medium/nineteen_eight_four_gutenberg_org_0100021.txt")
	if err != nil {
		slog.Error("Error reading file", "error", err)
		os.Exit(1)
	}

	img, _ := os.ReadFile("../images/example_image.jpg")
	encs := hex.EncodeToString(img)

	cursor := 0
	startMarkerWritten := false
	endMarkerWritten := false

	outfile := strings.Builder{}
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
		outfile.Write([]byte(fmt.Sprintf("%c", or)))
	}

	// Make a trailing tail of whitespace, not ideal, but works
	if cursor < len(encs) {
		for _, r := range encs[cursor:] {
			outfile.Write([]byte(fmt.Sprintf("%c", hexToUnicode[r])))
		}
	}
	os.WriteFile("out.txt", []byte(outfile.String()), 0644)
}
