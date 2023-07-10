package UI

// squareDrawing renders a Quadrilateral in a terminal by parsing the quadrilateral
// @parma ql: target rendered quadrilateral
// @return Whether the rendering is successful
func squareDrawing(ql *Quadrilateral) bool {

	style := ql.style
	cWidth, cHeight := getCanvasSize(ql)

	oldVolume := ql.volume

	if style.AutoSize {
		if ql.volume.Width == Auto {
			ql.volume.Width = cWidth
		}

		if ql.volume.Height == Auto {
			ql.volume.Height = cHeight
		}
	}

	parent := ql.parent
	position := ql.position
	volume := ql.volume

	cTop := 0
	cLeft := 0
	cRight := 0
	cBottom := 0
	yStart := position.totalY
	xStart := position.totalX
	yEnd := yStart + volume.Height
	xEnd := xStart + volume.Width

	basicsX := ""
	basicsY := ""

	globalBuf.WriteString(hiddenCursor)

	if parent != nil {
		pPosition, _ := parent.GetPosition()

		cTop = pPosition.totalY
		cLeft = pPosition.totalX
		cRight = cLeft + cWidth
		cBottom = cTop + cHeight

		if position.X > cWidth || position.Y > cHeight || position.X+volume.Width <= 0 || position.Y+volume.Height <= 0 { //不会在父节点中渲染
			return false
		}
	} else {
		cRight = SysWidth
		cBottom = SysHeight
	}

	if position.Type.Center != None { //Parse position.type.center
		nType := position.Type.Center
		if nType == PositionX || nType == PositionXY {
			xStart = confirmSquareCenter(cRight, volume.Width)
			xEnd = xStart + volume.Width
		}
		if nType == PositionY || nType == PositionXY {
			yStart = confirmSquareCenter(cBottom, volume.Height)
			yEnd = yStart + volume.Height
		}
	}

	if position.Type.Right != None { //Resolve position.type.right center and right exist at the same time, and the weight of right changes
		nType := position.Type.Right
		if nType == PositionX || nType == PositionXY {
			xEnd = confirmSquareEnd(position.X, cRight)
			xStart = xEnd - volume.Width
		}
		if nType == PositionY || nType == PositionXY {
			yEnd = confirmSquareEnd(position.Y, cBottom)
			yStart = yEnd - volume.Height
		}
	}

	qlYStart := yStart
	qlXStart := xStart
	qlYEnd := yEnd
	qlXEnd := xEnd

	yStart = confirmStartSquare(yStart, cTop) // Determine the final render position
	xStart = confirmStartSquare(xStart, cLeft)
	yEnd = confirmEndSquare(yEnd, cBottom)
	xEnd = confirmEndSquare(xEnd, cRight)

	if style.Color != "" {
		globalBuf.WriteString(style.Color)
	}

	if style.BackGroundColor != "" {
		globalBuf.WriteString(style.BackGroundColor)
	}

	if style.BorderColor != "" {
		basicsY = style.BorderColor
		basicsX = style.BorderColor
	}

	if style.BorderType != None { // Parsing the style
		switch style.BorderType {
		case ContinuousLine:
			basicsX += "-"
			basicsY += "|"
			break
		case DottedLine:
			basicsX += "."
			basicsY += "."
		}
	}

	if style.BorderColor != "" && style.Color != "" {
		basicsY += style.Color
		basicsX += style.Color
	}

	endLinePositionY := qlYEnd - 1
	endLinePositionX := qlXEnd - 1
	textIndex := 0
	for i := yStart; i < yEnd; i++ { //render y
		globalBuf.WriteString(setCursorPosition(uint32(xStart)+1, uint32(i)+1)) //set cursor
		for k := xStart; k < xEnd; k++ {                                        //render x
			if (i == qlYStart || i == endLinePositionY) && style.BorderType != None {
				globalBuf.WriteString(basicsX)
			} else if (k == qlXStart || k == endLinePositionX) && style.BorderType != None {
				globalBuf.WriteString(basicsY)
			} else if style.ShowText && textIndex < len(ql.text) {
				globalBuf.WriteByte(ql.text[textIndex])
				textIndex++
			} else {
				globalBuf.WriteByte(' ')
			}
		}
	}
	if style.AutoSize { // Re-place after dynamic calculation of width and height
		ql.volume = oldVolume
	}

	globalBuf.WriteString(closeAllProperties)

	return true
}

// confirmStartSquare auxiliary function, which assists the squareDrawing function to determine the rendering starting point
// @parma positionX: start of the current element, start: start of the renderable element
func confirmStartSquare(positionX int, start int) int {
	if (positionX < start || positionX <= 0) && start <= 0 {
		return 0
	} else if positionX < start && start > 0 {
		return start
	} else {
		return positionX
	}
}

// confirmEndSquare auxiliary function, which assists the squareDrawing function in determining the rendering endpoint
// @parma positionX: end point of the current element, start: end point of the render
func confirmEndSquare(positionX int, end int) int {
	if (positionX > end || positionX >= SysWidth) && end >= SysWidth {
		return SysWidth
	} else if positionX > end && end < SysWidth {
		return end
	} else {
		return positionX
	}
}

// confirmSquareCenter auxiliary function, which assists the squareDrawing function to determine the relative rendering starting point
// @parma pCount: offset of the current node, width: size of the canvas that can be rendered
func confirmSquareCenter(pCount, size int) int {
	return pCount/2 - size/2
}

// confirmSquareEnd auxiliary function, which assists the squareDrawing function to determine the relative End rendering starting point
// @parma pCount: offset of the current node, width: size of the canvas that can be rendered
func confirmSquareEnd(pCount, size int) int {
	return size - pCount
}

// getCanvasSize Gets the size of the canvas that can be rendered by the current node, usually the size of its parent
// @parma node: indicates the node to be obtained
// @return Wide, high
func getCanvasSize(node Node) (int, int) {
	parent, _ := node.GetParent()
	if parent == nil {
		return SysWidth, SysHeight
	}
	style, _ := parent.GetStyle()
	volume, _ := parent.GetVolume()
	if style.AutoSize {
		return getCanvasSize(parent)
	} else {
		return volume.Width, volume.Height
	}
}
