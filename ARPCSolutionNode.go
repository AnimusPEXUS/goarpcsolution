package goarpcsolution

import (
	"errors"
	"fmt"
	"time"

	"github.com/AnimusPEXUS/gojsonrpc2"
	"github.com/AnimusPEXUS/utils/anyutils"
	uuid "github.com/satori/go.uuid"
)

var _ ARPCSolutionCtlI = &ARPCSolutionNode{}

type ARPCSolutionNode struct {
	PushMessageToOutsideCB func(data []byte) error

	controller ARPCSolutionCtlI

	jrpc_node *gojsonrpc2.JSONRPC2Node

	debugName string
}

func NewARPCSolutionNode(
	controller ARPCSolutionCtlI,
) *ARPCSolutionNode {
	self := new(ARPCSolutionNode)
	self.debugName = "ARPCSolutionNode"
	self.controller = controller

	self.jrpc_node = gojsonrpc2.NewJSONRPC2Node()

	self.jrpc_node.OnRequestCB = self.PushMessageFromOutside

	return self
}

func (self *ARPCSolutionNode) SetDebugName(name string) {
	self.debugName = fmt.Sprintf("[%s]", name)
}

func (self *ARPCSolutionNode) GetDebugName() string {
	return self.debugName
}

func (self *ARPCSolutionNode) DebugPrintln(data ...any) {
	fmt.Println(append(append([]any{}, self.debugName), data...)...)
}

func (self *ARPCSolutionNode) DebugPrintfln(format string, data ...any) {
	fmt.Println(append(append([]any{}, self.debugName), fmt.Sprintf(format, data...))...)
}

// registers new call to counterpart and sends notification about this.
// use args to pass arguments. named and nameless arguments allowed to be intermixed.
// if multiple arguments with same name present only lattest one is accounted.
func (self *ARPCSolutionNode) CreateFuncCall(
	name string,
	args ...ARPCSolutionFuncArg,
) *ARPCSolutionFuncRes {
	return nil
}

// ----------------------------------------
// handle incomming messages
// ----------------------------------------

