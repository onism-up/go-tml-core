package demo

import (
	UI "github.com/onism-up/go-tml-core/tml"
)

func title(text string) UI.Node {
	node := UI.CreateQuadrilateral("title")

	node.SetText(text)

	node.SetVolume(UI.CanvasVolume{
		Width:  len(text),
		Height: 1,
	})

	node.SetPosition(UI.CanvasPosition{
		Y: 1,
	}, UI.CanvasPositionType{Center: UI.PositionX})

	return node
}

func displayNode(node UI.Node, display bool) UI.Node {
	style, _ := node.GetStyle()

	style.Display = display

	node.SetStyle(style)

	return node
}

func button(jumpNode UI.Node, text string, position UI.CanvasPosition, pType UI.CanvasPositionType) UI.Node {
	node := UI.CreateQuadrilateral("button")
	node.SetText(text)

	style, _ := node.GetStyle()

	style.BackGroundColor = UI.BlueBackGroundColor

	node.SetStyle(style)

	node.SetVolume(UI.CanvasVolume{
		Width:  20,
		Height: 5,
	})

	UI.Body.Insert(jumpNode)

	node.SetPosition(position, pType)

	displayNode(jumpNode, false)

	node.AddEventListener(UI.OnSelect, func(node UI.Node, origen UI.Node) {

		if origen == jumpNode {
			displayNode(jumpNode, false)
			UI.Select(UI.Body)
		} else {
			displayNode(jumpNode, true)
			UI.Select(jumpNode)
		}
	})

	node.AddEventListener(UI.OnKeyBord, func(node UI.Node, origen UI.Node) {

	})

	return node
}

func page() UI.Node {
	node := UI.CreateQuadrilateral("page")
	style, _ := node.GetStyle()
	position, _ := node.GetPosition()

	style.AutoSize = true

	position.ZIndex = 1

	node.SetStyle(style)
	node.SetPosition(position, position.Type)
	node.SetVolume(UI.CanvasVolume{Width: UI.Auto, Height: UI.Auto})

	return node
}

func box() UI.Node {
	node := UI.CreateQuadrilateral("box")

	style, _ := node.GetStyle()

	style.BackGroundColor = UI.BlackBackGroundColor
	style.BorderType = UI.DottedLine
	style.BorderColor = UI.GreenColor

	node.SetStyle(style)

	node.SetVolume(UI.CanvasVolume{Width: 10, Height: 5})

	node.SetText("Box")

	return node
}

func router1() UI.Node {
	node := page()
	style, _ := node.GetStyle()

	style.BackGroundColor = UI.WhiteBackGroundColor
	style.Color = UI.BlackColor

	node.SetStyle(style)

	boxNode := box()

	boxNode.SetPosition(UI.CanvasPosition{Y: 2, X: 1}, UI.CanvasPositionType{})

	node.AddEventListener(UI.OnKeyBord, func(node UI.Node, origen UI.Node) {
		keyBord, _ := node.GetKeyBord()
		if keyBord.Key == UI.KeyEsc {
			UI.Select(UI.Body)
			displayNode(node, false)
		}

		boxPosition, _ := boxNode.GetPosition()

		offsetX := 0
		offsetY := 0
		switch keyBord.Rune {
		case 'w':
			if boxPosition.Y > 3 {
				offsetY--
				boxNode.SetText("GoTop")
			}
			break
		case 's':
			if boxPosition.Y < UI.SysHeight-15 {
				offsetY++
				boxNode.SetText("GoBottom")
			}
		case 'a':
			if boxPosition.X > 1 {
				offsetX--
				boxNode.SetText("GoLeft")
			}
		case 'd':
			if boxPosition.X < UI.SysWidth-15 {
				offsetX++
				boxNode.SetText("GoRight")
			}
		default:
			boxNode.SetText("stop")
		}

		boxPosition.X += offsetX
		boxPosition.Y += offsetY

		boxNode.SetPosition(boxPosition)
	})

	node.Insert(title("The move button demo,check w,a,s,d move box, check esc close the page"), boxNode)

	return node
}

func text(output string, bgc, lightC string) UI.Node {
	lenOut := len(output)
	if lenOut > 0 {
		childStack := UI.NodeStack{}
		node := UI.CreateQuadrilateral("text")
		line := 0
		lastCut := 0
		for i, k := range output {
			gap := (i - lastCut)
			if k == '\n' || (gap != 0 && UI.SysWidth != 0 && gap%UI.SysWidth == 0) { //解析转义符
				child := textChild(output[lastCut:i], line)
				line++
				lastCut = i + 1
				style, _ := child.GetStyle()
				style.BackGroundColor = bgc
				child.SetStyle(style)
				childStack = append(childStack, child)
			}
		}

		style, _ := node.GetStyle()

		style.BackGroundColor = bgc

		style.AutoSize = true

		node.SetVolume(UI.CanvasVolume{Width: UI.Auto, Height: UI.Auto})

		node.SetStyle(style)

		index := -1

		node.Insert(childStack...)

		node.AddEventListener(UI.OnKeyBord, func(node UI.Node, origen UI.Node) {
			keyBord, _ := node.GetKeyBord()

			if keyBord.Key == UI.KeyArrowUp && index > 0 {
				index--
			} else if keyBord.Key == UI.KeyArrowDown && index < len(childStack) {
				index++
			} else if keyBord.Key == UI.KeyEsc {
				parent, _ := node.GetParent()
				UI.Select(parent)
			}
			for cIndex, cNode := range childStack {
				cStyle, _ := cNode.GetStyle()
				if cIndex == index {
					cStyle.BackGroundColor = lightC
				} else {
					cStyle.BackGroundColor = bgc
				}
				cNode.SetStyle(cStyle)
			}

		})

		return node
	} else {
		return nil
	}
}

