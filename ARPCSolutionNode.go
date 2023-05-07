package goarpcsolution

import (
	"errors"
	"fmt"
	"time"

	"github.com/AnimusPEXUS/gojsonrpc2"
	"github.com/AnimusPEXUS/utils/anyutils"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

// var _ ARPCSolutionCtlI = &ARPCSolutionNode{}

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

		call_id_uuid := uuid.FromStringOrNil(call_id)
		if uuid.Equal(uuid.Nil, call_id_uuid) {
			err_input = errors.New("invalid call_id")
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

		var response_on_uuid uuid.UUID
		if response_on_found {
			response_on_uuid = uuid.FromStringOrNil(response_on)
			if uuid.Equal(uuid.Nil, response_on_uuid) {
				err_input = errors.New("invalid response_on")
				break
			}
		}

		if !response_on_found {
			go self.controller.NewCall(
				call_id_uuid,
				uuid.Nil,
			)
		} else {
			go self.controller.NewCall(
				call_id_uuid,
				response_on_uuid,
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
			break
		}

		self.controller.NewBuffer(
			buffer_id_uuid,
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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
			err_input = errors.New("not found required parameter tarnsmission_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		tarnsmission_id_uuid := uuid.FromStringOrNil(tarnsmission_id)
		if uuid.Equal(uuid.Nil, tarnsmission_id_uuid) {
			err_input = errors.New("invalid tarnsmission_id")
			break
		}

		self.controller.NewTransmission(
			tarnsmission_id_uuid,
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

		port_id_uuid := uuid.FromStringOrNil(port_id)
		if uuid.Equal(uuid.Nil, port_id_uuid) {
			err_input = errors.New("invalid port_id")
			break
		}

		self.controller.NewSocket(
			port_id_uuid,
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

		call_id_uuid := uuid.FromStringOrNil(call_id)
		if uuid.Equal(uuid.Nil, call_id_uuid) {
			err_input = errors.New("invalid call_id")
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

		call_id_uuid := uuid.FromStringOrNil(call_id)
		if uuid.Equal(uuid.Nil, call_id_uuid) {
			err_input = errors.New("invalid call_id")
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
				call_id_uuid,
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

		call_id_uuid := uuid.FromStringOrNil(call_id)
		if uuid.Equal(uuid.Nil, call_id_uuid) {
			err_input = errors.New("invalid call_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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

		buffer_id_uuid := uuid.FromStringOrNil(buffer_id)
		if uuid.Equal(uuid.Nil, buffer_id_uuid) {
			err_input = errors.New("invalid buffer_id")
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
			err_input = errors.New("not found required parameter transmission_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		transmission_id_uuid := uuid.FromStringOrNil(transmission_id)
		if uuid.Equal(uuid.Nil, transmission_id_uuid) {
			err_input = errors.New("invalid transmission_id")
			break
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
			err_input = errors.New("not found required parameter listening_socket_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		listening_socket_id_uuid := uuid.FromStringOrNil(listening_socket_id)
		if uuid.Equal(uuid.Nil, listening_socket_id_uuid) {
			err_input = errors.New("invalid listening_socket_id")
			break
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
			err_input = errors.New("not found required parameter connected_socket_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		connected_socket_id_uuid := uuid.FromStringOrNil(connected_socket_id)
		if uuid.Equal(uuid.Nil, connected_socket_id_uuid) {
			err_input = errors.New("invalid connected_socket_id")
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
			err_input = errors.New("not found required parameter connected_socket_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		connected_socket_id_uuid := uuid.FromStringOrNil(connected_socket_id)
		if uuid.Equal(uuid.Nil, connected_socket_id_uuid) {
			err_input = errors.New("invalid connected_socket_id")
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

		connected_socket_id_uuid := uuid.FromStringOrNil(connected_socket_id)
		if uuid.Equal(uuid.Nil, connected_socket_id_uuid) {
			err_input = errors.New("invalid connected_socket_id")
			break
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
			err_input = errors.New("not found required parameter connected_socket_id")
			break
		}

		if err != nil {
			err_processing_internal = err
			break
		}

		connected_socket_id_uuid := uuid.FromStringOrNil(connected_socket_id)
		if uuid.Equal(uuid.Nil, connected_socket_id_uuid) {
			err_input = errors.New("invalid connected_socket_id")
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

		connected_socket_id_uuid := uuid.FromStringOrNil(connected_socket_id)
		if uuid.Equal(uuid.Nil, connected_socket_id_uuid) {
			err_input = errors.New("invalid connected_socket_id")
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

		connected_socket_id_uuid := uuid.FromStringOrNil(connected_socket_id)
		if uuid.Equal(uuid.Nil, connected_socket_id_uuid) {
			err_input = errors.New("invalid connected_socket_id")
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

func (self *ARPCSolutionNode) NewCall(
	call_id uuid.UUID,
	response_on uuid.UUID,
) error {
	msg := new(gojsonrpc2.Message)
	msg.Method = "NewCall"

	params := map[string]any{"call_id": call_id}

	if !uuid.Equal(response_on, uuid.Nil) {
		params["response_on"] = response_on
	}

	msg.Params = params

	return self.jrpc_node.SendNotification(msg)
}

func (self *ARPCSolutionNode) NewBuffer(
	buffer_id uuid.UUID,
) error {
	msg := new(gojsonrpc2.Message)
	msg.Method = "NewBuffer"

	msg.Params = map[string]any{"buffer_id": buffer_id.String()}

	return self.jrpc_node.SendNotification(msg)
}

func (self *ARPCSolutionNode) BufferUpdated(
	buffer_id uuid.UUID,
) error {
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferUpdated"

	msg.Params = map[string]any{"buffer_id": buffer_id.String()}

	return self.jrpc_node.SendNotification(msg)
}

func (self *ARPCSolutionNode) NewTransmission(
	tarnsmission_id uuid.UUID,
) error {
	msg := new(gojsonrpc2.Message)
	msg.Method = "NewTransmission"

	msg.Params = map[string]any{"tarnsmission_id": tarnsmission_id.String()}

	return self.jrpc_node.SendNotification(msg)
}

func (self *ARPCSolutionNode) NewSocket(
	listening_socket_id uuid.UUID,
) error {
	msg := new(gojsonrpc2.Message)
	msg.Method = "NewSocket"

	msg.Params = map[string]any{"listening_socket_id": listening_socket_id.String()}

	return self.jrpc_node.SendNotification(msg)
}

// ----------------------------------------
// Basic Calls
// ----------------------------------------

func (self *ARPCSolutionNode) subResultGetter01(
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
			return nil,
				false, false, errors.New("invalid message"), errors.New("protocol error")
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
	return nil,
		false, false, nil, errors.New("unknown error")
}

func (self *ARPCSolutionNode) CallGetList(
	response_timeout time.Duration,
) (
	buffer_id uuid.UUID,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallGetList"

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
		return uuid.Nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err :=
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if timedout || closed || result_err != nil || err != nil {
		buffer_id = uuid.Nil
		return
	}

	result_str, ok := result_any.(string)
	if !ok {
		return uuid.Nil,
			false, false, nil, errors.New("result must be uuid string")
	}

	result := uuid.FromStringOrNil(result_str)
	if uuid.Equal(uuid.Nil, result) {
		return uuid.Nil,
			false, false, nil, errors.New("invalid uuid")
	}

	return result, false, false, nil, nil

}

func (self *ARPCSolutionNode) CallGetArgCount(
	call_id uuid.UUID,
	response_timeout time.Duration,
) (
	res int,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallGetArgCount"
	msg.Params = map[string]any{"call_id": call_id.String()}

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

func (self *ARPCSolutionNode) CallGetArgValue(
	call_id uuid.UUID,
	index int,
	response_timeout time.Duration,
) (
	res *ARPCArgInfo,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallGetArgValue"
	msg.Params = map[string]any{
		"call_id":   call_id.String(),
		"arg_index": index,
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

	var result *ARPCArgInfo

	err = mapstructure.Decode(result_any, &result)
	if err != nil {
		return nil, false, false, nil, err
	}

	return result, false, false, nil, nil
}

func (self *ARPCSolutionNode) CallClose(
	call_id uuid.UUID,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

func (self *ARPCSolutionNode) BufferGetInfo(
	buffer_id uuid.UUID,
	response_timeout time.Duration,
) (
	info *ARPCBufferInfo,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferGetInfo"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.String(),
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

func (self *ARPCSolutionNode) BufferGetItemsCount(
	buffer_id uuid.UUID,
	response_timeout time.Duration,
) (
	count int,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferGetItemsCount"
	msg.Params = buffer_id.String()
	msg.Params = map[string]any{
		"buffer_id": buffer_id.String(),
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

func (self *ARPCSolutionNode) BufferGetItemsIds(
	buffer_id uuid.UUID,
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
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferGetItemsIds"
	msg.Params = buffer_id.String()
	msg.Params = map[string]any{
		"buffer_id": buffer_id.String(),
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

func (self *ARPCSolutionNode) BufferGetItemsByIds(
	buffer_id uuid.UUID,
	ids []string,
	response_timeout time.Duration,
) (
	buffer_items []*ARPCBufferItem,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferGetItemsByIds"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.String(),
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

func (self *ARPCSolutionNode) BufferGetItemsFirstTime(
	buffer_id uuid.UUID,
	response_timeout time.Duration,
) (
	time_ time.Time,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferGetItemsFirstTime"
	msg.Params = buffer_id.String()
	msg.Params = map[string]any{
		"buffer_id": buffer_id.String(),
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

func (self *ARPCSolutionNode) BufferGetItemsLastTime(
	buffer_id uuid.UUID,
	response_timeout time.Duration,
) (
	time_ time.Time,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferGetItemsIds"
	msg.Params = buffer_id.String()
	msg.Params = map[string]any{
		"buffer_id": buffer_id.String(),
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

func (self *ARPCSolutionNode) BufferSubscribeOnUpdatesNotification(
	buffer_id uuid.UUID,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferSubscribeOnUpdatesNotification"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.String(),
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

func (self *ARPCSolutionNode) BufferUnsubscribeFromUpdatesNotification(
	buffer_id uuid.UUID,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferUnsubscribeFromUpdatesNotification"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.String(),
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

func (self *ARPCSolutionNode) BufferGetIsSubscribedOnUpdatesNotification(
	buffer_id uuid.UUID,
	response_timeout time.Duration,
) (
	result bool,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferGetIsSubscribedOnUpdatesNotification"
	msg.Params = map[string]any{
		"buffer_id": buffer_id.String(),
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

func (self *ARPCSolutionNode) BufferGetListSubscribedUpdatesNotifications(
	response_timeout time.Duration,
) (
	buffer_id uuid.UUID,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "BufferGetListSubscribedUpdatesNotifications"

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
		return uuid.Nil, false, false, nil, err
	}

	result_any, timedout, closed, result_err, err =
		self.subResultGetter01(timedout_sig, closed_sig, msg_sig)

	if err != nil {
		return uuid.Nil, false, false, nil, err
	}

	return false, false, nil, nil
}

func (self *ARPCSolutionNode) BufferBinaryGetSize(
	buffer_id uuid.UUID,
	response_timeout time.Duration,
) (
	size int,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

func (self *ARPCSolutionNode) BufferBinaryGetSlice(
	buffer_id uuid.UUID,
	start_index, end_index int,
	response_timeout time.Duration,
) (
	data []byte,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

// ----------------------------------------
// Broadcasts
// ----------------------------------------

func (self *ARPCSolutionNode) TransmissionGetList(
	response_timeout time.Duration,
) (
	buffer_id uuid.UUID,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

func (self *ARPCSolutionNode) TransmissionGetInfo(
	transmission_id uuid.UUID,
	response_timeout time.Duration,
) (
	info *ARPCTransmissionInfo,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

// ----------------------------------------
// Sockets
// ----------------------------------------

func (self *ARPCSolutionNode) SocketGetList(
	response_timeout time.Duration,
) (
	buffer_id uuid.UUID,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

func (self *ARPCSolutionNode) SocketOpen(
	listening_socket_id uuid.UUID,
	response_timeout time.Duration,
) (
	connected_socket_id uuid.UUID,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

func (self *ARPCSolutionNode) SocketRead(
	connected_socket_id uuid.UUID,
	try_read_size int,
	response_timeout time.Duration,
) (
	b []byte,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

func (self *ARPCSolutionNode) SocketWrite(
	connected_socket_id uuid.UUID,
	b []byte,
	response_timeout time.Duration,
) (
	n int,
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

func (self *ARPCSolutionNode) SocketClose(
	connected_socket_id uuid.UUID,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"connected_socket_id": connected_socket_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

func (self *ARPCSolutionNode) SocketSetDeadline(
	connected_socket_id uuid.UUID,
	t time.Time,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

func (self *ARPCSolutionNode) SocketSetReadDeadline(
	connected_socket_id uuid.UUID,
	t time.Time,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}

func (self *ARPCSolutionNode) SocketSetWriteDeadline(
	connected_socket_id uuid.UUID,
	t time.Time,
	response_timeout time.Duration,
) (
	timedout bool,
	closed bool,
	result_err error,
	err error,
) {
	msg := new(gojsonrpc2.Message)
	msg.Method = "CallClose"
	msg.Params = map[string]any{
		"call_id": call_id.String(),
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

	if err != nil {
		return false, false, nil, err
	}

	return false, false, nil, nil
}
