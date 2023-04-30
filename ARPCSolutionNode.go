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

// some functions inherited from ARPCSolutionCtlI.
// respective function documentation - placed into interface
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

	self.jrpc_node.OnRequestCB = func(msg *gojsonrpc2.Message) (error, error) {
		e1, _, e3 := self.PushMessageFromOutside(msg)
		// TODO: maybe upward functions have to be upgraded to receive 3 errors
		return e1, e3
	}

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

// this function handles actual response to peer node by it self.
// results here are for internal use only.
// 1st value is for protocol (input) error;
// 2nd error is error which returned to peer;
// 3rd error is internal server error not intended to be returned to peer.
func (self *ARPCSolutionNode) PushMessageFromOutside(
	msg *gojsonrpc2.Message,
) (error, error, error) {

	if debug {
		self.DebugPrintln("PushMessageFromOutside()")
	}

	msg_id, msg_has_id := msg.GetId()

	if self.controller == nil {
		return nil,
			nil,
			errors.New("handling controller undefined")
	}

	var msg_par map[string]any

	msg_par, ok := (msg.Params).(map[string]any)
	if !ok {
		return errors.New("can't convert msg.Params to map[string]any"),
			errors.New("protocol error"),
			errors.New("protocol error")
	}

	var (
		// err_proto error
		result any = nil
		// todo: do something with err_code
		// err_code  int

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

		if !response_on_found {
			go self.controller.NewCall(
				uuid.FromStringOrNil(call_id),
				uuid.Nil,
			)
		} else {
			go self.controller.NewCall(
				uuid.FromStringOrNil(call_id),
				uuid.FromStringOrNil(response_on),
			)
		}
	case "NewBuffer":
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

		self.controller.NewBuffer(
			uuid.FromStringOrNil(buffer_id),
		)
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

		self.controller.BufferUpdated(
			uuid.FromStringOrNil(buffer_id),
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
			err_input = errors.New("not found required parameter tarnsmission_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		self.controller.NewTransmission(
			uuid.FromStringOrNil(tarnsmission_id),
		)
	case "NewSocket":
		port_id, not_found, err :=
			anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				true,
				"port_id",
			)

		if not_found {
			err_input = errors.New("not found required parameter port_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		self.controller.NewSocket(
			uuid.FromStringOrNil(port_id),
		)

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

		result, err_processing_not_internal, err_processing_internal =
			self.controller.CallGetArgCount(
				uuid.FromStringOrNil(call_id),
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

		arg_index, not_found, err := anyutils.TraverseObjectTree002_int(
			msg_par,
			true,
			true,
			"arg_index",
		)

		if not_found {
			err_input = errors.New("not found required parameter call_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		if arg_index < 0 {
			err_input = errors.New("arg_index must be >= 0")
			break
		}

		result, err_processing_not_internal, err_processing_internal =
			self.controller.CallGetArgValue(
				uuid.FromStringOrNil(call_id),
				arg_index,
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

		self.controller.CallClose(
			uuid.FromStringOrNil(call_id),
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

		result, err_processing_not_internal, err_processing_internal =
			self.controller.BufferGetInfo(
				uuid.FromStringOrNil(buffer_id),
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

		result, err_processing_not_internal, err_processing_internal =
			self.controller.BufferGetItemsCount(
				uuid.FromStringOrNil(buffer_id),
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
				uuid.FromStringOrNil(buffer_id),
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
				uuid.FromStringOrNil(buffer_id),
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

		result, err_processing_not_internal, err_processing_internal =
			self.controller.BufferGetItemsFirstTime(
				uuid.FromStringOrNil(buffer_id),
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

		result, err_processing_not_internal, err_processing_internal =
			self.controller.BufferGetItemsLastTime(
				uuid.FromStringOrNil(buffer_id),
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

		err_processing_not_internal, err_processing_internal =
			self.controller.BufferSubscribeOnUpdatesNotification(
				uuid.FromStringOrNil(buffer_id),
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

		err_processing_not_internal, err_processing_internal =
			self.controller.BufferUnsubscribeFromUpdatesNotification(
				uuid.FromStringOrNil(buffer_id),
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

		result, err_processing_not_internal, err_processing_internal =
			self.controller.BufferGetIsSubscribedOnUpdatesNotification(
				uuid.FromStringOrNil(buffer_id),
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

		result, err_processing_not_internal, err_processing_internal =
			self.controller.BufferGetIsSubscribedOnUpdatesNotification(
				uuid.FromStringOrNil(buffer_id),
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

		result, err_processing_not_internal, err_processing_internal =
			self.controller.BufferBinaryGetSize(
				uuid.FromStringOrNil(buffer_id),
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
				uuid.FromStringOrNil(buffer_id),
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
			err_input = errors.New("not found required parameter transmission_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		result, err_processing_not_internal, err_processing_internal =
			self.controller.TransmissionGetInfo(
				uuid.FromStringOrNil(transmission_id),
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
			err_input = errors.New("not found required parameter listening_socket_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		result, err_processing_not_internal, err_processing_internal =
			self.controller.SocketOpen(
				uuid.FromStringOrNil(listening_socket_id),
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
			err_input = errors.New("not found required parameter connected_socket_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		try_read_size, not_found, err := anyutils.TraverseObjectTree002_int(
			msg_par,
			true,
			true,
			"try_read_size",
		)

		if not_found {
			err_input = errors.New("not found required parameter try_read_size")
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
				uuid.FromStringOrNil(connected_socket_id),
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
			err_input = errors.New("not found required parameter connected_socket_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
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
				uuid.FromStringOrNil(connected_socket_id),
				b,
			)

	case "SocketClose":
		listening_socket_id, not_found, err :=
			anyutils.TraverseObjectTree002_string(
				msg_par,
				true,
				true,
				"listening_socket_id",
			)

		if not_found {
			err_input = errors.New("not found required parameter listening_socket_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		err_processing_not_internal, err_processing_internal =
			self.controller.SocketClose(
				uuid.FromStringOrNil(listening_socket_id),
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
			err_input = errors.New("not found required parameter connected_socket_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
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
				uuid.FromStringOrNil(connected_socket_id),
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
				uuid.FromStringOrNil(connected_socket_id),
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
				uuid.FromStringOrNil(connected_socket_id),
				t,
			)

		result = err_processing_not_internal == nil &&
			err_processing_internal == nil
	}

	if msg_has_id {
		err_processing_internal = self.methodReplyAction(
			msg_id,
			result,
			// todo: fix code
			10,
			err_input,
			err_processing_not_internal,
			err_processing_internal,
		)
		if err_processing_internal != nil {
			return nil, nil, err_processing_internal
		}
	}

	return err_input, err_processing_not_internal, err_processing_internal
}

// note: err_code used only if err_reply and/or err != nil.
// maybe it should be generated by methodReplyAction itself and shouldn't be
// provided by caller
func (self *ARPCSolutionNode) methodReplyAction(
	msg_id any,
	result any,
	err_code int,
	err_input error,
	err_processing_not_internal error,
	err_processing_internal error,
) error {
	msg := new(gojsonrpc2.Message)
	msg.SetId(msg_id)

	err_code = 10 // todo: find better value

	var err error

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

func (self *ARPCSolutionNode) NewCall(
	call_id uuid.UUID,
	response_on uuid.UUID,
)

func (self *ARPCSolutionNode) NewBuffer(id uuid.UUID)

func (self *ARPCSolutionNode) NewBroadcast(id uuid.UUID)

func (self *ARPCSolutionNode) NewSocket(id uuid.UUID)

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
// Sockets
// ----------------------------------------

func (self *ARPCSolutionNode) SocketsGetListeningIdList() (
	[]uuid.UUID,
	error,
)

func (self *ARPCSolutionNode) SocketOpen(id uuid.UUID) error

func (self *ARPCSolutionNode) SocketRead(id uuid.UUID, b []byte) (n int, err error)
func (self *ARPCSolutionNode) SocketWrite(id uuid.UUID, b []byte) (n int, err error)

func (self *ARPCSolutionNode) SocketClose(id uuid.UUID) error
