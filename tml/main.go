package tml

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"sync"
	"time"
)

// getWindowSize Gets the terminal window size
// @return The return value is width and Height of the terminal
func getWindowSize() (int, int) {
	fd := int(os.Stdout.Fd())
	w, h, err := terminal.GetSize(fd)

	if err != nil {
		panic(GetWindowSizeError + err.Error())
	}
	return w, h
}

// listenWindowSize Continuously listens to the serial port size of the terminal and triggers an event
// @parma fy: listening frequency. The base unit is 1ms
func listenWindowSize(fy int) {
	for true {
		width, height := getWindowSize()
		if (width != SysWidth || height != SysHeight) && Body != nil {
			SysWidth = width
			SysHeight = height
			Render()
			autoSizeChangeTrigger(Body, Body, true)

			time.Sleep(time.Millisecond * time.Duration(fy))
		}

	}
}

// listenKeyBord Continuously listens for global keyboard events
func listenKeyBord() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()
	for {
		char, key, err := keyboard.GetKey()

		if err != nil {
			panic(err.Error())
		}

		if SelectNode == nil {
			SelectNode = Body
		}

		SelectNode.setKeyBord(keyboard.KeyEvent{
			Key:  key,
			Rune: char,
			Err:  err,
		})

		if SelectNode.isUnMount() || !triggerEvent(SelectNode, OnKeyBord, SelectNode) { //尝试检测节点是否失效，如果失效会进行回退
			backSelect(SelectNode)
		}
	}
}

// Select the select a node to be used as the output node to listen for onKeyBord events. This node must listen for OnKeyBord events. Otherwise, the node automatically rolls back until the parent node has a listener or Body
// @parma node Selected node
func Select(node Node) {
	if node != nil {
		oldSelectNode := SelectNode
		SelectNode = node
		triggerEvent(node, OnSelect, oldSelectNode)
	}
}

// backSelect Attempts to roll back the cursor to a valid node
// @parma node Indicates the node to be rolled back
func backSelect(node Node) {
	nodeParent, _ := node.GetParent()
	if node != Body && nodeParent != nil {
		if nodeParent.isUnMount() {
			backSelect(nodeParent)
		} else {
			triggerEvent(nodeParent, OnSelect, node)
			SelectNode = nodeParent
		}
	} else if SelectNode != Body && nodeParent == nil {
		SelectNode = Body
	}

}

// Start initializes the service, using the framework's mandatory function, which does most of the initialization
// @parma autoFlash: Whether to automatically render, and the relatively timely response to the window during automatic rendering fy: rendering frequency and window size acquisition frequency, the base unit is ms keyBordEvent: whether to enable key event listening
func Start(fy uint16, keyBordEvent bool) {
	//确保未初始化
	if isInit {
		return
	}

	// 初始化name分类map
	nameLibrary = make(map[string]NodeStack)

	//初始化Body
	Body = CreateQuadrilateral(BodyName)
	Body.SetVolume(CanvasVolume{ // Only the width and height programs that display the declaration for adjustment have the right to change
		Width:  Auto,
		Height: Auto,
	})
	Body.SetStyle(CanvasStyle{AutoSize: true, Display: true, ShowText: true, Color: WhiteColor})

	Select(Body)

	if fy <= 0 { // Check the validity of the refresh frequency. If the refresh frequency does not comply with the rule, use the default value
		fy = SizeReLoadFrequency
	}

	fyParma := int(fy) // Converted to an int parameter

	// Query the window size every x ms
	go listenWindowSize(fyParma)

	if keyBordEvent {
		// Continuously listen for keyboard events
		go listenKeyBord()
	}

	// Modify the initialization signal
	isInit = true
}

// elementLoop renders a node and all its children, using recursion
// @parma The node tree that will be rendered
func elementLoop(node Node) {

	if node == Body {
		globalBuf.Reset() // Initialize the output file when re-rendering
	}

	renderResult := renderer(node) // Render node first to determine the node adaptability
	if !renderResult {             //父组件无法渲染则停止所有子组件的渲染
		return
	}

	zIndexTree := createZIndexTree(ZIndexRenderType)
	child, _ := node.GetChildren()
	zIndexTree.Init(child)
	childNode, ok := zIndexTree.GetNode()
	for ok {
		elementLoop(childNode)
		childNode, ok = zIndexTree.GetNode()
	}
}

// notRenderable is mainly used to detect whether the screen can be rendered
// @return The terminal cannot be used for rendering when the return value is true
func notRenderable() bool {
	return SysHeight < 0 || SysWidth < 0
}

