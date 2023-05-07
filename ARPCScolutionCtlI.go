package goarpcsolution

import (
	"net"
	"time"

	uuid "github.com/satori/go.uuid"
)

// about UUID is:
//
// it is adviced to keep in mind, what UUIDs of
// Calls, Buffers, Transmissions, available and connected sockets
// does not share space and resources can have identicall UUID, thou
// belonging to different things.
//
// once more:
// we have few groups (spaces) of separate IDs. those groups are:
// 1. Calls
// 2. Buffers
// 3. Transmissions
// 4. Available Sockets
// 5. Connected sockets
//
// so if You got, say, 'aaaa' UUID indicating an Available Socket,
// this doesn't mean 'you can't get same aaaa UUID for opened socket'

type ARPCSolutionCtlI interface {

	// ----------------------------------------
	// Notifications
	// ----------------------------------------

	// inform node about new call availability
	NewCall(
		call_id uuid.UUID,
		response_on uuid.UUID,
	)

	// inform node about new buffer availability
	NewBuffer(
		buffer_id uuid.UUID,
	)

	// you have to subscribe to bufffer updates to receive this
	// notifications (see BufferSubscribeOnUpdatesNotification())
	BufferUpdated(
		buffer_id uuid.UUID,
	)

	// inform node about new transmission availability
	NewTransmission(
		tarnsmission_id uuid.UUID,
	)

	// inform node about new socket availability
	NewSocket(
		listening_socket_id uuid.UUID,
	)

	// ----------------------------------------
	// Basic Calls
	// ----------------------------------------

	// generates list of active calls and returns it's buffer_id
	CallGetList() (
		buffer_id uuid.UUID,
		err_processing_not_internal, err_processing_internal error,
	)

	// returns int count of attributes of call
	CallGetArgCount(
		call_id uuid.UUID,
	) (
		res int,
		err_processing_not_internal, err_processing_internal error,
	)

	// returns ARPCArgInfo for selected arg of selected call
	CallGetArgValue(
		call_id uuid.UUID,
		index int,
	) (
		info *ARPCArgInfo,
		err_processing_not_internal, err_processing_internal error,
	)

	// inform node what call isn't longer needed
	CallClose(
		call_id uuid.UUID,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	// ----------------------------------------
	// Buffers
	// ----------------------------------------

	// get information about bufffer
	BufferGetInfo(
		buffer_id uuid.UUID,
	) (
		info *ARPCBufferInfo,
		err_processing_not_internal, err_processing_internal error,
	)

	// get count of items in buffer.
	// note: for binary buffers use BufferBinary* functions to get
	// bytes count and bytes slices
	BufferGetItemsCount(
		buffer_id uuid.UUID,
	) (
		count int,
		err_processing_not_internal, err_processing_internal error,
	)

	// get exact ids of buffer items using Buffer Item Specifiers
	BufferGetItemsIds(
		buffer_id uuid.UUID,
		first_spec, last_spec *ARPCBufferItemSpecifier,
		// TODO: do we need include_last?
		// include_last bool,
	) (
		ids []string,
		err_processing_not_internal, err_processing_internal error,
	)

	// returns Times of specified buffer items
	BufferGetItemsTimesByIds(
		buffer_id uuid.UUID,
		ids []string,
	) (
		times []time.Time,
		err_processing_not_internal, err_processing_internal error,
	)

	// retrive actual buffer items with payloads
	BufferGetItemsByIds(
		buffer_id uuid.UUID,
		ids []string,
	) (
		buffer_items []*ARPCBufferItem,
		err_processing_not_internal, err_processing_internal error,
	)

	// next two functions to enchance transmission updates rereival

	BufferGetItemsFirstTime(
		buffer_id uuid.UUID,
	) (
		time time.Time,
		err_processing_not_internal, err_processing_internal error,
	)

	BufferGetItemsLastTime(
		buffer_id uuid.UUID,
	) (
		time time.Time,
		err_processing_not_internal, err_processing_internal error,
	)

	// receive notifications on buffer updates
	BufferSubscribeOnUpdatesNotification(
		buffer_id uuid.UUID,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	BufferUnsubscribeFromUpdatesNotification(
		buffer_id uuid.UUID,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	BufferGetIsSubscribedOnUpdatesNotification(
		buffer_id uuid.UUID,
	) (
		bool,
		err_processing_not_internal, err_processing_internal error,
	)

	BufferGetListSubscribedUpdatesNotifications() (
		buffer_id uuid.UUID,
		err_processing_not_internal, err_processing_internal error,
	)

	// ---v--- BufferBinary ---v---

	// BufferBinary functions works only if buffer in binary mode

	BufferBinaryGetSize(
		buffer_id uuid.UUID,
	) (
		size int,
		err_processing_not_internal, err_processing_internal error,
	)

	BufferBinaryGetSlice(
		buffer_id uuid.UUID,
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
		buffer_id uuid.UUID,
		err_processing_not_internal, err_processing_internal error,
	)

	TransmissionGetInfo(
		transmission_id uuid.UUID,
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
		buffer_id uuid.UUID,
		err_processing_not_internal, err_processing_internal error,
	)

	SocketOpen(
		listening_socket_id uuid.UUID,
	) (
		connected_socket_id uuid.UUID,
		err_processing_not_internal, err_processing_internal error,
	)

	SocketRead(
		connected_socket_id uuid.UUID,
		try_read_size int,
	) (
		b []byte,
		err_processing_not_internal, err_processing_internal error,
	)

	SocketWrite(
		connected_socket_id uuid.UUID,
		b []byte,
	) (
		n int,
		err_processing_not_internal, err_processing_internal error,
	)

	SocketClose(
		connected_socket_id uuid.UUID,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	SocketSetDeadline(
		connected_socket_id uuid.UUID,
		t time.Time,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	SocketSetReadDeadline(
		connected_socket_id uuid.UUID,
		t time.Time,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	SocketSetWriteDeadline(
		connected_socket_id uuid.UUID,
		t time.Time,
	) (
		err_processing_not_internal, err_processing_internal error,
	)

	// get local representation of remote connected socket
	SocketGetConn(connected_socket_id uuid.UUID) (net.Conn, error)
}
