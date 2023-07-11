package tml

import (
	"errors"
	"github.com/eiannone/keyboard"
	"reflect"
	"sync"
)

// Quadrilateral a square Canvas
type Quadrilateral struct {
	Canvas                     //inherited struct
	keyboard.KeyEvent          //store the data generated by key bord events
	props             sync.Map //props that store custom information
}

func (ql *Quadrilateral) SetProps(key, value string) error {
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}

	ql.props.Store(key, value)
	return nil
}

func (ql *Quadrilateral) setKeyBord(event keyboard.KeyEvent) {
	ql.KeyEvent = event
}

func (ql *Quadrilateral) GetKeyBord() (keyboard.KeyEvent, error) {
	if ql.unMount {
		return ql.KeyEvent, errors.New(OperatingEmptyNodeError)
	}

	return ql.KeyEvent, nil
}

func (ql *Quadrilateral) SetText(text string) error {
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}
	if ql.text != text {
		Render()
		triggerEvent(ql, OnInput, ql)
	}
	ql.text = text
	return nil
}

func (ql *Quadrilateral) Insert(node ...Node) error {
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}
	ql.children = append(ql.children, node...)
	for _, childNode := range node {
		childNode.setParent(ql)
	}
	Render()
	return nil
}

func (ql *Quadrilateral) setParent(node Node) error {
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}

	oldParent := ql.parent
	if oldParent != nil && node != nil {
		oldStyle, _ := oldParent.GetStyle()
		newStyle, _ := node.GetStyle()
		if oldStyle.Display != newStyle.Display {
			displayEventTrigger(ql, node, newStyle.Display, true)
		}
	} else if oldParent != nil && node == nil {
		displayEventTrigger(ql, nil, false, true)
	} else if oldParent == nil && node != nil {
		newStyle, _ := node.GetStyle()
		if newStyle.Display {
			displayEventTrigger(ql, node, true, true)
		}
	}

	ql.parent = node

	Render()

	if node == nil {
		return errors.New(ParentNodeNil)
	}

	return nil
}

func (ql *Quadrilateral) SetVolume(volume CanvasVolume) error {
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}

	if ql.volume.Width != volume.Width || ql.volume.Height != volume.Height {
		triggerEvent(ql, OnSizeChange, ql)
		Render()
		for _, child := range ql.children { // Handle adaptive numeric events in child nodes
			autoSizeChangeTrigger(child, ql, true)
		}
	}

	ql.volume = volume
	return nil
}

func (ql *Quadrilateral) SetPosition(position CanvasPosition, pType ...CanvasPositionType) error {
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}

	if position.X != ql.position.X || position.Y != ql.position.Y { // Attempt to trigger an event
		triggerEvent(ql, OnMove, ql)
		for _, child := range ql.children { // Handle adaptive numeric events in child nodes
			autoSizeChangeTrigger(child, ql, true)
		}
	}

	for _, child := range ql.children {
		autoSizeChangeTrigger(child, ql, true)
	}

	oldPosition := ql.position

	ql.position = CanvasPosition{
		X:      position.X,
		Y:      position.Y,
		totalX: ql.position.totalX,
		totalY: ql.position.totalY,
		Type:   position.Type,
		ZIndex: position.ZIndex,
	}

	typeLen := len(pType)

	if typeLen > 0 {
		ql.position.Type = pType[typeLen-1]
	}

	if !reflect.DeepEqual(oldPosition, ql.position) {
		Render()
	}

	return nil
}

func (ql *Quadrilateral) SetStyle(style CanvasStyle) error {
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}
	if style.Display != ql.style.Display {
		displayEventTrigger(ql, ql, style.Display, true) // Trigger show and hide events and trigger recursively
	}

	oldStyle := ql.style
	ql.style = style
	if !reflect.DeepEqual(oldStyle, style) {
		Render()
	}

	return nil
}

func (ql *Quadrilateral) GetVolume() (CanvasVolume, error) {
	if ql.unMount {
		return ql.volume, errors.New(OperatingEmptyNodeError)
	}
	return ql.volume, nil
}

func (ql *Quadrilateral) GetPosition() (CanvasPosition, error) {
	if ql.unMount {
		return ql.position, errors.New(OperatingEmptyNodeError)
	}
	return ql.position, nil
}

func (ql *Quadrilateral) GetStyle() (CanvasStyle, error) {

	if ql.unMount {
		return ql.style, errors.New(OperatingEmptyNodeError)
	}
	return ql.style, nil
}

func (ql *Quadrilateral) GetProps() (map[string]string, error) {
	props := make(map[string]string)

	if ql.unMount {
		return props, errors.New(OperatingEmptyNodeError)
	}

	ql.props.Range(func(key, value any) bool {
		props[key.(string)] = value.(string)
		return true
	})

	return props, nil
}

func (ql *Quadrilateral) AddEventListener(event uint8, callback EventCallBack) error {
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}
	addEvent(ql, event, callback)
	return nil
}

func (ql *Quadrilateral) DeleteEventListener(event uint8, callback EventCallBack) error {
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}
	untieEvent(ql, event, callback)
	return nil
}

func (ql *Quadrilateral) GetChildren() ([]Node, error) {
	if ql.unMount {
		return ql.children, errors.New(OperatingEmptyNodeError)
	}
	return ql.children, nil
}

func (ql *Quadrilateral) GetParent() (Node, error) {
	if ql.unMount {
		return ql.parent, errors.New(OperatingEmptyNodeError)
	}
	return ql.parent, nil
}

func (ql *Quadrilateral) RemoveChildren(node Node) error {
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}
	attr, _ := node.GetAttr()
	for index, qlNode := range ql.children {
		if childAttr, _ := qlNode.GetAttr(); childAttr.Key == attr.Key {
			ql.children = append(ql.children[:index], ql.children[index+1:]...)
			if !qlNode.isUnMount() {
				qlNode.Remove()
			}
		}
	}
	return nil
}

func (ql *Quadrilateral) Remove() error { // Trigger show and hide events and trigger recursively
	if ql.unMount {
		return errors.New(OperatingEmptyNodeError)
	}

	if ql == Body {
		return errors.New(DeleteTopLeveNodeError)
	}
	nodeParent, _ := ql.GetParent()
	delNodeFromRenderStack(ql)
	delNodeFromNameIndex(ql)
	delNodeFromBase(ql)
	deleteEvent(ql)
	triggerEvent(ql, OnRemove, ql)
	if nodeParent != nil && !nodeParent.isUnMount() {
		nodeParent.RemoveChildren(ql)
	}
	ql.unMount = true
	return nil
}

func (ql *Quadrilateral) isUnMount() bool {
	return ql.unMount
}

func (ql *Quadrilateral) GetAttr() (CanvasAttr, error) {
	if ql.unMount {
		return CanvasAttr{}, errors.New(OperatingEmptyNodeError)
	}
	return CanvasAttr{
		Name: ql.name,
		Tag:  ql.tag,
		Key:  ql.key,
	}, nil
}
