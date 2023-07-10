# go-tml-core


![build](https://img.shields.io/badge/build-passing-green)
![edition](https://img.shields.io/badge/edition-v0.0.2-orange)
![Powered](https://img.shields.io/badge/Powered%20by-onism-blue)
![Status](https://img.shields.io/badge/Status-Test%20phase-red)

## Demonstration
![Demonstration](https://picdm.sunbangyan.cn/2023/07/10/nuo6go.gif)

## Peculiarity
- Low memory and resource usage
- Asynchronous graphics renderer
- Highly scalable
- Componentized components can be used
- More complete event support
- Window adaptive, position dynamic calculation

> Ensure that the target terminal supports VT100 before use,The output of the program is based on VT100.
## Start
Here shows part of the code, he can not run correctly, if you want to see all the code, please view or debug in the [demo](github.com/onism-up/go-tml-core/tree/main/demo)
```go
package main

import "github.com/onism-up/go-tml-core/tml"

func main(){
    tml.Start(100, true) //initialize
    
    style, _ := tml.Body.GetStyle()
    style.BackGroundColor = tml.YellowBackGroundColor
	tml.Body.SetStyle(style) //root node style Settings
    
    button1 := button(router1(), "MoveBox", tml.CanvasPosition{
    X: 10,
    }, tml.CanvasPositionType{Center: tml.PositionY}) //Create component
    
    button2 := button(/* Some information is omitted */)
    
    button3 := button(/* Some information is omitted */)
    
    buttonBase := []tml.Node{button1, button2, button3}
    
    index := -1
	
    tml.Body.AddEventListener(tml.OnKeyBord, func(node tml.Node, _ tml.Node) { //Listening event
        keyBord, _ := node.GetKeyBord()
        
        switch keyBord.Key {
        case tml.KeyArrowLeft:
            if index > 0 {
            index--
            }
            break   
        case tml.KeyArrowRight:
            if index < 2 {
            index++
            }
            break
        case tml.KeyEnter: //If transfer is selected, the transferred node will continue to listen for keyboard events
            if index <= 2 && index >= 0 {
                tml.Select(buttonBase[index])
            }
        }
        
        for buttonIndex, buttonNode := range buttonBase {
            buttonStyle, _ := buttonNode.GetStyle()
            
            if buttonIndex == index {
                buttonStyle.BackGroundColor = tml.PurpleBackGroundColor
                buttonStyle.BorderType = tml.ContinuousLine
                buttonStyle.BorderColor = tml.WhiteColor
            } else {
                buttonStyle.BackGroundColor = tml.BlueBackGroundColor
                buttonStyle.BorderType = tml.None
            }
            
            buttonNode.SetStyle(buttonStyle)
        }})
	//Insert all nodes into the top-level node to display
    tml.Body.Insert(title("<< demo: check left or right move light block >>"), button1, button2, button3)
	
    for true { //Prevent program termination
		
    }
}
```
## Customize a Node or component
At present, there is no wrapped component in tml, only a base component and a base node, if you want to create a new rendering node, you can try to create a new file in tml, and then add a node to the renderer, if you want to encapsulate a new component, you can try based on the base node, if you feel that your component or rendering node is good, You can try to push it to the project