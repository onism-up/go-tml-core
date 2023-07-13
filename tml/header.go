package tml

import (
	"github.com/eiannone/keyboard"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Universal constant
const (
	Auto                int    = math.MaxInt      //Display declaration run program adaptive
	SizeReLoadFrequency uint16 = 10               //Frequency of window query, in ms
	BodyName                   = "body"           //The Name of the top-level node is fixed, but the node whose Name is body may have more than body. When getting body, it is recommended to use the global variable instead of the GetNodeByName function
	QuadrilateralTag           = "quadrilateral"  //Tag of different types of nodes. Here, the tag is a quadrangle
	ZIndexRenderType    uint8  = 0                //Different types of renderers, which can get Node through the hierarchy
	PackageName                = "TMLRenderer"    //Package name, usually used for information printing
	RenderLazy                 = time.Microsecond //Asynchronous wait time
)

// Universal variable
var (
	SysWidth           int                  = 0     //Global width
	SysHeight          int                  = 0     //Global altitude
	isInit             bool                 = false //Whether to initialize
	quadrilateralStore sync.Map                     //Quadrilateral storage library Key:string value:Node
	nameLibrary        map[string]NodeStack         //A library of nodes sorted by name
	Body               Node                 = nil   //Root node
	zIndexStack        sync.Map                     //Render level Key:uint32 index level  value:NodeStack The node whose hierarchy is stored
	DefaultStyle       = CanvasStyle{       //Create a default style for the node
		Display:         true,
		BorderType:      None,
		AutoSize:        false,
		ShowText:        true,
		Color:           WhiteColor,
		BackGroundColor: BlackBackGroundColor,
	}
	globalBuf  = strings.Builder{}       //Final printed V100 data
	eventStore sync.Map                  //Event storage, to prevent thread conflicts, map using sync.map map[string]Event{}
	SelectNode Node                = nil //The currently selected node
)

// style correlation constant
const (
	None           uint8 = 0 //This style is usually not valid
	ContinuousLine uint8 = 1 //Continuous line
	DottedLine     uint8 = 2 //Dashed line
	PositionX      uint8 = 1 //Used for x centering in Position
	PositionY      uint8 = 2 //Used for y centering in Position
	PositionXY     uint8 = 3 //Used for x and y centering in Position

)

// Error information constant
const (
	OperatingEmptyNodeError           = "you are trying to operate a node that has been unmounted node"
	DeleteTopLeveNodeError            = "you are trying to delete a top-level node, which is not allowed to be destroyed"
	WrongCursorMovementDirectionError = "incorrect cursor movement. Please check the writing specification"
	RenderUninitializedNodeError      = "trying to render an uninitialized node"
	GetWindowSizeError                = "an attempt to get the window size failed, causing the framework to fail: "
	ParentNodeNil                     = "The parent node is nil, and setting the parent node to nil may cause the cursor to reset"
)

// VT100 exclusive
const (
	vT100Basics                = "\033["
	closeAllProperties         = "\033[0m" //Close all properties
	highlight                  = "\033[1m" //Set to highlight
	underline                  = "\033[4m" //Underline
	flicker                    = "\033[5m" //Flicker
	backDisplay                = "\033[7m" //Reverse display
	blanking                   = "\033[8m" //Blanking
	left                  byte = 'D'
	right                 byte = 'C'
	top                   byte = 'A'
	bottom                byte = 'B'
	clearScreen                = "\033[2J"   //Clear screen
	clearTheCursorEnd          = "\033[K"    //Clear the content from the cursor to the end of the line
	saveCursor                 = "\033[s"    //Save cursor position
	restoreCursor              = "\033[u"    //Restore cursor position
	hiddenCursor               = "\033[?25l" //Hide cursor
	showCursor                 = "\033[?25h" //Show cursor
	BlackColor                 = "\033[30m"  //Black color
	RedColor                   = "\033[31m"  //Red color
	GreenColor                 = "\033[32m"  //Green color
	YellowColor                = "\033[33m"  //Yellow color
	BlueColor                  = "\033[34m"  //Blue color
	PurpleColor                = "\033[35m"  //Purple color
	CyanColor                  = "\033[36m"  //Cyan color
	WhiteColor                 = "\033[37m"  //White color
	BlackBackGroundColor       = "\033[40m"  //Black background color
	RedBackGroundColor         = "\033[41m"  //Red background color
	GreenBackGroundColor       = "\033[42m"  //Green background color
	YellowBackGroundColor      = "\033[43m"  //Yellow background color
	BlueBackGroundColor        = "\033[44m"  //Blue background color
	PurpleBackGroundColor      = "\033[45m"  //Purple background color
	CyanBackGroundColor        = "\033[46m"  //Cyan background color
	WhiteBackGroundColor       = "\033[47m"  //White background color
)

// Event specific constant
const (
	OnMove       uint8 = 0
	OnSizeChange uint8 = 1
	OnShow       uint8 = 2
	OnHidden     uint8 = 3
	OnInput      uint8 = 4
	OnRemove     uint8 = 5
	OnSelect     uint8 = 6
	OnKeyBord    uint8 = 7
)

type Key uint16

// Keycode
const (
	KeyF1 keyboard.Key = 0xFFFF - iota
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	key_min // see terminfo
)

const (
	KeyCtrlTilde      keyboard.Key = 0x00
	KeyCtrl2          keyboard.Key = 0x00
	KeyCtrlSpace      keyboard.Key = 0x00
	KeyCtrlA          keyboard.Key = 0x01
	KeyCtrlB          keyboard.Key = 0x02
	KeyCtrlC          keyboard.Key = 0x03
	KeyCtrlD          keyboard.Key = 0x04
	KeyCtrlE          keyboard.Key = 0x05
	KeyCtrlF          keyboard.Key = 0x06
	KeyCtrlG          keyboard.Key = 0x07
	KeyBackspace      keyboard.Key = 0x08
	KeyCtrlH          keyboard.Key = 0x08
	KeyTab            keyboard.Key = 0x09
	KeyCtrlI          keyboard.Key = 0x09
	KeyCtrlJ          keyboard.Key = 0x0A
	KeyCtrlK          keyboard.Key = 0x0B
	KeyCtrlL          keyboard.Key = 0x0C
	KeyEnter          keyboard.Key = 0x0D
	KeyCtrlM          keyboard.Key = 0x0D
	KeyCtrlN          keyboard.Key = 0x0E
	KeyCtrlO          keyboard.Key = 0x0F
	KeyCtrlP          keyboard.Key = 0x10
	KeyCtrlQ          keyboard.Key = 0x11
	KeyCtrlR          keyboard.Key = 0x12
	KeyCtrlS          keyboard.Key = 0x13
	KeyCtrlT          keyboard.Key = 0x14
	KeyCtrlU          keyboard.Key = 0x15
	KeyCtrlV          keyboard.Key = 0x16
	KeyCtrlW          keyboard.Key = 0x17
	KeyCtrlX          keyboard.Key = 0x18
	KeyCtrlY          keyboard.Key = 0x19
	KeyCtrlZ          keyboard.Key = 0x1A
	KeyEsc            keyboard.Key = 0x1B
	KeyCtrlLsqBracket keyboard.Key = 0x1B
	KeyCtrl3          keyboard.Key = 0x1B
	KeyCtrl4          keyboard.Key = 0x1C
	KeyCtrlBackslash  keyboard.Key = 0x1C
	KeyCtrl5          keyboard.Key = 0x1D
	KeyCtrlRsqBracket keyboard.Key = 0x1D
	KeyCtrl6          keyboard.Key = 0x1E
	KeyCtrl7          keyboard.Key = 0x1F
	KeyCtrlSlash      keyboard.Key = 0x1F
	KeyCtrlUnderscore keyboard.Key = 0x1F
	KeySpace          keyboard.Key = 0x20
	KeyBackspace2     keyboard.Key = 0x7F
	KeyCtrl8          keyboard.Key = 0x7F
)

// cursorMovement Control cursor movement
// @parma direction: direction, see constants for details, step: the number of steps to move
// @return the final VT100 style
func cursorMovement(direction byte, step uint32) string {
	if direction != left && direction != right && direction != top && direction != bottom {
		throw(WrongCursorMovementDirectionError)
	}
	strBuff := strings.Builder{}
	strBuff.WriteString(vT100Basics)
	strBuff.WriteString(strconv.Itoa(int(step)))
	strBuff.WriteByte(direction)
	return strBuff.String()
}

// setCursorPosition sets the cursor position
// @parma x: x-axis position y: y-axis position
// @return the final VT100 style
func setCursorPosition(x, y uint32) string {
	strBuff := strings.Builder{}
	strBuff.WriteString(vT100Basics)
	strBuff.WriteString(strconv.Itoa(int(y)))
	strBuff.WriteByte(';')
	strBuff.WriteString(strconv.Itoa(int(x)))
	strBuff.WriteByte('H')
	return strBuff.String()
}

// CanvasPosition canvas position information, as well as its weight
type CanvasPosition struct {
	X      int                //x-axis position
	Y      int                //y-axis position
	totalX int                //including self x and x offset of all its parent nodes, any component should update this property when it is renderable, this property is crucial for most child nodes
	totalY int                //including self y and y offset of all its parent nodes
	Type   CanvasPositionType //different position types, such as center display
	ZIndex uint32             //weight, the higher the display priority, the higher the weight canvas will cover the lower weight (only valid at the same level)
}

// CanvasPositionType the type represented by position
type CanvasPositionType struct {
	Center uint8
	Right  uint8
}

type NodeStack []Node

// CanvasVolume describes the volume parameters of canvas
type CanvasVolume struct {
	Width  int
	Height int
}

// CanvasStyle describes the style
type CanvasStyle struct {
	Display         bool   //whether to display, not delete
	AutoSize        bool   //adaptive size, its size inherits from the parent element, and will be overwritten by valid values when volume's width\height is not equal to 0
	BorderType      uint8  //whether to display border, and
	BorderColor     string //border color
	Color           string //text color
	BackGroundColor string //background color
	ShowText        bool   //whether to display text
}

// Canvas  main body
type Canvas struct {
	name     string         //used for identification, but not unique
	tag      string         //tag
	key      string         //unique key
	position CanvasPosition //position and weight information
	volume   CanvasVolume   //volume information
	style    CanvasStyle    //style
	parent   Node           //parent node
	children []Node         //child nodes
	unMount  bool           //whether to unmount
	text     string         //the text data that Node needs to render, and should show as much as possible when it can be displayed

}

// CanvasAttr some key information of Canvas
type CanvasAttr struct {
	Name string //Node name
	Tag  string //Node type
	Key  string //Node unique identifier
}

// EventCallBack callback function type
type EventCallBack func(node Node, origen Node)

type Event = map[uint8][]EventCallBack

// Node the common interface that elements need to implement
type Node interface {
	Insert(node ...Node) error                                              //insert elements into the current node
	SetVolume(volume CanvasVolume) error                                    //set the volume by the Volume field
	SetPosition(position CanvasPosition, pType ...CanvasPositionType) error //set the Position field related properties
	SetStyle(style CanvasStyle) error                                       //set the Style related properties
	GetVolume() (CanvasVolume, error)                                       //return the current node's related volume
	GetPosition() (CanvasPosition, error)                                   //get the Position field related information
	GetStyle() (CanvasStyle, error)                                         //get the style
	GetProps() (map[string]string, error)                                   //get the Props
	SetProps(key, value string) error                                       //set the Props custom field
	GetKeyBord() (keyboard.KeyEvent, error)                                 //try to get the keyboard event value of the current node (only accurate when obtained in the event)
	AddEventListener(event uint8, callback EventCallBack) error             //add event listener
	DeleteEventListener(event uint8, callback EventCallBack) error          //delete event listener
	GetAttr() (CanvasAttr, error)                                           //get the Attr field related information
	RemoveChildren(node Node) error                                         //delete a child node
	GetParent() (Node, error)                                               //get the parent node
	GetChildren() ([]Node, error)                                           //get all child nodes
	isUnMount() bool                                                        //check if unmounted
	Remove() error                                                          //self-delete
	setParent(node Node) error                                              //set the parent node of Node
	SetText(text string) error                                              //set text
	setKeyBord(keyboard.KeyEvent)                                           //set key bord
}

// RenderZIndexTree A hierarchy tree generated from the RenderStack
type RenderZIndexTree struct {
	nodeBase   []Node
	originBase []Node
}

// GetNode Gets the next node
func (rt *RenderZIndexTree) GetNode() (Node, bool) {
	if len(rt.nodeBase) == 0 {
		return nil, false
	} else {
		var miniZIndexNode Node = nil
		var miniIndex = 0

		for index, node := range rt.nodeBase {
			nodePosition, _ := node.GetPosition()
			if miniZIndexNode == nil {
				miniZIndexNode = node
				miniIndex = index
			} else if position, _ := miniZIndexNode.GetPosition(); position.ZIndex > nodePosition.ZIndex {
				miniZIndexNode = node
				miniIndex = index
			}
		}
		rt.nodeBase = append(rt.nodeBase[:miniIndex], rt.nodeBase[miniIndex+1:]...)

		return miniZIndexNode, true
	}
}

func (rt *RenderZIndexTree) Init(nodes NodeStack) {
	rt.nodeBase = make(NodeStack, len(nodes))
	rt.originBase = nodes
	copy(rt.nodeBase, nodes)
}

// PaintingTree Interface of the painting tree
type PaintingTree interface {
	GetNode() (Node, bool) // Get the next node
	Init(NodeStack)        // initialize
}
