package goarpcsolution

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// about UUID is:
//
// it is adviced to keep in mind, what UUIDs of
// Calls, Buffers, Transmissions, available and connected ports
// does not share space and resources can have identicall UUID, thou
// belonging to different things.
//
// once more:
// we have few groups (spaces) of separate IDs. those groups are:
// 1. Calls
// 2. Buffers
// 3. Transmissions
// 4. Available Ports
// 5. Connected ports
//
// so if You got, say, 'aaaa' UUID indicating an Available Port,
// this doesn't mean 'you can't get same aaaa UUID for opened port'

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

	// inform node about new transmission availability
	NewTransmission(
		tarnsmission_id uuid.UUID,
	)

	// inform node about new port availability
	NewPort(
		port_id uuid.UUID,
	)

	// ----------------------------------------
	// Basic Calls
	// ----------------------------------------

	// generates list of active calls and returns it's buffer_id
	CallGetList() (
		buffer_id uuid.UUID,
		reply_err, err error,
	)

	// returns int count of attributes of call
	CallGetArgCount(
		call_id uuid.UUID,
	) (
		res int,
		reply_err, err error,
	)

	// returns ARPCArgInfo for selected arg of selected call
	CallGetArgInfo(
		call_id uuid.UUID,
		index int,
	) (
		info *ARPCArgInfo,
		reply_err, err error,
	)

	// inform node what call isn't longer needed
	CallClose(
		call_id uuid.UUID,
	) (
		reply_err, err error,
	)

	// ----------------------------------------
	// Buffers
	// ----------------------------------------

	// get information about bufffer
	BufferGetInfo(
		buffer_id uuid.UUID,
	) (
		info *ARPCBufferInfo,
		reply_err, err error,
	)

	// get count of items in buffer.
	// note: for binary buffers use BufferBinary* functions to get
	// bytes count and bytes slices
	BufferGetItemsCount(
		buffer_id uuid.UUID,
	) (
		count int,
		reply_err, err error,
	)

	// get exact ids of buffer items using Buffer Item Specifiers
	BufferGetItemsIds(
		buffer_id uuid.UUID,
		first_spec, last_spec *ARPCBufferItemSpecifier,
		include_last bool,
	) (
		ids []string,
		reply_err, err error,
	)

	// returns Times of specified buffer items
	BufferGetItemsTimesByIds(
		buffer_id uuid.UUID,
		ids []string,
	) (
		times []time.Time,
		reply_err, err error,
	)

	// retrive actual buffer items with payloads
	BufferGetItemsByIds(
		buffer_id uuid.UUID,
		ids []string,
	) (
		buffer_items []*ARPCBufferItem,
		reply_err, err error,
	)

	// next two functions to enchance transmission updates rereival

	BufferGetItemsFirstTime(
		buffer_id uuid.UUID,
	) (
		time time.Time,
		reply_err, err error,
	)

	BufferGetItemsLastTime(
		buffer_id uuid.UUID,
	) (
		time time.Time,
		reply_err, err error,
	)

	// receive notifications on buffer updates
	BufferSubscribeOnUpdatesNotification(
		buffer_id uuid.UUID,
	) (
		reply_err, err error,
	)

	BufferUnsubscribeFromUpdatesNotification(
		buffer_id uuid.UUID,
	) (
		reply_err, err error,
	)

	// ---v--- BufferBinary ---v---

	// BufferBinary functions works only if buffer in binary mode

	BufferBinaryGetSize(
		buffer_id uuid.UUID,
	) (
		size int,
		reply_err, err error,
	)

	BufferBinaryGetSlice(
		buffer_id uuid.UUID,
		start, end int,
	) (
		data []byte,
		reply_err, err error,
	)

	// ---^--- BufferBinary ---^---

	// ----------------------------------------
	// Transmissions
	// ----------------------------------------

	// like the CallGetList(), same here but for transmissions:
	// result is buffer with ids.
	TransmissionGetList() (
		buffer_id uuid.UUID,
		reply_err, err error,
	)

	TransmissionGetInfo(
		broadcast_id uuid.UUID,
	) (
		info *ARPCTransmissionInfo,
		reply_err, err error,
	)

	// ----------------------------------------
	// Ports
	// ----------------------------------------

	// just like the CallGetList() or TransmissionGetList(),
	// same here but for ports:
	// result is buffer with ids.
	// except: only ports available to be opened are listed -
	// already connected ports - are not returned
	PortGetList() (
		buffer_id uuid.UUID,
		reply_err, err error,
	)

	PortOpen(listening_port_id uuid.UUID) (
		connected_port_id uuid.UUID,
		reply_err, err error,
	)

	// TODO: not sure we need explicit PortRead() function.
	// maybe it should be assumed 'write' to the buffer whould be it's
	// 'read' from the other side
	PortRead(connected_port_id uuid.UUID, b []byte) (
		n int,
		reply_err, err error,
	)

	PortWrite(connected_port_id uuid.UUID, b []byte) (
		n int,
		reply_err, err error,
	)

	PortClose(connected_port_id uuid.UUID) (
		reply_err, err error,
	)
}
