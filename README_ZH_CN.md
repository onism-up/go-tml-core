# go-tml-core


![build](https://img.shields.io/badge/build-passing-green)
![edition](https://img.shields.io/badge/edition-v0.0.2-orange)
![Powered](https://img.shields.io/badge/Powered%20by-onism-blue)
![Status](https://img.shields.io/badge/Status-Test%20phase-red)

## 语言
[English](/README.md)

## 演示
![Demonstration](https://raw.githubusercontent.com/onism-up/go-tml-core/main/static/medium/demo.gif)

## 特性
- 相对低的资源和内存使用
- 异步的图形渲染
- 大小自适应，动态位置计算
- 可组件化编程
- 支持较完整的事件

> 使用tml之前请确保目标运行终端支持VT100
## 开始
这里只展示了部分代码段，不能被直接运行，完整代码参考[demo](/demo)
```go
package main

import "github.com/onism-up/go-tml-core/tml"

func main(){
    tml.Start(100, true) //初始化
    
    style, _ := tml.Body.GetStyle()
    style.BackGroundColor = tml.YellowBackGroundColor
	tml.Body.SetStyle(style) //设置顶级节点样式
    
    button1 := button(router1(), "MoveBox", tml.CanvasPosition{
    X: 10,
    }, tml.CanvasPositionType{Center: tml.PositionY}) //创建新节点
    
    button2 := button(/* 省略了一些代码 */)
    
    button3 := button(/* 省略了一些代码 */)
    
    buttonBase := []tml.Node{button1, button2, button3}
    
    index := -1
	
    tml.Body.AddEventListener(tml.OnKeyBord, func(node tml.Node, _ tml.Node) { //键盘事件监听
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
        case tml.KeyEnter: //如果你选择了一个节点，那么keybord事件将会转移到被选择的节点
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
    //向顶级节点中插入子节点，因为渲染是以顶级节点开始进行树状渲染	
    tml.Body.Insert(title("<< demo: check left or right move light block >>"), button1, button2, button3)
	
    for true { //防止进程关闭
		
    }
}
```
## 自定义组件/渲染节点
- 渲染节点：渲染节点和上述节点不同，渲染节点是负责将node渲染的节点，由一个解析函数进行分类分发，如果你想自定义渲染节点在tml文件夹中新建一个go文件并以渲染节点名称命名，然后在renderer函数中添加解析节点，你的渲染组件必须继承Node接口并参考基础渲染节点中函数的实现
- 自定义组件：目前tml没有封装任何的组件，如果向开发组件则可以基于任何已知的基础组件进行封装，利用事件系统做好选择组件的前进和后退
