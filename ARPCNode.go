package goarpcsolution

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AnimusPEXUS/gojsonrpc2"
	"github.com/AnimusPEXUS/gorecursionguard"
	"github.com/AnimusPEXUS/gouuidtools"
	"github.com/AnimusPEXUS/utils/anyutils"
	"github.com/mitchellh/mapstructure"
)

const ARPC_MSG_PREFIX_SIMPLE = "simple"
const ARPC_MSG_PREFIX_SIMPLE_PLUST_COLUMN = ARPC_MSG_PREFIX_SIMPLE + ":"
const ARPC_MSG_PREFIX_SIMPLE_PLUST_COLUMN_LEN = len(ARPC_MSG_PREFIX_SIMPLE_PLUST_COLUMN)

const ARPC_MSG_PREFIX_ARPC = "arpc"
const ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN = ARPC_MSG_PREFIX_ARPC + ":"
const ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN_LEN = len(ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN)

// var _ ARPCSolutionCtlI = &ARPCNode{}

// some functions inherited from ARPCSolutionCtlI.
// respective function documentation - placed into interface
type ARPCNode struct {
	PushMessageToOutsideCB func(data []byte) error

	controller ARPCNodeCtlI

	jrpc_node *gojsonrpc2.JSONRPC2Node

	debugName string

	debug bool

	stop_flag bool

	closeRecursionGuard *gorecursionguard.RecursionGuard
}

func NewARPCNode(
	controller ARPCNodeCtlI,
) *ARPCNode {
	self := new(ARPCNode)

	self.closeRecursionGuard = gorecursionguard.NewRecursionGuard(
		gorecursionguard.RGM_SilentReturn,
		nil,
	)

	self.debugName = "ARPCNode"
	self.controller = controller
	self.controller.SetNode(self)

	self.jrpc_node = gojsonrpc2.NewJSONRPC2Node()
	self.jrpc_node.PushMessageToOutsideCB = func(data []byte) error {
		if self.PushMessageToOutsideCB == nil {
			return errors.New("ARPCNode.PushMessageToOutsideCB == nil")
		}
		return self.PushMessageToOutsideCB(data)
	}

	self.jrpc_node.OnRequestCB =
		func(msg *gojsonrpc2.Message) (error, error) {
			if self.debug {
				self.DebugPrintfln("jrpc_node.OnRequestCB msg:", msg)
			}
			return self.handleJRPCNodeMessage(msg)
		}

	return self
}

func (self *ARPCNode) SetDebug(val bool) {
	self.debug = val
	if self.jrpc_node != nil {
		self.jrpc_node.SetDebug(val)
	}
}

func (self *ARPCNode) SetDebugName(name string) {
	self.debugName = fmt.Sprintf("[%s]", name)
}

func (self *ARPCNode) GetDebugName() string {
	return self.debugName
}

func (self *ARPCNode) DebugPrintln(data ...any) {
	fmt.Println(append(append([]any{}, self.debugName), data...)...)
}

func (self *ARPCNode) DebugPrintfln(format string, data ...any) {
	fmt.Println(
		append(
			append([]any{}, self.debugName),
			fmt.Sprintf(format, data...),
		)...,
	)
}

func (self *ARPCNode) nodeInvalidStateException() {
	if self.stop_flag || self.controller == nil {
		panic("node is in invalid state")
	}
}

func (self *ARPCNode) GetController() ARPCNodeCtlI {
	return self.controller
}

// func (self *ARPCNode) IsStopFlag() bool {
// 	return self.stop_flag
// }

// if controller is set - calls it's Close();
// if jrpc2 node is set - calls it's Close();
// sets this node into invalid state.
// node can't be reused after Close() and should be replaced.
func (self *ARPCNode) Close() {
	self.closeRecursionGuard.Do(
		func() {

			self.stop_flag = true

			if self.controller != nil {
				self.controller.Close()
				self.controller = nil
			}

			if self.jrpc_node != nil {
				self.jrpc_node.Close()
				self.jrpc_node = nil
			}
		},
	)

}

// ============ vvvvvvvvvvvvvvvvvvvv ============
// ------------ gojsonrpc2 functions ------------
// ============ vvvvvvvvvvvvvvvvvvvv ============

// sometimes this may be more convinient than ARPC Call

// note: this function always adds "s:" prefix to msg.Method
// note: just like in gojsonrpc2 this function doesn't check msg on
//
//	validity
func (self *ARPCNode) SendMessage(msg *gojsonrpc2.Message) error {
	self.nodeInvalidStateException()
	msg.Method = ARPC_MSG_PREFIX_SIMPLE_PLUST_COLUMN + msg.Method
	return self.jrpc_node.SendMessage(msg)
}

