package tml

import "fmt"

/*
Events that have been supported
	OnMove: Node displacement
	OnSizeChange: Node size change
	OnShow: Node are shown from hiding
	OnHidden: The node is hidden from display
	OnInput: node Text data change
	OnRemove: The node is deleted.
	OnSelect: The selected node can listen for keystroke events using the OnKeyBord event. If the node does not have a corresponding handler, the selected node will be invalid
	OnKeyBord: Triggered when the node is pressed and the keyboard is pressed
*/

// createEvent Add event
// @parma node: The node where the event is to be created
// @return Execution result
func createEvent(node Node) bool {
	attr, _ := node.GetAttr()
	_, ok := eventStore.Load(attr.Key)
	if !ok {
		eventStore.Store(attr.Key, Event{})
	}

	return !ok
}

// deleteEvent Delete event
// @parma node: The node for which the event is to be deleted
// @return Execution result
func deleteEvent(node Node) bool {
	attr, _ := node.GetAttr()
	_, ok := eventStore.Load(attr.Key)
	if ok {
		eventStore.Delete(attr.Key)
	}
	return ok
}

// addEvent Add event to Node
// @parma node: The node to which the event is to be added  eventName: Event type  callback: Event callback
// @return Execution result
func addEvent(node Node, eventName uint8, callback EventCallBack) bool {
	if callback == nil {
		return false
	}
	attr, _ := node.GetAttr()
	nodeEventAny, ok := eventStore.Load(attr.Key)

	nodeEvent := nodeEventAny.(Event)

	if ok {
		callbackStack, ok := nodeEvent[eventName]
		if ok {
			nodeEvent[eventName] = append(callbackStack, callback)
		} else {
			nodeEvent[eventName] = []EventCallBack{callback}
		}
		return true
	}

	return false

}

// untieEvent Unbind event
// @parma node: The node for which the event is to be unbound  eventName: Event type  callback: Event callback
// @return Execution result
func untieEvent(node Node, eventName uint8, callback EventCallBack) bool {
	if callback == nil {
		return false
	}
	attr, _ := node.GetAttr()
	nodeEventAny, ok := eventStore.Load(attr.Key)
	nodeEvent := nodeEventAny.(Event)
	callbackKey := fmt.Sprintf("%v", callback)
	if ok {
		callbackStack, ok := nodeEvent[eventName]
		if ok {
			for index, nodeCB := range callbackStack {
				if callbackKey == fmt.Sprintf("%v", nodeCB) {
					nodeEvent[eventName] = append(callbackStack[:index], callbackStack[index+1:]...)
				}
			}

			if len(nodeEvent[eventName]) == 0 {
				delete(nodeEvent, eventName)
			}

		}
	}
	return false
}

// triggerEvent Trigger event
// @parma node: The Node to which the event is to be triggered  eventName: Event type  origen: Event source node
// @return Execution result
func triggerEvent(node Node, eventName uint8, origen Node) bool {
	attr, _ := node.GetAttr()
	nodeEventAny, ok := eventStore.Load(attr.Key)

	nodeEvent := nodeEventAny.(Event)

	if ok {
		callbackStack, ok := nodeEvent[eventName]
		if ok {
			for _, callBack := range callbackStack {
				if callBack != nil {
					callBack(node, origen)
				}
			}
		}
		return ok
	}
	return false
}

// displayEventTrigger Node Display event notification
// @parma node: Target node  origen: Event source node  display: Event type  deep: Whether to enable in-depth notification
func displayEventTrigger(node Node, origen Node, display bool, deep bool) {
	key := OnShow
	if display == false {
		key = OnHidden
	}

	style, _ := node.GetStyle()

	if style.Display != false { //只有node可显的时候触发事件
		triggerEvent(node, key, origen)

		if deep {
			childNodes, _ := node.GetChildren()
			for _, childNode := range childNodes {
				displayEventTrigger(childNode, origen, display, deep)
			}
		}
	}

}

// autoSizeChangeTrigger Handles nodes with adaptive properties
// @parma node: Target node  origen: Event source node  deep: Whether to enable in-depth notification
func autoSizeChangeTrigger(node Node, origen Node, deep bool) {
	style, _ := node.GetStyle()
	position, _ := node.GetPosition()

	if style.AutoSize {
		triggerEvent(node, OnSizeChange, origen)
	}

	if position.Type.Center > None && position.Type.Center <= PositionXY {
		triggerEvent(node, OnMove, origen)
	}

	if !style.AutoSize && position.Type.Center == None {
		return
	}

	if deep {
		childNodes, _ := node.GetChildren()
		for _, childNode := range childNodes {
			autoSizeChangeTrigger(childNode, origen, deep)
		}
	}
}