func (self *ARPCSolutionNode) PushMessageFromOutside(
	msg *gojsonrpc2.Message,
) (error, error) {

	if debug {
		self.DebugPrintln("PushMessageFromOutside()")
	}

	if self.controller == nil {
		return nil,
			errors.New("handling controller undefined")
	}

	var msg_par map[string]any

	msg_par, ok := (msg.Params).(map[string]any)
	if !ok {
		return errors.New("can't convert msg.Params to map[string]any"),
			errors.New("protocol error")
	}

	msg_method := msg.Method

	msg_id, msg_has_id := msg.GetId()
	if !msg_has_id {

		switch msg_method {
		default:
			return errors.New("invalid call name"),
				errors.New("protocol error")
		case "NewCall":
			call_id, _, err := anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				false,
				"call_id",
			)
			if err != nil {
				return errors.New("new call_id required"),
					errors.New("protocol error")
			}

			response_on, response_on_found, err := anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				true,
				"response_on",
			)
			if err != nil {
				// TODO: better handling required
				return errors.New("possible protocol error"),
					errors.New("problem with response_on field")
			}

			if !response_on_found {
				self.NewCall(uuid.FromStringOrNil(call_id), uuid.Nil)
			} else {
				go self.controller.NewCall(
					uuid.FromStringOrNil(call_id),
					uuid.FromStringOrNil(response_on),
				)
			}
		case "NewBuffer":
			buffer_id, _, err := anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				false,
				"buffer_id",
			)
			if err != nil {
				return errors.New("new buffer id required"),
					errors.New("protocol error")
			}

			self.controller.NewBuffer(
				uuid.FromStringOrNil(buffer_id),
			)
		case "NewBroadcast":
			broadcast_id, _, err := anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				false,
				"broadcast_id",
			)
			if err != nil {
				return errors.New("new broadcast id required"),
					errors.New("protocol error")
			}

			self.controller.NewBroadcast(
				uuid.FromStringOrNil(broadcast_id),
			)
		case "NewPort":
			port_id, _, err := anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				false,
				"port_id",
			)
			if err != nil {
				return errors.New("new port id required"),
					errors.New("protocol error")
			}

			self.controller.NewPort(
				uuid.FromStringOrNil(port_id),
			)

		}

		// NOTE: notification doesn't expect to get any reply from us

		return nil, nil
	} else {
		var (
			// err_proto error
			result    any
			err_code  int
			err_reply error
			err_ret   error
		)
		switch msg_method {
		default:
			return errors.New("invalid method name"),
				errors.New("protocol error")
		case "CallGetList":
			result, err_reply, err_ret = self.controller.CallGetList()

		case "CallGetArgCount":
			call_id_str, _, err := anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				false,
				"call_id",
			)

			err_ret = err

			if err_ret != nil {
				goto return_result2
			}

			result, err_reply, err_ret = self.controller.CallGetArgCount(
				uuid.FromStringOrNil(call_id_str),
			)

		case "CallGetArgValue":
			call_id_str, _, err := anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				false,
				"call_id",
			)

			err_ret = err

			if err_ret != nil {
				goto return_result2
			}

			arg_index, _, err := anyutils.TraverseObjectTree002_uint(
				msg_par,
				true,
				false,
				"arg_index",
			)

			err_ret = err

			if err_ret != nil {
				goto return_result2
			}

			result, err_reply, err_ret = self.controller.CallGetArgValue(
				uuid.FromStringOrNil(call_id_str),
				arg_index,
			)

		case "CallClose":
			call_id_str, _, err := anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				false,
				"call_id",
			)

			err_ret = err

			if err_ret != nil {
				goto return_result2
			}

			err_reply, err_ret = self.controller.CallClose(
				uuid.FromStringOrNil(call_id_str),
			)

			err_ret = err

			if err_ret != nil {
				goto return_result2
			}

			result = "ok"

		case "BufferGetInfo":
		case "BufferGetItemsFirstTime":
		case "BufferGetItemsLastTime":
		case "BufferGetItemsByIds":
		case "BufferGetItemsStartingFromId":
		case "BufferGetItemsByFirstAndLastId":
		case "BufferGetItemsByFirstAndBeforeId":
		case "BufferGetItemsByFirstAndLastTime":
		case "BufferGetItemsByFirstAndBeforeTime":
		case "BufferSubscribeOnUpdatesIndicator":
		case "BroadcastGetIdList":
		case "BroadcastGetInfo":
		case "BroadcastSubscribeOnNew":
		case "PortGetListeningIdList":
		case "PortOpen":
		case "PortRead":
		case "PortWrite":
		case "PortClose":
		}

	return_result2:

		err_ret = self.methodReplyAction(
			msg_id,
			result,
			// todo: fix code
			1,
			err_reply,
			err_ret,
		)
		if err_ret != nil {
			return nil, err_ret
		}

		return nil, nil
		// if err_proto != nil {
		// 	return err_proto, errors.New("protocol error")
		// }

	}
}

// note: err_code used only if err_reply and/or err != nil.
// maybe it should be generated by methodReplyAction itself and shouldn't be
// provided by caller
func (self *ARPCSolutionNode) methodReplyAction(
	msg_id any,
	result any,
	err_code int,
	err_reply error,
	err error,
) error {
	msg := new(gojsonrpc2.Message)

	if err != nil {
		e := &gojsonrpc2.JSONRPC2Error{
			Code:    err_code,
			Message: err_reply.Error(),
		}
		msg.Error = e
		msg.SetId(msg_id)
		// note: intentionaly ignoring error from SendError()
		self.jrpc_node.SendError(msg)

		return err
	}

	if err_reply != nil {
		e := &gojsonrpc2.JSONRPC2Error{
			Code:    err_code,
			Message: err_reply.Error(),
		}
		msg.Error = e
		msg.SetId(msg_id)
		return self.jrpc_node.SendError(msg)
	}

	msg.Result = result

	return nil
}