// note: this function always adds "s:" prefix to msg.Method
// note: error if msg invalid
func (self *ARPCNode) SendRequest(
	msg *gojsonrpc2.Message,
	genid bool,
	unhandled bool,
	rh *gojsonrpc2.JSONRPC2NodeRespHandler,
	response_timeout time.Duration,
	request_id_hook *gojsonrpc2.JSONRPC2NodeNewRequestIdHook,
) (ret_any any, ret_err error) {
	self.nodeInvalidStateException()
	err := msg.IsInvalidError()
	if err != nil {
		return nil, err
	}
	if msg.Method != "" {
		msg.Method = ARPC_MSG_PREFIX_SIMPLE_PLUST_COLUMN + msg.Method
	}
	return self.jrpc_node.SendRequest(
		msg, genid, unhandled, rh, response_timeout, request_id_hook,
	)
}

// note: this function always adds "s:" prefix to msg.Method
// note: error if msg invalid
func (self *ARPCNode) SendNotification(msg *gojsonrpc2.Message) error {
	self.nodeInvalidStateException()
	err := msg.IsInvalidError()
	if err != nil {
		return err
	}
	if msg.Method != "" {
		msg.Method = ARPC_MSG_PREFIX_SIMPLE_PLUST_COLUMN + msg.Method
	}
	return self.jrpc_node.SendNotification(msg)
}

func (self *ARPCNode) SendResponse(msg *gojsonrpc2.Message) error {
	self.nodeInvalidStateException()
	err := msg.IsInvalidError()
	if err != nil {
		return err
	}

	return self.jrpc_node.SendResponse(msg)

}

func (self *ARPCNode) SendError(msg *gojsonrpc2.Message) error {
	self.nodeInvalidStateException()
	err := msg.IsInvalidError()
	if err != nil {
		return err
	}

	return self.jrpc_node.SendError(msg)
}

// ============ ^^^^^^^^^^^^^^^^^^^^ ============

// ----------------------------------------
// handle incomming messages
// ----------------------------------------

// #0 protocol violation - not critical for server running,
// #1 error - should be treated as server errors
func (self *ARPCNode) PushMessageFromOutside(data []byte) (error, error) {
	self.nodeInvalidStateException()
	return self.jrpc_node.PushMessageFromOutside(data)
}

func (self *ARPCNode) handleMessage_simple(
	msg *gojsonrpc2.Message,
) (error, error) {
	if self.debug {
		self.DebugPrintln("handleMessage_simple msg:", msg)
	}
	return self.controller.SimpleRequest(msg)
}

// func (self *ARPCNode) handleMessage_arpc(
// 	msg *gojsonrpc2.Message,
// ) (error, error) {
// 	return self.controller.SimpleRequest(msg)
// }

