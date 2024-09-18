# Image to Spreadsheet

## Quick start

### If you need to make an executable
```bash
go mod verify && go mod tidy && go build
```

### When you have the executable
```
imageToSpreadsheet <inputfile> <outputfile>
```
This will create a spreadsheet that can be opened by regular office software like OpenOffice or Office365.
It will also create a macro file with the extension ```.bas``` that can be used to view the image. Ask your resident spreadsheet tiger to help you in 'integrating' that macro into the file.

## About this method

This method deconstructs an image and then uses the RGB pixel data to create an sheet in a spreadsheet workbook where each cell represents the RGB values. 
It then creates an accompanying macro file that should be integrated into the sheet to be able to view the data.

### Advantages

Hiding data in a spreadsheet is a pretty natural thing to do. Spreadsheets are made for immense amounts of numbers and pixel color values are numbers.
The fact that the macro is separate makes it harder to determine whether it is an image or whether it is the sales forecast for your 320 subsidiaries.

The data can be copied over to any online spreadsheet as well, since it is only numbers, and can then be copied back into a spreadsheet with the macro.

### Disadvantages

In it's current form, it is of course quite easy to 'see' the image (just open the spreadsheet without the macro). 
This means that any scanner that can read a spreadsheet could quickly check if it is an image by rendering the data.
Also, if the image is known, a hash can be calculated over the data (and not the containing spreadsheet).

Rendering is horrendously slow. This is mitigated by the fact that someone can take a snapshot once it finally has rendered. 

If the recipient does not know the sender, then malware like code can be implemented and viewing the image will compromise your system.