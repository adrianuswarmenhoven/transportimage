
	' To use this macro with the data for the image, you must add this code to the spreadsheet as a macro (might take some time to run).
	' The data for the image must be in a sheet named "Data".
	' The spreadsheet must have a sheet called 'View' to display the image. You can just create a new sheet and name it 'View'.
	'
	' This viewer is for the image file ../images/example_image.jpg
	Sub CreateCanvas()
	Dim oDrawPage As Object
	Dim oShape As Object
	Dim x As Integer, y As Integer
	Dim canvasWidth As Integer, canvasHeight As Integer
	Dim cellSize As Integer
	Dim offSetX As Integer
	Dim offSetY As Integer

	Dim oSize As New com.sun.star.awt.Size
	Dim oPosition As New com.sun.star.awt.Point

	canvasWidth = 256 ' Number of cells wide
	canvasHeight = 256 ' Number of cells high
	cellSize = 40 ' Size of each cell in pixels
	
	offSetX=1000
	offSetY=1000
	
	oDrawPage = ThisComponent.getSheets().getByName("View").getDrawPage()
	dataSheet = ThisComponent.getSheets().getByName("Data")
	
	' Loop through the canvas size
	For x = 0 To canvasWidth - 1
		For y = 0 To canvasHeight - 1
			' Calculate cell positions in pixels
			Dim xPos As Integer
			Dim yPos As Integer
			xPos = x * cellSize
			yPos = y * cellSize
			
			' Create a rectangle shape (simulating a pixel)
			oShape = ThisComponent.createInstance("com.sun.star.drawing.RectangleShape")

			oDrawPage.add(oShape)

			oSize.Width = cellSize
			oSize.Height = cellSize

			oPosition.X = offSetX+(cellSize * x)
			oPosition.Y = offSetY+(cellSize * y)
			
    		Dim cellValue As Long
   			 Dim red As Long, green As Long, blue As Long
			
			cell = dataSheet.getCellByPosition(x,y)
            cellValue = cell.Value
                      ' Ensure the cell value is a valid 32-bit integer
            If IsNumeric(cellValue) Then
                cellValue = CLng(cellValue)
                
                ' Extract the RGB components from the 32-bit integer
                red = (cellValue And 16711680) \ &H10000
                if red>255 then
                   red = 255
                endif
                
                green = (cellValue And 65280) \ &H100
                if green>255 then
                   green=255
                endif
                
                blue = cellValue And 255
                if blue >255 then
                   blue=255
                endif
               
            End If
			With oShape
				.setSize(oSize)
				.setPosition(oPosition)
				.FillColor = RGB(red, green, blue) 
				.FillStyle = com.sun.star.drawing.FillStyle.SOLID ' or .NONE
				.FillTransparence = 0 ' 0-100 from opaque to transparent
				.LineStyle = com.sun.star.drawing.LineStyle.SOLID ' or .NONE
				.LineWidth = 1 
				.LineColor = RGB(red, green, blue)
			End With
		Next y
	Next x
End Sub