func textChild(output string, line int) UI.Node {
	node := UI.CreateQuadrilateral("textLine")
	node.SetVolume(UI.CanvasVolume{Height: 1, Width: UI.Auto})
	style, _ := node.GetStyle()
	position, _ := node.GetPosition()

	style.AutoSize = true
	position.Y = line
	node.SetStyle(style)
	node.SetPosition(position)
	node.SetText(output)
	return node
}

func router2() UI.Node {
	node := page()
	style, _ := node.GetStyle()
	jsonStr := "{\n  \"sites\": {\n    \"site\": [\n      {\n        \"id\": \"1\",\n        \"name\": \"runoob,\n        \"taps\": \"check up or down to move line height\"\n      },\n      {\n        \"id\": \"2\",\n        \"name\": \"runoobTools\",\n        \"url\": \"c.runoob.com\"\n      },\n      {\n        \"id\": \"3\",\n        \"name\": \"Google\",\n        \"taps\": \"if you want to quit,check esc\"\n      }\n    ]\n  }\n}\n"

	style.BackGroundColor = UI.RedBackGroundColor

	node.SetStyle(style)

	textNode := text(jsonStr, UI.RedBackGroundColor, UI.GreenBackGroundColor)

	textStyle, _ := textNode.GetStyle()

	textStyle.AutoSize = true

	textNode.SetStyle(textStyle)

	textNode.SetPosition(UI.CanvasPosition{Y: 3})

	node.AddEventListener(UI.OnSelect, func(node UI.Node, origen UI.Node) {
		if origen != textNode {
			UI.Select(textNode)
		} else {
			UI.Select(UI.Body)
			displayNode(node, false)
		}
	})

	node.Insert(title("Electronic Reader"), textNode)

	return node
}

func router3() UI.Node {
	node := page()
	style, _ := node.GetStyle()

	style.BackGroundColor = UI.GreenBackGroundColor

	style.BorderType = UI.DottedLine

	style.BorderColor = UI.BlackColor

	style.Color = UI.RedColor

	node.SetStyle(style)

	node.SetText("This is Page 3  tips: if you want backoff please check 1")

	node.AddEventListener(UI.OnKeyBord, func(node UI.Node, origen UI.Node) {
		keyBord, _ := node.GetKeyBord()
		node.SetText("you check: " + string(keyBord.Rune))
		if keyBord.Rune == '1' {
			UI.Select(UI.Body)
			displayNode(node, false)
		}
	})

	return node
}

func Run() {
	UI.Start(100, true)

	style, _ := UI.Body.GetStyle()
	style.BackGroundColor = UI.YellowBackGroundColor
	UI.Body.SetStyle(style)

	button1 := button(router1(), "MoveBox", UI.CanvasPosition{
		X: 10,
	}, UI.CanvasPositionType{Center: UI.PositionY})

	button2 := button(router2(), "Electronic Reader", UI.CanvasPosition{
		X: 35,
	}, UI.CanvasPositionType{Center: UI.PositionXY})

	button3 := button(router3(), "Awaiting development", UI.CanvasPosition{
		X: 10,
	}, UI.CanvasPositionType{Center: UI.PositionY, Right: UI.PositionX})

	buttonBase := []UI.Node{button1, button2, button3}

	index := -1

	UI.Body.AddEventListener(UI.OnKeyBord, func(node UI.Node, _ UI.Node) {
		keyBord, _ := node.GetKeyBord()

		switch keyBord.Key {
		case UI.KeyArrowLeft:
			if index > 0 {
				index--
			}
			break
		case UI.KeyArrowRight:
			if index < 2 {
				index++
			}
			break
		case UI.KeyEnter:
			if index <= 2 && index >= 0 {
				UI.Select(buttonBase[index])
			}
		}

		for buttonIndex, buttonNode := range buttonBase {
			buttonStyle, _ := buttonNode.GetStyle()

			if buttonIndex == index {
				buttonStyle.BackGroundColor = UI.PurpleBackGroundColor
				buttonStyle.BorderType = UI.ContinuousLine
				buttonStyle.BorderColor = UI.WhiteColor
			} else {
				buttonStyle.BackGroundColor = UI.BlueBackGroundColor
				buttonStyle.BorderType = UI.None
			}

			buttonNode.SetStyle(buttonStyle)
		}

	})

	UI.Body.Insert(title("<< demo: check left or right move light block >>"), button1, button2, button3)
	for true {

	}
}
