# Image to Whitespace

This method was inspired by Shell Company's [Poltergeist](https://github.com/Shell-Company/poltergeist) and the esoteric programming language [Whitespace](https://en.wikipedia.org/wiki/Whitespace_%28programming_language%29)

As a medium I am using the Gutenberg Project file [https://gutenberg.net.au/ebooks01/0100021.txt](https://gutenberg.net.au/ebooks01/0100021.txt)

Please support [Project Gutenberg](https://www.gutenberg.org)

## Quick start

### If you need to make an executable
```bash
go mod verify && go mod tidy && go build
```

### When you have the executable
```
imageToWhitespace <inputfile> <medium> <output>
```

Where the file that is pointed to by ```<medium>``` is a Unicode plain text file, preferrably a book, treatise or essay.

## About this method

Just as with the PCAP method, we cheat a bit by not deconstructing the image but rather reading the raw bytes and encoding them in the text.

We change every byte into a hex value and then use one of 16 selected Unicode whitespace characters to put into the output file.
The placement of that whitespace is whenever we see a normal ASCII whitespace.

The following whitespace characters are used:
```
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
```

These were selected 'by eye' so that when a file is rendered as plain text, the output looks roughly similar to the original.

### Advantages

You can hide it in other books. You can hide it in a 'mirror' of the Gutenberg Project (but mind you, a specialized scraper will find it).
It compresses pretty good (well, the text is, your whitespace characters still bulk it up).
Scanners will check the [MIME type](https://en.wikipedia.org/wiki/MIME) and see ```plain/text```. Most scanners, at the moment, will completely ignore ```plain/text``` because it can not harm the system (well, there have been instances where Unicode characters crashed systems, but that was just bugs, nothing structural).
Some form of plausible deniability; you can make a case that you want to read the book and did not know about the stego (just don't have the viewer on the same system)

### Disadvantages

The resulting file is large.
The resulting file can have a large trail of whitespace at the end if the payload is too large. A scanner could efficiently just check the last few bytes of a file and see if they contain nothing but Unicode whitespace.
Making a scanner for this is not hard; just check for the existence of the multiple types of whitespace and check if the word before and after them are valid words that normally would have been separated by a regular ````\u0020```` space.