// startElementLoop ElementLoop's startup function, which loops through node when rendered
// @parma fy: frequency of the loop
// @tips: discard function
func startElementLoop(fy int) {
	for true {
		Render()
		time.Sleep(time.Millisecond * time.Duration(fy))
	}
}

// renderDebounce Render anti-shake function, synchronization code after asynchronous rendering
var renderDebounce = debounce(RenderLazy)

// render render function
func render() {
	elementLoop(Body)
	print(globalBuf.String())
}

// Render Asynchronous rendering, which can be called manually
func Render() {
	if !notRenderable() {
		renderDebounce(render)
	}
}

// debounce anti-shake function
// @parma interval: indicates the anti-shake time
// @return is the final anti-shake function
func debounce(interval time.Duration) func(f func()) {
	var l sync.Mutex
	var timer *time.Timer

	return func(f func()) {
		l.Lock()
		defer l.Unlock()
		// Use lock to ensure that d.imer is stopped before updating.

		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(interval, f)
	}
}

// renderer Renderer that renders data as graphics
// @parma node: The node to be rendered
// @return node render the result, the result should come from different renderers, all renderers should return the correct rendering result for the renderer function
func renderer(node Node) bool {
	if notRenderable() || node.isUnMount() { //在窗口可视的情况下渲染，组件未销毁的情况下渲染
		return false
	}
	nodeAttr, _ := node.GetAttr()
	nodeStyle, _ := node.GetStyle()
	nodeVolume, _ := node.GetVolume()

	if nodeStyle.Display == false || (nodeStyle.AutoSize == false && (nodeVolume.Width <= 0 || nodeVolume.Height <= 0)) {
		return false
	}

	switch nodeAttr.Tag { // Determine the type of node and push the node to different parsers for classification rendering

	case QuadrilateralTag:
		return quadrilateralRender(node.(*Quadrilateral))
	}
	return false
}

// quadrilateralRender indicates deformation parsing render
// @parma ql: pointer to the quadrilateral struct
// @return returns the render result of the node, false if it does not render properly, so his children will not parse the render again
func quadrilateralRender(ql *Quadrilateral) bool {

	if (ql.volume.Width <= 0 || ql.volume.Height <= 0) && !ql.style.AutoSize {
		return false
	} else {
		ql.position.totalY = ql.position.Y
		ql.position.totalX = ql.position.X
		if ql.parent != nil {
			pPosition, _ := ql.parent.GetPosition()
			ql.position.totalY += pPosition.totalY
			ql.position.totalX += pPosition.totalX
		}

		return squareDrawing(ql)
	}
}

// CreateQuadrilateral Creates a quadrilateral Canvas
// @parma Name: Creates the name of the quadrilateral to mark the node, but does not force uniqueness
// @return creates a node that will be loaded into each global repository before it returns
func CreateQuadrilateral(name string) Node {
	key, _ := uuid.NewRandom()
	element := new(Quadrilateral)
	element.tag = QuadrilateralTag
	element.name = name
	element.key = fmt.Sprintf("%s", key)
	element.position.ZIndex = 0
	element.position.Y = 0
	element.position.X = 0
	element.unMount = false

	element.style = DefaultStyle

	loadNodeToBase(element)        // Loads the node into the sort repository
	lodeNodeToNameIndex(element)   // Loads a node into a repository classified by the Name field
	loadNodeToRenderStack(element) // Loads a node into the render-level repository
	createEvent(element)           // Create an event

	return element
}

// CreateNode Easy to create different types of nodes
// @parma nodeType: Type of Node name: Node name
// @return The final Node created will return nil if the type is incorrect
func CreateNode(nodeType string, name string) Node {
	var node Node = nil
	switch nodeType {
	case QuadrilateralTag:
		node = CreateQuadrilateral(name)
	}
	return node
}

// loadBase loads the node to the corresponding global repository
// @parma node: indicates the node to be loaded
func loadNodeToBase(node Node) {
	nodeAttr, _ := node.GetAttr()
	switch nodeAttr.Tag {
	case QuadrilateralTag:
		ql, ok := (node).(*Quadrilateral)
		if ok {
			quadrilateralStore.Store(ql.key, ql) // Stored in a global variable to facilitate the renderer to obtain information
		}
		break
	}
}

