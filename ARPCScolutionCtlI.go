package goarpcsolution

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type ARPCSolutionCtlI interface {
	// ----------------------------------------
	// Notifications
	// ----------------------------------------

	NewCall(
		call_id uuid.UUID,
		response_on uuid.UUID,
	)

	NewBuffer(id uuid.UUID)

	NewBroadcast(id uuid.UUID)

	NewPort(id uuid.UUID)

	// ----------------------------------------
	// Basic Calls
	// ----------------------------------------

	CallGetList() (buffer_id uuid.UUID, reply_err, err error)

	CallGetArgCount(call_id uuid.UUID) (
		res int,
		reply_err, err error,
	)

	CallGetArgValue(call_id uuid.UUID, index uint) (
		info *ARPCArgInfo,
		reply_err, err error,
	)

	CallClose(call_id uuid.UUID)

	// ----------------------------------------
	// Buffers
	// ----------------------------------------

	BufferGetInfo(id uuid.UUID) (
		info *ARPCBufferInfo,
		reply_err, err error,
	)

	BufferGetItemsFirstTime(id uuid.UUID) (
		time time.Time,
		reply_err, err error,
	)

	BufferGetItemsLastTime(id uuid.UUID) (
		time time.Time,
		reply_err, err error,
	)

	BufferGetItemsByIds(ids []uuid.UUID) (
		buffer_items []*ARPCBufferItem,
		err error,
	)

	BufferGetItemsStartingFromId(id uuid.UUID) (
		buffer_items []*ARPCBufferItem,
		reply_err, err error,
	)

	BufferGetItemsByFirstAndLastId(
		first_id uuid.UUID,
		last_id uuid.UUID,
	) (
		buffer_items []*ARPCBufferItem,
		reply_err, err error,
	)

	BufferGetItemsByFirstAndBeforeId(
		first_id uuid.UUID,
		last_id uuid.UUID,
	) (
		buffer_items []*ARPCBufferItem,
		reply_err, err error,
	)

	BufferGetItemsByFirstAndLastTime(
		first_time time.Time,
		last_time time.Time,
	) (
		buffer_items []*ARPCBufferItem,
		reply_err, err error,
	)

	BufferGetItemsByFirstAndBeforeTime(
		first_time time.Time,
		last_time time.Time,
	) (
		buffer_items []*ARPCBufferItem,
		reply_err, err error,
	)

	BufferSubscribeOnUpdatesIndicator(buffer_id uuid.UUID) (
		reply_err, err error,
	)

	// ----------------------------------------
	// Broadcasts
	// ----------------------------------------

	BroadcastGetIdList() (
		ids []uuid.UUID,
		reply_err, err error,
	)

	BroadcastGetInfo(broadcast_id uuid.UUID) (
		info *ARPCBroadcastInfo,
		reply_err, err error,
	)

	BroadcastSubscribeOnNew() error

	// ----------------------------------------
	// Ports
	// ----------------------------------------

	PortGetListeningIdList() (
		ids []uuid.UUID,
		reply_err, err error,
	)

	PortOpen(id uuid.UUID) (
		reply_err, err error,
	)

	PortRead(id uuid.UUID, b []byte) (
		n int,
		reply_err, err error,
	)

	PortWrite(id uuid.UUID, b []byte) (
		n int,
		reply_err, err error,
	)

	PortClose(id uuid.UUID) (
		reply_err, err error,
	)
}
