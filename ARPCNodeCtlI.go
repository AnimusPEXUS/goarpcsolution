package goarpcsolution

import (
	"net"
	"time"

	"github.com/AnimusPEXUS/gojsonrpc2"
	"github.com/AnimusPEXUS/gouuidtools"
)

type ARPCNodeCtlI interface {

	// node will call this func on node init and
	// pass pointer to self to controller
	SetNode(node *ARPCNode)

	// return true if node can't be used. if true, some Node
	// functions will panic on attemt to call them.
	// this is feature for developers.
	IsStopFlag() bool

	// this is just like OnRequestCB on JSONRPC2Node, with exception:
	// to get here, (on jsonrpc2 level) message should be sent with "s:"
	// prefix in method name. but this msg.Method will be with "s:" removed.
	SimpleRequest(msg *gojsonrpc2.Message) (error, error)

	// get local representation of remote connected socket
	SocketGetConn(connected_socket_id *gouuidtools.UUID) (net.Conn, error)

	Close()

	// =========== vvvvvvvvvvvvvvvvvvvvvvv ===========
	// ----------- node functions handlers -----------
	// =========== vvvvvvvvvvvvvvvvvvvvvvv ===========

	// ----------------------------------------
	// Notifications
	// ----------------------------------------

	// inform node about new call availability
	NewCall(
		call_id *gouuidtools.UUID,
		response_on *gouuidtools.UUID,
	)

	// inform node about new buffer availability
	// NewBuffer(
	// 	buffer_id *gouuidtools.UUID,
	// )

	// you have to subscribe to bufffer updates to receive this
	// notifications (see BufferSubscribeOnUpdatesNotification())
	BufferUpdated(
		buffer_id *gouuidtools.UUID,
	)

	// inform node about new transmission availability
	NewTransmission(
		transmission_id *gouuidtools.UUID,
	)

	// inform node about new socket availability
	// NewSocket(
	// 	listening_socket_id *gouuidtools.UUID,
	// )

	// ----------------------------------------
	// Basic Calls
	// ----------------------------------------

	// generates list of active calls and returns it's buffer_id
	CallGetList() (
		buffer_id *gouuidtools.UUID,
		err_processing_not_internal, err_processing_internal error,
	)

	// returns name of method which call is calling
	CallGetName(
		call_id *gouuidtools.UUID,
	) (
		name string,
		err_processing_not_internal, err_processing_internal error,
	)

	// returns int count of attributes of call
	CallGetArgCount(
		call_id *gouuidtools.UUID,
	) (
		res int,
		err_processing_not_internal, err_processing_internal error,
	)

	// returns ARPCArgInfo for selected arg of selected call
	CallGetArgValues(
		call_id *gouuidtools.UUID,
		first, last int,
	) (
		res []*ARPCArgInfo,
		err_processing_not_internal, err_processing_internal error,
	)

	// inform node what call isn't longer needed
	CallClose(
		call_id *gouuidtools.UUID,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	// ----------------------------------------
	// Buffers
	// ----------------------------------------

	// get information about bufffer
	BufferGetInfo(
		buffer_id *gouuidtools.UUID,
	) (
		info *ARPCCallForJSON,
		err_processing_not_internal, err_processing_internal error,
	)

	// get count of items in buffer.
	// note: for binary buffers use BufferBinary* functions to get
	// bytes count and bytes slices
	BufferGetItemsCount(
		buffer_id *gouuidtools.UUID,
	) (
		count int,
		err_processing_not_internal, err_processing_internal error,
	)

	// get exact ids of buffer items using Buffer Item Specifiers
	BufferGetItemsIds(
		buffer_id *gouuidtools.UUID,
		first_spec, last_spec *ARPCBufferItemSpecifier,
		// TODO: do we need include_last?
		// include_last bool,
	) (
		ids []string,
		err_processing_not_internal, err_processing_internal error,
	)

	// returns Times of specified buffer items
	BufferGetItemsTimesByIds(
		buffer_id *gouuidtools.UUID,
		ids []string,
	) (
		times []time.Time,
		err_processing_not_internal, err_processing_internal error,
	)

	// retrive actual buffer items with payloads
	BufferGetItemsByIds(
		buffer_id *gouuidtools.UUID,
		ids []string,
	) (
		buffer_items []*ARPCBufferItem,
		err_processing_not_internal, err_processing_internal error,
	)

	// next two functions to enchance transmission updates rereival

	BufferGetItemsFirstTime(
		buffer_id *gouuidtools.UUID,
	) (
		time_ time.Time,
		err_processing_not_internal, err_processing_internal error,
	)

	BufferGetItemsLastTime(
		buffer_id *gouuidtools.UUID,
	) (
		time_ time.Time,
		err_processing_not_internal, err_processing_internal error,
	)

	// receive notifications on buffer updates
	BufferSubscribeOnUpdatesNotification(
		buffer_id *gouuidtools.UUID,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	BufferUnsubscribeFromUpdatesNotification(
		buffer_id *gouuidtools.UUID,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	BufferGetIsSubscribedOnUpdatesNotification(
		buffer_id *gouuidtools.UUID,
	) (
		r bool,
		err_processing_not_internal, err_processing_internal error,
	)

	BufferGetListSubscribedUpdatesNotifications() (
		buffer_id *gouuidtools.UUID,
		err_processing_not_internal, err_processing_internal error,
	)

	// ---v--- BufferBinary ---v---

	// BufferBinary functions works only if buffer in binary mode

	BufferBinaryGetSize(
		buffer_id *gouuidtools.UUID,
	) (
		size int,
		err_processing_not_internal, err_processing_internal error,
	)

	BufferBinaryGetSlice(
		buffer_id *gouuidtools.UUID,
		start_index, end_index int,
	) (
		data []byte,
		err_processing_not_internal, err_processing_internal error,
	)

	// ---^--- BufferBinary ---^---

	// ----------------------------------------
	// Transmissions
	// ----------------------------------------

	// like the CallGetList(), same here but for transmissions:
	// result is buffer with ids.
	TransmissionGetList() (
		buffer_id *gouuidtools.UUID,
		err_processing_not_internal, err_processing_internal error,
	)

	TransmissionGetInfo(
		transmission_id *gouuidtools.UUID,
	) (
		info *ARPCTransmissionInfo,
		err_processing_not_internal, err_processing_internal error,
	)

	// ----------------------------------------
	// Sockets
	// ----------------------------------------

	// note: socket here - mimics Go's net.Conn interface

	// just like the CallGetList() or TransmissionGetList(),
	// same here but for sockets:
	// result is buffer with ids.
	// except: only sockets available to be opened are listed -
	// already connected sockets - are not returned
	SocketGetList() (
		buffer_id *gouuidtools.UUID,
		err_processing_not_internal, err_processing_internal error,
	)

	SocketOpen(
		listening_socket_id *gouuidtools.UUID,
	) (
		connected_socket_id *gouuidtools.UUID,
		err_processing_not_internal, err_processing_internal error,
	)

	SocketRead(
		connected_socket_id *gouuidtools.UUID,
		try_read_size int,
	) (
		b []byte,
		err_processing_not_internal, err_processing_internal error,
	)

	SocketWrite(
		connected_socket_id *gouuidtools.UUID,
		b []byte,
	) (
		n int,
		err_processing_not_internal, err_processing_internal error,
	)

	SocketClose(
		connected_socket_id *gouuidtools.UUID,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	SocketSetDeadline(
		connected_socket_id *gouuidtools.UUID,
		t time.Time,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	SocketSetReadDeadline(
		connected_socket_id *gouuidtools.UUID,
		t time.Time,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	SocketSetWriteDeadline(
		connected_socket_id *gouuidtools.UUID,
		t time.Time,
	) (
		err_processing_not_internal, err_processing_internal error,
	)
}