// this function handles actual response to peer node by it self.
// results here are for internal use only.
// 1st value is for protocol (input) error;
// 2nd other errors
func (self *ARPCNode) handleJRPCNodeMessage(
	msg *gojsonrpc2.Message,
) (error, error) {

	if self.debug {
		self.DebugPrintln("handleJRPCNodeMessage msg:", msg)
	}

	if self.controller == nil {
		// TODO: replace with panic?
		return nil,
			errors.New("handling controller undefined")
	}

	if strings.HasPrefix(msg.Method, ARPC_MSG_PREFIX_SIMPLE_PLUST_COLUMN) {
		if self.debug {
			self.DebugPrintln("handleJRPCNodeMessage msg handler: simple")
		}
		msg.Method = msg.Method[ARPC_MSG_PREFIX_SIMPLE_PLUST_COLUMN_LEN:]
		return self.handleMessage_simple(msg)
	}

	if strings.HasPrefix(msg.Method, ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN) {
		if self.debug {
			self.DebugPrintln("handleJRPCNodeMessage msg handler: arpc")
		}

		msg.Method = msg.Method[ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN_LEN:]

		msg_id, msg_has_id := msg.GetId()

		var msg_par map[string]any

		msg_par, ok := (msg.Params).(map[string]any)
		if !ok {
			return errors.New("can't convert msg.Params to map[string]any"),
				errors.New("protocol error")
		}

		var (
			// err_proto error
			result any = nil

			err_code int
			// input error (protocol error) (this is notification. include into log only)
			err_input error
			// error (this is notification. include into log only)
			err_processing_not_internal error
			// error to report in server log
			err_processing_internal error
		)

		// ------------ Notifications ------------

		switch msg.Method {
		default:
			err_code = int(gojsonrpc2.ProtocolErrorMethodNotFound)
			err_input = errors.New("invalid method name")
			err_processing_not_internal = errors.New("protocol error")
			err_processing_internal = errors.New("protocol error")

		case "NewCall":
			call_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"call_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter call_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			call_id_uuid, err := gouuidtools.NewUUIDFromString(call_id)
			if err != nil {
				err_input = err
				break
			}

			response_on, response_on_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"response_on",
				)
			if err != nil {
				err_input = errors.New("possible problem with response_on field")
				err_processing_internal = err
				break
			}

			var response_on_uuid *gouuidtools.UUID
			if response_on_found {
				response_on_uuid, err = gouuidtools.NewUUIDFromString(response_on)
				if err != nil {
					err_input = err
					break
				}
			}

			if !response_on_found {
				go self.controller.NewCall(
					call_id_uuid,
					nil,
				)
			} else {
				go self.controller.NewCall(
					call_id_uuid,
					response_on_uuid,
				)
			}

		// case "NewBuffer":
		// 	buffer_id, not_found, err :=
		// 		anyutils.TraverseObjectTree002_string(
		// 			msg_par,
		// 			true,
		// 			true,
		// 			"buffer_id",
		// 		)

		// 	if not_found {
		// 		err_input = errors.New("not found required parameter buffer_id")
		// 		break
		// 	}

		// 	if err != nil {
		// 		err_processing_internal = err
		// 		break
		// 	}

		// 	buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
		// 	if err != nil {
		// 		err_input = err
		// 		break
		// 	}

		// self.controller.NewBuffer(
		// 	buffer_id_uuid,
		// )

		case "BufferUpdated":
			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			self.controller.BufferUpdated(
				buffer_id_uuid,
			)

		case "NewTransmission":
			tarnsmission_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"tarnsmission_id",
				)

			if not_found {
				err_input =
					errors.New("not found required parameter tarnsmission_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			tarnsmission_id_uuid, err :=
				gouuidtools.NewUUIDFromString(tarnsmission_id)
			if err != nil {
				err_input = err
				break
			}

			self.controller.NewTransmission(
				tarnsmission_id_uuid,
			)

		// case "NewSocket":
		// 	port_id, not_found, err :=
		// 		anyutils.TraverseObjectTree002_string(
		// 			msg_par,
		// 			true,
		// 			true,
		// 			"port_id",
		// 		)

		// 	if not_found {
		// 		err_input = errors.New("not found required parameter port_id")
		// 		break
		// 	}

		// 	if err != nil {
		// 		err_processing_internal = err
		// 		break
		// 	}

		// 	port_id_uuid, err := gouuidtools.NewUUIDFromString(port_id)
		// 	if err != nil {
		// 		err_input = err
		// 		break
		// 	}

		// 	self.controller.NewSocket(
		// 		port_id_uuid,
		// 	)

		// ------------ Methods ------------

		case "CallGetList":
			result, err_processing_not_internal, err_processing_internal =
				self.controller.CallGetList()

		case "CallGetArgCount":
			call_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"call_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter call_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			call_id_uuid, err := gouuidtools.NewUUIDFromString(call_id)
			if err != nil {
				err_input = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.CallGetArgCount(
					call_id_uuid,
				)

		case "CallGetArgValue":
			call_id, not_found, err := anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				true,
				"call_id",
			)

			if not_found {
				err_input = errors.New("not found required parameter call_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			call_id_uuid, err := gouuidtools.NewUUIDFromString(call_id)
			if err != nil {
				err_input = err
				break
			}

			first, not_found, err := anyutils.TraverseObjectTree002_int(
				msg_par,
				true,
				true,
				"first",
			)

			if not_found {
				err_input = errors.New("not found required parameter 'first'")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			if first < 0 {
				err_input = errors.New("first must be >= 0")
				break
			}

			last, not_found, err := anyutils.TraverseObjectTree002_int(
				msg_par,
				true,
				true,
				"first",
			)

			if not_found {
				err_input = errors.New("not found required parameter 'last'")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			if last < first {
				err_input = errors.New("last must be >= first")
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.CallGetArgValues(
					call_id_uuid,
					first, last,
				)

		case "CallClose":
			call_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"call_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter call_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			call_id_uuid, err := gouuidtools.NewUUIDFromString(call_id)
			if err != nil {
				err_input = err
				break
			}

			self.controller.CallClose(
				call_id_uuid,
			)

		case "BufferGetInfo":
			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferGetInfo(
					buffer_id_uuid,
				)

		case "BufferGetItemsCount":
			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferGetItemsCount(
					buffer_id_uuid,
				)

		case "BufferGetItemsIds":
			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			// 1st spec

			first_spec_str, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"first_spec",
				)

			if not_found {
				err_input = errors.New("not found required parameter first_spec")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			first_spec, ty := NewARPCBufferItemSpecifierFromString(first_spec_str)
			if ty == ARPCBufferItemSpecifierTypeInvalid {
				err_input = errors.New("invalid value for first_spec")
				break
			}

			// 2nd spec

			last_spec_str, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"last_spec",
				)

			if not_found {
				err_input = errors.New("not found required parameter last_spec")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			last_spec, ty := NewARPCBufferItemSpecifierFromString(last_spec_str)
			if ty == ARPCBufferItemSpecifierTypeInvalid {
				err_input = errors.New("invalid value for last_spec")
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferGetItemsIds(
					buffer_id_uuid,
					first_spec,
					last_spec,
				)

		case "BufferGetItemsTimesByIds":
			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			ids, not_found, err :=
				anyutils.TraverseObjectTree002_str_list(
					msg_par,
					true,
					true,
					"ids",
				)

			if not_found {
				err_input = errors.New("not found required parameter ids")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferGetItemsTimesByIds(
					buffer_id_uuid,
					ids,
				)

		case "BufferGetItemsByIds":
			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			ids, not_found, err :=
				anyutils.TraverseObjectTree002_str_list(
					msg_par,
					true,
					true,
					"ids",
				)

			if not_found {
				err_input = errors.New("not found required parameter ids")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferGetItemsByIds(
					buffer_id_uuid,
					ids,
				)

		case "BufferGetItemsFirstTime":
			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferGetItemsFirstTime(
					buffer_id_uuid,
				)

		case "BufferGetItemsLastTime":

			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferGetItemsLastTime(
					buffer_id_uuid,
				)

		case "BufferSubscribeOnUpdatesNotification":

			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			err_processing_not_internal, err_processing_internal =
				self.controller.BufferSubscribeOnUpdatesNotification(
					buffer_id_uuid,
				)

			result = err_processing_not_internal == nil &&
				err_processing_internal == nil

		case "BufferUnsubscribeFromUpdatesNotification":

			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			err_processing_not_internal, err_processing_internal =
				self.controller.BufferUnsubscribeFromUpdatesNotification(
					buffer_id_uuid,
				)

			result = err_processing_not_internal == nil &&
				err_processing_internal == nil

		case "BufferGetIsSubscribedOnUpdatesNotification":

			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferGetIsSubscribedOnUpdatesNotification(
					buffer_id_uuid,
				)

		case "BufferGetListSubscribedUpdatesNotifications":

			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferGetIsSubscribedOnUpdatesNotification(
					buffer_id_uuid,
				)

		case "BufferBinaryGetSize":

			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferBinaryGetSize(
					buffer_id_uuid,
				)

		case "BufferBinaryGetSlice":
			buffer_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"buffer_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter buffer_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			buffer_id_uuid, err := gouuidtools.NewUUIDFromString(buffer_id)
			if err != nil {
				err_input = err
			}

			// 1st spec

			start_index, not_found, err :=
				anyutils.TraverseObjectTree002_int(
					msg_par,
					true,
					true,
					"start_index",
				)

			if not_found {
				err_input = errors.New("not found required parameter first_spec")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			// 2nd spec

			end_index, not_found, err :=
				anyutils.TraverseObjectTree002_int(
					msg_par,
					true,
					true,
					"end_index",
				)

			if not_found {
				err_input = errors.New("not found required parameter last_spec")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.BufferBinaryGetSlice(
					buffer_id_uuid,
					start_index,
					end_index,
				)

		case "TransmissionGetList":
			result, err_processing_not_internal, err_processing_internal =
				self.controller.TransmissionGetList()

		case "TransmissionGetInfo":
			transmission_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"transmission_id",
				)

			if not_found {
				err_input =
					errors.New("not found required parameter transmission_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			transmission_id_uuid, err :=
				gouuidtools.NewUUIDFromString(transmission_id)
			if err != nil {
				err_input = err
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.TransmissionGetInfo(
					transmission_id_uuid,
				)

		case "SocketGetList":
			result, err_processing_not_internal, err_processing_internal =
				self.controller.SocketGetList()

		case "SocketOpen":
			listening_socket_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"listening_socket_id",
				)

			if not_found {
				err_input =
					errors.New("not found required parameter listening_socket_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			listening_socket_id_uuid, err :=
				gouuidtools.NewUUIDFromString(listening_socket_id)
			if err != nil {
				err_input = err
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.SocketOpen(
					listening_socket_id_uuid,
				)

		case "SocketRead":
			connected_socket_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"connected_socket_id",
				)

			if not_found {
				err_input =
					errors.New("not found required parameter connected_socket_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			connected_socket_id_uuid, err :=
				gouuidtools.NewUUIDFromString(connected_socket_id)
			if err != nil {
				err_input = err
			}

			try_read_size, not_found, err := anyutils.TraverseObjectTree002_int(
				msg_par,
				true,
				true,
				"try_read_size",
			)

			if not_found {
				err_input =
					errors.New("not found required parameter try_read_size")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			if try_read_size < 0 {
				err_input = errors.New("try_read_size must be >= 0")
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.SocketRead(
					connected_socket_id_uuid,
					try_read_size,
				)

		case "SocketWrite":

			connected_socket_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"connected_socket_id",
				)

			if not_found {
				err_input =
					errors.New("not found required parameter connected_socket_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			connected_socket_id_uuid, err :=
				gouuidtools.NewUUIDFromString(connected_socket_id)
			if err != nil {
				err_input = err
			}

			b, not_found, err := anyutils.TraverseObjectTree002_byte_list(
				msg_par,
				true,
				true,
				"b",
			)

			if not_found {
				err_input = errors.New("not found required parameter b")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			result, err_processing_not_internal, err_processing_internal =
				self.controller.SocketWrite(
					connected_socket_id_uuid,
					b,
				)

		case "SocketClose":
			connected_socket_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"connected_socket_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter connected_socket_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			connected_socket_id_uuid, err :=
				gouuidtools.NewUUIDFromString(connected_socket_id)
			if err != nil {
				err_input = err
			}

			err_processing_not_internal, err_processing_internal =
				self.controller.SocketClose(
					connected_socket_id_uuid,
				)

			result = err_processing_not_internal == nil &&
				err_processing_internal == nil

		case "SocketSetDeadline":
			connected_socket_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"connected_socket_id",
				)

			if not_found {
				err_input =
					errors.New(
						"not found required parameter connected_socket_id",
					)
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			connected_socket_id_uuid, err :=
				gouuidtools.NewUUIDFromString(connected_socket_id)
			if err != nil {
				err_input = err
			}

			t_str, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"t",
				)

			if not_found {
				err_input = errors.New("not found required parameter t")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			t, err := time.Parse(time.RFC3339Nano, t_str)
			if err != nil {
				err_processing_internal = err
				break
			}

			err_processing_not_internal, err_processing_internal =
				self.controller.SocketSetDeadline(
					connected_socket_id_uuid,
					t,
				)

			result = err_processing_not_internal == nil &&
				err_processing_internal == nil

		case "SocketSetReadDeadline":
			connected_socket_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"connected_socket_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter connected_socket_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			connected_socket_id_uuid, err :=
				gouuidtools.NewUUIDFromString(connected_socket_id)
			if err != nil {
				err_input = err
			}

			t_str, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"t",
				)

			if not_found {
				err_input = errors.New("not found required parameter t")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			t, err := time.Parse(time.RFC3339Nano, t_str)
			if err != nil {
				err_processing_internal = err
				break
			}

			err_processing_not_internal, err_processing_internal =
				self.controller.SocketSetReadDeadline(
					connected_socket_id_uuid,
					t,
				)

			result = err_processing_not_internal == nil &&
				err_processing_internal == nil

		case "SocketSetWriteDeadline":
			connected_socket_id, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"connected_socket_id",
				)

			if not_found {
				err_input = errors.New("not found required parameter connected_socket_id")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			connected_socket_id_uuid, err :=
				gouuidtools.NewUUIDFromString(connected_socket_id)
			if err != nil {
				err_input = err
			}

			t_str, not_found, err :=
				anyutils.TraverseObjectTree002_string(
					msg_par,
					true,
					true,
					"t",
				)

			if not_found {
				err_input = errors.New("not found required parameter t")
				break
			}

			if err != nil {
				err_processing_internal = err
				break
			}

			t, err := time.Parse(time.RFC3339Nano, t_str)
			if err != nil {
				err_processing_internal = err
				break
			}

			err_processing_not_internal, err_processing_internal =
				self.controller.SocketSetWriteDeadline(
					connected_socket_id_uuid,
					t,
				)

			result = err_processing_not_internal == nil &&
				err_processing_internal == nil
		}

		if msg_has_id {
			err_processing_internal = self.methodReplyAction(
				msg_id,
				result,
				err_code,
				err_input,
				err_processing_not_internal,
				err_processing_internal,
			)
			if err_processing_internal != nil {
				return nil, err_processing_internal
			}
		}

		return err_input, err_processing_internal
	}

	return errors.New("invalid message format"), errors.New("input error")
}

// note: err_code used only if err_reply and/or err != nil.
// maybe it should be generated by methodReplyAction itself and shouldn't be
// provided by caller
func (self *ARPCNode) methodReplyAction(
	msg_id any,
	result any,
	err_code int,
	err_input error,
	err_processing_not_internal error,
	err_processing_internal error,
) error {
	msg := new(gojsonrpc2.Message)

	defer func() {
		if self.debug {
			self.DebugPrintln("methodReplyAction msg_id:", msg_id)
			self.DebugPrintln("methodReplyAction result:", result)
			self.DebugPrintln("methodReplyAction err_code:", err_code)
			self.DebugPrintln("methodReplyAction err_input:", err_input)
			self.DebugPrintln("methodReplyAction err_processing_not_internal:", err_processing_not_internal)
			self.DebugPrintln("methodReplyAction err_processing_internal:", err_processing_internal)
			self.DebugPrintln("methodReplyAction msg:", msg)
		}
	}()

	msg.SetId(msg_id)

	// err_code = 10 // todo: find better value

	var err error

	if err_input != nil {

	}

	if err_processing_internal != nil {
		e := &gojsonrpc2.JSONRPC2Error{
			Code:    int(gojsonrpc2.ProtocolErrorInternalError),
			Message: "internal server error",
		}
		msg.Error = e
		// note: intentionaly ignoring error from SendError()
		err = self.jrpc_node.SendError(msg)
		if err != nil {
			return err
		}

		return err_processing_internal
	}

	if err_processing_not_internal != nil {
		e := &gojsonrpc2.JSONRPC2Error{
			Code:    err_code,
			Message: err_processing_not_internal.Error(),
		}
		msg.Error = e
		return self.jrpc_node.SendError(msg)
	}

	if err_input != nil {
		if err_code != int(gojsonrpc2.ProtocolErrorMethodNotFound) {
			err_code = int(gojsonrpc2.ProtocolErrorInvalidParams)
		}
		e := &gojsonrpc2.JSONRPC2Error{
			Code:    err_code,
			Message: err_input.Error(),
		}
		msg.Error = e
		return self.jrpc_node.SendError(msg)
	}

	msg.Result = result
	err = self.jrpc_node.SendResponse(msg)
	if err != nil {
		return err
	}

	return nil
}

// ----------------------------------------
// notifications
// ----------------------------------------

func (self *ARPCNode) NewCall(
	call_id *gouuidtools.UUID,
	response_on *gouuidtools.UUID,
) error {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "NewCall"

	params := map[string]any{"call_id": call_id.Format()}

	if response_on != nil && !response_on.IsNil() {
		params["response_on"] = response_on
	}

	msg.Params = params

	return self.jrpc_node.SendNotification(msg)
}

// func (self *ARPCNode) NewBuffer(
// 	buffer_id *gouuidtools.UUID,
// ) error {
// 	msg := new(gojsonrpc2.Message)
// 	msg.Method =ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "NewBuffer"

// 	msg.Params = map[string]any{"buffer_id": buffer_id.Format()}

// 	return self.jrpc_node.SendNotification(msg)
// }

func (self *ARPCNode) BufferUpdated(
	buffer_id *gouuidtools.UUID,
) error {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferUpdated"

	msg.Params = map[string]any{"buffer_id": buffer_id.Format()}

	return self.jrpc_node.SendNotification(msg)
}

func (self *ARPCNode) NewTransmission(
	tarnsmission_id *gouuidtools.UUID,
) error {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "NewTransmission"

	msg.Params = map[string]any{
		"tarnsmission_id": tarnsmission_id.Format(),
	}

	return self.jrpc_node.SendNotification(msg)
}

// func (self *ARPCNode) NewSocket(
// 	listening_socket_id *gouuidtools.UUID,
// ) error {
// 	msg := new(gojsonrpc2.Message)
// 	msg.Method =ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "NewSocket"

// 	msg.Params = map[string]any{
// 		"listening_socket_id": listening_socket_id.Format(),
// 	}

// 	return self.jrpc_node.SendNotification(msg)
// }

// ----------------------------------------
// Basic Calls
// ----------------------------------------

func (self *ARPCNode) subResultGetter01(
	timedout_sig <-chan struct{},
	closed_sig <-chan struct{},
	msg_sig <-chan *gojsonrpc2.Message,
) (result any, timedout bool, closed bool,
	result_err error, err error) {

	select {
	case <-timedout_sig:
		return nil, true, false, nil, nil
	case <-closed_sig:
		return nil, false, true, nil, nil
	case res_msg := <-msg_sig:
		if res_msg.IsInvalid() {
			return nil, false, false,
				errors.New("invalid message"), errors.New("protocol error")
		}
		if !res_msg.HasResponseFields() {
			return nil,
				false, false, nil, errors.New("not a response")
		}
		if res_msg.IsError() {
			return nil,
				false, false, errors.New(res_msg.Error.Message), nil
		}

		return res_msg.Result,
			false, false, nil, errors.New("result must be uuid string")

	}
}

func (self *ARPCNode) CallGetList(
	response_timeout time.Duration,
) (
	buffer_id *gouuidtools.UUID,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "CallGetList"

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		buffer_id = nil
		return
	}

	result_str, ok := result_any.(string)
	if !ok {
		return nil,
			false, false, nil, errors.New("result must be uuid string")
	}

	result, err := gouuidtools.NewUUIDFromString(result_str)
	if err != nil {
		return
	}

	return result, false, false, nil, nil

}

func (self *ARPCNode) CallGetInfo(
	call_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	result *ARPCCallForJSON,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "CallGetInfo"
	msg.Params = map[string]any{"call_id": call_id.Format()}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		result = nil
		return
	}

	err = mapstructure.Decode(result_any, &result)
	if err != nil {
		return nil, false, false, nil, err
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) CallGetArgCount(
	call_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	res int,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "CallGetArgCount"
	msg.Params = map[string]any{"call_id": call_id.Format()}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return 0, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		res = 0
		return
	}

	result, ok := result_any.(int)
	if !ok {
		return 0,
			false, false, nil, errors.New("result must be uuid string")
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) CallGetArgValues(
	call_id *gouuidtools.UUID,
	first, last int,
	response_timeout time.Duration,
) (
	res []*ARPCArgInfo,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "CallGetArgValue"
	msg.Params = map[string]any{
		"call_id": call_id.Format(),
		"first":   first,
		"last":    last,
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		res = nil
		return
	}

	var result []*ARPCArgInfo

	err = mapstructure.Decode(result_any, &result)
	if err != nil {
		return nil, false, false, nil, err
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) CallClose(
	call_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return false, false, nil, err
	}

	_, timedout, closed, result_err, err =
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		return
	}

	return false, false, nil, nil
}

// ----------------------------------------
// Buffers
// ----------------------------------------

func (self *ARPCNode) BufferGetInfo(
	buffer_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	info *ARPCBufferInfo,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferGetInfo"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		info = nil
		return
	}

	var result *ARPCBufferInfo

	err = mapstructure.Decode(result_any, &result)
	if err != nil {
		return nil, false, false, nil, err
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) BufferGetItemsCount(
	buffer_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	count int,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferGetItemsCount"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return 0, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		count = 0
		return
	}

	result, ok := result_any.(int)
	if !ok {
		return 0,
			false, false, nil, errors.New("result must be uuid string")
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) BufferGetItemsIds(
	buffer_id *gouuidtools.UUID,
	first_spec, last_spec *ARPCBufferItemSpecifier,
	// TODO: do we need include_last?
	// include_last bool,
	response_timeout time.Duration,
) (
	ids []string,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferGetItemsIds"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		ids = nil
		return
	}

	result, ok := result_any.([]string)
	if !ok {
		return nil,
			false, false, nil, errors.New("result must be uuid string")
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) BufferGetItemsByIds(
	buffer_id *gouuidtools.UUID,
	ids []string,
	response_timeout time.Duration,
) (
	buffer_items []*ARPCBufferItem,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferGetItemsByIds"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.Format(),
		"ids":       ids,
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		buffer_items = nil
		return
	}

	result_any_slice, ok := result_any.([]any)
	if !ok {
		return nil,
			false, false, nil, errors.New("result must be uuid string")
	}

	result := make([]*ARPCBufferItem, 0)

	for _, i := range result_any_slice {
		var x *ARPCBufferItem
		x = nil
		err = mapstructure.Decode(i, &x)
		if err != nil {
			return nil, false, false, nil, err
		}
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) BufferGetItemsFirstTime(
	buffer_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	time_ time.Time,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferGetItemsFirstTime"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return time.Time{}, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		time_ = time.Time{}
		return
	}

	result_str, ok := result_any.(string)
	if !ok {
		return time.Time{},
			false, false, nil, errors.New("result must be RFC3339Nano time string")
	}

	result, err := time.Parse(result_str, time.RFC3339Nano)
	if err != nil {
		return time.Time{},
			false, false, nil, err
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) BufferGetItemsLastTime(
	buffer_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	time_ time.Time,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferGetItemsIds"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return time.Time{}, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		time_ = time.Time{}
		return
	}

	result_str, ok := result_any.(string)
	if !ok {
		return time.Time{},
			false, false, nil, errors.New("result must be RFC3339Nano time string")
	}

	result, err := time.Parse(result_str, time.RFC3339Nano)
	if err != nil {
		return time.Time{}, false, false, nil, err
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) BufferSubscribeOnUpdatesNotification(
	buffer_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferSubscribeOnUpdatesNotification"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return false, false, nil, err
	}

	_, timedout, closed, result_err, err =
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		return
	}

	return false, false, nil, nil
}

func (self *ARPCNode) BufferUnsubscribeFromUpdatesNotification(
	buffer_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferUnsubscribeFromUpdatesNotification"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return false, false, nil, err
	}

	_, timedout, closed, result_err, err =
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		return
	}

	return false, false, nil, nil
}

func (self *ARPCNode) BufferGetIsSubscribedOnUpdatesNotification(
	buffer_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	result bool,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferGetIsSubscribedOnUpdatesNotification"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return false, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		return
	}

	result, ok := result_any.(bool)
	if !ok {
		return false,
			false, false, nil, errors.New("result must be uuid string")
	}

	return false, false, false, nil, nil
}

func (self *ARPCNode) BufferGetListSubscribedUpdatesNotifications(
	response_timeout time.Duration,
) (
	buffer_id *gouuidtools.UUID,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferGetListSubscribedUpdatesNotifications"

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		buffer_id = nil
		return
	}

	result_str, ok := result_any.(string)
	if !ok {
		return nil,
			false, false, nil, errors.New("result must be uuid string")
	}

	result, err := gouuidtools.NewUUIDFromString(result_str)
	if err != nil {
		return
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) BufferBinaryGetSize(
	buffer_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	size int,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferBinaryGetSize"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return 0, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		size = 0
		return
	}

	result, ok := result_any.(int)
	if !ok {
		return 0, false, false, nil, errors.New("result must be uuid string")
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) BufferBinaryGetSlice(
	buffer_id *gouuidtools.UUID,
	start_index, end_index int,
	response_timeout time.Duration,
) (
	data []byte,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "BufferBinaryGetSlice"
	msg.Params = map[string]any{
		"buffer_id":   buffer_id.Format(),
		"start_index": start_index,
		"end_index":   end_index,
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	_, timedout, closed, result_err, err =
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		buffer_id = nil
		return
	}

	return nil, false, false, nil, nil
}

// ----------------------------------------
// Broadcasts
// ----------------------------------------

func (self *ARPCNode) TransmissionGetList(
	response_timeout time.Duration,
) (
	buffer_id *gouuidtools.UUID,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "TransmissionGetList"

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		buffer_id = nil
		return
	}

	result_str, ok := result_any.(string)
	if !ok {
		return nil,
			false, false, nil, errors.New("result must be uuid string")
	}

	result, err := gouuidtools.NewUUIDFromString(result_str)
	if err != nil {
		return
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) TransmissionGetInfo(
	transmission_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	info *ARPCTransmissionInfo,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "TransmissionGetInfo"
	msg.Params = map[string]any{
		"transmission_id": transmission_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		info = nil
		return
	}

	var result *ARPCTransmissionInfo

	err = mapstructure.Decode(result_any, &result)
	if err != nil {
		return nil, false, false, nil, err
	}

	return result, false, false, nil, nil
}

// ----------------------------------------
// Sockets
// ----------------------------------------

func (self *ARPCNode) SocketGetList(
	response_timeout time.Duration,
) (
	buffer_id *gouuidtools.UUID,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "SocketGetList"

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		buffer_id = nil
		return
	}

	result_str, ok := result_any.(string)
	if !ok {
		return nil,
			false, false, nil, errors.New("result must be uuid string")
	}

	result, err := gouuidtools.NewUUIDFromString(result_str)
	if err != nil {
		return
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) SocketOpen(
	listening_socket_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	connected_socket_id *gouuidtools.UUID,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "SocketOpen"
	msg.Params = map[string]any{
		"listening_socket_id": listening_socket_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		connected_socket_id = nil
		return
	}

	result_str, ok := result_any.(string)
	if !ok {
		return nil,
			false, false, nil, errors.New("result must be uuid string")
	}

	result, err := gouuidtools.NewUUIDFromString(result_str)
	if err != nil {
		return
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) SocketRead(
	connected_socket_id *gouuidtools.UUID,
	try_read_size int,
	response_timeout time.Duration,
) (
	b []byte,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "SocketRead"
	msg.Params = map[string]any{
		"connected_socket_id": connected_socket_id.Format(),
		"try_read_size":       try_read_size,
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		b = nil
		return
	}

	result, ok := result_any.([]byte)
	if !ok {
		return nil,
			false, false, nil, errors.New("result must be uuid string")
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) SocketWrite(
	connected_socket_id *gouuidtools.UUID,
	b []byte,
	response_timeout time.Duration,
) (
	n int,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "SocketWrite"
	msg.Params = map[string]any{
		"connected_socket_id": connected_socket_id.Format(),
		"b":                   b,
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return 0, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		n = 0
		return
	}

	result, ok := result_any.(int)
	if !ok {
		return 0,
			false, false, nil, errors.New("result must be uuid string")
	}

	return result, false, false, nil, nil
}

func (self *ARPCNode) SocketClose(
	connected_socket_id *gouuidtools.UUID,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	self.nodeInvalidStateException()
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "CallClose"
	msg.Params = map[string]any{
		"connected_socket_id": connected_socket_id.Format(),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return false, false, nil, err
	}

	_, timedout, closed, result_err, err =
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		return
	}

	return false, false, nil, nil
}

func (self *ARPCNode) SocketSetDeadline(
	connected_socket_id *gouuidtools.UUID,
	t time.Time,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "SocketSetDeadline"
	msg.Params = map[string]any{
		"connected_socket_id": connected_socket_id.Format(),
		"t":                   t.Format(time.RFC3339Nano),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return false, false, nil, err
	}

	_, timedout, closed, result_err, err =
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		return
	}

	return false, false, nil, nil
}

func (self *ARPCNode) SocketSetReadDeadline(
	connected_socket_id *gouuidtools.UUID,
	t time.Time,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "SocketSetReadDeadline"
	msg.Params = map[string]any{
		"connected_socket_id": connected_socket_id.Format(),
		"t":                   t.Format(time.RFC3339Nano),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return false, false, nil, err
	}

	_, timedout, closed, result_err, err =
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		return
	}

	return false, false, nil, nil
}

func (self *ARPCNode) SocketSetWriteDeadline(
	connected_socket_id *gouuidtools.UUID,
	t time.Time,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = ARPC_MSG_PREFIX_ARPC_PLUS_COLUMN + "SocketSetWriteDeadline"
	msg.Params = map[string]any{
		"connected_socket_id": connected_socket_id.Format(),
		"t":                   t.Format(time.RFC3339Nano),
	}

	timedout_sig, closed_sig, msg_sig, rh :=
		gojsonrpc2.NewChannelledJSONRPC2NodeRespHandler()

	_, err = self.jrpc_node.SendRequest(
		msg,
		true,
		false,
		rh,
		response_timeout,
		nil,
	)
	if err != nil {
		return false, false, nil, err
	}

	_, timedout, closed, result_err, err =
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		return
	}

	return false, false, nil, nil
}
