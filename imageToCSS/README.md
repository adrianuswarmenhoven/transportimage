# Image to CSS

## Quick start

### If you need to make an executable
```bash
go mod verify && go mod tidy && go build
```

### When you have the executable
```
imageToCSS <inputfile> <outputfile>
```
This will create an HTML page that uses HTML elements to designate the pixels and CSS to color them.
You can view the HTML page in any browser (mind you, rendering is quite slow)

The input file must be an image (preferrably JPG or PNG).
The output file should have the extension .html so you can doubleclick it and view it in your browser.

## About this method

This method deconstructs an image and then uses the RGB pixel data to create HTML elements which are colored using CSS.

There is a limited 'intelligence' in that pixels of the exact same color have the same CSS.

### Advantages

Given that this is HTML and any MIME scanner will tell you that it is HTML, it probably will not be scanned by any client side or server side scanner.
The fact that it does not use the IMG tag (preventing the DOM to have an image element) also helps towards the fact that any monitoring software does not see any images being transferred or even rendered.
For detection, this slowness helps in making detection 'expensive'; a never-seen-before file must be rendered first (see disadvantages... it is slow) and then be assessed before a hash or other discerning marker can be gotten.
Furthermore, a recipient does not need any special software to view it; any browser will suffice.

### Disadvantages

The resulting file is huge. Not just a little bit bigger than the original, but huge. This is somewhat mitigated by the fact that people send big(ger) files via instant messaging and webpage sizes have gone up considerably in the last decades.
However, encoding large images this way is not a good idea. This is mitigated by the fact that we can now send relative low-res images and upscale them via AI. They might not be the same as the origin image, but to some, that does not matter.

Rendering is slow. Not just a little bit slower, but really, slow. Try it in your browser (get the example file from [output/image_to_css.html](../output/image_to_css.html)). This is mitigated by the fact that once it is rendered a snapshot can be taken to save as an image file. 

It is also highly detectable once a scanner knows what it is looking for: huge CSS lists, elements that are only 1 pixel wide, no text content etc...

If the recipient does not know the sender, then malware like code can be implemented and viewing the image will compromise your system.