// loadNodeToRenderStack loads the node to the index of the corresponding level. The load process attempts to de-load the node
// @parma node: indicates the node to be loaded
func loadNodeToRenderStack(node Node) {
	position, _ := node.GetPosition()
	attr, _ := node.GetAttr()
	nodeStackAny, ok := zIndexStack.Load(position.ZIndex)
	if ok {
		nodeStack := nodeStackAny.(NodeStack)
		for _, stackNode := range nodeStack {
			snAttr, _ := stackNode.GetAttr()
			if snAttr.Key == attr.Key {
				return
			}
		}
		nodeStack = append(nodeStack, node)
		zIndexStack.Store(position.ZIndex, nodeStack)
	} else {
		zIndexStack.Store(position.ZIndex, NodeStack{node})
	}
}

// delNodeFromRenderStack Deletes a Node from the rendering stack
// @parma node: indicates the node to be deleted
func delNodeFromRenderStack(node Node) {
	position, _ := node.GetPosition()
	attr, _ := node.GetAttr()

	nodeStackAny, ok := zIndexStack.Load(position.ZIndex)

	if ok {
		nodeStack := nodeStackAny.(NodeStack)
		for index, stackNode := range nodeStack {
			snAttr, _ := stackNode.GetAttr()
			if snAttr.Key == attr.Key {
				newNodeStack := append(nodeStack[:index], nodeStack[index+1:]...) //删除元素
				zIndexStack.Store(position.ZIndex, newNodeStack)
				break
			}
		}
	}
}

// getElementFromBase searches for an element in the global repository
// @parma Key: element key nodeType: element type. You can enter * to obtain it from any repository
// @return returns the Node type. The detailed type of the node needs to be determined by calling GetAttr. This method returns the uniform type
func getElementFromBase(key string) (Node, bool) {
	node, ok := getQuadrilateralElementFromBase(key)
	if ok {
		return node, ok
	}

	return node, ok
}

// delNodeFromBase Deletes the node from the corresponding global repository
// @parma node: indicates the node to be deleted
func delNodeFromBase(node Node) {
	nodeAttr, _ := node.GetAttr()
	switch nodeAttr.Tag {
	case QuadrilateralTag:
		ql, ok := (node).(*Quadrilateral)
		if ok {
			quadrilateralStore.Delete(ql.key) //Save to global variables to facilitate the renderer to obtain information
		}
		break
	}
}

// getQuadrilateralElementFromBase find an element in the quadrangle warehouse
// @parma Key: The key to find
// @return Indicates whether the search result is successful
func getQuadrilateralElementFromBase(key string) (*Quadrilateral, bool) {
	value, ok := quadrilateralStore.Load(key)
	ql, ok := value.(*Quadrilateral)
	return ql, ok
}

// lodeNodeToNameIndex stores elements in groups with their name field
// @parma node: indicates the node to be loaded
// @return The result of loading
func lodeNodeToNameIndex(node Node) bool {

	attr, _ := node.GetAttr()
	base, ok := nameLibrary[attr.Name]
	if ok {
		for _, value := range base { //去重
			snAttr, _ := value.GetAttr()
			if snAttr.Key == attr.Key {
				return false
			}
		}
		nameLibrary[attr.Name] = append(base, node)
	} else {
		nameLibrary[attr.Name] = NodeStack{node}
	}
	return true
}

// delNodeFromNameIndex removes the element with its name field
// @parma node: indicates the node to be deleted
func delNodeFromNameIndex(node Node) {
	attr, _ := node.GetAttr()
	base, ok := nameLibrary[attr.Name]
	if ok {
		for index, value := range base {
			snAttr, _ := value.GetAttr()
			if snAttr.Key == attr.Key {
				nameLibrary[attr.Name] = append(base[:index], base[index+1:]...)
			}
		}
	}

}

// GetNodeByName Tries to get a Node by Name
// @parma Name: node.Name
// @return The obtained node stack and the obtained result
func GetNodeByName(name string) (NodeStack, bool) {
	nameIndex, ok := nameLibrary[name]
	return nameIndex, ok
}

// createZIndexTree attempts to create a hierarchy tree, allowing different types of hierarchy trees to be selected, and defaulting to RenderTree if there are more supported types
// @parma typeType: type of the painting tree
// @return The return value is a hierarchy tree
func createZIndexTree(treeType uint8) PaintingTree {
	switch treeType {
	case ZIndexRenderType:
		return new(RenderZIndexTree)
	default:
		return new(RenderZIndexTree)
	}
}

// throw is used to throw an error internally
// @parma Error message
func throw(msg string) {
	str := strings.Builder{}
	str.WriteString("<<")
	str.WriteString(PackageName)
	str.WriteString(">>")
	str.WriteByte(' ')
	str.WriteString(msg)

	panic(str.String())
}

// Remove Used to destroy nodes in batches
func Remove(nodes ...Node) {
	for _, node := range nodes {
		node.Remove() // Try whether it is destroyed or not
	}
}