// ----------------------------------------
// notifications
// ----------------------------------------

func (self *ARPCSolutionNode) NewCall(
	call_id uuid.UUID,
	response_on uuid.UUID,
)

func (self *ARPCSolutionNode) NewBuffer(id uuid.UUID)

func (self *ARPCSolutionNode) NewBroadcast(id uuid.UUID)

func (self *ARPCSolutionNode) NewPort(id uuid.UUID)

// ----------------------------------------
// Basic Calls
// ----------------------------------------

func (self *ARPCSolutionNode) CallGetList() (buffer_id uuid.UUID, err error)

func (self *ARPCSolutionNode) CallGetArgCount(call_id uuid.UUID) (int, error)

func (self *ARPCSolutionNode) CallGetArgValue(call_id uuid.UUID, index int) (
	*ARPCArgInfo,
	error,
)

func (self *ARPCSolutionNode) CallClose(call_id uuid.UUID)

// ----------------------------------------
// Buffers
// ----------------------------------------

func (self *ARPCSolutionNode) GetBufferInfo(id uuid.UUID) (*ARPCBufferInfo, error)

func (self *ARPCSolutionNode) GetBufferItemsFirstTime(id uuid.UUID) (time.Time, error)

func (self *ARPCSolutionNode) GetBufferItemsLastTime(id uuid.UUID) (time.Time, error)

func (self *ARPCSolutionNode) GetBufferItemsByIds(ids []uuid.UUID) (
	buffer_items []*ARPCBufferItem,
	err error,
)

func (self *ARPCSolutionNode) GetBufferItemsStartingFromId(id uuid.UUID) (
	buffer_items []*ARPCBufferItem,
	err error,
)

func (self *ARPCSolutionNode) GetBufferItemsByFirstLastId(
	first_id uuid.UUID,
	last_id uuid.UUID,
) (
	buffer_items []*ARPCBufferItem,
	err error,
)

func (self *ARPCSolutionNode) GetBufferItemsByFirstLastIdExcludingly(
	first_id uuid.UUID,
	last_id uuid.UUID,
) (
	buffer_items []*ARPCBufferItem,
	err error,
)

func (self *ARPCSolutionNode) GetBufferItemsByFirstLastTime(
	first_time time.Time,
	last_time time.Time,
) (
	buffer_items []*ARPCBufferItem,
	err error,
)

func (self *ARPCSolutionNode) GetBufferItemsByFirstLastTimeExcludingly(
	first_time time.Time,
	last_time time.Time,
) (
	buffer_items []*ARPCBufferItem,
	err error,
)

func (self *ARPCSolutionNode) BufferSubscribeOnUpdatesIndicator(buffer_id uuid.UUID) error

// ----------------------------------------
// Broadcasts
// ----------------------------------------

func (self *ARPCSolutionNode) BroadcastIdList() (
	[]uuid.UUID,
	error,
)

func (self *ARPCSolutionNode) GetBroadcastInfo(broadcast_id uuid.UUID) (
	*ARPCBroadcastInfo,
	error,
)

func (self *ARPCSolutionNode) SubscribeOnNewBroadcasts() error

// ----------------------------------------
// Ports
// ----------------------------------------

func (self *ARPCSolutionNode) PortsGetListeningIdList() (
	[]uuid.UUID,
	error,
)

func (self *ARPCSolutionNode) PortOpen(id uuid.UUID) error

func (self *ARPCSolutionNode) PortRead(id uuid.UUID, b []byte) (n int, err error)
func (self *ARPCSolutionNode) PortWrite(id uuid.UUID, b []byte) (n int, err error)

func (self *ARPCSolutionNode) PortClose(id uuid.UUID) error
